package controller

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/kartikx04/chat/pkg"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type UseCase interface {
	SignUp(context.Context) error
	SignIn(context.Context) (string, error)
}

func init() {
	pkg.OAuthgolang = &oauth2.Config{
		RedirectURL:  pkg.LoadFile("REDIRECT_URL"),
		ClientID:     pkg.LoadFile("CLIENT_ID"),
		ClientSecret: pkg.LoadFile("CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
}

func GoogleSignOn(res http.ResponseWriter, req *http.Request) {
	state, err := pkg.GenerateStateToken()
	if err != nil {
		slog.ErrorContext(req.Context(), "failed to generate oauth state token", "error", err)
		http.Error(res, "internal server error", http.StatusInternalServerError)
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

	authURL := pkg.OAuthgolang.AuthCodeURL(state)
	slog.InfoContext(req.Context(), "oauth redirect initiated")
	http.Redirect(res, req, authURL, http.StatusTemporaryRedirect)
}

func Callback(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	state := req.FormValue("state")
	promptParam := req.URL.Query().Get("prompt")

	stateCookie, cookieErr := req.Cookie("oauth_state")

	if cookieErr != nil {
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

	token, err := pkg.OAuthgolang.Exchange(ctx, code)
	if err != nil {
		slog.ErrorContext(req.Context(), "failed to exchange token", "error", err)
		http.Error(res, "internal server error", http.StatusInternalServerError)
		return
	}

	// SECURITY: tokenBytes is not encrypted
	tokenBytes, err := json.Marshal(token)
	if err != nil {
		slog.Error("failed to marshal token payload", "error", err)
		http.Error(res, "internal server error", http.StatusInternalServerError)
		return
	}
	encoded := base64.RawURLEncoding.EncodeToString(tokenBytes)

	http.SetCookie(res, &http.Cookie{
		Name:     "oauth_token_raw",
		Value:    encoded,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	slog.Info("callback succeeded, redirecting to /home")
	http.Redirect(res, req, "/home", http.StatusFound)
}
