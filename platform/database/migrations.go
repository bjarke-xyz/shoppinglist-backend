package database

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(runDownMigrations bool) error {
	db, err := OpenDBConnection()
	if err != nil {
		return err
	}
	driver, err := postgres.WithInstance(db.ItemQueries.DB.DB, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file://platform/migrations", "ShoppingList", driver)
	if err != nil {
		return err
	}

	// err = m.Steps(2)
	// if err != nil {
	// 	return err
	// }

	if runDownMigrations {
		err = m.Down()
		if err != nil && err != migrate.ErrNoChange {
			return err
		}
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
