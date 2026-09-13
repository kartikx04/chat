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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProfile(t *testing.T) {
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

	err := userRepo.Create(app.Context, &user)
	require.NoError(t, err)

	sessionID := "test-session-123"

	session := models.Session{
		ID:         uuid.New(),
		SessionID:  sessionID,
		UserID:     userID,
		ExpiresAt:  time.Now().Add(72 * time.Hour),
		CreatedAt:  time.Now(),
		LastUsedAt: time.Now(),
	}

	err = app.TestDB.Create(&session).Error
	require.NoError(t, err)

	// Act
	req := httptest.NewRequest(
		http.MethodGet,
		"/public/profile",
		nil,
	)

	req.AddCookie(&http.Cookie{
		Name:  "session_id",
		Value: sessionID,
	})

	rec := httptest.NewRecorder()

	app.Router.ServeHTTP(rec, req)

	// Assert
	assert.Equal(t, http.StatusOK, rec.Code)

	assert.Equal(
		t,
		"application/json",
		rec.Header().Get("Content-Type"),
	)

	var profile domain.Profile

	err = json.NewDecoder(rec.Body).Decode(&profile)
	require.NoError(t, err)

	assert.Equal(t, "Kartik", profile.Name)
	assert.Equal(t, "kartik@example.com", profile.Email)
}
