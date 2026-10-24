package usecase

import (
	"context"
	"fmt"

	"github.com/kartikx04/chat/internal/domain"
	"github.com/kartikx04/chat/internal/models"
)

type userByIDFinder interface {
	GetByID(c context.Context, id string) (models.User, error)
}

type profileUsecase struct {
	userRepo userByIDFinder
}

func NewProfileUsecase(userRepo userByIDFinder) domain.ProfileUsecase {
	return &profileUsecase{
		userRepo: userRepo,
	}
}

func (pu *profileUsecase) GetProfileByID(c context.Context, id string) (*domain.Profile, error) {
	user, err := pu.userRepo.GetByID(c, id)
	if err != nil {
		return nil, fmt.Errorf("get user by ID for profile: %w", err)
	}

	return &domain.Profile{Name: user.Username, Email: user.Email}, nil
}
