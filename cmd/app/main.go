package main

import (
	"log/slog"
	"os"

	"github.com/kartikx04/chat/cmd/app/config"
	"github.com/kartikx04/chat/internal/app"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		slog.Error("error creating config", "error", err)
		os.Exit(1)
	}

	app.Run(cfg)
}
