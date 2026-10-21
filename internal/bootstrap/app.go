package bootstrap

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/kartikx04/chat/cmd/app/config"
	"github.com/kartikx04/chat/internal/api"
	"github.com/kartikx04/chat/internal/api/route"
	"github.com/kartikx04/chat/internal/database"
	applogger "github.com/kartikx04/chat/internal/logger"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func Run(cfg *config.App) {
	var err error

	logger := applogger.Init(cfg.Server.Env)

	db, err := database.Init(cfg)
	if err != nil {
		slog.Error("database initilization error", "error", err)
		os.Exit(1)
	}

	router := chi.NewRouter()

	oauthConfig := &oauth2.Config{
		RedirectURL:  cfg.GoogleAuth.RedirectURL,
		ClientID:     cfg.GoogleAuth.ClientID,
		ClientSecret: cfg.GoogleAuth.ClientSecret,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}

	route.Register(cfg, db, router, logger, oauthConfig)

	server := api.NewHTTPServer(cfg, router, logger)

	logger.Info("server running", "addr", server.Addr)

	// ListenAndServe blocks forever
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("server error", "error", err)
		os.Exit(1)
	}
}
