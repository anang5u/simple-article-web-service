package query

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"simple-ddd-cqrs/config"
	"simple-ddd-cqrs/pkg/domain"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

const articleRedisKeyPrefix = "article:"

// ArticleQueryHandler is an ArticleQueryHandler
type ArticleQueryHandler interface {
	GetListArticle(filters ...map[string]string) ([]*domain.ArticleModel, error)
	GetArticleByID(ID int) (*domain.ArticleModel, error)
	WithRedis(redis *redis.Client) *articleQueryHandler
	GetArticleFromCache(ctx context.Context, ID int) *domain.ArticleModel
	StoreArticleIntoCache(ctx context.Context, articles []*domain.ArticleModel) error
}

// articleQueryHandler
type articleQueryHandler struct {
	articleRepo domain.ArticleRepository
	redis       *redis.Client
}

// NewArticleQueryHandler is NewArticleQueryHandler
func NewArticleQueryHandler(articleRepo domain.ArticleRepository) ArticleQueryHandler {
	return &articleQueryHandler{
		articleRepo: articleRepo,
	}
}

// WithRedis
func (h *articleQueryHandler) WithRedis(redis *redis.Client) *articleQueryHandler {
	h.redis = redis
	return h
}

// GetListArticle
func (h *articleQueryHandler) GetListArticle(filters ...map[string]string) ([]*domain.ArticleModel, error) {
	articles, err := h.articleRepo.Get(filters...)
	if err != nil {
		return nil, err
	}

	// store articles into cache
	// background process
	go func() {
		if err := h.StoreArticleIntoCache(context.Background(), articles); err != nil {
			log.Println("Error on GetListArticle while store articles into cache. ", err)
		}
	}()

	return articles, nil
}

// GetArticleByID
func (h *articleQueryHandler) GetArticleByID(ID int) (*domain.ArticleModel, error) {
	ctx := context.Background()

	// Check if the item is available in the Redis cache
	cachedArticle := h.GetArticleFromCache(ctx, ID)
	if cachedArticle != nil {
		return cachedArticle, nil
	}

	article, err := h.articleRepo.GetByID(ID)
	if err != nil {
		return nil, err
	}

	// store article into cache
	// background process
	go func() {
		articles := []*domain.ArticleModel{}
		articles = append(articles, article)
		if err := h.StoreArticleIntoCache(ctx, articles); err != nil {
			log.Println("Error on GetArticleByID while store article into cache. ", err)
		}
	}()

	return article, nil
}

// GetArticleFromCache
func (h *articleQueryHandler) GetArticleFromCache(ctx context.Context, ID int) *domain.ArticleModel {
	if h.redis == nil {
		return nil
	}

	redisKey := fmt.Sprintf("%s%d", articleRedisKeyPrefix, ID)

	// Check if the item is available in the Redis cache
	cachedItemJSON, err := h.redis.Get(ctx, redisKey).Result()
	if err == nil {
		// Item found in the cache, unmarshal JSON to domain.Item
		var cachedArticle domain.ArticleModel
		err := json.Unmarshal([]byte(cachedItemJSON), &cachedArticle)
		if err != nil {
			log.Println("Error while GetArticleFromCache Unmarshal, ", err)
			return nil
		}

		// debug for development only
		if config.Get("APP_ENVIRONMENT") == "development" {
			log.Println("cached article found with id ", ID)
		}

		return &cachedArticle
	}

	// error redis
	if err != redis.Nil {
		log.Println("Error while Get from redis, ", err)
	}

	return nil
}

// StoreArticleIntoCache
func (h *articleQueryHandler) StoreArticleIntoCache(ctx context.Context, articles []*domain.ArticleModel) error {
	if h.redis == nil {
		return nil
	}

	for _, article := range articles {
		redisKey := fmt.Sprintf("%s%d", articleRedisKeyPrefix, article.ID)

		// Check if the article is still available in the Redis cache
		cachedArticle := h.GetArticleFromCache(ctx, article.ID)
		if cachedArticle != nil {
			continue // skip for caching
		}

		// Marshal the item to JSON and store it in the Redis cache
		itemJSON, err := json.Marshal(article)
		if err != nil {
			return err
		}

		// Set the item in the Redis cache with an expiration time
		expirationTime, err := strconv.Atoi(config.Get("CACHE_ARTICLE_EXP_TIME"))
		if err != nil {
			log.Println("Error strconv.Atoi CACHE_ARTICLE_EXP_TIME. ", err)
			expirationTime = 1
		}
		h.redis.Set(ctx, redisKey, itemJSON, (time.Duration(expirationTime) * time.Minute))

		// debug for development only
		if config.Get("APP_ENVIRONMENT") == "development" {
			log.Println("article stored into cache. id ", article.ID)
		}
	}

	return nil
}
