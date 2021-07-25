package controller

import (
	"ShoppingList-Backend/internal/pkg/list"
	"ShoppingList-Backend/pkg/application"
	"ShoppingList-Backend/pkg/middleware"
	"ShoppingList-Backend/pkg/server"
	"ShoppingList-Backend/pkg/utils"
	"database/sql"

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
func GetLists(app *application.Application) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		appUser := middleware.GetAppUser(c)

		lists, err := app.Queries.List.GetLists(appUser)
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
func GetDefaultList(app *application.Application) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		appUser := middleware.GetAppUser(c)
		defaultList, err := app.Queries.List.GetDefaultList(appUser)
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
func CreateList(app *application.Application) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		addList := &list.AddList{}
		if err := c.BodyParser(addList); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(server.HTTPError{
				Status: fiber.StatusBadRequest,
				Error:  err.Error(),
			})
		}

		validate := utils.NewValidator()

		appUser := middleware.GetAppUser(c)

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

		listId, err := app.Queries.List.CreateList(listToCreate)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(server.HTTPError{
				Status: fiber.StatusInternalServerError,
				Error:  err.Error(),
			})
		}

		createdList, err := app.Queries.List.GetList(listId, appUser)
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
func UpdateList(app *application.Application) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
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

		appUser := middleware.GetAppUser(c)
		foundList, err := app.Queries.List.GetList(id, appUser)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(server.HTTPError{
				Status: fiber.StatusNotFound,
				Error:  err.Error(),
			})
		}

		foundList.Name = addList.Name

		if err := app.Queries.List.UpdateList(foundList); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(server.HTTPError{
				Status: fiber.StatusInternalServerError,
				Error:  err.Error(),
			})
		}

		updatedList, err := app.Queries.List.GetList(id, appUser)
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
func SetDefaultList(app *application.Application) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(server.HTTPError{
				Status: fiber.StatusBadRequest,
				Error:  err.Error(),
			})
		}

		appUser := middleware.GetAppUser(c)
		foundList, err := app.Queries.List.GetList(id, appUser)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(server.HTTPError{
				Status: fiber.StatusNotFound,
				Error:  err.Error(),
			})
		}

		defaultList, err := app.Queries.List.SetDefaultList(appUser, foundList)
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
func AddItemToList(app *application.Application) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
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

		appUser := middleware.GetAppUser(c)
		foundList, err := app.Queries.List.GetList(listId, appUser)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(server.HTTPError{
				Status: fiber.StatusNotFound,
				Error:  err.Error(),
			})
		}

		item, err := app.Queries.Item.GetItem(itemId)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(server.HTTPError{
				Status: fiber.StatusNotFound,
				Error:  err.Error(),
			})
		}

		listItem, err := app.Queries.List.AddItemToList(foundList, item)
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
func UpdateListItem(app *application.Application) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
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

		appUser := middleware.GetAppUser(c)
		if _, err := app.Queries.List.GetList(listId, appUser); err != nil {
			return c.Status(fiber.StatusNotFound).JSON(server.HTTPError{
				Status: fiber.StatusNotFound,
				Error:  err.Error(),
			})
		}

		listItem, err := app.Queries.List.GetListItem(listItemId)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(server.HTTPError{
				Status: fiber.StatusNotFound,
				Error:  err.Error(),
			})
		}

		listItem.Crossed = updateListItem.Crossed

		if err := app.Queries.List.UpdateListItem(listItem); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(server.HTTPError{
				Status: fiber.StatusInternalServerError,
				Error:  err.Error(),
			})
		}

		return c.JSON(list.ListItemResponse{
			Data: listItem,
		})
	}
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
func RemoveItemFromList(app *application.Application) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
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

		appUser := middleware.GetAppUser(c)
		if _, err := app.Queries.List.GetList(listId, appUser); err != nil {
			return c.Status(fiber.StatusNotFound).JSON(server.HTTPError{
				Status: fiber.StatusNotFound,
				Error:  err.Error(),
			})
		}

		if err := app.Queries.List.RemoveItemFromList(listItemId); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(server.HTTPError{
				Status: fiber.StatusInternalServerError,
				Error:  err.Error(),
			})
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
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
func DeleteList(app *application.Application) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
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

		appUser := middleware.GetAppUser(c)
		foundList, err := app.Queries.List.GetList(list.ID, appUser)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(server.HTTPError{
				Status: fiber.StatusNotFound,
				Error:  err.Error(),
			})
		}

		if err := app.Queries.List.DeleteList(foundList); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(server.HTTPError{
				Status: fiber.StatusInternalServerError,
				Error:  err.Error(),
			})
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}
