package controller

import (
	"log/slog"
	"net/http"

	"github.com/kartikx04/chat/cmd/app/config"
	"github.com/kartikx04/chat/internal/domain"

	"github.com/kartikx04/chat/pkg"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func InitOAuth(authCfg config.GoogleAuth) {
	pkg.OAuthgolang = &oauth2.Config{
		RedirectURL:  authCfg.RedirectURL,
		ClientID:     authCfg.ClientID,
		ClientSecret: authCfg.ClientSecret,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
}

type OAuthController struct {
	AuthUseCase domain.AuthUseCase
	Env         string
}

func (ac *OAuthController) GoogleSignOn(res http.ResponseWriter, req *http.Request) {
	state, authURL, err := ac.AuthUseCase.InitiateGoogleOAuth(req.Context())
	if err != nil {
		slog.WarnContext(req.Context(), "error signing in", "error", err)
		http.Error(res, "internal server error", http.StatusBadRequest)
		return
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

func (ac *OAuthController) Callback(res http.ResponseWriter, req *http.Request) {
	state := req.URL.Query().Get("state")
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

	code := req.URL.Query().Get("code")
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
		MaxAge:   72 * 60 * 60,
		HttpOnly: true,
		Secure:   ac.Env == "production",
		SameSite: http.SameSiteLaxMode,
	})

	http.Redirect(res, req, "/auth/success", http.StatusSeeOther)
}
