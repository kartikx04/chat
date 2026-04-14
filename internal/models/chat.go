package models

import (
	"time"

	"github.com/google/uuid"
)

type Chat struct {
	Id            uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	FromId        uuid.UUID `gorm:"type:uuid;" json:"from_id"`
	ToId          uuid.UUID `gorm:"type:uuid;" json:"to_id"`
	Message       string    `gorm:"type:text;" json:"message"`
	CreatedAt     time.Time `gorm:"autoCreateTime;" json:"created_at"`
	CreatedAtUnix int64     `json:"created_at_unix"`
}

type Message struct {
	Type   string `json:"type"`
	User   string `json:"user,omitempty"`
	UserId string `json:"user_id,omitempty"`
	Chat   Chat   `json:"chat,omitempty"`
}
