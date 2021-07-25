package config

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func FiberConfig(cfg *Config) fiber.Config {
	return fiber.Config{
		ReadTimeout:           time.Second * time.Duration(cfg.ServerReadTimeout),
		DisableStartupMessage: true,
	}
}
