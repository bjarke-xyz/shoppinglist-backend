package server

import (
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func StartWithGracefulShutdown(a *fiber.App) {
	idleConnsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt) // Catch OS signals
		<-sigint

		if err := a.Shutdown(); err != nil {
			// Error from closing listeners, or context tiemout
			zap.S().Errorw("Server is not shutting down!", "Reason", err)
		}

		close(idleConnsClosed)
	}()
}

func Start(a *fiber.App) {
	zap.S().Infow("Api started", "SERVER_URL", os.Getenv("SERVER_URL"))
	if err := a.Listen(os.Getenv("SERVER_URL")); err != nil {
		zap.S().Errorw("Server is not running!", "Reason", err)
	}
}

type HTTPError struct {
	Status int    `json:"status"`
	Error  string `json:"error"`
}
