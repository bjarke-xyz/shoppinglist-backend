package lists

import (
	"ShoppingList-Backend/pkg/application"
	"ShoppingList-Backend/pkg/middleware"
	"ShoppingList-Backend/pkg/server"
	pkgWebsocket "ShoppingList-Backend/pkg/websocket"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func WsOnListChanges(hub *pkgWebsocket.Hub, app *application.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("w's type is %T\n", w)
		upgrader := websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			zap.S().Errorf("could not upgrade websocket: %v", err)
			return
		}

		user := middleware.UserFromContext(r.Context())
		zap.S().Infow("new websocket connection", "user", user)

		defaultListAssociation, cErr := app.Controllers.List.GetDefaultList(user)
		if cErr != nil {
			conn.WriteJSON(server.HTTPError{
				Status: cErr.StatusCode,
				Error:  cErr.Err.Error(),
			})
			return
		}

		client := pkgWebsocket.NewClient(hub, conn)
		client.SessionInfo["id"] = defaultListAssociation.ListID
		client.ReadWritePump()

	}
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
