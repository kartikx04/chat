package models

type GoogleUser struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}
