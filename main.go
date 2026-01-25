package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/kartikx04/chat/controllers"
)

func health(w http.ResponseWriter, r *http.Request) {
	log.Printf("ok")
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", controllers.RenderPage)
	mux.HandleFunc("/google-sso", controllers.GoogleSignOn)
	mux.HandleFunc("/callback", controllers.Callback)

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
