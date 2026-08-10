package route

import (
	"context"
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/kartikx04/chat/internal/api/controller"
	"gorm.io/gorm"
)

// http method grouping and permissions

func Register(ctx context.Context, db *gorm.DB, r *chi.Mux, logger *slog.Logger) {
	r.Get("/health", controller.Health)

	r.Get("/auth/google-sso", controller.GoogleSignOn)
	r.Get("/auth/google/callback", controller.Callback)

	NewProfileRouter(ctx, db, r, logger)
}
