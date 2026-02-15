package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID        string `gorm:"primaryKey" json:"id"`
	Auth0ID   string `gorm:"unique" json:"auth0_id"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	Picture   string `json:"picture"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
}
