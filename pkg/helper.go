package pkg

import (
	"log/slog"
	"time"

	"github.com/goombaio/namegenerator"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		slog.Error(".env file not loaded", "error", err)
	} else {
		slog.Info(".env file loaded")
	}
}

func GenerateUniqueName() string {
	seed := time.Now().UTC().UnixNano()
	nameGenerator := namegenerator.NewNameGenerator(seed)

	name := nameGenerator.Generate()
	return name
}
