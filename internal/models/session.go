package models

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	SessionID  string    `gorm:"unique" json:"session_id"`
	UserID     uuid.UUID `gorm:"unique" json:"user_id"`
	ExpiresAt  time.Time `json:"expires_at"`
	CreatedAt  time.Time `json:"created_at"`
	LastUsedAt time.Time `json:"last_used_at"`
}
