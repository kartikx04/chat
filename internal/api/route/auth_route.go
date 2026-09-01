package route

import (
	"log/slog"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/kartikx04/chat/cmd/app/config"
	"github.com/kartikx04/chat/internal/api/controller"
	"github.com/kartikx04/chat/internal/repository"
	useCase "github.com/kartikx04/chat/internal/usecase"
	"github.com/kartikx04/chat/pkg"
	"gorm.io/gorm"
)

func NewAuthRouter(cfg *config.App, timeout time.Duration, db *gorm.DB, r chi.Router, logger *slog.Logger) {
	controller.InitOAuth(cfg.GoogleAuth)

	ur := repository.NewUserRepository(db, logger)
	sr := repository.NewSessionRepository(db, logger)
	lc := &controller.OAuthController{
		AuthUseCase: useCase.NewAuthUseCase(pkg.OAuthgolang, ur, sr, timeout),
		Env:         cfg.Server.Env,
	}
	r.Get("/auth/google-sso", lc.GoogleSignOn)
	r.Get("/auth/google/callback", lc.Callback)
}
