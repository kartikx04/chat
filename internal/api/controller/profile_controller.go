package controller

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/kartikx04/chat/internal/auth"
	"github.com/kartikx04/chat/internal/domain"
)

type ProfileController struct {
	ProfileUseCase domain.ProfileUseCase
}

func (pc *ProfileController) Profile(res http.ResponseWriter, req *http.Request) {
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

	profile, err := pc.ProfileUseCase.GetProfileByID(req.Context(), claims.UserID)
	if err != nil {
		slog.WarnContext(req.Context(), "unable to fetch id")
		http.Error(res, "internal server error", http.StatusInternalServerError)
	}

	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(profile)
}
