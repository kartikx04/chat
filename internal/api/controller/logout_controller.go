package controller

import (
	"log/slog"
	"net/http"
)

func Logout(res http.ResponseWriter, req *http.Request) {
	http.SetCookie(res, &http.Cookie{
		Name:     "userSession",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
	http.SetCookie(res, &http.Cookie{
		Name:     "oauth_token_raw",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
	slog.InfoContext(req.Context(), "user logged out")
	http.Redirect(res, req, "/", http.StatusSeeOther)
}
