package route

import (
	"log/slog"
	"time"

	"github.com/kartikx04/chat/cmd/app/config"
	"github.com/kartikx04/chat/internal/api/controller"
	"github.com/kartikx04/chat/internal/repository"
	useCase "github.com/kartikx04/chat/internal/usecase"
	"gorm.io/gorm"
)

func NewProfileHandler(cfg *config.App, timeout time.Duration, db *gorm.DB, logger *slog.Logger) *controller.ProfileController {
	ur := repository.NewUserRepository(db, logger)
	sr := repository.NewSessionRepository(db, logger)
	return &controller.ProfileController{
		ProfileUseCase: useCase.NewProfileUseCase(timeout, ur),
		SessionRepo:    sr,
	}
}
