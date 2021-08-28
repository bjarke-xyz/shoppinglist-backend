package main

import (
	"ShoppingList-Backend/cmd/api/router"
	"ShoppingList-Backend/internal/pkg/migration"
	"ShoppingList-Backend/pkg/application"
	"ShoppingList-Backend/pkg/config"
	"ShoppingList-Backend/pkg/logger"
	"ShoppingList-Backend/pkg/server"
	"log"

	"github.com/gorilla/mux"
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

	if cfg.MigrateOnStartup {
		migration.DoMigration("up", cfg.GetDBConnStr())
	}

	r := mux.NewRouter().StrictSlash(true)

	router.SwaggerRoute(app)
	router.PrivateRoutes(app, r)

	server.Start(app.Cfg, r)
}
