package config

import (
	"github.com/kartikx04/chat/pkg"
)

type Database struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func GetDatabase() *Database {
	return &Database{
		Host:     pkg.LoadFile("DB_HOST"),
		Port:     pkg.LoadFile("DB_PORT"),
		User:     pkg.LoadFile("DB_USER"),
		Password: pkg.LoadFile("DB_PASSWORD"),
		DBName:   pkg.LoadFile("DB_NAME"),
		SSLMode:  pkg.LoadFile("DB_SSLMODE"),
	}
}
