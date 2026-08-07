package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kartikx04/chat/internal/api/route"
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

	route.Register(chi)

	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		c.Handler(chi).ServeHTTP(w, req)
	})

	port := pkg.LoadFile("SERVER_PORT")
	return &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: handler,
	}
}
