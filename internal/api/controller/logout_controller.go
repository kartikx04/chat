package controller

import (
	"log/slog"
	"net/http"

	"github.com/kartikx04/chat/internal/domain"
	"github.com/kartikx04/chat/pkg"
)

type LogoutController struct {
	AuthUseCase domain.AuthUseCase
}

func (ac *LogoutController) Logout(res http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("session_id")
	if err != nil {
		slog.InfoContext(req.Context(), "logout requested without session")
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	err = ac.AuthUseCase.Logout(req.Context(), cookie.Value)
	if err != nil {
		slog.ErrorContext(req.Context(), "failed to logout", "error", err)
		http.Error(res, "internal server error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(res, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   pkg.LoadFile("ENV") == "production",
		SameSite: http.SameSiteLaxMode,
	})
	slog.InfoContext(req.Context(), "user logged out")
	http.Redirect(res, req, "/", http.StatusSeeOther)
}
