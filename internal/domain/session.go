package domain

import (
	"context"

	"github.com/kartikx04/chat/internal/models"
)

type SessionRepository interface {
	Create(c context.Context, sess *models.Sessions) error
}
