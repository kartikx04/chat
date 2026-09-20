package route

import (
	"log/slog"
	"time"

	"github.com/kartikx04/chat/cmd/app/config"
	"github.com/kartikx04/chat/internal/api/controller"
	"github.com/kartikx04/chat/internal/repository"
	useCase "github.com/kartikx04/chat/internal/usecase"
	"github.com/kartikx04/chat/pkg"
	"gorm.io/gorm"
)

func NewAuthHandler(cfg *config.App, timeout time.Duration, db *gorm.DB, logger *slog.Logger) *controller.OAuthController {
	controller.InitOAuth(cfg.GoogleAuth)

	ur := repository.NewUserRepository(db, logger)
	sr := repository.NewSessionRepository(db, logger)
	return &controller.OAuthController{
		AuthUseCase: useCase.NewAuthUseCase(pkg.OAuthgolang, ur, sr, timeout),
		Env:         cfg.Server.Env,
	}
}
