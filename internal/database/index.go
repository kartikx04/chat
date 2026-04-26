package database

import (
	"fmt"
	"log/slog"
	"os"

	applogger "github.com/kartikx04/chat/internal/logger"
	"github.com/kartikx04/chat/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

var DB *gorm.DB

func InitDB(cfg Config) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode,
	)

	env := os.Getenv("ENV")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: applogger.NewSlogGormLogger(env),
	})
	if err != nil {
		slog.Error("db initialize error", "error", err)
		os.Exit(1)
	}

	DB = db

	if err = DB.AutoMigrate(&models.Chat{}, &models.Users{}); err != nil {
		slog.Error("migration failed", "error", err)
		os.Exit(1)
	}

	slog.Info("database connected", "host", cfg.Host, "name", cfg.DBName)
}
