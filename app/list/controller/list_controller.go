package controller

import (
	"ShoppingList-Backend/app/list"
	"ShoppingList-Backend/pkg/server"
	"ShoppingList-Backend/pkg/utils"
	"ShoppingList-Backend/platform/database"
	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GetLists func gets all lists for user
// @Description Get all lists for user
// @Summary get all lists for user
// @Tags lists
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} list.ListsResponse
// @Failure 500 {object} server.HTTPError
// @Failure 404 {object} server.HTTPError
// @Router /api/v1/lists [get]
func GetLists(c *fiber.Ctx) error {
	db, err := database.OpenDBConnection()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(server.HTTPError{
			Status: fiber.StatusInternalServerError,
			Error:  err.Error(),
		})
	}
	appUser := server.GetAppUser(c)
	lists, err := db.GetLists(appUser)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(server.HTTPError{
			Status: fiber.StatusNotFound,
			Error:  err.Error(),
		})
	}

	return c.JSON(list.ListsResponse{
		Data: lists,
	})
}

// GetDefaultList func Get the user's default list
// @Description Get the user's default list
// @Summary Get the user's default list
// @Tags lists
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} list.DefaultListResponse
// @Failure 500 {object} server.HTTPError
// @Failure 404 {object} server.HTTPError
// @Router /api/v1/lists/default [get]
func GetDefaultList(c *fiber.Ctx) error {
	db, err := database.OpenDBConnection()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(server.HTTPError{
			Status: fiber.StatusInternalServerError,
			Error:  err.Error(),
		})
	}

	appUser := server.GetAppUser(c)
	defaultList, err := db.GetDefaultList(appUser)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(server.HTTPError{
				Status: fiber.StatusNotFound,
				Error:  err.Error(),
			})
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(server.HTTPError{
				Status: fiber.StatusInternalServerError,
				Error:  err.Error(),
			})
		}
	}

	return c.JSON(list.DefaultListResponse{
		Data: defaultList,
	})
}

// CreateList func Create new list
// @Description Create new list
// @Summary Create new list
// @Tags lists
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param list body list.AddList true "Add list"
// @Success 200 {object} list.ListResponse
// @Failure 500 {object} server.HTTPError
// @Failure 400 {object} server.HTTPError
// @Router /api/v1/lists [post]
func CreateList(c *fiber.Ctx) error {
	addList := &list.AddList{}
	if err := c.BodyParser(addList); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(server.HTTPError{
			Status: fiber.StatusBadRequest,
			Error:  err.Error(),
		})
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(server.HTTPError{
			Status: fiber.StatusInternalServerError,
			Error:  err.Error(),
		})
	}

	validate := utils.NewValidator()

	appUser := server.GetAppUser(c)

	listToCreate := list.List{
		ID:      uuid.New(),
		Name:    addList.Name,
		OwnerID: appUser.ID,
	}

	if err := validate.Struct(addList); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(server.HTTPError{
			Status: fiber.StatusBadRequest,
			Error:  err.Error(),
		})
	}

	listId, err := db.CreateList(listToCreate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(server.HTTPError{
			Status: fiber.StatusInternalServerError,
			Error:  err.Error(),
		})
	}

	createdList, err := db.GetList(listId, appUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(server.HTTPError{
			Status: fiber.StatusInternalServerError,
			Error:  err.Error(),
		})
	}

	return c.JSON(list.ListResponse{
		Data: createdList,
	})
}

