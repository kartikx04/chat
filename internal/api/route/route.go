package route

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/kartikx04/chat/cmd/app/config"
	"github.com/kartikx04/chat/internal/database"
	"gorm.io/gorm"
)

// http method grouping and permissions

func Register(cfg *config.App, timeout time.Duration, db *gorm.DB, r *chi.Mux, logger *slog.Logger) {
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

	r.Get("/auth/success", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Login successful"))
	})

	NewAuthRouter(cfg, timeout, db, r, logger)
	NewLogoutRouter(cfg, timeout, db, r, logger)
	NewProfileRouter(cfg, timeout, db, r, logger)
}
