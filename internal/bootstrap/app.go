package bootstrap

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/kartikx04/chat/cmd/app/config"
	"github.com/kartikx04/chat/internal/api"
	"github.com/kartikx04/chat/internal/api/route"
	"github.com/kartikx04/chat/internal/database"
	applogger "github.com/kartikx04/chat/internal/logger"
)

func Run(cfg *config.App) {
	var err error

	logger := applogger.Init(cfg.Server.Env)

	db, err := database.Init(cfg)
	if err != nil {
		slog.Error("database initilization error", "error", err)
		os.Exit(1)
	}

	chi := chi.NewRouter()

	timeout := time.Duration(3) * time.Second
	route.Register(cfg, timeout, db, chi, logger)

	server := api.NewHTTPServer(cfg, chi)

	slog.Info("server running", "addr", server.Addr)

	// ListenAndServe blocks forever
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("server error", "error", err)
		os.Exit(1)
	}
}
