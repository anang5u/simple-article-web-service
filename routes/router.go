package routes

import (
	"simple-ddd-cqrs/controller"

	"github.com/gofiber/fiber/v2"
)

// App struct to hold dependencies for handlers
type App struct {
	articleController *controller.ArticleController
}

func Handle(app *fiber.App, articleController *controller.ArticleController) {
	// Init App with controller
	appConfig := App{
		articleController: articleController,
	}

	// Route POST create new article
	app.Post("/articles", appConfig.articleController.CreateArticle)

	// Route GET list article
	app.Get("/articles", appConfig.articleController.GetListArticle)

	// Route GET article by id
	// Add param name for improve SEO Optimization
	app.Get("/articles/:id/:name", appConfig.articleController.GetArticleByID)
}
