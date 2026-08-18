package domain

import (
	"context"
)

type Profile struct {
	Name  string
	Email string
}

type ProfileUseCase interface {
	GetProfileById(c context.Context, Id string) (*Profile, error)
}
