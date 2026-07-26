package pkg

import (
	"context"
	"crypto/rand"
	"errors"
	"io"
	"net/http"

	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
)

var (
	OAuthgolang *oauth2.Config
	Store       *sessions.CookieStore
)

func getStore() *sessions.CookieStore {
	if Store == nil {
		Store = sessions.NewCookieStore([]byte(LoadFile("TOKEN_SECRET")))
		Store.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   3600,
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
		}
	}
	return Store
}

// GetSession is a public function to safely get a session
func GetSession(req *http.Request, sessionName string) (*sessions.Session, error) {
	return getStore().Get(req, sessionName)
}

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

// GetUserData validates verification request and returns data of verified google user
func GetUserData(state, code, tokenCode string) ([]byte, error) {
	// compares the generated token string to the token retrieved from the parsed URL
	if state != tokenCode {
		return nil, errors.New("invalid user")
	}

	// converts authorization code into a token
	token, err := OAuthgolang.Exchange(context.Background(), code)
	if err != nil {
		return nil, err
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, err
	}

	// this is done to prevent memory leakage
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	// returns data of verified google user
	return data, nil
}
