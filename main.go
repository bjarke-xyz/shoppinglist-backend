package main

import (
	"ShoppingList-Backend/pkg/config"
	"ShoppingList-Backend/pkg/middleware"
	"ShoppingList-Backend/pkg/routing"
	"ShoppingList-Backend/pkg/server"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
		return
	}

	config := config.FiberConfig()

	app := fiber.New(config)

	middleware.FiberMiddleware(app)

	routing.SwaggerRoute(app)
	routing.PrivateRoutes(app)

	server.Start(app)
}
