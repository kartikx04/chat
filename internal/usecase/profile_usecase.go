package useCase

import (
	"context"

	"github.com/kartikx04/chat/internal/domain"
	"github.com/kartikx04/chat/internal/repository"
)

type profileUseCase struct {
	ctx      context.Context
	userRepo *repository.UserRepository
}

func NewProfileUseCase(ctx context.Context, userRepo *repository.UserRepository) domain.ProfileUseCase {
	return &profileUseCase{
		ctx:      ctx,
		userRepo: userRepo,
	}
}

func (pu *profileUseCase) GetProfileById(c context.Context, id string) (*domain.Profile, error) {
	user, err := pu.userRepo.GetUserById(id)
	if err != nil {
		return nil, err
	}

	return &domain.Profile{Name: user.Username, Email: user.Email}, nil
}
