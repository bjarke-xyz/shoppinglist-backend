package routes

import (
	"ShoppingList-Backend/app/controllers"
	"ShoppingList-Backend/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

func PrivateRoutes(a *fiber.App) {
	route := a.Group("/api/v1")

	items := route.Group("/items", middleware.JWTProtected())

	items.Get("/", controllers.GetItems)
}
