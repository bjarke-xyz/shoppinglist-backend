package server

import (
	"ShoppingList-Backend/pkg/application"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func Start(a *fiber.App, app *application.Application) {
	zap.S().Infow("Api started", "SERVER_URL", app.Cfg.GetServerUrl())
	if err := a.Listen(app.Cfg.GetServerUrl()); err != nil {
		zap.S().Errorw("Server is not running!", "Reason", err)
	}
}

type HTTPError struct {
	Status int    `json:"status"`
	Error  string `json:"error"`
}
