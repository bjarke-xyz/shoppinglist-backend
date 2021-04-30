package queries

import (
	"ShoppingList-Backend/app/models"

	"github.com/jmoiron/sqlx"
)

type ItemQueries struct {
	*sqlx.DB
}

func (q *ItemQueries) GetItems(ownerID string) ([]models.Item, error) {
	items := []models.Item{}

	query := `SELECT * FROM items WHERE owner_id = $1 AND deleted_at IS NULL`

	err := q.Select(&items, query, ownerID)

	if err != nil {
		return items, err
	}

	return items, nil
}

func (q *ItemQueries) GetItem(id uint) (models.Item, error) {
	item := models.Item{}

	query := `SELECT * FROM items WHERE id = $1 AND deleted_at IS NULL`

	err := q.Get(&item, query, id)
	if err != nil {
		return item, err
	}

	return item, err
}

func (q *ItemQueries) CreateItem(item *models.Item) (uint, error) {
	query := `INSERT INTO items (name, owner_id) VALUES ($1, $2) returning ID`

	res, err := q.Exec(query, item.Name, item.OwnerID)
	if err != nil {
		return 0, err
	}

	insertedId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint(insertedId), nil
}

func (q *ItemQueries) UpdateItem(item *models.Item) error {
	query := `UPDATE items SET updated_at = NOW(), name = $2 WHERE id = $1`
	_, err := q.Exec(query, item.ID, item.Name)
	if err != nil {
		return err
	}
	return nil
}

func (q ItemQueries) DeleteItem(item *models.Item) error {
	query := `UPDATE items SET deleted_at = NOW() WHERE id = $1`
	_, err := q.Exec(query, item.ID)
	if err != nil {
		return err
	}
	return nil
}
