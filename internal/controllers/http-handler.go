package controllers

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/kartikx04/chat/internal/database"
	"github.com/kartikx04/chat/internal/ws"
	"github.com/kartikx04/chat/pkg"
	"github.com/rs/cors"
)

func NewHTTPServer() *http.Server {
	r := http.NewServeMux()

	c := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:3000",
			"https://chat-0rnj.onrender.com",
			"https://banterrr.vercel.app",
		},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "Cookie"},
		ExposedHeaders:   []string{"Set-Cookie"},
		AllowCredentials: true,
	})

	r.HandleFunc("/health", func(w http.ResponseWriter, req *http.Request) {
		// Check DB
		sqlDB, err := database.DB.DB()
		if err != nil || sqlDB.Ping() != nil {
			slog.ErrorContext(req.Context(), "health: db unreachable")
			http.Error(w, `{"status":"error","db":false}`, http.StatusServiceUnavailable)
			return
		}

		// Check Redis
		if err := database.PingRedis(); err != nil {
			slog.ErrorContext(req.Context(), "health: redis unreachable", "error", err)
			http.Error(w, `{"status":"error","redis":false}`, http.StatusServiceUnavailable)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","db":true,"redis":true}`))
	})

	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWs(ws.HubInstance, w, r)
	})

	r.HandleFunc("/google-sso", GoogleSignOn)
	r.HandleFunc("/auth/google/callback", Callback)
	r.HandleFunc("/me", Me)
	r.HandleFunc("/contacts", contactListHandler)
	r.HandleFunc("/chat-history", chatHistoryHandler)
	r.HandleFunc("/add-contact", addContactHandler)
	r.HandleFunc("/verify-contact", verifyContactHandler)

	port := pkg.LoadFile("SERVER_PORT")

	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/ws" {
			r.ServeHTTP(w, req)
			return
		}
		c.Handler(r).ServeHTTP(w, req)
	})

	return &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: handler,
	}
}
