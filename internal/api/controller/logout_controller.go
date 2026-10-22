package controller

import (
	"log/slog"
	"net/http"

	"github.com/kartikx04/chat/internal/domain"
)

type LogoutController struct {
	authUsecase domain.AuthUsecase
	env         string
	logger      *slog.Logger
}

func NewLogoutController(authUsecase domain.AuthUsecase, env string, logger *slog.Logger) *LogoutController {
	return &LogoutController{
		authUsecase: authUsecase,
		env:         env,
		logger:      logger,
	}
}

func (lc *LogoutController) Logout(res http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("session_id")
	if err != nil {
		lc.logger.InfoContext(req.Context(), "logout requested without session")
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	err = lc.authUsecase.Logout(req.Context(), cookie.Value)
	if err != nil {
		lc.logger.ErrorContext(req.Context(), "failed to logout", "error", err)
		http.Error(res, "internal server error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(res, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   lc.env == "production",
		SameSite: http.SameSiteLaxMode,
	})
	lc.logger.InfoContext(req.Context(), "user logged out")
	http.Redirect(res, req, "/", http.StatusSeeOther)
}
