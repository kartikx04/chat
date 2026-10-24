package controller

import (
	"log/slog"
	"net/http"

	"github.com/kartikx04/chat/internal/domain"
)

type OAuthController struct {
	authUsecase domain.AuthUsecase
	env         string
	logger      *slog.Logger
}

func NewOAuthController(authUsecase domain.AuthUsecase, env string, logger *slog.Logger) *OAuthController {
	return &OAuthController{
		authUsecase: authUsecase,
		env:         env,
		logger:      logger,
	}
}

func (ac *OAuthController) StartGoogleOAuth(res http.ResponseWriter, req *http.Request) {
	state, authURL, err := ac.authUsecase.InitiateGoogleOAuth(req.Context())
	if err != nil {
		ac.logger.WarnContext(req.Context(), "error signing in", "error", err)
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

	ac.logger.InfoContext(req.Context(), "oauth redirect initiated")
	http.Redirect(res, req, authURL, http.StatusTemporaryRedirect)
}

func (ac *OAuthController) CompleteGoogleOAuth(res http.ResponseWriter, req *http.Request) {
	state := req.URL.Query().Get("state")
	promptParam := req.URL.Query().Get("prompt")

	stateCookie, err := req.Cookie("oauth_state")

	if err != nil {
		ac.logger.WarnContext(req.Context(), "oauth_state cookie missing on callback",
			"likely_silent_auth", promptParam == "none",
		)
		http.Redirect(res, req, "/auth/google-sso", http.StatusFound)
		return
	}

	if state != stateCookie.Value {
		ac.logger.ErrorContext(req.Context(), "oauth_state MISMATCH")
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
		ac.logger.ErrorContext(req.Context(), "code not found in callback")
		http.Error(res, "internal server error", http.StatusBadRequest)
		return
	}

	token, err := ac.authUsecase.CompleteGoogleOAuth(req.Context(), code)
	if err != nil {
		ac.logger.ErrorContext(req.Context(), "failed to finalize google auth", "error", err)
		http.Error(res, "internal server error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(res, &http.Cookie{
		Name:     "session_id",
		Value:    token,
		Path:     "/",
		MaxAge:   72 * 60 * 60,
		HttpOnly: true,
		Secure:   ac.env == "production",
		SameSite: http.SameSiteLaxMode,
	})

	http.Redirect(res, req, "/auth/success", http.StatusSeeOther)
}
