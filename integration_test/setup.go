package integrationtest

import (
	"context"
	"fmt"
	"log/slog"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/kartikx04/chat/cmd/app/config"
	"github.com/kartikx04/chat/internal/api/route"
	"github.com/kartikx04/chat/internal/database"
	applogger "github.com/kartikx04/chat/internal/logger"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/gorm"
)

type TestApp struct {
	Config  *config.App
	TestDB  *gorm.DB
	Router  chi.Router
	Logger  *slog.Logger
	Context context.Context
}

func Setup(t *testing.T) *TestApp {
	t.Helper()
	cfg, err := config.New()
	if err != nil {
		t.Fatalf("failed to parse env variables: %v", err)
	}

	if err := validateTestConfig(cfg); err != nil {
		t.Fatal(err)
	}

	logger := applogger.Init(cfg.Server.Env)

	testDB, err := database.Init(cfg)
	if err != nil {
		t.Fatalf("failed to initialise database: %v", err)
	}

	router := chi.NewRouter()

	timeout := 3 * time.Second

	oauthConfig := &oauth2.Config{
		RedirectURL:  cfg.GoogleAuth.RedirectURL,
		ClientID:     cfg.GoogleAuth.ClientID,
		ClientSecret: cfg.GoogleAuth.ClientSecret,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}

	route.Register(cfg, timeout, testDB, router, logger, oauthConfig)

	app := &TestApp{
		Config:  cfg,
		TestDB:  testDB,
		Router:  router,
		Logger:  logger,
		Context: context.Background(),
	}

	t.Cleanup(func() {
		sqlDB, err := testDB.DB()
		if err != nil {
			t.Logf("failed to get sql DB: %v", err)
			return
		}

		query := `
            DO $$
            DECLARE
                r RECORD;
            BEGIN
                FOR r IN (
                    SELECT tablename
                    FROM pg_tables
                    WHERE schemaname = 'public'
                      AND tablename != 'schema_migrations'
                ) LOOP
                    EXECUTE 'TRUNCATE TABLE '
                        || quote_ident(r.tablename)
                        || ' CASCADE';
                END LOOP;
            END $$;
        `

		if _, err := sqlDB.Exec(query); err != nil {
			t.Logf("failed to clean database: %v", err)
		}

		sqlDB.Close()
	})

	return app
}

func validateTestConfig(cfg *config.App) error {
	// Replace these with your actual config fields.
	if cfg.Server.Env != "integration-test" {
		return fmt.Errorf(
			"integration tests require APPENV=integration-test",
		)
	}

	if cfg.Database.DBName != "chat_test" {
		return fmt.Errorf(
			"integration tests require DB_NAME=chat_test",
		)
	}

	return nil
}
