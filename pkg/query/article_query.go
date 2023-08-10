package query

import "simple-ddd-cqrs/pkg/domain"

// ArticleQueryHandler is an ArticleQueryHandler
type ArticleQueryHandler interface {
	GetListArticle(filters ...map[string]string) ([]*domain.ArticleModel, error)
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

// GetArticleByID
func (h *articleQueryHandler) GetListArticle(filters ...map[string]string) ([]*domain.ArticleModel, error) {
	articles, err := h.articleRepo.Get(filters...)
	if err != nil {
		return nil, err
	}

	return articles, nil
}
