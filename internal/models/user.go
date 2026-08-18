package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	OAuthID   string    `gorm:"unique" json:"oauth_id"`
	Email     string    `gorm:"unique" json:"email"`
	Username  string    `gorm:"unique" json:"username"`
	Picture   string    `json:"picture"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
