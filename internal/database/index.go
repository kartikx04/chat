package database

import (
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
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
	err := runMigrations(dsn, cfg.Database.DBName)
	if err := runMigrations(dsn, cfg.Database.DBName); err != nil {
		slog.Error("error running migration", "error", err)
		return nil, fmt.Errorf("run migrations for database %q: %w", cfg.Database.DBName, err)
	}

	// Then open GORM connection for the app
	db, err := gorm.Open(gormpostgres.Open(dsn), &gorm.Config{
		Logger: applogger.NewSlogGormLogger(cfg.Server.Env),
	})
	if err != nil {
		slog.Error("db initialize error", "error", err)
		return nil, fmt.Errorf("open gorm connection for database %q: %w", cfg.Database.DBName, err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		slog.Error("failed to get underlying sql.DB", "error", err)
		return nil, fmt.Errorf("return sql db for database %q: %w", cfg.Database.DBName, err)
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	slog.Info("db pool configured", "max_open", 25, "max_idle", 10)

	DB = db
	slog.Info("database connected", "host", cfg.Database.Host, "name", cfg.Database.DBName)
	return DB, nil
}

func runMigrations(dsn, dbName string) error {
	// golang-migrate needs a *sql.DB, not GORM
	sqlDB, err := sql.Open("postgres", dsn)
	if err != nil {
		slog.Error("failed to open db for migrations", "error", err)
		return fmt.Errorf("open database connection: %w", err)
	}
	defer sqlDB.Close()

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		slog.Error("failed to create migration driver", "error", err)
		return fmt.Errorf("create postgres migration driver: %w", err)
	}

	// normal path
	migrationsPath := "file://internal/database/migrations"
	if flag.Lookup("test.v") != nil || flag.Lookup("test.run") != nil {
		// test path
		migrationsPath = "file://../internal/database/migrations"
	}

	m, err := migrate.NewWithDatabaseInstance(
		migrationsPath,
		dbName,
		driver,
	)
	if err != nil {
		slog.Error("failed to create migrator", "error", err)
		return fmt.Errorf("load migrations from %q: %w", migrationsPath, err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		slog.Error("migration failed", "error", err)
		return fmt.Errorf("apply migrations: %w", err)
	}

	version, dirty, _ := m.Version()
	if dirty {
		slog.Error("database is in dirty migration state — manual fix required", "version", version)
		return fmt.Errorf("read migration version: %w", err)
	}

	slog.Info("migrations applied", "version", version)

	return nil
}