// UpdateList func Update list
// @Description Update list
// @Summary Update list
// @Tags lists
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path string true "List ID"
// @Param list body list.AddList true "Update list"
// @Success 200 {object} list.ListResponse
// @Failure 500 {object} server.HTTPError
// @Failure 404 {object} server.HTTPError
// @Failure 400 {object} server.HTTPError
// @Router /api/v1/lists/{id} [put]
func UpdateList(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(server.HTTPError{
			Status: fiber.StatusBadRequest,
			Error:  err.Error(),
		})
	}

	addList := &list.AddList{}
	if err := c.BodyParser(addList); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(server.HTTPError{
			Status: fiber.StatusBadRequest,
			Error:  err.Error(),
		})
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(server.HTTPError{
			Status: fiber.StatusInternalServerError,
			Error:  err.Error(),
		})
	}

	appUser := server.GetAppUser(c)
	foundList, err := db.GetList(id, appUser)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(server.HTTPError{
			Status: fiber.StatusNotFound,
			Error:  err.Error(),
		})
	}

	foundList.Name = addList.Name

	if err := db.UpdateList(foundList); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(server.HTTPError{
			Status: fiber.StatusInternalServerError,
			Error:  err.Error(),
		})
	}

	updatedList, err := db.GetList(id, appUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(server.HTTPError{
			Status: fiber.StatusInternalServerError,
			Error:  err.Error(),
		})
	}

	return c.JSON(list.ListResponse{
		Data: updatedList,
	})
}

// SetDefaultList func set default list
// @Description set default list
// @Summary set default list
// @Tags lists
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path string true "List ID"
// @Success 200 {object} list.ListResponse
// @Failure 500 {object} server.HTTPError
// @Failure 404 {object} server.HTTPError
// @Failure 400 {object} server.HTTPError
// @Router /api/v1/lists/{id}/default [put]
func SetDefaultList(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(server.HTTPError{
			Status: fiber.StatusBadRequest,
			Error:  err.Error(),
		})
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(server.HTTPError{
			Status: fiber.StatusInternalServerError,
			Error:  err.Error(),
		})
	}

	appUser := server.GetAppUser(c)
	foundList, err := db.GetList(id, appUser)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(server.HTTPError{
			Status: fiber.StatusNotFound,
			Error:  err.Error(),
		})
	}

	defaultList, err := db.SetDefaultList(appUser, foundList)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(server.HTTPError{
			Status: fiber.StatusInternalServerError,
			Error:  err.Error(),
		})
	}

	return c.JSON(list.DefaultListResponse{
		Data: defaultList,
	})
}

// AddItemToList func Add item to list
// @Description Add item to list
// @Summary Add item to list
// @Tags lists, items
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param list-id path string true "List ID"
// @Param item-id path string true "Item ID"
// @Success 200 {object} list.ListItemResponse
// @Failure 500 {object} server.HTTPError
// @Failure 404 {object} server.HTTPError
// @Failure 400 {object} server.HTTPError
// @Router /api/v1/lists/{list-id}/{item-id} [patch]
func AddItemToList(c *fiber.Ctx) error {
	log.Println("hej")
	listIdStr := c.Params("id")
	listId, err := uuid.Parse(listIdStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(server.HTTPError{
			Status: fiber.StatusBadRequest,
			Error:  err.Error(),
		})
	}

	itemIdStr := c.Params("itemId")
	itemId, err := uuid.Parse(itemIdStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(server.HTTPError{
			Status: fiber.StatusBadRequest,
			Error:  err.Error(),
		})
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(server.HTTPError{
			Status: fiber.StatusInternalServerError,
			Error:  err.Error(),
		})
	}

	appUser := server.GetAppUser(c)
	foundList, err := db.GetList(listId, appUser)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(server.HTTPError{
			Status: fiber.StatusNotFound,
			Error:  err.Error(),
		})
	}

	item, err := db.GetItem(itemId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(server.HTTPError{
			Status: fiber.StatusNotFound,
			Error:  err.Error(),
		})
	}

	listItem, err := db.AddItemToList(foundList, item)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(server.HTTPError{
			Status: fiber.StatusInternalServerError,
			Error:  err.Error(),
		})
	}

	return c.JSON(list.ListItemResponse{
		Data: listItem,
	})
}

