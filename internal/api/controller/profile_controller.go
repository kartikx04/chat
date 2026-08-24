package controller

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/kartikx04/chat/internal/domain"
)

type ProfileController struct {
	ProfileUseCase domain.ProfileUseCase
	SessionRepo    domain.SessionRepository
}

func (pc *ProfileController) Profile(res http.ResponseWriter, req *http.Request) {
	Cookie, err := req.Cookie("session_id")
	if err != nil {
		slog.Warn("error fetching cookie with value session_id", "error", err)
		http.Redirect(res, req, "auth failed", http.StatusForbidden)
		return
	}

	sessionID := Cookie.Value
	session, err := pc.SessionRepo.GetBySessionID(req.Context(), sessionID)
	if err != nil {
		slog.Warn("error fetching session id from database", "error", err)
		http.Redirect(res, req, "auth failed", http.StatusBadRequest)
		return
	}

	if session.SessionID != sessionID {
		slog.Warn("session not found in database", "error", err)
		http.Redirect(res, req, "auth failed", http.StatusBadRequest)
		return
	}
	slog.Info("user session found and user authenticated")

	profile, err := pc.ProfileUseCase.GetProfileByID(req.Context(), session.UserID.String())
	if err != nil {
		slog.WarnContext(req.Context(), "unable to fetch id")
		http.Error(res, "internal server error", http.StatusInternalServerError)
	}

	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(profile)
}
