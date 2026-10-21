package domain

import (
	"context"
)

type AuthUsecase interface {
	InitiateGoogleOAuth(ctx context.Context) (string, string, error)
	CompleteGoogleOAuth(ctx context.Context, code string) (string, error)

	Logout(ctx context.Context, cookie string) error
}
