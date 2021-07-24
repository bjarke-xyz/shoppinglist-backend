package database

import (
	"ShoppingList-Backend/app/item"
	"ShoppingList-Backend/app/list"

	"github.com/jmoiron/sqlx"
)

type Queries struct {
	Item *item.ItemQueries
	List *list.ListQueries
	db   *sqlx.DB
}

type DBConnection interface {
	OpenDBConnection() (*Queries, error)
	Close() error
}

func OpenDBConnection() (*Queries, error) {
	db, err := PostgreSQLConnection()
	if err != nil {
		return nil, err
	}

	return &Queries{
		Item: &item.ItemQueries{DB: db},
		List: &list.ListQueries{DB: db},
		db:   db,
	}, nil
}

func (q *Queries) Close() error {
	return q.db.Close()
}
