package domain_test

import (
	"simple-ddd-cqrs/pkg/domain"
	"simple-ddd-cqrs/service"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

// Test Article Create
func Test_ArticleCreate(t *testing.T) {
	db, mock := service.NewDBMock()
	repo := domain.CreateArticleRepository(db)
	defer db.Close()

	var (
		author  = "Author test"
		title   = "Title test"
		body    = "Body test"
		created = time.Now()
	)

	query := "INSERT INTO articles \\(author, title, body, created\\) VALUES \\(\\$1, \\$2, \\$3, \\$4\\)"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(author, title, body, created).WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.Create(author, title, body, created)

	assert.NoError(t, err)
}

// Test Article Get
func Test_ArticleGet(t *testing.T) {
	db, mock := service.NewDBMock()
	repo := domain.CreateArticleRepository(db)
	defer db.Close()

	var (
		ID      = 1
		author  = "Author test"
		title   = "Title test"
		body    = "Body test"
		created = time.Now()
	)

	query := "SELECT id, author, title, body, created FROM articles"

	rows := sqlmock.NewRows([]string{"id", "author", "title", "body", "created"}).
		AddRow(ID, author, title, body, created)

	mock.ExpectQuery(query).WillReturnRows(rows)

	articles, err := repo.Get()

	assert.NotEmpty(t, articles)
	assert.NoError(t, err)
	assert.Len(t, articles, 1)
}

// Test Get Article ByID
func Test_GetArticleByID(t *testing.T) {
	db, mock := service.NewDBMock()
	repo := domain.CreateArticleRepository(db)
	defer db.Close()

	var (
		ID      = 1
		author  = "Author test"
		title   = "Title test"
		body    = "Body test"
		created = time.Now()
	)

	query := "SELECT id, author, title, body, created FROM articles WHERE id = \\$1"

	rows := sqlmock.NewRows([]string{"id", "author", "title", "body", "created"}).
		AddRow(ID, author, title, body, created)

	mock.ExpectQuery(query).WithArgs(ID).WillReturnRows(rows)

	article, err := repo.GetByID(ID)
	assert.NotNil(t, article)
	assert.NoError(t, err)
}
