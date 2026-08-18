package useCase

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"time"

	"github.com/kartikx04/chat/internal/domain"
	"github.com/kartikx04/chat/internal/models"
	"github.com/kartikx04/chat/pkg"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

type authUseCase struct {
	oauthConfig    *oauth2.Config
	userRepo       domain.UserRepository
	contextTimeout time.Duration
}

type GoogleUser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	VerifiedEmail bool   `json:"verified_email"`
	Picture       string `json:"picture"`
}

func NewAuthUseCase(oauthConfig *oauth2.Config, userRepo domain.UserRepository, timeout time.Duration) domain.AuthUseCase {
	return &authUseCase{
		oauthConfig:    oauthConfig,
		userRepo:       userRepo,
		contextTimeout: timeout,
	}
}

func (au *authUseCase) InitiateGoogleOAuth(ctx context.Context) (string, string, error) {
	state, err := pkg.GenerateStateToken()
	if err != nil {
		slog.ErrorContext(ctx, "failed to generate oauth state token", "error", err)
		return "", "", err
	}

	authURL := pkg.OAuthgolang.AuthCodeURL(state)
	slog.InfoContext(ctx, "oauth redirect initiated")
	return state, authURL, nil
}

func (au *authUseCase) FinaliseGoogleAuth(ctx context.Context, code string) (string, error) {
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

	var googleUser GoogleUser

	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		return "", err
	}

	user, err := au.userRepo.GetByEmail(ctx, googleUser.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {

			user = models.User{
				OAuthID:  googleUser.ID,
				Email:    googleUser.Email,
				Username: pkg.GenerateUniqueName(),
				Picture:  googleUser.Picture,
				Role:     "user",
			}

			if err := au.userRepo.Create(ctx, &user); err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	}

	sessionID := "asdf"
	// sessionID := uuid.NewString()

	// session := models.Session{
	// 	ID:           uuid.NewString(),
	// 	UserID:       user.ID,
	// 	RefreshToken: token.RefreshToken,
	// 	ExpiresAt:    token.Expiry,
	// }

	// if err := au.sessionRepo.Create(ctx, &session); err != nil {
	// 	return "", err
	// }

	return sessionID, nil
}
