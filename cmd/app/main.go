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
	env := pkg.LoadFile("ENV")
	applogger.Init(env)

	slog.Info("server starting", "env", env, "port", pkg.LoadFile("SERVER_PORT"))

	dbConfig := config.GetDatabase()
	database.InitDB(dbConfig)

	server := controllers.NewHTTPServer()

	slog.Info("server running", "addr", server.Addr)

	// ListenAndServe blocks forever
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("server error", "error", err)
		os.Exit(1)
	}
}
