package api

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kartikx04/chat/cmd/app/config"
	"github.com/kartikx04/chat/internal/api/middleware"
	"github.com/rs/cors"
)

func NewHTTPServer(cfg *config.App, router *chi.Mux, logger *slog.Logger) *http.Server {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "Cookie"},
		ExposedHeaders:   []string{"Set-Cookie"},
		AllowCredentials: true,
	})

	handler := middleware.RequestID(logger)(c.Handler(router))

	return &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: handler,
	}
}
