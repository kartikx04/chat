package controllers

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kartikx04/chat/internal/database"
	"github.com/kartikx04/chat/pkg"
	"github.com/rs/cors"
)

func NewHTTPServer() *http.Server {
	chi := chi.NewRouter()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "Cookie"},
		ExposedHeaders:   []string{"Set-Cookie"},
		AllowCredentials: true,
	})

	chi.Get("/health", Health)

	chi.Get("/auth/google-sso", GoogleSignOn)
	chi.Get("/auth/google/callback", Callback)

	chi.Get("/home", Home)

	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		c.Handler(chi).ServeHTTP(w, req)
	})

	port := pkg.LoadFile("SERVER_PORT")
	return &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: handler,
	}
}

func Health(w http.ResponseWriter, req *http.Request) {
	// Check DB
	sqlDB, err := database.DB.DB()
	if err != nil || sqlDB.Ping() != nil {
		slog.ErrorContext(req.Context(), "health: db unreachable")
		http.Error(w, `{"status":"error","db":false}`, http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
