package models

import (
	"time"

	"github.com/google/uuid"
)

type Chat struct {
	Id        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	FromID    uuid.UUID `gorm:"type:uuid;" json:"from_id"`
	ToID      uuid.UUID `gorm:"type:uuid;" json:"to_id"`
	Message   string    `gorm:"type:text;" json:"message"`
	CreatedAt time.Time `gorm:"autoCreateTime;" json:"created_at"`
}
