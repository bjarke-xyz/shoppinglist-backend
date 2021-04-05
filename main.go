package main

import (
	"log"
	"os"

	"ShoppingList-Backend/controller"
	_ "ShoppingList-Backend/docs"
	"ShoppingList-Backend/model"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

// @title Shopping List API
// @version 1.0
// @description API for the Shopping List application
// @BasePath /api/v1
func main() {
	r := gin.Default()

	db, err := model.Init()
	if err != nil {
		log.Fatal(err)
	}

	err = godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	jwksURL := os.Getenv("jwks_uri")
	keycloakURL := os.Getenv("keycloak_url")

	c := controller.NewController(db)

	r.Use(controller.CORSMiddleware())
	v1 := r.Group("api/v1")
	{
		items := v1.Group("/items")
		{
			items.Use(controller.JWTAuthorized(jwksURL, keycloakURL))
			items.GET("", c.FindItems)
			items.POST("", c.CreateItem)
			items.PUT("/:id", c.UpdateItem)
			items.DELETE("/:id", c.DeleteItem)
		}

		lists := v1.Group("/lists")
		{
			lists.Use(controller.JWTAuthorized(jwksURL, keycloakURL))
			lists.GET("", c.FindLists)
			lists.POST("", c.CreateList)
			lists.PATCH("/add/:listId/:itemId", c.AddItemToList)
			lists.PATCH("/remove/:listId/:itemId", c.RemoveItemFromList)
			lists.PUT("/:id", c.UpdateList)
			lists.DELETE("/:id", c.DeleteList)
		}
	}

	// /swagger/index.html
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run()
}
