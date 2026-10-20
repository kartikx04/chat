package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/kartikx04/chat/internal/models"
	"github.com/kartikx04/chat/pkg"
)

type sessionFinder interface {
	FindByToken(ctx context.Context, token string) (models.Session, error)
}

func Auth(sessionRepo sessionFinder, logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			log := pkg.RequestLogger(req.Context(), logger)

			cookie, err := req.Cookie("session_id")
			if err != nil {
				log.WarnContext(req.Context(), "error fetching cookie with value session_id", "error", err)
				http.Redirect(res, req, "Unauthorized", http.StatusUnauthorized)
				return
			}

			token := cookie.Value
			session, err := sessionRepo.FindByToken(req.Context(), token)
			if err != nil {
				log.WarnContext(req.Context(), "error fetching token from session", "error", err)
				http.Redirect(res, req, "Unauthorized", http.StatusUnauthorized)
				return
			}

			if time.Now().After(session.ExpiresAt) {
				log.WarnContext(req.Context(), "authentication failed: session expired")
				http.Error(res, "unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := pkg.AttachUserID(req.Context(), session.UserID.String())

			authenticatedLog := pkg.RequestLogger(ctx, logger)
			authenticatedLog.InfoContext(ctx, "user session authenticated")

			next.ServeHTTP(res, req.WithContext(ctx))
		})
	}
}
