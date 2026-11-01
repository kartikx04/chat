package controller

import (
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogoutController(t *testing.T) {
	tests := []struct {
		name          string
		env           string
		cookie        bool
		err           error
		wantStatus    int
		wantCallCount int
		wantSecure    bool
	}{
		{
			name:          "missing cookie redirects without logout",
			env:           "development",
			cookie:        false,
			err:           nil,
			wantStatus:    http.StatusSeeOther,
			wantCallCount: 0,
			wantSecure:    false,
		},
		{
			name:          "logout failure returns 500",
			env:           "development",
			cookie:        true,
			err:           errors.New("delete failed"),
			wantStatus:    http.StatusInternalServerError,
			wantCallCount: 1,
			wantSecure:    false,
		},
		{
			name:          "successful production logout clears secure cookie",
			env:           "production",
			cookie:        true,
			err:           nil,
			wantStatus:    http.StatusSeeOther,
			wantCallCount: 1,
			wantSecure:    true,
		},
		{
			name:          "successful development logout clears unsecure cookie",
			env:           "development",
			cookie:        true,
			err:           nil,
			wantStatus:    http.StatusSeeOther,
			wantCallCount: 1,
			wantSecure:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auth := &fakeAuthUsecase{
				logoutErr: tt.err,
			}
			controller := NewLogoutController(auth, tt.env, slog.New(slog.Handler(slog.DiscardHandler)))
			req := httptest.NewRequest(http.MethodGet, "/logout", nil)
			if tt.cookie {
				req.AddCookie(&http.Cookie{Name: "session_id", Value: "session-123"})
			}
			recorder := httptest.NewRecorder()
			controller.Logout(recorder, req)

			if got, want := recorder.Code, tt.wantStatus; got != want {
				t.Fatalf("status: got %v, want %v", got, want)
			}
			if got, want := auth.logoutCallCount, tt.wantCallCount; got != want {
				t.Errorf("logout call count: got %d, want %d", got, want)
			}
			if got, want := auth.gotLogoutToken, "session-123"; tt.wantCallCount == 1 && got != want {
				t.Errorf("token: got %q, want %q", got, want)
			}
			if tt.wantStatus == http.StatusSeeOther {
				if got, want := recorder.Header().Get("Location"), "/"; got != want {
					t.Errorf("lcoation: got %q, want %q", got, want)
				}
			}

			if tt.wantCallCount == 1 && tt.err == nil {
				var cleared *http.Cookie
				for _, cookie := range recorder.Result().Cookies() {
					if cookie.Name == "session_id" {
						cleared = cookie
					}
				}
				if cleared == nil || cleared.MaxAge != -1 || cleared.Secure != tt.wantSecure {
					t.Errorf("unexpected clear cookie: %#v", cleared)
				}
			}
		})
	}
}
