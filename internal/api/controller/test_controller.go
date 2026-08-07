package controller

import (
	"log/slog"
	"net/http"

	"github.com/kartikx04/chat/internal/database"
)

func Health(w http.ResponseWriter, req *http.Request) {
	sqlDB, err := database.DB.DB()
	if err != nil || sqlDB.Ping() != nil {
		slog.ErrorContext(req.Context(), "health: db unreachable")
		http.Error(w, `{"status":"error","db":false}`, http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
