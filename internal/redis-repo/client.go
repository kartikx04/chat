package redisrepo

import (
	"context"
	"fmt"
	"log/slog"
	"os"

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
			os.Getenv("REDIS_HOST"),
			os.Getenv("REDIS_PORT"),
		)
		client = redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       0,
		})
	}

	if _, err := client.Ping(context.Background()).Result(); err != nil {
		slog.Error("redis connection failed", "error", err)
		os.Exit(1)
	}

	redisClient = client
	slog.Info("redis connected",
		"host", os.Getenv("REDIS_HOST"),
		"port", os.Getenv("REDIS_PORT"),
	)
}

func Close() error {
	slog.Info("redis closing connection")
	return redisClient.Close()
}
