package query

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"simple-ddd-cqrs/pkg/domain"

	"github.com/go-redis/redis/v8"
)

const articleRedisKeyPrefix = "article:"

// ArticleQueryHandler is an ArticleQueryHandler
type ArticleQueryHandler interface {
	GetListArticle(filters ...map[string]string) ([]*domain.ArticleModel, error)
	GetArticleByID(ID int) (*domain.ArticleModel, error)
	WithRedis(redis *redis.Client) *articleQueryHandler
	GetArticleFromCache(ID int) *domain.ArticleModel
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

	return articles, nil
}

// GetArticleByID
func (h *articleQueryHandler) GetArticleByID(ID int) (*domain.ArticleModel, error) {
	// Check if the item is available in the Redis cache
	cachedArticle := h.GetArticleFromCache(ID)
	if cachedArticle != nil {
		return cachedArticle, nil
	}

	article, err := h.articleRepo.GetByID(ID)
	if err != nil {
		return nil, err
	}

	return article, nil
}

// GetArticleFromCache
func (h *articleQueryHandler) GetArticleFromCache(ID int) *domain.ArticleModel {
	if h.redis == nil {
		return nil
	}

	ctx := context.Background()
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

		return &cachedArticle
	}

	// error redis
	if err != redis.Nil {
		log.Println("Error while Get from redis, ", err)
	}

	return nil
}
