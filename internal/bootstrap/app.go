package bootstrap

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/kartikx04/chat/cmd/app/config"
	"github.com/kartikx04/chat/internal/api"
	"github.com/kartikx04/chat/internal/database"
	applogger "github.com/kartikx04/chat/internal/logger"
	"github.com/kartikx04/chat/pkg"
)

func Run(cfg *config.App) {
	env := pkg.LoadFile("ENV")

	applogger.Init(env)

	database.InitDB(cfg)

	server := api.NewHTTPServer(cfg)

	slog.Info("server running", "addr", server.Addr)

	// ListenAndServe blocks forever
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("server error", "error", err)
		os.Exit(1)
	}
}
