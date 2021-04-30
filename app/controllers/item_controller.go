package controllers

import (
	"ShoppingList-Backend/app/models"
	"ShoppingList-Backend/platform/database"
	"log"

	"github.com/gofiber/fiber/v2"
)

func GetItems(c *fiber.Ctx) error {
	db, err := database.OpenDBConnection()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	identityUser := c.Locals("user").(models.IdentityUser)
	log.Printf("%v", identityUser)
	items, err := db.GetItems(identityUser.ID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": items,
	})
}
