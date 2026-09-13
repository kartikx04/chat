package integrationtest

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/kartikx04/chat/cmd/app/config"
	"github.com/kartikx04/chat/internal/api/route"
	"github.com/kartikx04/chat/internal/database"
	applogger "github.com/kartikx04/chat/internal/logger"
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

	logger := applogger.Init(cfg.Server.Env)

	testDB, err := database.Init(cfg)
	if err != nil {
		t.Fatalf("failed to initialise database: %v", err)
	}

	router := chi.NewRouter()

	timeout := 3 * time.Second
	route.Register(cfg, timeout, testDB, router, logger)

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
