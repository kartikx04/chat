package domain

import (
	"context"
)

type Profile struct {
	Name  string
	Email string
}

type ProfileUsecase interface {
	GetProfileByID(c context.Context, ID string) (*Profile, error)
}
