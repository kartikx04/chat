package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/kartikx04/chat/internal/domain"
	"github.com/kartikx04/chat/internal/models"
	"github.com/kartikx04/chat/pkg"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

type authUsecase struct {
	oauthConfig *oauth2.Config
	userRepo    domain.UserRepository
	sessionRepo domain.SessionRepository
}

type GoogleUser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	VerifiedEmail bool   `json:"verified_email"`
	Picture       string `json:"picture"`
}

func NewAuthUsecase(oauthConfig *oauth2.Config, userRepo domain.UserRepository, sessionRepo domain.SessionRepository) domain.AuthUsecase {
	return &authUsecase{
		oauthConfig: oauthConfig,
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}

func (au *authUsecase) InitiateGoogleOAuth(ctx context.Context) (string, string, error) {
	state, err := pkg.GenerateStateToken()
	if err != nil {
		slog.ErrorContext(ctx, "failed to generate oauth state token", "error", err)
		return "", "", err
	}

	authURL := au.oauthConfig.AuthCodeURL(state)
	slog.InfoContext(ctx, "oauth redirect initiated")
	return state, authURL, nil
}

func (au *authUsecase) CompleteGoogleOAuth(ctx context.Context, code string) (string, error) {
	token, err := au.oauthConfig.Exchange(ctx, code)
	if err != nil {
		return "", err
	}

	client := au.oauthConfig.Client(ctx, token)

	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("google user info returned: %s", resp.Status)
	}
	var googleUser GoogleUser

	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		return "", err
	}
	if googleUser.ID == "" || googleUser.Email == "" || !googleUser.VerifiedEmail {
		return "", fmt.Errorf("google returned an incomplete or unverified identity")
	}
	if googleUser.ID == "" || googleUser.Email == "" || !googleUser.VerifiedEmail {
		return "", fmt.Errorf("google returned an incomplete or unverified identity")
	}

	user, err := au.userRepo.GetByOAuthID(ctx, googleUser.ID)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return "", err
		}

		_, emailErr := au.userRepo.GetByEmail(ctx, googleUser.Email)
		if emailErr == nil {
			return "", domain.ErrAccountLinkRequired
		}
		if !errors.Is(emailErr, gorm.ErrRecordNotFound) {
			return "", emailErr
		}

		now := time.Now()
		user = models.User{
			ID:        uuid.New(),
			OAuthID:   googleUser.ID,
			Email:     googleUser.Email,
			Username:  pkg.GenerateUniqueName(),
			Picture:   googleUser.Picture,
			Role:      "user",
			CreatedAt: now,
			UpdatedAt: now,
		}

		if err := au.userRepo.Create(ctx, &user); err != nil {
			return "", err
		}
	}

	sessionToken := pkg.GenerateSecureID()
	session := models.Session{
		ID:         uuid.New(),
		Token:      sessionToken,
		UserID:     user.ID,
		ExpiresAt:  time.Now().Add(72 * time.Hour),
		CreatedAt:  time.Now(),
		LastUsedAt: time.Now(),
	}
	if err := au.sessionRepo.Create(ctx, &session); err != nil {
		return "", err
	}

	return sessionToken, nil
}

func (au *authUsecase) Logout(ctx context.Context, cookie string) error {
	return au.sessionRepo.DeleteByToken(ctx, cookie)
}
