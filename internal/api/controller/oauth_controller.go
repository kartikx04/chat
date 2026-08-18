package controller

import (
	"log/slog"
	"net/http"

	"github.com/kartikx04/chat/internal/domain"
	"github.com/kartikx04/chat/pkg"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func init() {
	pkg.OAuthgolang = &oauth2.Config{
		RedirectURL:  pkg.LoadFile("REDIRECT_URL"),
		ClientID:     pkg.LoadFile("CLIENT_ID"),
		ClientSecret: pkg.LoadFile("CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
}

type AuthController struct {
	AuthUseCase domain.AuthUseCase
}

func (ac *AuthController) GoogleSignOn(res http.ResponseWriter, req *http.Request) {

	state, authURL, err := ac.AuthUseCase.InitiateGoogleOAuth(req.Context())
	if err != nil {
		slog.WarnContext(req.Context(), "error signing in", "error", err)
		http.Error(res, "internal server error", http.StatusBadRequest)
	}

	http.SetCookie(res, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Path:     "/",
		MaxAge:   300,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false,
	})

	slog.InfoContext(req.Context(), "oauth redirect initiated")
	http.Redirect(res, req, authURL, http.StatusTemporaryRedirect)
}

func (ac *AuthController) Callback(res http.ResponseWriter, req *http.Request) {
	state := req.FormValue("state")
	promptParam := req.URL.Query().Get("prompt")

	stateCookie, err := req.Cookie("oauth_state")

	if err != nil {
		slog.Warn("oauth_state cookie missing on callback",
			"likely_silent_auth", promptParam == "none",
		)
		http.Redirect(res, req, "/auth/google-sso", http.StatusFound)
		return
	}

	if state != stateCookie.Value {
		slog.Error("oauth_state MISMATCH - possible CSRF or stale link",
			"query_state", state,
			"cookie_state", stateCookie.Value,
		)
		http.Error(res, "internal server error", http.StatusBadRequest)
		return
	}

	http.SetCookie(res, &http.Cookie{
		Name:     "oauth_state",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	code := req.FormValue("code")
	if code == "" {
		slog.ErrorContext(req.Context(), "code not found in callback")
		http.Error(res, "internal server error", http.StatusBadRequest)
		return
	}

	sessionID, err := ac.AuthUseCase.FinaliseGoogleAuth(req.Context(), code)
	if err != nil {
		slog.ErrorContext(req.Context(), "failed to finalize google auth", "error", err)
		http.Error(res, "internal server error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(res, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	http.Redirect(res, req, pkg.LoadFile("FRONTEND_URL")+"/home", http.StatusFound)
}
