package main

import (
	"ShoppingList-Backend/cmd/api/router"
	"ShoppingList-Backend/pkg/application"
	"ShoppingList-Backend/pkg/config"
	"ShoppingList-Backend/pkg/logger"
	"ShoppingList-Backend/pkg/middleware"
	"ShoppingList-Backend/pkg/server"
	"log"

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
	}
	cfg := config.New()
	logger.SetLogs(zap.DebugLevel, cfg.LogFormat)

	app, err := application.Get(cfg)
	if err != nil {
		zap.S().Fatalf("Database error: %v", err)
	}

	fiberConfig := config.FiberConfig(app.Cfg)

	fiberApp := fiber.New(fiberConfig)

	middleware.FiberMiddleware(fiberApp)

	router.SwaggerRoute(fiberApp, app)
	router.PrivateRoutes(fiberApp, app)

	server.Start(fiberApp, app)
}
