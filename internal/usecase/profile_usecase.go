package usecase

import (
	"context"
	"fmt"

	"github.com/kartikx04/chat/internal/domain"
)

type profileUseCase struct {
	ctx      context.Context
	userRepo domain.UserRepository
}

func NewProfileUseCase(userRepo domain.UserRepository) domain.ProfileUseCase {
	return &profileUseCase{
		userRepo: userRepo,
	}
}

func (pu *profileUseCase) GetProfileByID(c context.Context, id string) (*domain.Profile, error) {
	user, err := pu.userRepo.GetByID(c, id)
	if err != nil {
		return nil, fmt.Errorf("get user by ID for profile: %w", err)
	}

	return &domain.Profile{Name: user.Username, Email: user.Email}, nil
}
