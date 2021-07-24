package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func FiberMiddleware(a *fiber.App) {
	a.Use(
		cors.New(),
		requestid.New(),
		// logger.New(),
		ZapLogger(),
	)
}
