package main

import (
	"ShoppingList-Backend/pkg/configs"
	"ShoppingList-Backend/pkg/middleware"
	"ShoppingList-Backend/pkg/routes"
	"ShoppingList-Backend/pkg/utils"
	"ShoppingList-Backend/platform/database"
	"flag"
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

	migrationsDownFlag := flag.Bool("migrations-down", false, "Set to true, to run down migrations first")
	flag.Parse()
	err = database.RunMigrations(*migrationsDownFlag)
	if err != nil {
		log.Fatalf("Could not run migrations: %v", err)
		return
	}

	config := configs.FiberConfig()

	app := fiber.New(config)

	middleware.FiberMiddleware(app)

	routes.SwaggerRoute(app)
	routes.PrivateRoutes(app)

	utils.StartServer(app)
}
