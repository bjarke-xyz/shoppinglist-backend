package main

import (
	"ShoppingList-Backend/pkg/config"
	"ShoppingList-Backend/pkg/logger"

	"github.com/golang-migrate/migrate/v4"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	godotenv.Load()
	cfg := config.New()
	logger.SetLogs(zap.DebugLevel, cfg.LogFormat)

	log := zap.S()

	direction := cfg.Migrate
	if direction != "down" && direction != "up" {
		log.Fatalf("-migrate accepts [up, down] values only")
	}

	m, err := migrate.New("file://db/migrations", cfg.GetDBConnStr())
	if err != nil {
		log.Fatalf("Error getting migration files: %v", err)
	}

	if direction == "up" {
		if err := m.Up(); err != nil {
			if err == migrate.ErrNoChange {
				log.Infof("%v", err)
			} else {
				log.Fatalf("Failed to migrate up: %v", err)
			}
		}
	}
	if direction == "down" {
		if err := m.Down(); err != nil {
			if err == migrate.ErrNoChange {
				log.Infof("%v", err)
			} else {
				log.Fatalf("Failed to migrate down: %v", err)
			}
		}
	}
}
