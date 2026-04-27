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
	IsSelf        bool      `gorm:"-" json:"is_self"`
}

type Message struct {
	Type   string `json:"type"`
	User   string `json:"user,omitempty"`
	UserId string `json:"user_id,omitempty"`
	To     string `json:"to,omitempty"`
	Chat   Chat   `json:"chat,omitempty"`
}

type ContactList struct {
	Id           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	LastActivity int64     `json:"last_activity"`
}
