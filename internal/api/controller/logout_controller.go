package controller

import (
	"log/slog"
	"net/http"

	"github.com/kartikx04/chat/internal/domain"
)

type LogoutController struct {
	AuthUseCase domain.AuthUseCase
	Env         string
}

func (lc *LogoutController) Logout(res http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("session_id")
	if err != nil {
		slog.InfoContext(req.Context(), "logout requested without session")
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	err = lc.AuthUseCase.Logout(req.Context(), cookie.Value)
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
		Secure:   lc.Env == "production",
		SameSite: http.SameSiteLaxMode,
	})
	slog.InfoContext(req.Context(), "user logged out")
	http.Redirect(res, req, "/", http.StatusSeeOther)
}
