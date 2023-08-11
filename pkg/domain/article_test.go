package domain_test

import (
	"simple-ddd-cqrs/pkg/domain"
	"simple-ddd-cqrs/service"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

var (
	ID      = 1
	author  = "Author test"
	title   = "Title test"
	name    = "title-test"
	body    = "Body test"
	created = time.Now()
)

// Test Article Create
func Test_ArticleCreate(t *testing.T) {
	db, mock := service.NewDBMock()
	repo := domain.CreateArticleRepository(db)
	defer db.Close()

	query := "INSERT INTO articles \\(author, title, name, body, created\\) VALUES \\(\\$1, \\$2, \\$3, \\$4, \\$5\\)"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(author, title, name, body, created).WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.Create(author, title, name, body, created)

	assert.NoError(t, err)
}

// Test Article Get
func Test_ArticleGet(t *testing.T) {
	db, mock := service.NewDBMock()
	repo := domain.CreateArticleRepository(db)
	defer db.Close()

	query := "SELECT id, author, title, name, body, created FROM articles"

	rows := sqlmock.NewRows([]string{"id", "author", "title", "name", "body", "created"}).
		AddRow(ID, author, title, name, body, created)

	mock.ExpectQuery(query).WillReturnRows(rows)

	articles, err := repo.Get()

	assert.NotEmpty(t, articles)
	assert.NoError(t, err)
	assert.Len(t, articles, 1)
}

func Test_ArticleGetWithFilter(t *testing.T) {
	db, mock := service.NewDBMock()
	repo := domain.CreateArticleRepository(db)
	defer db.Close()

	query := "SELECT id, author, title, name, body, created FROM articles"

	rows := sqlmock.NewRows([]string{"id", "author", "title", "name", "body", "created"}).
		AddRow(ID, author, title, name, body, created)

	mock.ExpectQuery(query).WillReturnRows(rows)

	filter := map[string]string{
		"author": "Author test",
	}
	articles, err := repo.Get(filter)

	assert.NotEmpty(t, articles)
	assert.NoError(t, err)
	assert.Len(t, articles, 1)
}

// Test Get Article ByID
func Test_GetArticleByID(t *testing.T) {
	db, mock := service.NewDBMock()
	repo := domain.CreateArticleRepository(db)
	defer db.Close()

	query := "SELECT id, author, title, name, body, created FROM articles WHERE id = \\$1"

	rows := sqlmock.NewRows([]string{"id", "author", "title", "name", "body", "created"}).
		AddRow(ID, author, title, name, body, created)

	mock.ExpectQuery(query).WithArgs(ID).WillReturnRows(rows)

	article, err := repo.GetByID(ID)
	assert.NotNil(t, article)
	assert.NoError(t, err)
}

// Test buildFilter Article
func Test_BuildFilterValueArticle(t *testing.T) {
	repo := domain.CreateArticleRepository(nil)

	filter := map[string]string{
		"title":  "title search",
		"body":   "body search",
		"author": "author 1",
	}
	sFilter, values := repo.BuildFilterValues(filter)

	assert.Contains(t, sFilter, "AND")
	assert.Len(t, values, 3)
}
