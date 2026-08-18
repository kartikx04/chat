package useCase

import (
	"context"
	"time"

	"github.com/kartikx04/chat/internal/domain"
)

type profileUseCase struct {
	ctx            context.Context
	userRepo       domain.UserRepository
	ContextTimeout time.Duration
}

func NewProfileUseCase(timeout time.Duration, userRepo domain.UserRepository) domain.ProfileUseCase {
	return &profileUseCase{
		ContextTimeout: timeout,
		userRepo:       userRepo,
	}
}

func (pu *profileUseCase) GetProfileById(c context.Context, id string) (*domain.Profile, error) {
	user, err := pu.userRepo.GetByID(c, id)
	if err != nil {
		return nil, err
	}

	return &domain.Profile{Name: user.Username, Email: user.Email}, nil
}
