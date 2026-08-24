package domain

import (
	"context"

	"github.com/kartikx04/chat/internal/models"
)

type SessionRepository interface {
	Create(c context.Context, sess *models.Session) error
	GetByID(ctx context.Context, id string) (models.Session, error)
	GetBySessionID(ctx context.Context, id string) (models.Session, error)
	DeleteByID(ctx context.Context, id string) error
}
