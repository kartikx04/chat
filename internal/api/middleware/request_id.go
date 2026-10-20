package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/kartikx04/chat/pkg"

	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func RequestID(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			requestID := uuid.NewString()

			ctx := pkg.AttachRequestID(req.Context(), requestID)
			req = req.WithContext(ctx)

			res.Header().Set("X-Request-ID", requestID)

			recorder := chimiddleware.NewWrapResponseWriter(res, req.ProtoMajor)
			startedAt := time.Now()

			next.ServeHTTP(res, req)

			pkg.RequestLogger(ctx, logger).InfoContext(
				ctx,
				"http request completed",
				"method", req.Method,
				"path", req.URL.Path,
				"status", recorder.Status(),
				"duration_ms", time.Since(startedAt).Milliseconds(),
			)
		})
	}
}
