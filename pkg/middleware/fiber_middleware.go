package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/websocket/v2"
)

func FiberMiddleware(a *fiber.App) {
	a.Use(
		cors.New(),
		requestid.New(),
		NewZapLogger(),
	)
}

func IsWebSocketAllowed() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("wsPath", c.Path())
			c.Locals("wsAllowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	}
}
