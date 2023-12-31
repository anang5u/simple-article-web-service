package controller

import (
	"net/http"
	"simple-ddd-cqrs/pkg/command"
	"simple-ddd-cqrs/pkg/query"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

// ArticleController
type ArticleController struct {
	articleCommandHandler command.ArticleCommandHandler
	articleQueryHandler   query.ArticleQueryHandler
}

// NewArticleController
func NewArticleController(articleCommandHandler command.ArticleCommandHandler, articleQueryHandler query.ArticleQueryHandler) *ArticleController {
	return &ArticleController{
		articleCommandHandler: articleCommandHandler,
		articleQueryHandler:   articleQueryHandler,
	}
}

// CreateArticle
func (ctr *ArticleController) CreateArticle(c *fiber.Ctx) error {
	var req struct {
		Author string `json:"author"`
		Title  string `json:"title"`
		Body   string `json:"body"`
	}
	created := time.Now()

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	article, err := ctr.articleCommandHandler.CreateArticle(req.Author, req.Title, req.Body, created)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(article)
}

// GetListArticle
func (ctr *ArticleController) GetListArticle(c *fiber.Ctx) error {
	// filter by search title, body
	// filter by author
	searchTitle := c.Query("title") // search title
	searchBody := c.Query("body")   // search body
	byAuthor := c.Query("author")   // filter by author

	filters := map[string]string{}
	if len(searchTitle) > 0 {
		filters["title"] = searchTitle
	}
	if len(searchBody) > 0 {
		filters["body"] = searchBody
	}
	if len(byAuthor) > 0 {
		filters["author"] = byAuthor
	}

	article, err := ctr.articleQueryHandler.GetListArticle(filters)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(article)
}

// GetArticleByID
func (ctr *ArticleController) GetArticleByID(c *fiber.Ctx) error {
	articleID := c.Params("id")
	ID, err := strconv.Atoi(articleID)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	article, err := ctr.articleQueryHandler.GetArticleByID(ID)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(article)
}
