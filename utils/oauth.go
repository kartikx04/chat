package utils

import (
	"crypto/rand"

	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
)

var (
	OAuthgolang *oauth2.Config
	Store       = sessions.NewCookieStore([]byte(LoadFile("TOKEN_SECRET")))
)

const (
	tokenSet    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	tokenLength = 15
)

// GenerateRandomString generates a random string of the specified length(15).
func TokenString() (string, error) {
	charsetLength := len(tokenSet)

	randomBytes := make([]byte, tokenLength)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	for i := 0; i < tokenLength; i++ {
		randomBytes[i] = tokenSet[int(randomBytes[i])%charsetLength]
	}

	return string(randomBytes), nil
}
