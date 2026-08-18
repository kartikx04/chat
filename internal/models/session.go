package models

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID           string
	UserID       uuid.UUID
	RefreshToken string
	ExpiresAt    time.Time
	CreatedAt    time.Time
}
