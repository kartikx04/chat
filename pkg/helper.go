package pkg

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func InitEnv() {
	// Load .env only in development (won't error if file doesn't exist)
	_ = godotenv.Load()
}

func LoadFile(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("Missing required environment variable: %s", key)
	}
	return val
}
