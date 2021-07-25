package application

import (
	"ShoppingList-Backend/internal/pkg/item"
	"ShoppingList-Backend/internal/pkg/list"
	"ShoppingList-Backend/internal/pkg/queries"
	"ShoppingList-Backend/pkg/config"
	"ShoppingList-Backend/pkg/db"
)

type Application struct {
	Cfg     *config.Config
	Queries *queries.Queries
}

func Get(cfg *config.Config) (*Application, error) {
	db, err := db.Get(cfg.GetDBConnStr())
	if err != nil {
		return nil, err
	}

	queries := &queries.Queries{
		Item: &item.ItemQueries{
			DB: db.Client,
		},
		List: &list.ListQueries{
			DB: db.Client,
		},
	}

	return &Application{
		Cfg:     cfg,
		Queries: queries,
	}, nil
}
