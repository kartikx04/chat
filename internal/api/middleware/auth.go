package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/kartikx04/chat/internal/domain"
	"github.com/kartikx04/chat/pkg"
)

func Auth(sessionRepo domain.SessionRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			Cookie, err := req.Cookie("session_id")
			if err != nil {
				slog.WarnContext(req.Context(), "error fetching cookie with value session_id", "error", err)
				http.Redirect(res, req, "Unauthorized", http.StatusUnauthorized)
				return
			}

			sessionID := Cookie.Value
			session, err := sessionRepo.GetBySessionID(req.Context(), sessionID)
			if err != nil {
				slog.WarnContext(req.Context(), "error fetching session id from database", "error", err)
				http.Redirect(res, req, "Unauthorized", http.StatusUnauthorized)
				return
			}

			if time.Now().After(session.ExpiresAt) {
				slog.WarnContext(req.Context(), "authentication failed: session expired")
				http.Error(res, "unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := pkg.AttachUserID(req.Context(), session.UserID.String())
			next.ServeHTTP(res, req.WithContext(ctx))
		})
	}
}
