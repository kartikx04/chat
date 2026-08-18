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
	"github.com/kartikx04/chat/pkg"
)

func Run(cfg *config.App) {
	env := pkg.LoadFile("ENV")

	logger := applogger.Init(env)
	db := database.Init(cfg)

	chi := chi.NewRouter()

	timeout := time.Duration(3) * time.Second
	route.Register(timeout, db, chi, logger)

	server := api.NewHTTPServer(cfg, chi)

	slog.Info("server running", "addr", server.Addr)

	// ListenAndServe blocks forever
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("server error", "error", err)
		os.Exit(1)
	}
}
