package redisrepo

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

// Initialize Redis Connection
func InitRedis() {
	var client *redis.Client
	var err error

	redisURL := os.Getenv("REDIS_URL")

	if redisURL != "" {
		// ✅ Preferred: use full URL (Render internal connection)
		opt, err := redis.ParseURL(redisURL)
		if err != nil {
			log.Fatalf("❌ Failed to parse Redis URL: %v", err)
		}
		client = redis.NewClient(opt)
	} else {
		// fallback (not recommended for Render)
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

	// Test connection
	_, err = client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("❌ Failed to connect to Redis: %v", err)
	}

	redisClient = client
	log.Println("✅ Redis connection established")
}
