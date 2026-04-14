package models

import (
	"github.com/google/uuid"
)

type Chat struct {
	Id        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	FromId    uuid.UUID `gorm:"type:uuid;" json:"from_id"`
	ToId      uuid.UUID `gorm:"type:uuid;" json:"to_id"`
	Message   string    `gorm:"type:text;" json:"message"`
	CreatedAt int64     `gorm:"autoCreateTime;" json:"created_at"`
}
