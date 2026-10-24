package route

import (
	"log/slog"

	"github.com/kartikx04/chat/internal/api/controller"
	"github.com/kartikx04/chat/internal/repository"
	"github.com/kartikx04/chat/internal/usecase"
	"gorm.io/gorm"
)

func NewProfileController(db *gorm.DB, logger *slog.Logger) *controller.ProfileController {
	ur := repository.NewUserRepository(db, logger)
	return controller.NewProfileController(usecase.NewProfileUsecase(ur), logger)
}
