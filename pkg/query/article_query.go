package query

import "simple-ddd-cqrs/pkg/domain"

// ArticleQueryHandler is an ArticleQueryHandler
type ArticleQueryHandler interface {
	GetListArticle(filters ...map[string]string) ([]*domain.ArticleModel, error)
	GetArticleByID(ID int) (*domain.ArticleModel, error)
}

// articleQueryHandler
type articleQueryHandler struct {
	articleRepo domain.ArticleRepository
}

// NewArticleQueryHandler is NewArticleQueryHandler
func NewArticleQueryHandler(articleRepo domain.ArticleRepository) ArticleQueryHandler {
	return &articleQueryHandler{
		articleRepo: articleRepo,
	}
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
	article, err := h.articleRepo.GetByID(ID)
	if err != nil {
		return nil, err
	}

	return article, nil
}
