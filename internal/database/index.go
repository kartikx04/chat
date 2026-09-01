package database

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/kartikx04/chat/cmd/app/config"
	applogger "github.com/kartikx04/chat/internal/logger"
	_ "github.com/lib/pq"
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init(cfg *config.App) (*gorm.DB, error) {

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Database.Host, cfg.Database.User, cfg.Database.Password, cfg.Database.DBName, cfg.Database.Port, cfg.Database.SSLMode,
	)

	// Run migrations first using raw sql.DB
	runMigrations(dsn, cfg.Database.DBName)

	// Then open GORM connection for the app
	db, err := gorm.Open(gormpostgres.Open(dsn), &gorm.Config{
		Logger: applogger.NewSlogGormLogger(cfg.Server.Env),
	})
	if err != nil {
		slog.Error("db initialize error", "error", err)
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		slog.Error("failed to get underlying sql.DB", "error", err)
		return nil, err
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	slog.Info("db pool configured", "max_open", 25, "max_idle", 10)

	DB = db
	slog.Info("database connected", "host", cfg.Database.Host, "name", cfg.Database.DBName)
	return DB, nil
}

func runMigrations(dsn, dbName string) {
	// golang-migrate needs a *sql.DB, not GORM
	sqlDB, err := sql.Open("postgres", dsn)
	if err != nil {
		slog.Error("failed to open db for migrations", "error", err)
		os.Exit(1)
	}
	defer sqlDB.Close()

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		slog.Error("failed to create migration driver", "error", err)
		os.Exit(1)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/database/migrations",
		dbName,
		driver,
	)
	if err != nil {
		slog.Error("failed to create migrator", "error", err)
		os.Exit(1)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		slog.Error("migration failed", "error", err)
		os.Exit(1)
	}

	version, dirty, _ := m.Version()
	if dirty {
		slog.Error("database is in dirty migration state — manual fix required",
			"version", version,
		)
		os.Exit(1)
	}

	slog.Info("migrations applied", "version", version)
}
