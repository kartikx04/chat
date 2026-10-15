package integrationtest

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/kartikx04/chat/internal/domain"
	"github.com/kartikx04/chat/internal/models"
	"github.com/kartikx04/chat/internal/repository"
)

func TestProfile(t *testing.T) {
	var err error

	app := Setup(t)

	userID := uuid.New()
	OAuthID := uuid.New()

	user := models.User{
		ID:        userID,
		OAuthID:   OAuthID.String(),
		Email:     "kartik@example.com",
		Username:  "Kartik",
		Picture:   "asdf",
		Role:      "user",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	userRepo := repository.NewUserRepository(app.TestDB, app.Logger)

	err = userRepo.Create(app.Context, &user)
	if err != nil {
		t.Fatalf("create user: %v", err)
	}

	token := "test-session-token-123"

	session := models.Session{
		ID:         uuid.New(),
		Token:      token,
		UserID:     userID,
		ExpiresAt:  time.Now().Add(72 * time.Hour),
		CreatedAt:  time.Now(),
		LastUsedAt: time.Now(),
	}

	sessionRepo := repository.NewSessionRepository(app.TestDB, app.Logger)
	err = sessionRepo.Create(app.Context, &session)
	if err != nil {
		t.Fatalf("create session: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/v1/users/me", nil)

	req.AddCookie(&http.Cookie{
		Name:  "session_id",
		Value: token,
	})

	rec := httptest.NewRecorder()

	app.Router.ServeHTTP(rec, req)

	if got, want := rec.Code, http.StatusOK; got != want {
		t.Errorf("want %v, got: %v", want, got)
	}

	if got, want := rec.Header().Get("Content-Type"), "application/json"; got != want {
		t.Errorf("wanted content type: %v, got: %v", want, got)
	}

	var profile domain.Profile

	err = json.NewDecoder(rec.Body).Decode(&profile)
	if err != nil {
		t.Errorf("decode body: %v", err)
	}

	if got, want := profile.Name, "Kartik"; got != want {
		t.Errorf("user name = %v, want %v", got, want)
	}
	if got, want := profile.Email, "kartik@example.com"; got != want {
		t.Errorf("user email = %v want %v", got, want)
	}
}
