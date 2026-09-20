package route

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/kartikx04/chat/cmd/app/config"
	"github.com/kartikx04/chat/internal/api/middleware"
	"github.com/kartikx04/chat/internal/database"
	"github.com/kartikx04/chat/internal/repository"
	"gorm.io/gorm"
)

// http method grouping and permissions

func Register(cfg *config.App, timeout time.Duration, db *gorm.DB, r *chi.Mux, logger *slog.Logger) {
	authHandler := NewAuthHandler(cfg, timeout, db, logger)
	profileHandler := NewProfileHandler(cfg, timeout, db, logger)
	logoutHandler := NewLogoutHandler(cfg, timeout, db, logger)
	sessionRepo := repository.NewSessionRepository(db, logger)

	// public
	r.Get("/auth/success", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Login successful"))
	})

	r.Get("/health", func(w http.ResponseWriter, req *http.Request) {
		sqlDB, err := database.DB.DB()
		if err != nil || sqlDB.Ping() != nil {
			slog.ErrorContext(req.Context(), "health: db unreachable")
			http.Error(w, `{"status":"error","db":false}`, http.StatusServiceUnavailable)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	})

	r.Get("/auth/google-sso", authHandler.GoogleSignOn)
	r.Get("/auth/google/callback", authHandler.Callback)

	// protected
	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth(sessionRepo))

		r.Post("/logout", logoutHandler.Logout)
		r.Get("/profile", profileHandler.Profile)
	})
}
