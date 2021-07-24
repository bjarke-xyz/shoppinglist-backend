package database

import (
	"fmt"
	"os"
	"strconv"

	"github.com/jmoiron/sqlx"

	_ "github.com/jackc/pgx/v4/stdlib" // PostgreSQL pgx driver
)

func PostgreSQLConnection() (*sqlx.DB, error) {
	maxConn, _ := strconv.Atoi(os.Getenv("DB_MAX_CONNECTIONS"))
	maxIdleConn, _ := strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONNECTIONS"))
	// maxLifetimeConn, _ := strconv.Atoi(os.Getenv("DB_MAX_LIFETIME_CONNECTIONS"))

	db, err := sqlx.Open("pgx", os.Getenv("DB_SERVER"))
	if err != nil {
		return nil, fmt.Errorf("error, not connected to database, %w", err)
	}

	db.SetMaxOpenConns(maxConn)     // Default is 0 (Unlimited)
	db.SetMaxIdleConns(maxIdleConn) // Default is 2
	// db.SetConnMaxLifetime(time.Duration(time.Hour)) // 0, connections are reused forever

	return db, nil
}