package route

import (
	"log/slog"
	"time"

	"github.com/kartikx04/chat/cmd/app/config"
	"github.com/kartikx04/chat/internal/api/controller"
	"github.com/kartikx04/chat/internal/repository"
	"github.com/kartikx04/chat/internal/usecase"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

func NewOAuthController(cfg *config.App, timeout time.Duration, db *gorm.DB, logger *slog.Logger, oauthConfig *oauth2.Config) *controller.OAuthController {
	ur := repository.NewUserRepository(db, logger)
	sr := repository.NewSessionRepository(db, logger)
	return &controller.OAuthController{
		AuthUseCase: usecase.NewAuthUseCase(oauthConfig, ur, sr, timeout),
		Env:         cfg.Server.Env,
	}
}
