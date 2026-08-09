package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kartikx04/chat/cmd/app/config"
	"github.com/kartikx04/chat/internal/api/route"
	"github.com/rs/cors"
)

func NewHTTPServer(cfg *config.App) *http.Server {
	chi := chi.NewRouter()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "Cookie"},
		ExposedHeaders:   []string{"Set-Cookie"},
		AllowCredentials: true,
	})

	route.Register(chi)

	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		c.Handler(chi).ServeHTTP(w, req)
	})

	return &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: handler,
	}
}
