package usecase

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/google/uuid"

	"github.com/kartikx04/chat/internal/domain"
	"github.com/kartikx04/chat/internal/models"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

type fakeAuthUserRepository struct {
	createCallCount int
	createdUser     models.User
	createErr       error

	user         models.User
	gotCallCount int
	gotErr       error

	oauthUser         models.User
	oauthGetCallCount int
	oauthGetErr       error
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
	f.gotCallCount++
	return f.user, f.gotErr
}

func (f *fakeAuthUserRepository) GetByID(_ context.Context, _ string) (models.User, error) {
	return models.User{}, nil
}

func (f *fakeAuthUserRepository) GetByOAuthID(_ context.Context, _ string) (models.User, error) {
	f.oauthGetCallCount++
	return f.oauthUser, f.oauthGetErr
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
			return jsonResponse(req, http.StatusOK, `{"id":"google-123", "email":"person@example.com", "name":"person", "verified_email":true, "picture":"https://example.test/p.png"}`), nil
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

func TestCompleteGoogleOAuthIdentityResolution(t *testing.T) {
	existingUser := models.User{
		ID:      uuid.New(),
		OAuthID: "google-123",
		Email:   "person@example.com",
	}

	tests := []struct {
		name             string
		oauthUser        models.User
		oauthGetErr      error
		emailUser        models.User
		emailGetErr      error
		sessionCreateErr error
		wantErr          error
		wantCreate       int
		wantEmailLookup  int
		wantSession      int
		wantSessionUser  uuid.UUID
	}{
		{
			name:            "uses existing provider identity",
			oauthUser:       existingUser,
			wantSession:     1,
			wantSessionUser: existingUser.ID,
		},
		{
			name:            "creates user for new provider and email",
			oauthGetErr:     gorm.ErrRecordNotFound,
			emailGetErr:     gorm.ErrRecordNotFound,
			wantCreate:      1,
			wantEmailLookup: 1,
			wantSession:     1,
		},
		{
			name:            "rejects email owned by another identity",
			oauthGetErr:     gorm.ErrRecordNotFound,
			emailUser:       models.User{ID: uuid.New(), OAuthID: "other-id", Email: "person@example.com"},
			wantErr:         domain.ErrAccountLinkRequired,
			wantEmailLookup: 1,
		},
		{
			name:             "returns session creation failure",
			oauthUser:        existingUser,
			sessionCreateErr: errors.New("session creation failed"),
			wantErr:          errors.New("session creation failed"),
			wantSession:      1,
			wantSessionUser:  existingUser.ID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := &fakeAuthUserRepository{
				oauthUser:   tt.oauthUser,
				oauthGetErr: tt.oauthGetErr,
				user:        tt.emailUser,
				gotErr:      tt.emailGetErr,
			}
			sessionRepo := &fakeAuthSessionRepository{createErr: tt.sessionCreateErr}
			auth := NewAuthUsecase(authTestOAuthConfig(), userRepo, sessionRepo)

			_, err := auth.CompleteGoogleOAuth(authTestContext(http.StatusOK), "code-123")
			if tt.wantErr != nil && err == nil {
				t.Fatal("expected an error")
			}
			if tt.wantErr == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantErr == domain.ErrAccountLinkRequired && !errors.Is(err, tt.wantErr) {
				t.Fatalf("error = %v, want %v", err, tt.wantErr)
			}
			if got, want := userRepo.createCallCount, tt.wantCreate; got != want {
				t.Errorf("user create calls = %d, want %d", got, want)
			}
			if got, want := userRepo.gotCallCount, tt.wantEmailLookup; got != want {
				t.Errorf("email lookup calls = %d, want %d", got, want)
			}
			if got, want := sessionRepo.createCallCount, tt.wantSession; got != want {
				t.Errorf("session create calls = %d, want %d", got, want)
			}
			if tt.wantSessionUser != uuid.Nil && sessionRepo.createdSession.UserID != tt.wantSessionUser {
				t.Errorf("session user ID = %v, want %v", sessionRepo.createdSession.UserID, tt.wantSessionUser)
			}
		})
	}
}
