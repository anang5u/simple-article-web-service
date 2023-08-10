package command_test

import (
	"simple-ddd-cqrs/pkg/command"
	"simple-ddd-cqrs/pkg/domain"
	"simple-ddd-cqrs/service"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func Test_ArticleCommandCreate(t *testing.T) {
	db, mock := service.NewDBMock()
	articleRepo := domain.CreateArticleRepository(db)
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

	articleCommandHandler := command.NewArticleCommandHandler(articleRepo)

	article, err := articleCommandHandler.CreateArticle(author, title, body, created)

	assert.NoError(t, err)
	assert.Equal(t, article.Author, "Author test")
}
