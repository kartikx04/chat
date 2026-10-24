package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/kartikx04/chat/internal/domain"
	"github.com/kartikx04/chat/internal/models"
)

type fakeUserByIDFinder struct {
	user  models.User
	gotID string

	callCount int
	err       error
}

func (f *fakeUserByIDFinder) GetByID(_ context.Context, id string) (models.User, error) {
	f.gotID = id
	f.callCount++

	return f.user, f.err
}

func TestProfileUsecaseGetProfileByID(t *testing.T) {
	userID := uuid.NewString()
	repositoryError := errors.New("user was not found")

	tests := []struct {
		name            string
		user            models.User
		wantProfile     domain.Profile
		repositoryError error
		wantError       error
	}{
		{
			name: "maps user to profile",
			user: models.User{
				Username: "kartik",
				Email:    "kartik@example.com",
			},
			wantProfile: domain.Profile{
				Name:  "kartik",
				Email: "kartik@example.com",
			},
		},
		{
			name:            "repository returns error",
			repositoryError: repositoryError,
			wantError:       repositoryError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			finder := &fakeUserByIDFinder{
				user: tt.user,
				err:  tt.repositoryError,
			}
			profileService := NewProfileUsecase(finder)

			gotProfile, gotError := profileService.GetProfileByID(context.Background(), userID)

			if got, want := finder.gotID, userID; got != want {
				t.Errorf("user ID passed to repository: got %v, want %v", got, want)
			}

			if got, want := finder.callCount, 1; got != want {
				t.Errorf("repository calls: got %v, want %v", got, want)
			}

			if tt.wantError != nil {
				if !errors.Is(gotError, tt.wantError) {
					t.Errorf("error: got %v, want %v", gotError, tt.wantError)
				}
				if gotProfile != nil {
					t.Errorf("profile: got %#v, want nil", gotProfile)
				}
				return
			}

			if gotError != nil {
				t.Fatalf("GetProfileByID returned unexpected error: %v", gotError)
			}

			if gotProfile == nil {
				t.Fatalf("profile: got nil, want a profile")
			}

			if got, want := *gotProfile, tt.wantProfile; got != want {
				t.Errorf("profiel: got %v, want %v", got, want)
			}
		})
	}

}
