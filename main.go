package main

import (
	"ShoppingList-Backend/pkg/configs"
	"ShoppingList-Backend/pkg/middleware"
	"ShoppingList-Backend/pkg/routes"
	"ShoppingList-Backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

func main() {
	config := configs.FiberConfig()

	app := fiber.New(config)

	middleware.FiberMiddleware(app)

	routes.SwaggerRoute(app)
	routes.PrivateRoutes(app)

	utils.StartServer(app)
}
