package controller

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/kartikx04/chat/internal/domain"
	"github.com/kartikx04/chat/pkg"
)

type ProfileController struct {
	profileUseCase domain.ProfileUseCase
	logger         *slog.Logger
}

func NewProfileController(profileUseCase domain.ProfileUseCase, logger *slog.Logger) *ProfileController {
	return &ProfileController{
		profileUseCase: profileUseCase,
		logger:         logger,
	}
}

func (pc *ProfileController) Profile(res http.ResponseWriter, req *http.Request) {
	userID, ok := pkg.GetUserID(req.Context())
	if !ok {
		pc.logger.WarnContext(req.Context(), "user ID is missing from request context")
		http.Error(res, "unauthorized", http.StatusUnauthorized)
		return
	}

	profile, err := pc.profileUseCase.GetProfileByID(req.Context(), userID)
	if err != nil {
		pc.logger.ErrorContext(req.Context(), "unable to load profile", "error", err)
		http.Error(res, "unable to load profile", http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(res).Encode(profile); err != nil {
		pc.logger.ErrorContext(req.Context(), "unable to encode profile response", "error", err)
	}
}
