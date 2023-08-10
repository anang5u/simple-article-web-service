package query_test

import (
	"simple-ddd-cqrs/pkg/domain"
	"simple-ddd-cqrs/pkg/query"
	"simple-ddd-cqrs/service"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func Test_ArticleQueryGet(t *testing.T) {
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

func Test_GetQueryArticleByID(t *testing.T) {
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
