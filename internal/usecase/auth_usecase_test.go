package usecase

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/kartikx04/chat/internal/models"
	"golang.org/x/oauth2"
)

type fakeAuthUserRepository struct {
	createCallCount int
	createdUser     models.User
	createErr       error

	user         models.User
	getCallCount int
	getErr       error
}

func (f *fakeAuthUserRepository) Create(_ context.Context, user *models.User) error {
	f.createCallCount++
	f.createdUser = *user
	return f.createErr
}

func (f *fakeAuthUserRepository) List(_ context.Context) ([]models.User, error) {
	return nil, nil
}

func (f *fakeAuthUserRepository) GetByEmail(_ context.Context, _ string) (models.User, error) {
	f.getCallCount++
	return f.user, f.getErr
}

func (f *fakeAuthUserRepository) GetByID(_ context.Context, _ string) (models.User, error) {
	return models.User{}, nil
}

type fakeAuthSessionRepository struct {
	createCallCount int
	createdSession  models.Session
	createErr       error

	deleteCallCount int
	deleteToken     string
	deleteErr       error
}

func (f *fakeAuthSessionRepository) Create(_ context.Context, sess *models.Session) error {
	f.createCallCount++
	f.createdSession = *sess
	return f.createErr
}

func (f *fakeAuthSessionRepository) GetByID(_ context.Context, id string) (models.Session, error) {
	return models.Session{}, nil
}

func (f *fakeAuthSessionRepository) FindByToken(_ context.Context, id string) (models.Session, error) {
	return models.Session{}, nil
}

func (f *fakeAuthSessionRepository) DeleteByToken(_ context.Context, token string) error {
	f.deleteCallCount++
	f.deleteToken = token
	return f.deleteErr
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func jsonResponse(req *http.Request, status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Status:     http.StatusText(status),
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(body)),
		Request:    req,
	}
}

func authTestTransport(tokenStatus int) http.RoundTripper {
	return roundTripFunc(func(req *http.Request) (*http.Response, error) {
		switch req.URL.Host + req.URL.Path {
		case "oauth.test/token":
			if tokenStatus != http.StatusOK {
				return jsonResponse(req, http.StatusBadRequest, `{"error": "invalid grant"}`), nil
			}
			return jsonResponse(req, http.StatusOK, `{"access_token": "access-123", "token_type": "Bearer", "expires_in": 3600}`), nil
		case "www.googleapis.com/oauth2/v2/userinfo":
			return jsonResponse(req, http.StatusOK, `{"id": "google-123", "email": "person@example.com", "name": "person", "picture": "https://example.test/p.png"}`), nil
		default:
			return nil, fmt.Errorf("unexpected oauth http request: %s", req.URL)
		}
	})
}

func authTestContext(tokenStatus int) context.Context {
	return context.WithValue(
		context.Background(),
		oauth2.HTTPClient,
		&http.Client{Transport: authTestTransport(tokenStatus)},
	)
}

func authTestOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		RedirectURL:  "https://app.example.test/callback",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://oauth.test/authorize",
			TokenURL: "https://oauth.test/token",
		},
		Scopes: []string{"profile"},
	}
}
