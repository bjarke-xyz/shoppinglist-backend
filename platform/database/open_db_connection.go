package database

import "ShoppingList-Backend/app/queries"

type Queries struct {
	*queries.ItemQueries
	*queries.ListQueries
}

func OpenDBConnection() (*Queries, error) {
	db, err := PostgreSQLConnection()
	if err != nil {
		return nil, err
	}

	return &Queries{
		ItemQueries: &queries.ItemQueries{DB: db},
		ListQueries: &queries.ListQueries{DB: db},
	}, nil
}
