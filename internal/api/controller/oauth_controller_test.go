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
				t.Fatalf("init call count: got %d, want %d", got, want)
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

func TestCompleteGoogleOAuth(t *testing.T) {
	tests := []struct {
		name        string
		query       string
		cookieState string
		addCookie   bool
		completeErr error
		token       string

		wantLocation          string
		wantStatus            int
		wantCompleteCallCount int
	}{
		{
			name:         "redirects when state cookie is missing",
			query:        "?state=state-123&code=code-123",
			wantStatus:   http.StatusFound,
			wantLocation: "/auth/google-sso",
		},
		{
			name:                  "rejects state mismatch",
			query:                 "?state=wrong&code=code-123",
			cookieState:           "state-123",
			addCookie:             true,
			wantStatus:            http.StatusBadRequest,
			wantCompleteCallCount: 0,
		},
		{
			name:                  "rejects missing code",
			query:                 "?state=state-123",
			cookieState:           "state-123",
			addCookie:             true,
			wantStatus:            http.StatusBadRequest,
			wantCompleteCallCount: 0,
		},
		{
			name:                  "returns error when OAuth completion fails",
			query:                 "?state=state-123&code=code-123",
			cookieState:           "state-123",
			addCookie:             true,
			completeErr:           errors.New("exchange failed"),
			wantStatus:            http.StatusInternalServerError,
			wantCompleteCallCount: 1,
		},
		{
			name:                  "sets session cookie and redirects",
			query:                 "?state=state-123&code=code-123",
			cookieState:           "state-123",
			addCookie:             true,
			token:                 "session-123",
			wantStatus:            http.StatusSeeOther,
			wantLocation:          "/auth/success",
			wantCompleteCallCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auth := &fakeAuthUsecase{
				token:       tt.token,
				completeErr: tt.completeErr,
			}
			controller := NewOAuthController(auth, "production", slog.New(slog.Handler(slog.DiscardHandler)))
			req := httptest.NewRequest(http.MethodGet, "/auth/google/callback"+tt.query, nil)
			if tt.addCookie {
				req.AddCookie(&http.Cookie{Name: "oauth_state", Value: tt.cookieState})
			}
			recorder := httptest.NewRecorder()
			controller.CompleteGoogleOAuth(recorder, req)

			if got, want := recorder.Code, tt.wantStatus; got != want {
				t.Errorf("status: got %v, want %v", got, want)
			}
			if got, want := auth.completeCallCount, tt.wantCompleteCallCount; got != want {
				t.Errorf("complete call count: got %d, want %d", got, want)
			}
			if got, want := recorder.Header().Get("Location"), tt.wantLocation; tt.wantLocation != "" && got != want {
				t.Errorf("location: got %q, want %q", got, want)
			}
			if got, want := auth.gotCode, "code-123"; tt.wantCompleteCallCount == 1 && got != want {
				t.Errorf("code passed to usecase: got %q, want %q", got, want)
			}

			if tt.token != "" {
				var sessionCookie *http.Cookie
				for _, cookie := range recorder.Result().Cookies() {
					if cookie.Name == "session_id" {
						sessionCookie = cookie
					}
				}

				if sessionCookie == nil {
					t.Fatal("session_id cookie was not set")
				}
				if sessionCookie.Value != tt.token || !sessionCookie.Secure || !sessionCookie.HttpOnly {
					t.Errorf("unexpected session cookie: %#v", sessionCookie)
				}
			}
		})
	}
}
