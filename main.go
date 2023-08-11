package main

import (
	"fmt"

	"github.com/joho/godotenv"
)

func init() {
	// Load the .env file in the current directory
	godotenv.Load()
}

func main() {
	fmt.Println("A simple article web service with DDD-CQRS")
}
