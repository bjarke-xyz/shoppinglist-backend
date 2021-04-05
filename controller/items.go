package controller

import (
	"ShoppingList-Backend/model"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// FindItems godoc
// @Summary Find all items
// @Produce json
// @Success 200 {object} []model.Item
// @Router /items [get]
func (c *Controller) FindItems(ctx *gin.Context) {

	user := c.getOrCreateUser(ctx)

	var items []model.Item
	c.db.Where(&model.Item{OwnerID: user.ID}).Find(&items)

	c.ok(ctx, items)
}

// CreateItem godoc
// @Summary Create item
// @Accept json
// @Produce json
// @Param item body model.AddItem true "Add item"
// @Success 200 {object} model.Item
// @Failure 400 {object} controller.HttpError
// @Failure 500 {object} controller.HttpError
// @Router /items [post]
func (c *Controller) CreateItem(ctx *gin.Context) {

	var addItem model.AddItem
	if err := ctx.ShouldBindJSON(&addItem); err != nil {
		c.error(ctx, http.StatusBadRequest, err)
		return
	}
	user := c.getOrCreateUser(ctx)
	var existingItems []model.Item
	c.db.Where(&model.Item{OwnerID: user.ID}).Find(&existingItems)
	if err := addItem.Validation(existingItems); err != nil {
		c.error(ctx, http.StatusBadRequest, err)
		return
	}

	item := model.Item{Name: addItem.Name, OwnerID: user.ID}

	err := c.db.Create(&item).Error
	if err != nil {
		c.error(ctx, http.StatusInternalServerError, err)
		return
	}

	c.ok(ctx, item)
}

// UpdateItem godoc
// @Summary Update item
// @Accept json
// @Produce json
// @Param  id path int true "Item ID"
// @Param item body model.UpdateItem true "Update item"
// @Success 200 {object} model.Item
// @Failure 400 {object} controller.HttpError
// @Failure 500 {object} controller.HttpError
// @Router /items/{id} [put]
func (c *Controller) UpdateItem(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.error(ctx, http.StatusBadRequest, err)
		return
	}

	var updateItem model.UpdateItem
	if err := ctx.ShouldBindJSON(&updateItem); err != nil {
		c.error(ctx, http.StatusBadRequest, err)
		return
	}
	user := c.getOrCreateUser(ctx)
	var existingItems []model.Item
	c.db.Where(&model.Item{OwnerID: user.ID}).Find(&existingItems)
	if err := updateItem.Validation(existingItems); err != nil {
		c.error(ctx, http.StatusBadRequest, err)
		return
	}

	itemInDb := model.Item{}
	if err := c.db.Where(&model.Item{OwnerID: user.ID}).First(&itemInDb, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.error(ctx, http.StatusNotFound, err)
		} else {
			c.error(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	itemInDb.Name = updateItem.Name

	if err := c.db.Save(&itemInDb).Error; err != nil {
		c.error(ctx, http.StatusInternalServerError, err)
		return
	}

	c.ok(ctx, itemInDb)
}

// DeleteItem godoc
// @Summary delete an item
// @Accept json
// @Produce json
// @Param id path int true "Item ID"
// @Success 204
// @Failure 404 {object} controller.HttpError
// @Failure 500 {object} controller.HttpError
// @Router /items/{id} [delete]
func (c *Controller) DeleteItem(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.error(ctx, http.StatusBadRequest, err)
		return
	}

	user := c.getOrCreateUser(ctx)

	item := model.Item{}
	if err := c.db.Where(&model.Item{OwnerID: user.ID}).First(&item, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.error(ctx, http.StatusNotFound, err)
		} else {
			c.error(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	c.db.Delete(&item)
	c.db.Delete(&model.ListItem{ItemID: item.ID})
	c.status(ctx, http.StatusNoContent, nil)
}
