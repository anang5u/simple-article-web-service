package domain

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
)

// Article model
type ArticleModel struct {
	ID      int       `json:"id"`
	Author  string    `json:"author"`
	Title   string    `json:"title"`
	Name    string    `json:"name"` // for improve SEO Optimization
	Body    string    `json:"body"`
	Created time.Time `json:"created"`
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

// ArticleRepository
type ArticleRepository interface {
	Create(author, title, name, body string, created time.Time) error
	Get(filters ...map[string]string) ([]*ArticleModel, error)
	GetByID(ID int) (*ArticleModel, error)
	BuildFilterValues(filter map[string]string) (string, []interface{})
}

// article
type article struct {
	DB *sql.DB
}

// CreateArticleRepository
func CreateArticleRepository(db *sql.DB) ArticleRepository {
	return &article{
		DB: db,
	}
}

// Create perform to create new article
func (r *article) Create(author, title, name, body string, created time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), cTX_TIMEOUT*time.Second)
	defer cancel()

	query := "INSERT INTO articles (author, title, name, body, created) VALUES ($1, $2, $3, $4, $5)"
	stmt, err := r.DB.PrepareContext(ctx, query)
	if err != nil {
		log.Println("Error while Create PrepareContext article: ", err)
		return errCreateArticle
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, author, title, name, body, created)
	if err != nil {
		log.Println("Error while Create ExecContext article: ", err)
		return errCreateArticle
	}

	return nil
}

// Get perform to get list of articles
func (r *article) Get(filters ...map[string]string) ([]*ArticleModel, error) {
	articles := make([]*ArticleModel, 0)

	ctx, cancel := context.WithTimeout(context.Background(), cTX_TIMEOUT*time.Second)
	defer cancel()

	// filter article by
	sCond := ""
	var filterValues []interface{}
	if len(filters) > 0 {
		sCond, filterValues = r.BuildFilterValues(filters[0])
	}

	query := fmt.Sprintf(`
		SELECT 
			id, 
			author, 
			title, 
			name,
			body, 
			created 
		FROM articles %s 
		ORDER BY created 
		DESC LIMIT %d`, sCond, aRTICLES_PER_PAGE)

	rows, err := r.DB.QueryContext(ctx, query, filterValues...)
	if err != nil {
		log.Printf("Error while Get query article. Err: %s Query: %s\n", err.Error(), query)
		return nil, errArticleNotFound
	}
	defer rows.Close()

	for rows.Next() {
		result := new(ArticleModel)
		err = rows.Scan(
			&result.ID,
			&result.Author,
			&result.Title,
			&result.Name,
			&result.Body,
			&result.Created,
		)

		if err != nil {
			log.Println("Error while Fetch rows for article: ", err)
			return nil, errArticleNotFound
		}

		articles = append(articles, result)
	}

	// fix empty article
	if len(articles) == 0 {
		return nil, errArticleNotFound
	}

	return articles, nil
}

// GetByID
func (r *article) GetByID(ID int) (*ArticleModel, error) {
	result := new(ArticleModel)

	ctx, cancel := context.WithTimeout(context.Background(), cTX_TIMEOUT*time.Second)
	defer cancel()

	err := r.DB.QueryRowContext(ctx, "SELECT id, author, title, name, body, created FROM articles WHERE id = $1", ID).
		Scan(
			&result.ID,
			&result.Author,
			&result.Title,
			&result.Name,
			&result.Body,
			&result.Created,
		)
	if err != nil {
		log.Printf("Error while get article by ID: %d, Err: %s\n", ID, err.Error())
		return nil, errArticleNotFound
	}

	return result, nil
}

// BuildFilterValues
func (r *article) BuildFilterValues(filter map[string]string) (string, []interface{}) {
	sFilter := ""
	var filterValues []interface{}
	num := 1

	// check filter is appear
	if len(filter) > 0 {
		sFilter = "WHERE"
	}

	// filter by keyword search title
	if title, ok := filter["title"]; ok {
		sFilter = fmt.Sprintf("%s LOWER(title) LIKE '%%' || $%d || '%%' AND", sFilter, num)
		filterValues = append(filterValues, strings.ToLower(title))

		num++
	}

	// filter by keyword search body
	if body, ok := filter["body"]; ok {
		sFilter = fmt.Sprintf("%s LOWER(body) LIKE '%%' || $%d || '%%' AND", sFilter, num)
		filterValues = append(filterValues, strings.ToLower(body))

		num++
	}

	// filter by author
	if author, ok := filter["author"]; ok {
		sFilter = fmt.Sprintf("%s LOWER(author) = $%d AND", sFilter, num)
		filterValues = append(filterValues, strings.ToLower(author))

		num++
	}

	// trim end of string sFilter AND
	sFilter = strings.TrimRight(sFilter, "AND")

	return sFilter, filterValues
}
