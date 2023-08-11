package routes

import (
	"simple-ddd-cqrs/controller"
	"simple-ddd-cqrs/pkg/command"
	"simple-ddd-cqrs/pkg/domain"
	"simple-ddd-cqrs/pkg/query"
	"simple-ddd-cqrs/service"

	"github.com/gofiber/fiber/v2"
)

// App struct to hold dependencies for handlers
type App struct {
	articleController *controller.ArticleController
}

func Handle(app *fiber.App) {
	// create article repository
	articleRepo := domain.CreateArticleRepository(service.GetDBConnection())

	// Init command and query handlers
	articleCommandHandler := command.NewArticleCommandHandler(articleRepo)
	articleQueryHandler := query.NewArticleQueryHandler(articleRepo)

	// Init controller
	articleController := controller.NewArticleController(articleCommandHandler, articleQueryHandler)

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
