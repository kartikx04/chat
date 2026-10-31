package controller

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

type fakeAuthUsecase struct {
	state         string
	authURL       string
	initErr       error
	initCallCount int

	token             string
	gotCode           string
	completeErr       error
	completeCallCount int

	gotLogoutToken  string
	logoutErr       error
	logoutCallCount int
}

func (f *fakeAuthUsecase) InitiateGoogleOAuth(_ context.Context) (string, string, error) {
	f.initCallCount++
	return f.state, f.authURL, f.initErr
}

func (f *fakeAuthUsecase) CompleteGoogleOAuth(_ context.Context, code string) (string, error) {
	f.completeCallCount++
	f.gotCode = code

	return f.token, f.completeErr
}

func (f *fakeAuthUsecase) Logout(_ context.Context, token string) error {
	f.logoutCallCount++
	f.gotLogoutToken = token
	return f.logoutErr
}

func TestInititiateGoogleOAuth(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
	}{
		{
			name:       "redirects to provider",
			err:        nil,
			wantStatus: http.StatusTemporaryRedirect,
		},
		{
			name:       "returns error when inititialisation fails",
			err:        errors.New("failed"),
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auth := &fakeAuthUsecase{
				state:   "state-123",
				authURL: "https://accounts.example.test/auth",
				initErr: tt.err,
			}
			controller := NewOAuthController(auth, "development", slog.New(slog.Handler(slog.DiscardHandler)))
			req := httptest.NewRequest(http.MethodGet, "/auth/google-sso", nil)
			recorder := httptest.NewRecorder()

			controller.StartGoogleOAuth(recorder, req)

			if got, want := recorder.Code, tt.wantStatus; got != want {
				t.Fatalf("status: got %v, want %v", got, want)
			}

			if got, want := auth.initCallCount, 1; got != want {
				t.Fatalf("init call count: got %v, want %v", got, want)
			}

			if tt.err != nil {
				return
			}

			if got, want := recorder.Header().Get("Location"), auth.authURL; got != want {
				t.Errorf("location: got %q, want %q", got, want)
			}

			cookies := recorder.Result().Cookies()
			if len(cookies) == 0 || cookies[0].Name != "oauth_state" || cookies[0].Value != auth.state {
				t.Errorf("oauth_state cookie not set as expected: %#v", cookies)
			}
		})
	}
}
