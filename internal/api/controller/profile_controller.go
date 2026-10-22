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
	log := pkg.RequestLogger(req.Context(), pc.logger)

	userID, ok := pkg.GetUserID(req.Context())
	if !ok {
		log.WarnContext(req.Context(), "unable to fetch user id from context")
		http.Error(res, "auth failed", http.StatusBadRequest)
		return
	}

	profile, err := pc.profileUseCase.GetProfileByID(req.Context(), userID)
	if err != nil {
		log.WarnContext(req.Context(), "unable to fetch id")
		http.Error(res, "auth failed", http.StatusBadRequest)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(profile)
}
