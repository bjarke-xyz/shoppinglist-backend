package controller

import (
	"ShoppingList-Backend/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type Controller struct {
	db *gorm.DB
}

type HttpError struct {
	Code  int    `json:"code" example:"400"`
	Error string `json:"error" example:"Bad request"`
}

func NewController(db *gorm.DB) *Controller {
	return &Controller{db: db}
}

func (c *Controller) ok(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

func (c *Controller) status(ctx *gin.Context, code int, data interface{}) {
	ctx.JSON(code, gin.H{"data": data})
}

func (c *Controller) error(ctx *gin.Context, code int, err error) {
	e := HttpError{
		Code:  code,
		Error: err.Error(),
	}
	ctx.JSON(code, e)
}

func (c *Controller) getOrCreateUser(ctx *gin.Context) model.IdentityUser {
	user, exists := ctx.Get("user")
	if !exists {
		// TODO: Create user
	}
	return user.(model.IdentityUser)
}
