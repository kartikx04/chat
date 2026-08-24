package route

import (
	"log/slog"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/kartikx04/chat/internal/api/controller"
	"github.com/kartikx04/chat/internal/repository"
	useCase "github.com/kartikx04/chat/internal/usecase"
	"github.com/kartikx04/chat/pkg"
	"gorm.io/gorm"
)

func NewLogoutRouter(timeout time.Duration, db *gorm.DB, r chi.Router, logger *slog.Logger) {
	ur := repository.NewUserRepository(db, logger)
	sr := repository.NewSessionRepository(db, logger)
	dc := &controller.LogoutController{
		AuthUseCase: useCase.NewAuthUseCase(pkg.OAuthgolang, ur, sr, timeout),
	}
	r.Get("/public/profile", dc.Logout)
}
