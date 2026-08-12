package useCase

import (
	"context"

	"github.com/kartikx04/chat/internal/domain"
	"github.com/kartikx04/chat/internal/models"
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

func (pu *profileUseCase) GetProfileById(c context.Context, Id string) (*models.Users, error) {
	// TODO: set ctx timeout duration

	return nil, nil
}
