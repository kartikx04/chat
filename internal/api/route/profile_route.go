package route

import (
	"log/slog"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/kartikx04/chat/internal/api/controller"
	"github.com/kartikx04/chat/internal/repository"
	useCase "github.com/kartikx04/chat/internal/usecase"
	"gorm.io/gorm"
)

func NewProfileRouter(timeout time.Duration, db *gorm.DB, r chi.Router, logger *slog.Logger) {
	ur := repository.NewUserRepository(db, logger)
	dc := &controller.ProfileController{
		ProfileUseCase: useCase.NewProfileUseCase(timeout, ur),
	}
	r.Get("/public/profile", dc.Profile)
}
