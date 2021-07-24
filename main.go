package main

import (
	"ShoppingList-Backend/pkg/config"
	"ShoppingList-Backend/pkg/middleware"
	"ShoppingList-Backend/pkg/routing"
	"ShoppingList-Backend/pkg/server"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
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

	var logFormat string
	if os.Getenv("APP_ENV") == "production" {
		logFormat = server.LOGFORMAT_JSON
	} else {
		logFormat = server.LOGFORMAT_CONSOLE
	}
	server.SetLogs(zap.DebugLevel, logFormat)

	config := config.FiberConfig()

	app := fiber.New(config)

	middleware.FiberMiddleware(app)

	routing.SwaggerRoute(app)
	routing.PrivateRoutes(app)

	server.Start(app)
}
