package controllers

import (
	"ShoppingList-Backend/app/models"
	"ShoppingList-Backend/pkg/utils"
	"ShoppingList-Backend/platform/database"
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
// @Success 200 {object} models.ListsResponse
// @Failure 500 {object} models.HTTPError
// @Failure 404 {object} models.HTTPError
// @Router /api/v1/lists [get]
func GetLists(c *fiber.Ctx) error {
	db, err := database.OpenDBConnection()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.HTTPError{
			Status: fiber.StatusInternalServerError,
			Error:  err.Error(),
		})
	}
	identityUser := c.Locals("user").(models.IdentityUser)
	log.Printf("%v", identityUser)
	lists, err := db.GetLists(identityUser.ID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.HTTPError{
			Status: fiber.StatusNotFound,
			Error:  err.Error(),
		})
	}

	return c.JSON(models.ListsResponse{
		Data: lists,
	})
}

// CreateList func Create new list
// @Description Create new list
// @Summary Create new list
// @Tags lists
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param list body models.AddList true "Add list"
// @Success 200 {object} models.ListResponse
// @Failure 500 {object} models.HTTPError
// @Failure 400 {object} models.HTTPError
// @Router /api/v1/lists [post]
func CreateList(c *fiber.Ctx) error {
	addList := &models.AddList{}
	if err := c.BodyParser(addList); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.HTTPError{
			Status: fiber.StatusBadRequest,
			Error:  err.Error(),
		})
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.HTTPError{
			Status: fiber.StatusInternalServerError,
			Error:  err.Error(),
		})
	}

	validate := utils.NewValidator()

	identityUser := c.Locals("user").(models.IdentityUser)

	list := models.List{
		ID:        uuid.New(),
		Name:      addList.Name,
		IsDefault: &addList.IsDefault,
		OwnerID:   identityUser.ID,
	}

	if err := validate.Struct(addList); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.HTTPError{
			Status: fiber.StatusBadRequest,
			Error:  err.Error(),
		})
	}

	listId, err := db.CreateList(list)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.HTTPError{
			Status: fiber.StatusInternalServerError,
			Error:  err.Error(),
		})
	}

	createdList, err := db.GetList(listId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.HTTPError{
			Status: fiber.StatusInternalServerError,
			Error:  err.Error(),
		})
	}

	return c.JSON(models.ListResponse{
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
// @Param list body models.AddList true "Update list"
// @Success 200 {object} models.ListResponse
// @Failure 500 {object} models.HTTPError
// @Failure 404 {object} models.HTTPError
// @Failure 400 {object} models.HTTPError
// @Router /api/v1/lists/{id} [put]
func UpdateList(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.HTTPError{
			Status: fiber.StatusBadRequest,
			Error:  err.Error(),
		})
	}

	addList := &models.AddList{}
	if err := c.BodyParser(addList); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.HTTPError{
			Status: fiber.StatusBadRequest,
			Error:  err.Error(),
		})
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.HTTPError{
			Status: fiber.StatusInternalServerError,
			Error:  err.Error(),
		})
	}

	identityUser := c.Locals("user").(models.IdentityUser)
	foundList, err := db.GetList(id)
	if err != nil && foundList.OwnerID != identityUser.ID {
		return c.Status(fiber.StatusNotFound).JSON(models.HTTPError{
			Status: fiber.StatusNotFound,
			Error:  err.Error(),
		})
	}

	foundList.Name = addList.Name
	foundList.IsDefault = &addList.IsDefault

	if err := db.UpdateList(foundList); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.HTTPError{
			Status: fiber.StatusInternalServerError,
			Error:  err.Error(),
		})
	}

	updatedList, err := db.GetList(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.HTTPError{
			Status: fiber.StatusInternalServerError,
			Error:  err.Error(),
		})
	}

	return c.JSON(models.ListResponse{
		Data: updatedList,
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
// @Success 200 {object} models.ListItemResponse
// @Failure 500 {object} models.HTTPError
// @Failure 404 {object} models.HTTPError
// @Failure 400 {object} models.HTTPError
// @Router /api/v1/lists/{list-id}/{item-id} [patch]
func AddItemToList(c *fiber.Ctx) error {
	log.Println("hej")
	listIdStr := c.Params("id")
	listId, err := uuid.Parse(listIdStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.HTTPError{
			Status: fiber.StatusBadRequest,
			Error:  err.Error(),
		})
	}

	itemIdStr := c.Params("itemId")
	itemId, err := uuid.Parse(itemIdStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.HTTPError{
			Status: fiber.StatusBadRequest,
			Error:  err.Error(),
		})
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.HTTPError{
			Status: fiber.StatusInternalServerError,
			Error:  err.Error(),
		})
	}

	identityUser := c.Locals("user").(models.IdentityUser)
	list, err := db.GetList(listId)
	if err != nil && list.OwnerID != identityUser.ID {
		return c.Status(fiber.StatusNotFound).JSON(models.HTTPError{
			Status: fiber.StatusNotFound,
			Error:  err.Error(),
		})
	}

	item, err := db.GetItem(itemId)
	if err != nil && list.OwnerID != identityUser.ID {
		return c.Status(fiber.StatusNotFound).JSON(models.HTTPError{
			Status: fiber.StatusNotFound,
			Error:  err.Error(),
		})
	}

	listItem, err := db.AddItemToList(list, item)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.HTTPError{
			Status: fiber.StatusInternalServerError,
			Error:  err.Error(),
		})
	}

	return c.JSON(models.ListItemResponse{
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
// @Param list body models.UpdateListItem true "Update list item"
// @Success 200 {object} models.ListItemResponse
// @Failure 500 {object} models.HTTPError
// @Failure 404 {object} models.HTTPError
// @Failure 400 {object} models.HTTPError
// @Router /api/v1/lists/{list-id}/{list-item-id} [put]
func UpdateListItem(c *fiber.Ctx) error {
	listIdStr := c.Params("id")
	listId, err := uuid.Parse(listIdStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.HTTPError{
			Status: fiber.StatusBadRequest,
			Error:  err.Error(),
		})
	}

	listItemIdStr := c.Params("listItemId")
	listItemId, err := uuid.Parse(listItemIdStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.HTTPError{
			Status: fiber.StatusBadRequest,
			Error:  err.Error(),
		})
	}

	updateListItem := &models.UpdateListItem{}
	if err := c.BodyParser(updateListItem); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.HTTPError{
			Status: fiber.StatusBadRequest,
			Error:  err.Error(),
		})
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.HTTPError{
			Status: fiber.StatusInternalServerError,
			Error:  err.Error(),
		})
	}

	identityUser := c.Locals("user").(models.IdentityUser)
	list, err := db.GetList(listId)
	if err != nil && list.OwnerID != identityUser.ID {
		return c.Status(fiber.StatusNotFound).JSON(models.HTTPError{
			Status: fiber.StatusNotFound,
			Error:  err.Error(),
		})
	}

	listItem, err := db.GetListItem(listItemId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.HTTPError{
			Status: fiber.StatusNotFound,
			Error:  err.Error(),
		})
	}

	listItem.Crossed = updateListItem.Crossed

	if err := db.UpdateListItem(listItem); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.HTTPError{
			Status: fiber.StatusInternalServerError,
			Error:  err.Error(),
		})
	}

	return c.JSON(models.ListItemResponse{
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
// @Failure 500 {object} models.HTTPError
// @Failure 404 {object} models.HTTPError
// @Failure 400 {object} models.HTTPError
// @Router /api/v1/lists/{list-id}/{list-item-id} [delete]
func RemoveItemFromList(c *fiber.Ctx) error {

	listIdStr := c.Params("id")
	listId, err := uuid.Parse(listIdStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.HTTPError{
			Status: fiber.StatusBadRequest,
			Error:  err.Error(),
		})
	}

	listItemIdStr := c.Params("listItemId")
	listItemId, err := uuid.Parse(listItemIdStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.HTTPError{
			Status: fiber.StatusBadRequest,
			Error:  err.Error(),
		})
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.HTTPError{
			Status: fiber.StatusInternalServerError,
			Error:  err.Error(),
		})
	}

	identityUser := c.Locals("user").(models.IdentityUser)
	list, err := db.GetList(listId)
	if err != nil && list.OwnerID != identityUser.ID {
		return c.Status(fiber.StatusNotFound).JSON(models.HTTPError{
			Status: fiber.StatusNotFound,
			Error:  err.Error(),
		})
	}

	if err := db.RemoveItemFromList(listItemId); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.HTTPError{
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
// @Failure 500 {object} models.HTTPError
// @Failure 400 {object} models.HTTPError
// @Router /api/v1/lists/{id} [delete]
func DeleteList(c *fiber.Ctx) error {

	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.HTTPError{
			Status: fiber.StatusBadRequest,
			Error:  err.Error(),
		})
	}

	list := &models.List{
		ID: id,
	}

	validate := utils.NewValidator()

	if err := validate.StructPartial(list, "id"); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.HTTPError{
			Status: fiber.StatusBadRequest,
			Error:  err.Error(),
		})
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.HTTPError{
			Status: fiber.StatusInternalServerError,
			Error:  err.Error(),
		})
	}

	foundList, err := db.GetList(list.ID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.HTTPError{
			Status: fiber.StatusNotFound,
			Error:  err.Error(),
		})
	}

	identityUser := c.Locals("user").(models.IdentityUser)
	if foundList.OwnerID != identityUser.ID {
		return c.Status(fiber.StatusNotFound).JSON(models.HTTPError{
			Status: fiber.StatusNotFound,
			Error:  err.Error(),
		})
	}

	if err := db.DeleteList(foundList); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.HTTPError{
			Status: fiber.StatusInternalServerError,
			Error:  err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
