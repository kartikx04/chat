package redisrepo

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"os"

	"github.com/kartikx04/chat/internal/database"
	"github.com/kartikx04/chat/pkg"
	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func InitRedis() {
	var client *redis.Client

	redisURL := os.Getenv("REDIS_URL")

	if redisURL != "" {
		opt, err := redis.ParseURL(redisURL)
		if err != nil {
			slog.Error("failed to parse redis url", "error", err)
			os.Exit(1)
		}
		client = redis.NewClient(opt)
	} else {
		addr := fmt.Sprintf("%s:%s",
			pkg.LoadFile("REDIS_HOST"),
			pkg.LoadFile("REDIS_PORT"),
		)
		client = redis.NewClient(&redis.Options{
			Addr:      addr,
			Username:  pkg.LoadFile("REDIS_USER"),
			Password:  pkg.LoadFile("REDIS_PASSWORD"),
			DB:        0,
			TLSConfig: &tls.Config{},
		})
	}

	if _, err := client.Ping(context.Background()).Result(); err != nil {
		slog.Error("redis connection failed", "error", err)
		os.Exit(1)
	}

	redisClient = client

	database.PingRedis = func() error {
		return redisClient.Ping(context.Background()).Err()
	}

	slog.Info("redis connected",
		"host", os.Getenv("REDIS_HOST"),
		"port", os.Getenv("REDIS_PORT"),
	)
}

func Close() error {
	slog.Info("redis closing connection")
	return redisClient.Close()
}
