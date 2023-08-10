package query_test

import (
	"simple-ddd-cqrs/pkg/domain"
	"simple-ddd-cqrs/pkg/query"
	"simple-ddd-cqrs/service"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
)

func TestArticleQueryHandler_GetListArticle(t *testing.T) {
	db, mock := service.NewDBMock()
	articleRepo := domain.CreateArticleRepository(db)
	defer db.Close()

	var (
		ID      = 1
		author  = "Author test"
		title   = "Title test"
		body    = "Body test"
		created = time.Now()
	)

	queryStmt := "SELECT id, author, title, body, created FROM articles"

	rows := sqlmock.NewRows([]string{"id", "author", "title", "body", "created"}).
		AddRow(ID, author, title, body, created)

	mock.ExpectQuery(queryStmt).WillReturnRows(rows)

	articleQueryHandler := query.NewArticleQueryHandler(articleRepo)

	articles, err := articleQueryHandler.GetListArticle()

	assert.NotEmpty(t, articles)
	assert.NoError(t, err)
	assert.Len(t, articles, 1)
}

func TestArticleQueryHandler_GetArticleByID(t *testing.T) {
	db, mock := service.NewDBMock()
	articleRepo := domain.CreateArticleRepository(db)
	defer db.Close()

	var (
		ID      = 1
		author  = "Author test"
		title   = "Title test"
		body    = "Body test"
		created = time.Now()
	)

	queryStmt := "SELECT id, author, title, body, created FROM articles WHERE id = \\$1"

	rows := sqlmock.NewRows([]string{"id", "author", "title", "body", "created"}).
		AddRow(ID, author, title, body, created)

	mock.ExpectQuery(queryStmt).WithArgs(ID).WillReturnRows(rows)

	articleQueryHandler := query.NewArticleQueryHandler(articleRepo)

	article, err := articleQueryHandler.GetArticleByID(ID)
	assert.NotNil(t, article)
	assert.NoError(t, err)
}

func TestArticleQueryHandler_GetArticleFromCache(t *testing.T) {
	// mock Redis client
	client, mock := redismock.NewClientMock()

	articleQueryHandler := query.NewArticleQueryHandler(nil).WithRedis(client)

	// Mock Redis response for cached item
	mock.ExpectGet("article:123").
		SetVal(`{"ID": 123, "Author": "Test Author", "Title": "Test Title", "Body": "Test Body", "Created": "2023-07-29T15:00:00Z"}`)

	cachedArticle := articleQueryHandler.GetArticleFromCache(123)
	assert.NotNil(t, cachedArticle)
	assert.Equal(t, 123, cachedArticle.ID)
	assert.Equal(t, "Test Author", cachedArticle.Author)
	// ...assert other fields as needed

	// Test other scenarios...
}
