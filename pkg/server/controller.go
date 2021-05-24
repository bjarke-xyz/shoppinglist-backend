package server

import (
	"ShoppingList-Backend/app/user"

	"github.com/gofiber/fiber/v2"
)

func GetAppUser(ctx *fiber.Ctx) user.AppUser {
	appUser := ctx.Locals("user").(user.AppUser)
	return appUser
}
