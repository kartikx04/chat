package controller

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/kartikx04/chat/internal/domain"
	"github.com/kartikx04/chat/pkg"
)

type ProfileController struct {
	ProfileUseCase domain.ProfileUseCase
}

func (pc *ProfileController) Profile(res http.ResponseWriter, req *http.Request) {
	userID, ok := pkg.GetUserID(req.Context())
	if !ok {
		slog.WarnContext(req.Context(), "unable to fetch user id from context")
		http.Error(res, "auth failed", http.StatusBadRequest)
		return
	}

	profile, err := pc.ProfileUseCase.GetProfileByID(req.Context(), userID)
	if err != nil {
		slog.WarnContext(req.Context(), "unable to fetch id")
		http.Error(res, "auth failed", http.StatusBadRequest)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(profile)
}
