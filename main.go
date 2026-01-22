package main

import (
	"log"
	"net/http"
)

func health(w http.ResponseWriter, r *http.Request) {
	log.Printf("ok")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", health)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
