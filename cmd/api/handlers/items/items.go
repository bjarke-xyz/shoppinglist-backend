package items

import (
	"ShoppingList-Backend/internal/pkg/item"
	"ShoppingList-Backend/pkg/application"
	"ShoppingList-Backend/pkg/middleware"
	"ShoppingList-Backend/pkg/server"
	"ShoppingList-Backend/pkg/utils"

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
// @Success 200 {object} item.ItemsResponse
// @Failure 500 {object} server.HTTPError
// @Failure 404 {object} server.HTTPError
// @Router /api/v1/items [get]
func GetItems(app *application.Application) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		appUser := middleware.GetAppUser(c)
		items, err := app.Queries.Item.GetItems(appUser.ID)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(server.HTTPError{
				Status: fiber.StatusNotFound,
				Error:  err.Error(),
			})
		}

		return c.JSON(item.ItemsResponse{
			Data: items,
		})
	}
}

// CreateItem func Create new item
// @Description Create new item
// @Summary Create new item
// @Tags items
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param item body item.AddItem true "Add item"
// @Success 200 {object} item.ItemResponse
// @Failure 500 {object} server.HTTPError
// @Failure 400 {object} server.HTTPError
// @Router /api/v1/items [post]
func CreateItem(app *application.Application) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		addItem := &item.AddItem{}
		if err := c.BodyParser(addItem); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(server.HTTPError{
				Status: fiber.StatusBadRequest,
				Error:  err.Error(),
			})
		}

		validate := utils.NewValidator()

		appUser := middleware.GetAppUser(c)

		itemToCreate := &item.Item{
			ID:      uuid.New(),
			Name:    addItem.Name,
			OwnerID: appUser.ID,
		}

		if err := validate.Struct(addItem); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(server.HTTPError{
				Status: fiber.StatusBadRequest,
				Error:  err.Error(),
			})
		}

		itemId, err := app.Queries.Item.CreateItem(itemToCreate)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(server.HTTPError{
				Status: fiber.StatusInternalServerError,
				Error:  err.Error(),
			})
		}

		createdItem, err := app.Queries.Item.GetItem(itemId)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(server.HTTPError{
				Status: fiber.StatusInternalServerError,
				Error:  err.Error(),
			})
		}

		return c.JSON(item.ItemResponse{
			Data: createdItem,
		})
	}
}

// UpdateItem func Update item
// @Description Update item
// @Summary Update item
// @Tags items
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path string true "Item ID"
// @Param item body item.AddItem true "Update item"
// @Success 200 {object} item.ItemResponse
// @Failure 500 {object} server.HTTPError
// @Failure 404 {object} server.HTTPError
// @Failure 400 {object} server.HTTPError
// @Router /api/v1/items/{id} [put]
func UpdateItem(app *application.Application) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(server.HTTPError{
				Status: fiber.StatusBadRequest,
				Error:  err.Error(),
			})
		}

		addItem := &item.AddItem{}
		if err := c.BodyParser(addItem); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(server.HTTPError{
				Status: fiber.StatusBadRequest,
				Error:  err.Error(),
			})
		}

		foundItem, err := app.Queries.Item.GetItem(id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(server.HTTPError{
				Status: fiber.StatusNotFound,
				Error:  err.Error(),
			})
		}

		foundItem.Name = addItem.Name

		if err := app.Queries.Item.UpdateItem(&foundItem); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(server.HTTPError{
				Status: fiber.StatusInternalServerError,
				Error:  err.Error(),
			})
		}

		updatedItem, err := app.Queries.Item.GetItem(id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(server.HTTPError{
				Status: fiber.StatusInternalServerError,
				Error:  err.Error(),
			})
		}

		return c.JSON(item.ItemResponse{
			Data: updatedItem,
		})
	}
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
// @Failure 500 {object} server.HTTPError
// @Failure 400 {object} server.HTTPError
// @Router /api/v1/items/{id} [delete]
func DeleteItem(app *application.Application) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(server.HTTPError{
				Status: fiber.StatusBadRequest,
				Error:  err.Error(),
			})
		}

		item := &item.Item{
			ID: id,
		}

		validate := utils.NewValidator()

		if err := validate.StructPartial(item, "id"); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(server.HTTPError{
				Status: fiber.StatusBadRequest,
				Error:  err.Error(),
			})
		}

		foundItem, err := app.Queries.Item.GetItem(item.ID)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(server.HTTPError{
				Status: fiber.StatusNotFound,
				Error:  err.Error(),
			})
		}

		appUser := middleware.GetAppUser(c)
		if foundItem.OwnerID != appUser.ID {
			return c.Status(fiber.StatusNotFound).JSON(server.HTTPError{
				Status: fiber.StatusNotFound,
				Error:  err.Error(),
			})
		}

		if err := app.Queries.Item.DeleteItem(&foundItem); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(server.HTTPError{
				Status: fiber.StatusInternalServerError,
				Error:  err.Error(),
			})
		}

		return c.SendStatus(fiber.StatusNoContent)

	}

}
