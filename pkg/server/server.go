package server

import (
	"ShoppingList-Backend/pkg/config"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"go.uber.org/zap"
)

type Server struct{}

func Start(cfg *config.Config, r *mux.Router) {

	handler := cors.New(cors.Options{
		AllowedOrigins: []string{"http://slv4.bjarke.xyz", "http://localhost:3000", "http://localhost:4200"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		Debug:          false,
	}).Handler(r)
	srv := http.Server{
		Addr:         cfg.GetServerUrl(),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      handler,
	}

	zap.S().Infow("Api started", "SERVER_URL", cfg.GetServerUrl())
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

func (s *Server) RespondError(w http.ResponseWriter, r *http.Request, status int, err error) {
	s.Respond(w, r, status, HTTPError{
		Status: status,
		Error:  err.Error(),
	})
}

func (s *Server) Decode(w http.ResponseWriter, r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

type HTTPError struct {
	Status int    `json:"status"`
	Error  string `json:"error"`
}
