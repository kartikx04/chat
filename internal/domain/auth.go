package domain

import (
	"context"
)

type AuthUseCase interface {
	InitiateGoogleOAuth(ctx context.Context) (string, string, error)
	FinaliseGoogleAuth(ctx context.Context, code string) (string, error)

	Logout(ctx context.Context, cookie string) error
}
