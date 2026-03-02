package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	AuthOId   string    `gorm:"unique" json:"auth_o_id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Picture   string    `json:"picture"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

type OAuthData struct {
	Id            string `gorm:"primaryKey" json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Picture       string `json:"picture"`
}
