package routes_test

import (
	"database/sql/driver"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"simple-ddd-cqrs/controller"
	"simple-ddd-cqrs/pkg/command"
	"simple-ddd-cqrs/pkg/domain"
	"simple-ddd-cqrs/pkg/query"
	"simple-ddd-cqrs/routes"
	"simple-ddd-cqrs/service"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

type anyTime struct{}

func (anyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

// Test Route POST /articles - Create new article
func TestHandle_PostArticle(t *testing.T) {
	// init DB Mock
	db, mock := service.NewDBMock()
	defer db.Close()

	// Init article repository
	articleRepo := domain.CreateArticleRepository(db)

	articleMock := struct {
		Author string
		Title  string
		Name   string
		Body   string
	}{
		Author: "Author Test",
		Title:  "Title Test",
		Name:   "title-test",
		Body:   "Body test",
	}

	query := regexp.QuoteMeta(`INSERT INTO articles`)
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(articleMock.Author, articleMock.Title, articleMock.Name, articleMock.Body, anyTime{}).WillReturnResult(sqlmock.NewResult(1, 1))

	// Init COMMAND handlers
	articleCommandHandler := command.NewArticleCommandHandler(articleRepo)
	articleController := controller.NewArticleController(articleCommandHandler, nil)

	// Membuat instance Fiber app
	app := fiber.New()

	// Memasang handler ke dalam Fiber app
	routes.Handle(app, articleController)

	// Req payload
	payload, _ := json.Marshal(articleMock)

	// Test Route POST /articles
	req := httptest.NewRequest("POST", "/articles", strings.NewReader(string(payload)))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)

	assert.NoError(t, err)

	// read and unmarshal respons JSON
	var createdArticle domain.ArticleModel
	err = json.NewDecoder(resp.Body).Decode(&createdArticle)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	assert.NoError(t, err)
	assert.Equal(t, articleMock.Author, createdArticle.Author)
	assert.Equal(t, articleMock.Title, createdArticle.Title)
	assert.Equal(t, articleMock.Name, createdArticle.Name)
	// ... assert lainnya sesuai kebutuhan
}

// Test Route GET /articles - get list articles
func TestHandle_GetListArticle(t *testing.T) {
	// init DB Mock
	db, mock := service.NewDBMock()
	defer db.Close()

	var (
		ID      = 1
		author  = "Author Test"
		title   = "Title Test"
		name    = "title-test"
		body    = "Body Test"
		created = time.Now()
	)
	queryStmt := "SELECT id, author, title, name, body, created FROM articles"
	rows := sqlmock.NewRows([]string{"id", "author", "title", "name", "body", "created"}).
		AddRow(ID, author, title, name, body, created)

	mock.ExpectQuery(queryStmt).WillReturnRows(rows)

	// Init article repository
	articleRepo := domain.CreateArticleRepository(db)

	// Init QUERY handlers
	articleQueryHandler := query.NewArticleQueryHandler(articleRepo)
	articleController := controller.NewArticleController(nil, articleQueryHandler)

	// Membuat instance Fiber app
	app := fiber.New()

	// Memasang handler ke dalam Fiber app
	routes.Handle(app, articleController)

	// Test Route POST /articles
	req := httptest.NewRequest("GET", "/articles", nil)
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// read and unmarshal respons JSON
	var articles []domain.ArticleModel
	err = json.NewDecoder(resp.Body).Decode(&articles)
	assert.NoError(t, err)
	// ... assert lainnya sesuai kebutuhan
}
