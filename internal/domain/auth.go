package domain

import (
	"context"
	"errors"
)

var ErrAccountLinkRequired = errors.New("account linking required")

type AuthUsecase interface {
	InitiateGoogleOAuth(ctx context.Context) (string, string, error)
	CompleteGoogleOAuth(ctx context.Context, code string) (string, error)

	Logout(ctx context.Context, cookie string) error
}
