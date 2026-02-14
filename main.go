package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/kartikx04/chat/models"
	"github.com/kartikx04/chat/routes"
	"github.com/kartikx04/chat/utils"
)

func health(w http.ResponseWriter, r *http.Request) {
	log.Printf("ok")
}

func main() {
	mux := http.NewServeMux()

	// routes
	routes.UIRoutes(mux)
	routes.AuthRoutes(mux)

	//config for database
	config := models.Config{
		Host:     utils.LoadFile("DB_HOST"),
		Port:     utils.LoadFile("DB_PORT"),
		User:     utils.LoadFile("DB_USER"),
		Password: utils.LoadFile("DB_PASSWORD"),
		DBName:   utils.LoadFile("DB_NAME"),
		SSLMode:  utils.LoadFile("DB_SSLMODE"),
	}

	// Initialize DB
	models.InitDB(config)

	// server instance
	srv := &http.Server{
		Handler:      mux,
		Addr:         "127.0.0.1:8007",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Printf("listening on port: %s\n", srv.Addr)

	// test route health
	mux.HandleFunc("/health", health)

	// running server
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
