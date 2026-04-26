package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

// Use a typed key to avoid context collisions with other packages
type contextKey string

const RequestIDKey contextKey = "request_id"

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Respect upstream request ID if present (e.g. from load balancer)
		id := r.Header.Get("X-Request-ID")
		if id == "" {
			id = uuid.New().String()
		}

		// Stamp it on the response so client can correlate
		w.Header().Set("X-Request-ID", id)

		// Put it in context so all downstream code can grab it
		ctx := context.WithValue(r.Context(), RequestIDKey, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Helper — call this anywhere you have a context
func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(RequestIDKey).(string); ok {
		return id
	}
	return ""
}
