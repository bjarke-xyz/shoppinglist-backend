package server

import (
	"ShoppingList-Backend/pkg/config"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Server struct{}

func Start(a *fiber.App, cfg *config.Config, r *mux.Router) {
	zap.S().Infow("Api started", "SERVER_URL", cfg.GetServerUrl())
	srv := http.Server{
		Addr:         cfg.GetServerUrl(),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}
	if err := srv.ListenAndServe(); err != nil {
		zap.S().Errorf("API could not start: %v", err)
	}
}

func (s *Server) Respond(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	w.WriteHeader(status)
	if data != nil {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			zap.S().Errorw("Error encoding data", "error", err, "data", data)
		}
	}
}

type HTTPError struct {
	Status int    `json:"status"`
	Error  string `json:"error"`
}
