package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/kartikx04/chat/internal/models"
	"github.com/kartikx04/chat/pkg"
)

type fakeSessionFinder struct {
	session  models.Session
	gotToken string
	calls    int
	err      error
}

func (f *fakeSessionFinder) FindByToken(_ context.Context, token string) (models.Session, error) {
	f.gotToken = token
	f.calls += 1

	return f.session, f.err
}

func TestAuth(t *testing.T) {
	t.Run("adds user ID to context for valid session", func(t *testing.T) {
		userID := uuid.New()
		wantToken := "valid-session-token"
		wantUserID := userID.String()

		sessionFinder := &fakeSessionFinder{
			session: models.Session{
				UserID:    userID,
				ExpiresAt: time.Now().Add(time.Hour),
			},
		}

		var nextHandlerCalled bool
		var gotUserID string
		var contextGotUserID bool

		next := http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			nextHandlerCalled = true
			gotUserID, contextGotUserID = pkg.GetUserID(req.Context())

			res.WriteHeader(http.StatusNoContent)
		})

		req := httptest.NewRequest(http.MethodGet, "/v1/users/me", nil)
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: wantToken,
		})

		recorder := httptest.NewRecorder()

		Auth(sessionFinder)(next).ServeHTTP(recorder, req)

		if got, want := recorder.Code, http.StatusNoContent; got != want {
			t.Errorf("status: got %v, want %v", got, want)
		}

		if !nextHandlerCalled {
			t.Error("next handler was not called")
		}

		if !contextGotUserID {
			t.Error("user ID was not added to context")
		}

		if got, want := gotUserID, wantUserID; got != want {
			t.Errorf("userID: got %v, want %v", got, want)
		}

		if got, want := sessionFinder.gotToken, wantToken; got != want {
			t.Errorf("token: got %v, want %v", got, want)
		}

		if got, want := sessionFinder.calls, 1; got != want {
			t.Errorf("repository calls: got %v, want %v", got, want)
		}
	})
}
