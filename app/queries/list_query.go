package queries

import (
	"ShoppingList-Backend/app/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ListQueries struct {
	*sqlx.DB
}

func (q *ListQueries) getListItems(listIds []uuid.UUID) ([]models.ListItem, error) {
	listItems := []models.ListItem{}
	query, args, err := sqlx.In(`SELECT * FROM list_item WHERE list_id IN (?)`, listIds)
	if err != nil {
		return listItems, err
	}

	query = q.Rebind(query)
	err = q.Select(&listItems, query, args...)
	if err != nil {
		return listItems, err
	}

	return listItems, nil
}

func (q *ListQueries) getItems(itemIds []uuid.UUID) ([]models.Item, error) {
	items := []models.Item{}
	query, args, err := sqlx.In(`SELECT * FROM items WHERE id IN (?)`, itemIds)
	if err != nil {
		return items, err
	}

	query = q.Rebind(query)
	err = q.Select(&items, query, args...)
	if err != nil {
		return items, err
	}
	return items, nil
}

func (q *ListQueries) populateWithItems(lists []models.List) error {
	// Get ListItems
	listIds := make([]uuid.UUID, len(lists))
	for _, list := range lists {
		listIds = append(listIds, list.ID)
	}
	listItems, err := q.getListItems(listIds)
	if err != nil {
		return err
	}
	listItemsByListId := make(map[uuid.UUID][]models.ListItem)
	for _, listItem := range listItems {
		_, present := listItemsByListId[listItem.ListID]
		if present {
			listItemsByListId[listItem.ListID] = append(listItemsByListId[listItem.ListID], listItem)
		} else {
			listItemsByListId[listItem.ListID] = []models.ListItem{listItem}
		}
	}
	for i, list := range lists {
		lists[i].Items = listItemsByListId[list.ID]
	}

	// Get ListItem.Item
	itemIds := make([]uuid.UUID, len(lists))
	for _, list := range lists {
		for _, listItem := range list.Items {
			itemIds = append(itemIds, listItem.ItemID)
		}
	}
	items, err := q.getItems(itemIds)
	if err != nil {
		return err
	}
	itemsById := make(map[uuid.UUID]models.Item)
	for _, item := range items {
		itemsById[item.ID] = item
	}
	for i := range lists {
		for j, item := range lists[i].Items {
			lists[i].Items[j].Item = itemsById[item.ItemID]
		}
	}

	return nil
}

func (q *ListQueries) GetLists(ownerID string) ([]models.List, error) {
	lists := []models.List{}

	query := `SELECT * FROM lists WHERE owner_id = $1 AND deleted_at IS NULL ORDER BY created_at ASC`

	err := q.Select(&lists, query, ownerID)

	if err != nil {
		return lists, err
	}

	q.populateWithItems(lists)

	for i := range lists {
		if lists[i].Items == nil {
			lists[i].Items = make([]models.ListItem, 0)
		}
	}

	return lists, nil
}

func (q *ListQueries) GetList(id uuid.UUID) (models.List, error) {
	list := models.List{}

	query := `SELECT * FROM lists WHERE id = $1 AND deleted_at IS NULL LIMIT 1`

	err := q.Get(&list, query, id)
	if err != nil {
		return list, err
	}

	lists := []models.List{list}
	q.populateWithItems(lists)

	if lists[0].Items == nil {
		lists[0].Items = make([]models.ListItem, 0)
	}

	return lists[0], err
}

func (q *ListQueries) CreateList(list models.List) (uuid.UUID, error) {
	if *list.IsDefault {
		_, err := q.Exec(`UPDATE lists SET is_default = false, updated_at = NOW() WHERE owner_id = $1`, list.OwnerID)
		if err != nil {
			return uuid.Nil, err
		}
	}

	query := `INSERT INTO lists (id, name, is_default, owner_id) VALUES ($1, $2, $3, $4)`
	_, err := q.Exec(query, list.ID, list.Name, list.IsDefault, list.OwnerID)
	if err != nil {
		return uuid.Nil, err
	}

	return list.ID, nil
}

func (q *ListQueries) UpdateList(list models.List) error {
	if *list.IsDefault {
		_, err := q.Exec(`UPDATE lists SET is_default = false, updated_at = NOW() WHERE owner_id = $1`, list.OwnerID)
		if err != nil {
			return err
		}
	}
	query := `UPDATE lists SET updated_at = NOW(), name = $2, is_default = $3 WHERE id = $1`
	_, err := q.Exec(query, list.ID, list.Name, list.IsDefault)
	if err != nil {
		return err
	}
	return nil
}

func (q *ListQueries) AddItemToList(list models.List, item models.Item) (models.ListItem, error) {
	listItem := models.ListItem{ID: uuid.New()}
	query := `INSERT INTO list_item (id, list_id, item_id) VALUES ($1, $2, $3)`
	_, err := q.Exec(query, listItem.ID, list.ID, item.ID)
	if err != nil {
		return listItem, err
	}
	fetchQuery := `SELECT * FROM list_item WHERE id = $1`
	err = q.Get(&listItem, fetchQuery, listItem.ID)
	if err != nil {
		return listItem, err
	}
	listItem.Item = item
	return listItem, nil
}

func (q *ListQueries) UpdateListItem(listItem models.ListItem) error {
	query := `UPDATE list_item SET updated_at = NOW(), crossed = $1 WHERE id = $2`
	_, err := q.Exec(query, listItem.Crossed, listItem.ID)
	if err != nil {
		return err
	}
	return nil
}

func (q *ListQueries) GetListItem(id uuid.UUID) (models.ListItem, error) {
	listItem := models.ListItem{}
	query := `SELECT * FROM list_item where id = $1`
	err := q.Get(&listItem, query, id)
	if err != nil {
		return listItem, err
	}
	item := models.Item{}
	itemQuery := `SELECT * FROM items WHERE id = $1`
	err = q.Get(&item, itemQuery, listItem.ItemID)
	if err != nil {
		return listItem, err
	}
	listItem.Item = item
	return listItem, nil
}

func (q *ListQueries) RemoveItemFromList(id uuid.UUID) error {
	query := `DELETE FROM list_item WHERE id = $1`
	_, err := q.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func (q *ListQueries) DeleteList(list models.List) error {
	query := `UPDATE lists SET deleted_at = NOW() WHERE id = $1`
	_, err := q.Exec(query, list.ID)
	if err != nil {
		return err
	}
	return nil
}
