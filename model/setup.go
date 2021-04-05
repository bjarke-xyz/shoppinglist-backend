package model

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
)

func Init() (*gorm.DB, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	appEnv := os.Getenv("app_env")

	username := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbName := os.Getenv("db_name")
	dbHost := os.Getenv("db_host")
	dbPort := os.Getenv("db_port")

	dbUri := fmt.Sprintf("host=%s user=%s dbname=%s port=%s sslmode=disable password=%s", dbHost, username, dbName, dbPort, password) //Build connection string
	// fmt.Println(dbUri)

	conn, err := gorm.Open("postgres", dbUri)
	if err != nil {
		return nil, err
	}

	if appEnv == "dev" {
		conn.LogMode(true)
	}

	conn.Debug().AutoMigrate(&ListItem{})
	conn.Debug().AutoMigrate(&Item{})
	conn.Debug().AutoMigrate(&List{})

	return conn, nil
}
