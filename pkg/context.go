package pkg

import (
	"context"
	"log/slog"
)

type contextKey int

const (
	userIDKey contextKey = iota
	requestIDKey
)

func AttachRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey, requestID)
}

func GetRequestID(ctx context.Context) (string, bool) {
	requestID, ok := ctx.Value(requestIDKey).(string)
	return requestID, ok
}

func AttachUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func GetUserID(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userIDKey).(string)
	return userID, ok
}

func RequestLogger(ctx context.Context, logger *slog.Logger) *slog.Logger {
	if requestID, ok := GetRequestID(ctx); ok {
		logger = logger.With("request_id", requestID)
	}

	if userID, ok := GetUserID(ctx); ok {
		logger = logger.With("user_id", userID)
	}

	return logger
}
