package controllers

import (
	"ShoppingList-Backend/app/models"
	"ShoppingList-Backend/pkg/utils"
	"ShoppingList-Backend/platform/database"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GetItems func gets all items for user
// @Description Get all items for user
// @Summary get all items for user
// @Tags items
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} models.ItemsResponse
// @Failure 500 {object} models.HTTPError
// @Failure 404 {object} models.HTTPError
// @Router /api/v1/items [get]
func GetItems(c *fiber.Ctx) error {
	db, err := database.OpenDBConnection()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.HTTPError{
			Status: fiber.StatusInternalServerError,
			Error:  err.Error(),
		})
	}

	identityUser := c.Locals("user").(models.IdentityUser)
	log.Printf("%v", identityUser)
	items, err := db.GetItems(identityUser.ID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.HTTPError{
			Status: fiber.StatusNotFound,
			Error:  err.Error(),
		})
	}

	return c.JSON(models.ItemsResponse{
		Data: items,
	})
}

// CreateItem func Create new item
// @Description Create new item
// @Summary Create new item
// @Tags items
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param item body models.AddItem true "Add item"
// @Success 200 {object} models.ItemResponse
// @Failure 500 {object} models.HTTPError
// @Failure 400 {object} models.HTTPError
// @Router /api/v1/items [post]
func CreateItem(c *fiber.Ctx) error {
	addItem := &models.AddItem{}
	if err := c.BodyParser(addItem); err != nil {
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

	item := &models.Item{
		ID:      uuid.New(),
		Name:    addItem.Name,
		OwnerID: identityUser.ID,
	}

	if err := validate.Struct(addItem); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.HTTPError{
			Status: fiber.StatusBadRequest,
			Error:  err.Error(),
		})
	}

	itemId, err := db.CreateItem(item)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.HTTPError{
			Status: fiber.StatusInternalServerError,
			Error:  err.Error(),
		})
	}

	createdItem, err := db.GetItem(itemId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.HTTPError{
			Status: fiber.StatusInternalServerError,
			Error:  err.Error(),
		})
	}

	return c.JSON(models.ItemResponse{
		Data: createdItem,
	})
}

// UpdateItem func Update item
// @Description Update item
// @Summary Update item
// @Tags items
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path string true "Item ID"
// @Param item body models.AddItem true "Update item"
// @Success 200 {object} models.ItemResponse
// @Failure 500 {object} models.HTTPError
// @Failure 404 {object} models.HTTPError
// @Failure 400 {object} models.HTTPError
// @Router /api/v1/items/{id} [put]
func UpdateItem(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.HTTPError{
			Status: fiber.StatusBadRequest,
			Error:  err.Error(),
		})
	}

	addItem := &models.AddItem{}
	if err := c.BodyParser(addItem); err != nil {
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

	foundItem, err := db.GetItem(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.HTTPError{
			Status: fiber.StatusNotFound,
			Error:  err.Error(),
		})
	}

	foundItem.Name = addItem.Name

	if err := db.UpdateItem(&foundItem); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.HTTPError{
			Status: fiber.StatusInternalServerError,
			Error:  err.Error(),
		})
	}

	updatedItem, err := db.GetItem(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.HTTPError{
			Status: fiber.StatusInternalServerError,
			Error:  err.Error(),
		})
	}

	return c.JSON(models.ItemResponse{
		Data: updatedItem,
	})
}

// DeleteItem func Delete item
// @Description Delete item
// @Summary Delete item
// @Tags items
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path string true "Item ID"
// @Success 204 {string} status "ok"
// @Failure 500 {object} models.HTTPError
// @Failure 400 {object} models.HTTPError
// @Router /api/v1/items/{id} [delete]
func DeleteItem(c *fiber.Ctx) error {

	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.HTTPError{
			Status: fiber.StatusBadRequest,
			Error:  err.Error(),
		})
	}

	item := &models.Item{
		ID: id,
	}

	validate := utils.NewValidator()

	if err := validate.StructPartial(item, "id"); err != nil {
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

	foundItem, err := db.GetItem(item.ID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.HTTPError{
			Status: fiber.StatusNotFound,
			Error:  err.Error(),
		})
	}

	identityUser := c.Locals("user").(models.IdentityUser)
	if foundItem.OwnerID != identityUser.ID {
		return c.Status(fiber.StatusNotFound).JSON(models.HTTPError{
			Status: fiber.StatusNotFound,
			Error:  err.Error(),
		})
	}

	if err := db.DeleteItem(&foundItem); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.HTTPError{
			Status: fiber.StatusInternalServerError,
			Error:  err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
