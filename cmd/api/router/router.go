package router

import (
	itemsHandler "ShoppingList-Backend/cmd/api/handlers/items"
	listsHandler "ShoppingList-Backend/cmd/api/handlers/lists"
	"ShoppingList-Backend/pkg/application"
	"ShoppingList-Backend/pkg/middleware"
	"net/http"

	socketio "github.com/googollee/go-socket.io"
	"github.com/gorilla/mux"
	"go.uber.org/zap"

	// http-swagger middleware
	httpSwagger "github.com/swaggo/http-swagger"

	// docs are generated by Swag CLI
	_ "ShoppingList-Backend/api"
)

func SocketIoRoutes(app *application.Application, r *mux.Router) {
	app.SocketIo.OnConnect("/", func(c socketio.Conn) error {
		c.SetContext("")
		zap.S().Infow("connected:", "id", c.ID())
		return nil
	})
	app.SocketIo.OnEvent("/", "message", func(s socketio.Conn, msg string) {
		zap.S().Infow("message", "msg", msg)
		s.Emit("reply", "have "+msg)
	})
	r.HandleFunc("/socket.io/", app.SocketIo.ServeHTTP)
}

func PrivateRoutes(app *application.Application, r *mux.Router) {
	apiV1 := r.PathPrefix("/api/v1").Subrouter()

	// Items
	items := apiV1.PathPrefix("/items").Subrouter()
	items.Use(middleware.JWTProtected(app.Cfg))
	items.HandleFunc("", itemsHandler.GetItems(app)).Methods("GET")
	items.HandleFunc("", itemsHandler.CreateItem(app)).Methods("POST")
	items.HandleFunc("/{id}", itemsHandler.UpdateItem(app)).Methods("PUT")
	items.HandleFunc("/{id}", itemsHandler.DeleteItem(app)).Methods("DELETE")

	// Lists
	lists := apiV1.PathPrefix("/lists").Subrouter()
	lists.Use(middleware.JWTProtected(app.Cfg))

	lists.HandleFunc("", listsHandler.GetLists(app)).Methods("GET")
	lists.HandleFunc("/default", listsHandler.GetDefaultList(app)).Methods("GET")
	lists.HandleFunc("", listsHandler.CreateList(app)).Methods("POST")
	lists.HandleFunc("/{id}", listsHandler.UpdateList(app)).Methods("PUT")
	lists.HandleFunc("/{id}/default", listsHandler.SetDefaultList(app)).Methods("PUT")
	lists.HandleFunc("/{id}", listsHandler.DeleteList(app)).Methods("DELETE")
	lists.HandleFunc("/{id}/items/crossed", listsHandler.ClearCrossedListItems(app)).Methods("DELETE")
	lists.HandleFunc("/{id}/items/{itemId}", listsHandler.AddItemToList(app)).Methods("POST")
	lists.HandleFunc("/{id}/items/{listItemId}", listsHandler.UpdateListItem(app)).Methods("PUT")
	lists.HandleFunc("/{id}/items/{listItemId}", listsHandler.RemoveItemFromList(app)).Methods("DELETE")

	// // SSE
	// sse := apiV1.PathPrefix("/sse").Subrouter()

	// sseTicket := sse.PathPrefix("/ticket").Subrouter()
	// sseTicket.Use(middleware.JWTProtected(app.Cfg))
	// sseTicket.HandleFunc("/", handlers.CreateSseTicket(app)).Methods("POST")

	// sse.HandleFunc("/events", handlers.SseEvents(app)).Methods("GET")
}

func SwaggerRoute(app *application.Application, r *mux.Router) {
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, r.URL.Opaque+"/swagger/index.html", http.StatusMovedPermanently)
	}).Methods("GET")
	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("#swagger-ui"),
	))
}
