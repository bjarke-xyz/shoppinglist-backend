package utils

import (
	"log"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v2"
)

func StartServerWithGracefulShutdown(a *fiber.App) {
	idleConnsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt) // Catch OS signals
		<-sigint

		if err := a.Shutdown(); err != nil {
			// Error from closing listeners, or context tiemout
			log.Printf("Oops... Server is not shutting down! Reason: %v", err)
		}

		close(idleConnsClosed)
	}()
}

func StartServer(a *fiber.App) {
	if err := a.Listen(os.Getenv("SERVER_URL")); err != nil {
		log.Printf("Oops... Server is not running! Reason: %v", err)
	}
}
