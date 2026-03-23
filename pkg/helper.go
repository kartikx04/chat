package pkg

import (
	"os"

	"github.com/joho/godotenv"
)

func InitEnv() {
	// Load .env only in local/dev
	if os.Getenv("ENV") == "" || os.Getenv("ENV") == "development" {
		_ = godotenv.Load()
	}
}
func LoadFile(key string) string {
	if os.Getenv("ENV") != "production" {
		_ = godotenv.Load()
	}
	return os.Getenv(key)
}
