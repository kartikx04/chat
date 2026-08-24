package route

import (
	"log/slog"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/kartikx04/chat/internal/api/controller"
	"gorm.io/gorm"
)

// http method grouping and permissions

func Register(timeout time.Duration, db *gorm.DB, r *chi.Mux, logger *slog.Logger) {
	r.Get("/health", controller.Health)

	NewAuthRouter(timeout, db, r, logger)
	NewLogoutRouter(timeout, db, r, logger)
	NewProfileRouter(timeout, db, r, logger)
}
