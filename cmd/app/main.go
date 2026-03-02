package main

import (
	"github.com/kartikx04/chat/internal/controllers"
	"github.com/kartikx04/chat/internal/database"
	"github.com/kartikx04/chat/pkg"
)

func main() {
	//config for database
	config := database.Config{
		Host:     pkg.LoadFile("DB_HOST"),
		Port:     pkg.LoadFile("DB_PORT"),
		User:     pkg.LoadFile("DB_USER"),
		Password: pkg.LoadFile("DB_PASSWORD"),
		DBName:   pkg.LoadFile("DB_NAME"),
		SSLMode:  pkg.LoadFile("DB_SSLMODE"),
	}

	// Initialize DB
	database.InitDB(config)

	controllers.StartHTTPServer()
}
