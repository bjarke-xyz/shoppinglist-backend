package routing

import (
	itemController "ShoppingList-Backend/app/item/controller"
	listController "ShoppingList-Backend/app/list/controller"
	"ShoppingList-Backend/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

func PrivateRoutes(a *fiber.App) {
	route := a.Group("/api/v1")

	items := route.Group("/items", middleware.JWTProtected())
	items.Get("/", itemController.GetItems)
	items.Post("/", itemController.CreateItem)
	items.Put("/:id", itemController.UpdateItem)
	items.Delete("/:id", itemController.DeleteItem)

	lists := route.Group("/lists", middleware.JWTProtected())
	lists.Get("/", listController.GetLists)
	lists.Get("/default", listController.GetDefaultList)
	lists.Post("/", listController.CreateList)
	lists.Put("/:id", listController.UpdateList)
	lists.Put("/:id/default", listController.SetDefaultList)
	lists.Delete("/:id", listController.DeleteList)
	lists.Post("/:id/items/:itemId", listController.AddItemToList)
	lists.Put("/:id/items/:listItemId", listController.UpdateListItem)
	lists.Delete("/:id/items/:listItemId", listController.RemoveItemFromList)

}
