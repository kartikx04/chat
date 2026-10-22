package middleware

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"

	"github.com/kartikx04/chat/pkg"
)

func TestRequestID(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))

	var gotRequestID string
	var contextContainRequestID bool

	next := http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		gotRequestID, contextContainRequestID = pkg.GetRequestID(req.Context())
		res.WriteHeader(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	recorder := httptest.NewRecorder()

	RequestID(logger)(next).ServeHTTP(recorder, req)

	if got, want := recorder.Code, http.StatusNoContent; got != want {
		t.Errorf("status: got %v, want %v", got, want)
	}

	if !contextContainRequestID {
		t.Error("RequestID was not added to context")
	}

	if _, err := uuid.Parse(gotRequestID); err != nil {
		t.Errorf("request ID is not a valid UUID: %q", gotRequestID)
	}

	if got, want := recorder.Header().Get("X-Request-ID"), gotRequestID; got != want {
		t.Errorf("requestID: got %v, want %v", got, want)
	}
}
