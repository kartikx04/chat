package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	var filename string

	appEnv := os.Getenv("APPENV")

	switch appEnv {
	case "development":
		slog.Info("Running in DEVELOPMENT mode.")
		filename = ".env.dev"
	case "production":
		slog.Info("Running in PRODUCTION mode.")
		filename = ".env"
	case "integration-test":
		slog.Info("Running in Test mode.")
		filename = ".env.test"
	default:
		slog.Warn("No matching APPENV found. Proceeding with system environment variables only.", "APPENV", appEnv)
		return
	}

	if filename != "" {
		if err := godotenv.Load(filename); err != nil {
			slog.Error("env file not loaded", "file", filename, "error", err)
		} else {
			slog.Info("env file loaded", "file", filename)
		}
	}
}
