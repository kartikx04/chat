package pkg

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

func InitEnv() {
	_ = godotenv.Load()
}

func LoadFile(key string) string {
	val := os.Getenv(key)
	if val == "" {
		slog.Error("missing required environment variable", "key", key)
		os.Exit(1)
	}
	return val
}
