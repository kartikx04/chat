package integrationtest

import (
	"context"
	"fmt"
	"log/slog"
	"os"
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
		t.Fatalf("load integration-test configuration: %v", err)
	}

	if err := validateTestConfig(cfg); err != nil {
		t.Fatalf("refusing to run destructive integration-test setup: %v", err)
	}

	logger := applogger.Init(cfg.Server.Env)

	testDB, err := database.Init(cfg)
	if err != nil {
		t.Fatalf("initialize integration database %q at %s:%s: %v", cfg.Database.DBName, cfg.Database.Host, cfg.Database.Port, err)
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
			t.Errorf("integration-test cleanup: get database handle for %q: %v", cfg.Database.DBName, err)
			return
		}

		defer func() {
			if err := sqlDB.Close(); err != nil {
				t.Errorf("integration-test cleanup: close database %q: %v", cfg.Database.DBName, err)
			}
		}()

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
			t.Errorf("integration-test cleanup: truncate tables in %q: %v", cfg.Database.DBName, err)
		}
	})

	return app
}

func validateTestConfig(cfg *config.App) error {
	if appEnv := os.Getenv("APPENV"); appEnv != "integration-test" {
		return fmt.Errorf("APPENV must be %q, got %q", "integration-test", appEnv)
	}

	if cfg.Server.Env != "integration-test" {
		return fmt.Errorf("integration tests require APPENV=integration-test")
	}

	if cfg.Database.DBName != "chat_test" {
		return fmt.Errorf("integration tests require DB_NAME=chat_test")
	}

	return nil
}
