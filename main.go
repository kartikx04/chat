package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/kartikx04/chat/controllers"
	"github.com/kartikx04/chat/models"
	"github.com/kartikx04/chat/utils"
)

func health(w http.ResponseWriter, r *http.Request) {
	log.Printf("ok")
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", controllers.RenderPage)
	mux.HandleFunc("/google-sso", controllers.GoogleSignOn)
	mux.HandleFunc("/callback", controllers.Callback)
	mux.HandleFunc("/home", controllers.Home)

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

	srv := &http.Server{
		Handler:      mux,
		Addr:         "127.0.0.1:8007",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Printf("listening on port: %s\n", srv.Addr)

	mux.HandleFunc("/health", health)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
