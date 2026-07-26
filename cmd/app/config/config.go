package config

import (
	"log/slog"

	"github.com/joho/godotenv"
)

type App struct {
	Database   database
	Server     server
	JWT        jwt
	GoogleAuth googleAuth
}

type database struct {
	Host     string `env:"DB_HOST,required"`
	Port     string `env:"DB_PORT,required"`
	User     string `env:"DB_USER,required"`
	Password string `env:"DB_PASSWORD,required"`
	DBName   string `env:"DB_NAME,required"`
	SSLMode  string `env:"DB_SSLMODE,required"`
}

type server struct {
	port string `env:"SERVER_PORT,required"`
	env  string `env:"ENV" envDefault:"development"`
}

type jwt struct {
	secret string `env:"JWT_SECRET,required"`
}

type googleAuth struct {
	clientId     string `env:"CLIENT_ID,required"`
	clientSecret string `env:"CLIENT_SECRET,required"`
	redirectURL  string `env:"REDIRECT_URL,required"`
	tokenSecret  string `env:"TOKEN_SECRET,required"`
}

func New() (*App, error) {
	cfg := &App{}
	if err := godotenv.Load("../../.env"); err != nil {
		return nil, err
	}
	slog.Info(".env loaded")
	return cfg, nil
}

// func GetDatabase() *Database {
// 	return &Database{
// 		Host:     pkg.LoadFile("DB_HOST"),
// 		Port:     pkg.LoadFile("DB_PORT"),
// 		User:     pkg.LoadFile("DB_USER"),
// 		Password: pkg.LoadFile("DB_PASSWORD"),
// 		DBName:   pkg.LoadFile("DB_NAME"),
// 		SSLMode:  pkg.LoadFile("DB_SSLMODE"),
// 	}
// }
