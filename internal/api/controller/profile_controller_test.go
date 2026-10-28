package controller

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kartikx04/chat/internal/domain"
	"github.com/kartikx04/chat/pkg"
)

type fakeProfileByIDFinder struct {
	profile *domain.Profile
	gotID   string

	callCount int
	err       error
}

func (f *fakeProfileByIDFinder) GetProfileByID(_ context.Context, id string) (*domain.Profile, error) {
	f.gotID = id
	f.callCount++

	return f.profile, f.err
}

func TestProfileControllerReturnsProfile(t *testing.T) {
	wantUserID := "user-id-123"
	wantProfile := domain.Profile{
		Name:  "kartik",
		Email: "kartik@example.com",
	}
	finder := &fakeProfileByIDFinder{
		profile: &wantProfile,
		gotID:   wantUserID,
	}
	logger := slog.New(slog.Handler(slog.DiscardHandler))

	profileService := NewProfileController(finder, logger)

	req := httptest.NewRequest(http.MethodGet, "/v1/users/me", nil)
	req = req.WithContext(pkg.AttachUserID(req.Context(), wantUserID))
	recorder := httptest.NewRecorder()

	profileService.Profile(recorder, req)

	if got, want := recorder.Code, http.StatusOK; got != want {
		t.Errorf("code: got %v, want %v", got, want)
	}

	if got, want := recorder.Header().Get("Content-Type"), "application/json"; got != want {
		t.Errorf("content type: got %v, want %v", got, want)
	}

	if got, want := finder.callCount, 1; got != want {
		t.Errorf("usecase calls: got %v, want %v", got, want)
	}

	if got, want := finder.gotID, wantUserID; got != want {
		t.Errorf("user ID: got %v, want %v", got, want)
	}

	var gotProfile domain.Profile
	if err := json.NewDecoder(recorder.Body).Decode(&gotProfile); err != nil {
		t.Fatalf("decode response json: %v", err)
	}

	if got, want := gotProfile, wantProfile; got != want {
		t.Errorf("profile: got %v, want %v", got, want)
	}
}
