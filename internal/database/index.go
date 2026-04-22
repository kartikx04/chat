package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/kartikx04/chat/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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

	// Default to Warn level
	logLevel := logger.Warn

	// Check ENV after initialization
	if os.Getenv("ENV") == "development" {
		logLevel = logger.Info
	}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logLevel,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		log.Fatal("db initialize error:", err)
	}

	DB = db

	// Migrate
	err = DB.AutoMigrate(&models.Chat{}, &models.Users{})
	if err != nil {
		log.Fatal("migration failed:", err)
	}

	fmt.Println("✅ Database connected")
}
