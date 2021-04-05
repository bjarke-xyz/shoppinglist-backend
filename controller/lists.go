package controller

import (
	"ShoppingList-Backend/model"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// FindLists godoc
// @Summary Find all lists
// @Produce json
// @Success 200 {object} []model.List
// @Router /lists [get]
func (c *Controller) FindLists(ctx *gin.Context) {
	user := c.getOrCreateUser(ctx)

	var lists []model.List
	c.db.Where(&model.List{OwnerID: user.ID}).Find(&lists)

	for i, list := range lists {
		var listItems []model.ListItem
		err := c.db.Model(&list).Association("Items").Find(&listItems).Error
		if err != nil {
			c.error(ctx, http.StatusInternalServerError, err)
			return
		}
		for j, item := range listItems {
			var assocItem model.Item
			err = c.db.Model(&item).Association("Item").Find(&assocItem).Error
			if err != nil {
				c.error(ctx, http.StatusInternalServerError, err)
				return
			}
			listItems[j].Item = assocItem
		}
		lists[i].Items = make([]model.ListItem, 0)
		lists[i].Items = append(lists[i].Items, listItems...)
	}

	c.ok(ctx, lists)
}

// CreateList godoc
// @Summary Create list
// @Accept json
// @Produce json
// @Param list body model.AddList true "Add list"
// @Success 200 {object} model.List
// @Failure 400 {object} controller.HttpError
// @Failure 500 {object} controller.HttpError
// @Router /lists [post]
func (c *Controller) CreateList(ctx *gin.Context) {
	var addList model.AddList
	if err := ctx.ShouldBindJSON(&addList); err != nil {
		c.error(ctx, http.StatusBadRequest, err)
		return
	}

	user := c.getOrCreateUser(ctx)

	var existingLists []model.List
	c.db.Where(&model.List{OwnerID: user.ID}).Find(&existingLists)

	if err := addList.Validation(existingLists); err != nil {
		c.error(ctx, http.StatusBadRequest, err)
		return
	}

	list := model.List{Name: addList.Name, OwnerID: user.ID, Default: len(existingLists) == 0}

	err := c.db.Create(&list).Error
	if err != nil {
		c.error(ctx, http.StatusInternalServerError, err)
		return
	}

	c.ok(ctx, list)
}

// UpdateList godoc
// @Summary Update list
// @Accept json
// @Produce json
// @Param  id path int true "List ID"
// @Param list body model.UpdateList true "Update list"
// @Success 200 {object} model.List
// @Failure 400 {object} controller.HttpError
// @Failure 500 {object} controller.HttpError
// @Router /lists/{id} [put]
func (c *Controller) UpdateList(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.error(ctx, http.StatusBadRequest, err)
		return
	}

	var updateList model.UpdateList
	if err := ctx.ShouldBindJSON(&updateList); err != nil {
		c.error(ctx, http.StatusBadRequest, err)
		return
	}

	log.Println(updateList)

	user := c.getOrCreateUser(ctx)

	var existingLists []model.List
	c.db.Where(&model.List{OwnerID: user.ID}).Find(&existingLists)

	if err := updateList.Validation(existingLists); err != nil {
		c.error(ctx, http.StatusBadRequest, err)
		return
	}

	listInDb := model.List{}
	if err := c.db.Where(&model.List{OwnerID: user.ID}).First(&listInDb, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.error(ctx, http.StatusNotFound, err)
		} else {
			c.error(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	listInDb.Name = updateList.Name
	listInDb.Default = updateList.Default

	if updateList.Default {
		if err := c.db.Model(&model.List{}).Where("owner_id = ?", user.ID).Update("default", "false").Error; err != nil {
			c.error(ctx, http.StatusInternalServerError, err)
			return
		}
	}

	if err := c.db.Save(&listInDb).Error; err != nil {
		c.error(ctx, http.StatusInternalServerError, err)
		return
	}

	c.ok(ctx, listInDb)
}

// AddItemToList godoc
// @Summary Add an item to a list
// @Accept json
// @Produce json
// @Param listId path int true "List ID"
// @Param itemId path int true "Item ID"
// @Success 200 {object} model.ListItem
// @Failure 404 {object} controller.HttpError
// @Failure 500 {object} controller.HttpError
// @Router /lists/add/{listId}/{itemId}  [put]
func (c *Controller) AddItemToList(ctx *gin.Context) {
	listIdStr := ctx.Param("listId")
	listId, err := strconv.Atoi(listIdStr)
	if err != nil {
		c.error(ctx, http.StatusBadRequest, err)
		return
	}

	user := c.getOrCreateUser(ctx)

	listInDb := model.List{}
	if err := c.db.Where(&model.List{OwnerID: user.ID}).First(&listInDb, listId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.error(ctx, http.StatusNotFound, err)
		} else {
			c.error(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	itemIdStr := ctx.Param("itemId")
	itemId, err := strconv.Atoi(itemIdStr)
	if err != nil {
		c.error(ctx, http.StatusBadRequest, err)
		return
	}
	itemInDb := model.Item{}
	if err := c.db.Where(&model.Item{OwnerID: user.ID}).First(&itemInDb, itemId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.error(ctx, http.StatusNotFound, err)
		} else {
			c.error(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	listItem := model.ListItem{ListID: listInDb.ID, ItemID: itemInDb.ID}

	err = c.db.Save(&listItem).Error
	if err != nil {
		c.error(ctx, http.StatusInternalServerError, err)
		return
	}

	c.ok(ctx, listItem)
}

// RemoveItemFromList godoc
// @Summary Remove an item from a list
// @Accept json
// @Produce json
// @Param listId path int true "List ID"
// @Param itemId path int true "Item ID"
// @Success 204
// @Failure 404 {object} controller.HttpError
// @Failure 500 {object} controller.HttpError
// @Router /lists/remove/{listId}/{itemId}  [put]
func (c *Controller) RemoveItemFromList(ctx *gin.Context) {
	listIdStr := ctx.Param("listId")
	listId, err := strconv.Atoi(listIdStr)
	if err != nil {
		c.error(ctx, http.StatusBadRequest, err)
		return
	}
	itemIdStr := ctx.Param("itemId")
	itemId, err := strconv.Atoi(itemIdStr)
	if err != nil {
		c.error(ctx, http.StatusBadRequest, err)
		return
	}

	c.db.Delete(&model.ListItem{ListID: uint(listId), ItemID: uint(itemId)})
	c.status(ctx, http.StatusNoContent, nil)
}

// DeleteList godoc
// @Summary delete a list
// @Accept json
// @Produce json
// @Param id path int true "List ID"
// @Success 204
// @Failure 404 {object} controller.HttpError
// @Failure 500 {object} controller.HttpError
// @Router /lists/{id} [delete]
func (c *Controller) DeleteList(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.error(ctx, http.StatusBadRequest, err)
		return
	}

	user := c.getOrCreateUser(ctx)

	list := model.List{}
	if err := c.db.Where(&model.List{OwnerID: user.ID}).First(&list, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.error(ctx, http.StatusNotFound, err)
		} else {
			c.error(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	c.db.Delete(&list)
	c.db.Delete(&model.ListItem{ListID: list.ID})

	c.status(ctx, http.StatusNoContent, nil)
}
