package handlers

import (
	"ShoppingList-Backend/internal/pkg/common"
	"ShoppingList-Backend/pkg/application"
	"ShoppingList-Backend/pkg/middleware"
	"fmt"
	"net/http"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func sseTicketKey(keyprefix string, ticket string) string {
	key := fmt.Sprintf("%v.API:SSETICKET:%v", keyprefix, ticket)
	return key
}

// CreateSseTicket func creates a new ticket for subscribing to sse events
// @Description creates a new ticket for subscribing to sse events
// @Summary create new ticket
// @Tags sse
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} common.Response{data=string}
// @Failure 500 {object} server.HTTPError
// @Failure 404 {object} server.HTTPError
// @Router /api/v1/sse/ticket [post]
func CreateSseTicket(app *application.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		redisConn := app.Redis.Get()
		defer redisConn.Close()
		uuid, err := uuid.NewUUID()
		if err != nil {
			app.Srv.RespondError(w, r, 500, fmt.Errorf("could not create ticket: %w", err))
			return
		}
		ticket := uuid.String()
		key := sseTicketKey(app.Cfg.GetRedisPrefix(), ticket)
		value := middleware.UserFromContext(r.Context()).ID
		_, err = redisConn.Do("SETEX", key, (30 * time.Second).Seconds(), value)
		if err != nil {
			app.Srv.RespondError(w, r, 500, fmt.Errorf("could not store ticket: %w", err))
			return
		}
		app.Srv.Respond(w, r, 200, common.Response{
			Data: ticket,
		})
	}
}

func SseEvents(app *application.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userTicket := r.URL.Query().Get("ticket")
		if userTicket == "" {
			app.Srv.RespondError(w, r, 400, fmt.Errorf("no ticket supplied"))
			return
		}
		key := sseTicketKey(app.Cfg.GetRedisPrefix(), userTicket)

		redisConn := app.Redis.Get()
		defer redisConn.Close()
		userId, err := redis.String(redisConn.Do("GET", key))
		if err != nil || userId == "" {
			app.Srv.RespondError(w, r, 500, fmt.Errorf("could not get ticket: %w", err))
			return
		}

		zap.S().Infow("tickets", "redisTicket", userId, "userTicket", userTicket)

		ctx := middleware.SetContextUser(r.Context(), userId)

		app.SseBroker.Handle(w, r.WithContext(ctx))
	}
}
