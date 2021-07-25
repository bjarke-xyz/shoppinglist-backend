package item

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ItemQueries struct {
	DB *sqlx.DB
}

func (q *ItemQueries) GetItems(ownerID string) ([]Item, error) {
	items := []Item{}

	query := `SELECT * FROM items WHERE owner_id = $1 AND deleted_at IS NULL`

	err := q.DB.Select(&items, query, ownerID)

	if err != nil {
		return items, err
	}

	return items, nil
}

func (q *ItemQueries) GetItem(id uuid.UUID) (Item, error) {
	item := Item{}

	query := `SELECT * FROM items WHERE id = $1 AND deleted_at IS NULL`

	err := q.DB.Get(&item, query, id)
	if err != nil {
		return item, err
	}

	return item, err
}

func (q *ItemQueries) CreateItem(item *Item) (uuid.UUID, error) {
	existingItem := Item{}
	fetchQuery := `SELECT * FROM items WHERE name = $1 AND owner_id = $2`
	err := q.DB.Get(&existingItem, fetchQuery, item.Name, item.OwnerID)
	if err == nil {
		return existingItem.ID, nil
	}

	query := `INSERT INTO items (id, name, owner_id) VALUES ($1, $2, $3)`

	_, err = q.DB.Exec(query, item.ID, item.Name, item.OwnerID)
	if err != nil {
		return uuid.Nil, err
	}

	return item.ID, nil
}

func (q *ItemQueries) UpdateItem(item *Item) error {
	query := `UPDATE items SET updated_at = NOW(), name = $2 WHERE id = $1`
	_, err := q.DB.Exec(query, item.ID, item.Name)
	if err != nil {
		return err
	}
	return nil
}

func (q *ItemQueries) DeleteItem(item *Item) error {
	query := `UPDATE items SET deleted_at = NOW() WHERE id = $1`
	_, err := q.DB.Exec(query, item.ID)
	if err != nil {
		return err
	}
	return nil
}
