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

	err := q.Get(&items, query, ownerID)

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
	query := `INSERT INTO items VALUES ($1, $2, $3, $4, $5) returning ID`

	// TODO: return id
	res, err := q.Exec(query, item.CreatedAt, item.UpdatedAt, item.DeletedAt, item.Name, item.OwnerID)
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
	query := `UPDATE items SET updated_at = $2, name = $3 WHERE id = $1`
	_, err := q.Exec(query, item.ID, item.UpdatedAt, item.Name)
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
