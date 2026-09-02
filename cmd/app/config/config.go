package config

import (
	"log/slog"

	"github.com/caarlos0/env/v11"
)

type App struct {
	Database   Database
	Server     Server
	GoogleAuth GoogleAuth
}

type Database struct {
	Host           string `env:"DB_HOST,required"`
	Port           string `env:"DB_PORT,required"`
	User           string `env:"DB_USER,required"`
	Password       string `env:"DB_PASSWORD,required"`
	DBName         string `env:"DB_NAME,required"`
	SSLMode        string `env:"DB_SSLMODE,required"`
	ContextTimeout int    `env:"DB_CONTEXT_TIMEOUT" envDefault:"5"`
}

type Server struct {
	Port string `env:"SERVER_PORT,required"`
	Env  string `env:"ENV" envDefault:"development"`
}

type GoogleAuth struct {
	ClientID     string `env:"CLIENT_ID,required"`
	ClientSecret string `env:"CLIENT_SECRET,required"`
	RedirectURL  string `env:"REDIRECT_URL,required"`
}

func New() (*App, error) {
	LoadEnv()

	cfg := &App{}
	err := env.Parse(cfg)
	if err != nil {
		return nil, err
	}
	slog.Info("config loaded")

	return cfg, nil
}
