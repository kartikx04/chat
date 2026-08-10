package bootstrap

import (
	"context"
	"log/slog"
	"net/http"
	"os"

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

	ctx := context.Background()
	logger := applogger.Init(env)
	db := database.Init(cfg)

	chi := chi.NewRouter()
	route.Register(ctx, db, chi, logger)

	server := api.NewHTTPServer(cfg, chi)

	slog.Info("server running", "addr", server.Addr)

	// ListenAndServe blocks forever
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("server error", "error", err)
		os.Exit(1)
	}
}