// UpdateListItem func Update list item
// @Description Update list item
// @Summary Update list item
// @Tags lists, items
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param list-id path string true "List ID"
// @Param list-item-id path string true "List-Item ID"
// @Param list body list.UpdateListItem true "Update list item"
// @Success 200 {object} list.ListItemResponse
// @Failure 500 {object} server.HTTPError
// @Failure 404 {object} server.HTTPError
// @Failure 400 {object} server.HTTPError
// @Router /api/v1/lists/{list-id}/{list-item-id} [put]
func UpdateListItem(c *fiber.Ctx) error {
	listIdStr := c.Params("id")
	listId, err := uuid.Parse(listIdStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(server.HTTPError{
			Status: fiber.StatusBadRequest,
			Error:  err.Error(),
		})
	}

	listItemIdStr := c.Params("listItemId")
	listItemId, err := uuid.Parse(listItemIdStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(server.HTTPError{
			Status: fiber.StatusBadRequest,
			Error:  err.Error(),
		})
	}

	updateListItem := &list.UpdateListItem{}
	if err := c.BodyParser(updateListItem); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(server.HTTPError{
			Status: fiber.StatusBadRequest,
			Error:  err.Error(),
		})
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(server.HTTPError{
			Status: fiber.StatusInternalServerError,
			Error:  err.Error(),
		})
	}

	appUser := server.GetAppUser(c)
	if _, err := db.GetList(listId, appUser); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(server.HTTPError{
			Status: fiber.StatusNotFound,
			Error:  err.Error(),
		})
	}

	listItem, err := db.GetListItem(listItemId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(server.HTTPError{
			Status: fiber.StatusNotFound,
			Error:  err.Error(),
		})
	}

	listItem.Crossed = updateListItem.Crossed

	if err := db.UpdateListItem(listItem); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(server.HTTPError{
			Status: fiber.StatusInternalServerError,
			Error:  err.Error(),
		})
	}

	return c.JSON(list.ListItemResponse{
		Data: listItem,
	})
}

// RemoveItemFromList func Remove item from list
// @Description Remove item from list
// @Summary Remove item from list
// @Tags lists, items
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param list-id path string true "List ID"
// @Param list-item-id path string true "List-Item ID"
// @Success 204 {string} status "ok"
// @Failure 500 {object} server.HTTPError
// @Failure 404 {object} server.HTTPError
// @Failure 400 {object} server.HTTPError
// @Router /api/v1/lists/{list-id}/{list-item-id} [delete]
func RemoveItemFromList(c *fiber.Ctx) error {

	listIdStr := c.Params("id")
	listId, err := uuid.Parse(listIdStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(server.HTTPError{
			Status: fiber.StatusBadRequest,
			Error:  err.Error(),
		})
	}

	listItemIdStr := c.Params("listItemId")
	listItemId, err := uuid.Parse(listItemIdStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(server.HTTPError{
			Status: fiber.StatusBadRequest,
			Error:  err.Error(),
		})
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(server.HTTPError{
			Status: fiber.StatusInternalServerError,
			Error:  err.Error(),
		})
	}

	appUser := server.GetAppUser(c)
	if _, err := db.GetList(listId, appUser); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(server.HTTPError{
			Status: fiber.StatusNotFound,
			Error:  err.Error(),
		})
	}

	if err := db.RemoveItemFromList(listItemId); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(server.HTTPError{
			Status: fiber.StatusInternalServerError,
			Error:  err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// DeleteList func Delete list
// @Description Delete list
// @Summary Delete list
// @Tags lists
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path string true "List ID"
// @Success 204 {string} status "ok"
// @Failure 500 {object} server.HTTPError
// @Failure 400 {object} server.HTTPError
// @Router /api/v1/lists/{id} [delete]
func DeleteList(c *fiber.Ctx) error {

	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(server.HTTPError{
			Status: fiber.StatusBadRequest,
			Error:  err.Error(),
		})
	}

	list := &list.List{
		ID: id,
	}

	validate := utils.NewValidator()

	if err := validate.StructPartial(list, "id"); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(server.HTTPError{
			Status: fiber.StatusBadRequest,
			Error:  err.Error(),
		})
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(server.HTTPError{
			Status: fiber.StatusInternalServerError,
			Error:  err.Error(),
		})
	}

	appUser := server.GetAppUser(c)
	foundList, err := db.GetList(list.ID, appUser)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(server.HTTPError{
			Status: fiber.StatusNotFound,
			Error:  err.Error(),
		})
	}

	if err := db.DeleteList(foundList); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(server.HTTPError{
			Status: fiber.StatusInternalServerError,
			Error:  err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
