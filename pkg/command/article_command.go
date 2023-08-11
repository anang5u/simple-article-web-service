package command

import (
	"simple-ddd-cqrs/helper"
	"simple-ddd-cqrs/pkg/domain"
	"time"
)

// ArticleCommandHandler
type ArticleCommandHandler interface {
	CreateArticle(author, title, body string, created time.Time) (*domain.ArticleModel, error)
}

// articleCommandHandler
type articleCommandHandler struct {
	articleRepo domain.ArticleRepository
}

// NewArticleCommandHandler
func NewArticleCommandHandler(articleRepo domain.ArticleRepository) ArticleCommandHandler {
	return &articleCommandHandler{
		articleRepo: articleRepo,
	}
}

// CreateArticle Handler
func (h *articleCommandHandler) CreateArticle(author, title, body string, created time.Time) (*domain.ArticleModel, error) {
	article := &domain.ArticleModel{
		Author:  author,
		Title:   title,
		Body:    body,
		Created: created,
	}
	articleName := helper.Slugify(title)

	err := h.articleRepo.Create(article.Author, article.Title, articleName, article.Body, article.Created)
	if err != nil {
		return nil, err
	}

	return article, nil
}
