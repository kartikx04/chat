package main

import (
	"context"
	"log/slog"
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

	// Run server in goroutine so it doesn't block signal handling
	go controllers.StartHTTPServer()

	// Block until SIGINT or SIGTERM
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit

	slog.Info("shutdown signal received", "signal", sig.String())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := redisrepo.Close(); err != nil {
		slog.Error("redis close error", "error", err)
	}

	<-ctx.Done()
	slog.Info("server stopped")
}
