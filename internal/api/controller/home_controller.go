package controller

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/kartikx04/chat/internal/auth"
)

func Home(res http.ResponseWriter, req *http.Request) {
	authHeader := req.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		slog.WarnContext(req.Context(), "no auth header")
		http.Error(res, "unauthorized", http.StatusUnauthorized)
		return
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := auth.ValidateToken(tokenStr)
	if err != nil {
		slog.WarnContext(req.Context(), "me: invalid token", "error", err)
		http.Error(res, "unauthorized", http.StatusUnauthorized)
		return
	}

	slog.DebugContext(req.Context(), "me: identity resolved",
		"user_id", claims.UserID,
		"username", claims.Username,
	)

	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(map[string]string{
		"id":       claims.UserID,
		"username": claims.Username,
		"email":    claims.Email,
	})
}
