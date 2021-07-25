package queries

import (
	"ShoppingList-Backend/internal/pkg/item"
	"ShoppingList-Backend/internal/pkg/list"
)

type Queries struct {
	Item *item.ItemQueries
	List *list.ListQueries
}
