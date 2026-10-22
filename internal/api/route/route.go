package route

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kartikx04/chat/cmd/app/config"
	"github.com/kartikx04/chat/internal/api/middleware"
	"github.com/kartikx04/chat/internal/database"
	"github.com/kartikx04/chat/internal/repository"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

// http method grouping and permissions

func Register(cfg *config.App, db *gorm.DB, r *chi.Mux, logger *slog.Logger, oauthConfig *oauth2.Config) {
	oauthHandler := NewOAuthController(cfg, db, logger, oauthConfig)
	profileHandler := NewProfileController(db, logger)
	logoutHandler := NewLogoutController(cfg, db, logger, oauthConfig)
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

	r.Get("/auth/google-sso", oauthHandler.StartGoogleOAuth)
	r.Get("/auth/google/callback", oauthHandler.CompleteGoogleOAuth)

	// protected
	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth(sessionRepo, logger))

		r.Post("/logout", logoutHandler.Logout)
		r.Get("/v1/users/me", profileHandler.Profile)
	})
}
