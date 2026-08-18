package pkg

import (
	"log/slog"
	"os"
	"time"

	"github.com/goombaio/namegenerator"
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

func GenerateUniqueName() string {
	seed := time.Now().UTC().UnixNano()
	nameGenerator := namegenerator.NewNameGenerator(seed)

	name := nameGenerator.Generate()
	return name
}
