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
	items.Post("/", controllers.CreateItem)
	items.Put("/:id", controllers.UpdateItem)
	items.Delete("/:id", controllers.DeleteItem)

	lists := route.Group("/lists", middleware.JWTProtected())
	lists.Get("/", controllers.GetLists)
	lists.Post("/", controllers.CreateList)
	lists.Put("/:id", controllers.UpdateList)
	lists.Delete("/:id", controllers.DeleteList)
	lists.Post("/:id/items/:itemId", controllers.AddItemToList)
	lists.Put("/:id/items/:listItemId", controllers.UpdateListItem)
	lists.Delete("/:id/items/:listItemId", controllers.RemoveItemFromList)

}
