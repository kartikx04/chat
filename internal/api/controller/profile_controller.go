package controller

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/kartikx04/chat/internal/domain"
	"github.com/kartikx04/chat/pkg"
)

type ProfileController struct {
	profileUsecase domain.ProfileUsecase
	logger         *slog.Logger
}

func NewProfileController(profileUsecase domain.ProfileUsecase, logger *slog.Logger) *ProfileController {
	return &ProfileController{
		profileUsecase: profileUsecase,
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

	profile, err := pc.profileUsecase.GetProfileByID(req.Context(), userID)
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
