package database

import "ShoppingList-Backend/app/queries"

type Queries struct {
	*queries.ItemQueries
}

func OpenDBConnection() (*Queries, error) {
	db, err := PostgreSQLConnection()
	if err != nil {
		return nil, err
	}

	return &Queries{
		ItemQueries: &queries.ItemQueries{DB: db},
	}, nil
}
