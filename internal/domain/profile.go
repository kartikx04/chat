package domain

import (
	"context"

	"github.com/kartikx04/chat/internal/models"
)

type ProfileUseCase interface {
	GetProfileById(c context.Context, Id string) (*models.Users, error)
}
