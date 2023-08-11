package main

import (
	"fmt"
	"simple-ddd-cqrs/config"
	"simple-ddd-cqrs/controller"
	"simple-ddd-cqrs/pkg/command"
	"simple-ddd-cqrs/pkg/domain"
	"simple-ddd-cqrs/pkg/query"
	"simple-ddd-cqrs/routes"
	"simple-ddd-cqrs/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func init() {
	// Load the .env file in the current directory
	godotenv.Load()
}

func main() {
	fmt.Println("A simple article web service with DDD-CQRS")

	// Init Fiber app
	app := fiber.New()

	// Middleware CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// create article repository
	articleRepo := domain.CreateArticleRepository(service.GetDBConnection())

	// Init command and query handlers
	articleCommandHandler := command.NewArticleCommandHandler(articleRepo)
	articleQueryHandler := query.NewArticleQueryHandler(articleRepo)

	// Init controller
	articleController := controller.NewArticleController(articleCommandHandler, articleQueryHandler)

	// Route handler
	routes.Handle(app, articleController)

	// Run server
	err := app.Listen(":" + config.Get("APP_PORT"))
	if err != nil {
		panic(err)
	}

}
