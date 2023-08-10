package domain

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
)

// Article model
type ArticleModel struct {
	ID      int
	Author  string
	Title   string
	Body    string
	Created time.Time
}

// const default environment
const (
	cTX_TIMEOUT       = 5  // db context timeout in second
	aRTICLES_PER_PAGE = 20 // list articles fetch per page
)

// error message
var (
	errCreateArticle   = errors.New("create article failed")
	errArticleNotFound = errors.New("article not found")
)

// articleRepository
type articleRepository interface {
	Create(author, title, body string, created time.Time) error
	Get() ([]*ArticleModel, error)
	GetByID(ID int) (*ArticleModel, error)
}

// article
type article struct {
	DB *sql.DB
}

// CreateArticleRepository
func CreateArticleRepository(db *sql.DB) articleRepository {
	return &article{
		DB: db,
	}
}

// Create perform to create new article
func (r *article) Create(author, title, body string, created time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), cTX_TIMEOUT*time.Second)
	defer cancel()

	query := "INSERT INTO articles (author, title, body, created) VALUES ($1, $2, $3, $4)"
	stmt, err := r.DB.PrepareContext(ctx, query)
	if err != nil {
		log.Println("Error while Create PrepareContext article: ", err)
		return errCreateArticle
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, author, title, body, created)
	if err != nil {
		log.Println("Error while Create ExecContext article: ", err)
		return errCreateArticle
	}

	return nil
}

// Get perform to get list of articles
func (r *article) Get() ([]*ArticleModel, error) {
	articles := make([]*ArticleModel, 0)

	ctx, cancel := context.WithTimeout(context.Background(), cTX_TIMEOUT*time.Second)
	defer cancel()

	query := fmt.Sprintf(`SELECT id, author, title, body, created FROM articles ORDER BY created DESC LIMIT %d`, aRTICLES_PER_PAGE)
	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		log.Println("Error while Get query article: ", err)
		return nil, errArticleNotFound
	}
	defer rows.Close()

	for rows.Next() {
		result := new(ArticleModel)
		err = rows.Scan(
			&result.ID,
			&result.Author,
			&result.Title,
			&result.Body,
			&result.Created,
		)

		if err != nil {
			log.Println("Error while Fetch rows for article: ", err)
			return nil, errArticleNotFound
		}

		articles = append(articles, result)
	}

	return articles, nil
}

// GetByID
func (r *article) GetByID(ID int) (*ArticleModel, error) {
	result := new(ArticleModel)

	ctx, cancel := context.WithTimeout(context.Background(), cTX_TIMEOUT*time.Second)
	defer cancel()

	err := r.DB.QueryRowContext(ctx, "SELECT id, author, title, body, created FROM articles WHERE id = $1", ID).
		Scan(
			&result.ID,
			&result.Author,
			&result.Title,
			&result.Body,
			&result.Created,
		)
	if err != nil {
		log.Println("Error while article by ID: ", err)
		return nil, errArticleNotFound
	}
	return result, nil

}
