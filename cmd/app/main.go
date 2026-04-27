package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kartikx04/chat/internal/controllers"
	"github.com/kartikx04/chat/internal/database"
	applogger "github.com/kartikx04/chat/internal/logger"
	redisrepo "github.com/kartikx04/chat/internal/redis-repo"
	"github.com/kartikx04/chat/internal/ws"
	"github.com/kartikx04/chat/pkg"
)

func main() {
	pkg.InitEnv()

	env := os.Getenv("ENV")
	applogger.Init(env)

	slog.Info("server starting", "env", env, "port", os.Getenv("SERVER_PORT"))

	config := database.Config{
		Host:     pkg.LoadFile("DB_HOST"),
		Port:     pkg.LoadFile("DB_PORT"),
		User:     pkg.LoadFile("DB_USER"),
		Password: pkg.LoadFile("DB_PASSWORD"),
		DBName:   pkg.LoadFile("DB_NAME"),
		SSLMode:  pkg.LoadFile("DB_SSLMODE"),
	}

	database.InitDB(config)
	redisrepo.InitRedis()
	redisrepo.CreateChatIndex()
	ws.InitHub()

	server := controllers.NewHTTPServer()

	// Start server in background
	go func() {
		slog.Info("server running", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	slog.Info("shutdown signal received", "signal", sig.String())

	// 10 seconds for in-flight HTTP requests to finish
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Step 1 — stop accepting new HTTP requests, drain existing ones
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("http server shutdown error", "error", err)
	}
	slog.Info("http server stopped")

	// Step 2 — close all WebSocket connections
	ws.HubInstance.Shutdown()
	slog.Info("websocket hub stopped")

	// Step 3 — close Redis
	if err := redisrepo.Close(); err != nil {
		slog.Error("redis close error", "error", err)
	}
	slog.Info("redis closed")

	slog.Info("server stopped cleanly")
}
