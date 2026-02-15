package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Id        string `gorm:"primaryKey" json:"id"`
	Auth0Id   string `gorm:"unique" json:"auth0_id"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	Picture   string `json:"picture"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
}

type OAuthData struct {
	Id            string `gorm:"primaryKey" json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Picture       string `json:"picture"`
}
