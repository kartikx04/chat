package pkg

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadFile(key string) string {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatal("error loading .env file")
	}
	return os.Getenv(key)
}
