package lists

import (
	"ShoppingList-Backend/pkg/application"
	"ShoppingList-Backend/pkg/middleware"
	"ShoppingList-Backend/pkg/server"
	pkgWebsocket "ShoppingList-Backend/pkg/websocket"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func WsOnListChanges(hub *pkgWebsocket.Hub, app *application.Application) fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {

		appUser := middleware.WsGetAppUser(c)
		zap.S().Infow("new websocket connection", "user", appUser)
		defaultListAssociation, err := app.Queries.List.GetDefaultList(appUser)
		if err != nil {
			c.WriteJSON(server.HTTPError{
				Status: fiber.StatusNotFound,
				Error:  err.Error(),
			})
			return
		}

		client := pkgWebsocket.NewClient(hub, c)
		client.SessionInfo["id"] = defaultListAssociation.ListID
		client.ReadWritePump()
	})
}

func sessionInfoHasListId(listId uuid.UUID) pkgWebsocket.ClientFilter {
	return func(sessionInfo pkgWebsocket.SessionInfo) bool {
		id, ok := sessionInfo["id"]
		if !ok {
			zap.S().Infow("ID not in session", "id", id)
			return false
		}

		id, ok = id.(uuid.UUID)
		if !ok {
			return false
		}

		zap.S().Infow("Ids:", "id", id, "listId", listId)

		return id == listId
	}
}
