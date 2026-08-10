package route

import (
	"context"
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/kartikx04/chat/internal/api/controller"
	"github.com/kartikx04/chat/internal/repository"
	useCase "github.com/kartikx04/chat/internal/usecase"
	"gorm.io/gorm"
)

func NewProfileRouter(ctx context.Context, db *gorm.DB, r chi.Router, logger *slog.Logger) {
	ur := repository.NewUserRepository(db, logger)
	dc := &controller.ProfileController{
		ProfileUseCase: useCase.NewProfileUseCase(ctx, ur),
	}
	r.Get("/public/Profile", dc.Profile)
}
