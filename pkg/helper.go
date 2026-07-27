package pkg

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		slog.Error(".env file not loaded", "error", err)
	} else {
		slog.Info(".env file loaded")
	}
}

func LoadFile(key string) string {
	val := os.Getenv(key)
	if val == "" {
		slog.Error("missing required environment variable", "key", key)
		os.Exit(1)
	}
	return val
}
