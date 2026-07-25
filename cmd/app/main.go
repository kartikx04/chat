package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/kartikx04/chat/cmd/app/config"
	"github.com/kartikx04/chat/internal/controllers"
	"github.com/kartikx04/chat/internal/database"
	applogger "github.com/kartikx04/chat/internal/logger"
	"github.com/kartikx04/chat/pkg"
)

func main() {
	pkg.InitEnv()

	env := pkg.LoadFile("ENV")
	applogger.Init(env)

	slog.Info("server starting", "env", env, "port", pkg.LoadFile("SERVER_PORT"))

	dbConfig := config.Database()
	database.InitDB(*dbConfig)

	server := controllers.NewHTTPServer()

	// Start server in background
	slog.Info("server running", "addr", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("server error", "error", err)
		os.Exit(1)
	}
}
