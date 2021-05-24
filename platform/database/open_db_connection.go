package database

import (
	"ShoppingList-Backend/app/item"
	"ShoppingList-Backend/app/list"
)

type Queries struct {
	*item.ItemQueries
	*list.ListQueries
}

func OpenDBConnection() (*Queries, error) {
	db, err := PostgreSQLConnection()
	if err != nil {
		return nil, err
	}

	return &Queries{
		ItemQueries: &item.ItemQueries{DB: db},
		ListQueries: &list.ListQueries{DB: db},
	}, nil
}
