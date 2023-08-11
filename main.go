package main

import (
	"fmt"
	"simple-ddd-cqrs/config"
	"simple-ddd-cqrs/routes"

	"github.com/gofiber/fiber/v2"
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

	// Route handler
	routes.Handle(app)

	// Run server
	err := app.Listen(":" + config.Get("APP_PORT"))
	if err != nil {
		panic(err)
	}

}
