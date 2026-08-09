package config

import (
	"log"
	"log/slog"

	"github.com/caarlos0/env/v11"
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
	Port string `env:"SERVER_PORT,required"`
	Env  string `env:"ENV" envDefault:"development"`
}

type jwt struct {
	Secret string `env:"JWT_SECRET,required"`
}

type googleAuth struct {
	ClientId     string `env:"CLIENT_ID,required"`
	ClientSecret string `env:"CLIENT_SECRET,required"`
	RedirectURL  string `env:"REDIRECT_URL,required"`
	TokenSecret  string `env:"TOKEN_SECRET,required"`
}

func New() (*App, error) {
	cfg := &App{}
	err := env.Parse(cfg)
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}
	slog.Info(".env loaded")
	return cfg, nil
}
