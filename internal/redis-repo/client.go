package redisrepo

import (
	"context"
	"log"

	"github.com/kartikx04/chat/pkg"
	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

// Initialize Redis Connection
func InitRedis() *redis.Client {
	// server instance

	conn := redis.NewClient(&redis.Options{
		Addr:     pkg.LoadFile("REDIS_HOST"),
		Password: pkg.LoadFile("REDIS_PASSWORD"),
		DB:       0,
	})

	_, err := conn.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("❌ Failed to connect to Redis: %v", err)
	}
	log.Println("✅ Redis connection established")

	return conn
}
