package controller

import (
	"log/slog"
	"net/http"

	"github.com/kartikx04/chat/internal/domain"
	"github.com/kartikx04/chat/pkg"
)

type OAuthController struct {
	authUseCase domain.AuthUsecase
	env         string
	logger      *slog.Logger
}

func NewOAuthController(authUseCase domain.AuthUsecase, env string, logger *slog.Logger) *OAuthController {
	return &OAuthController{
		authUseCase: authUseCase,
		env:         env,
		logger:      logger,
	}
}

func (ac *OAuthController) StartGoogleOAuth(res http.ResponseWriter, req *http.Request) {
	log := pkg.RequestLogger(req.Context(), ac.logger)

	state, authURL, err := ac.authUseCase.InitiateGoogleOAuth(req.Context())
	if err != nil {
		log.WarnContext(req.Context(), "error signing in", "error", err)
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

	log.InfoContext(req.Context(), "oauth redirect initiated")
	http.Redirect(res, req, authURL, http.StatusTemporaryRedirect)
}

func (ac *OAuthController) CompleteGoogleOAuth(res http.ResponseWriter, req *http.Request) {
	log := pkg.RequestLogger(req.Context(), ac.logger)

	state := req.URL.Query().Get("state")
	promptParam := req.URL.Query().Get("prompt")

	stateCookie, err := req.Cookie("oauth_state")

	if err != nil {
		log.WarnContext(req.Context(), "oauth_state cookie missing on callback",
			"likely_silent_auth", promptParam == "none",
		)
		http.Redirect(res, req, "/auth/google-sso", http.StatusFound)
		return
	}

	if state != stateCookie.Value {
		log.ErrorContext(req.Context(), "oauth_state MISMATCH")
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
		log.ErrorContext(req.Context(), "code not found in callback")
		http.Error(res, "internal server error", http.StatusBadRequest)
		return
	}

	token, err := ac.authUseCase.CompleteGoogleOAuth(req.Context(), code)
	if err != nil {
		log.ErrorContext(req.Context(), "failed to finalize google auth", "error", err)
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
