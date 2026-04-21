package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/kartikx04/chat/pkg"
	"github.com/rs/cors"
)

func StartHTTPServer() {
	r := http.NewServeMux()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("ok")
	})

	r.HandleFunc("/google-sso", GoogleSignOn)
	r.HandleFunc("/callback", Callback)
	r.HandleFunc("/", RenderPage)
	r.HandleFunc("/home", Home)

	r.HandleFunc("/contacts", contactListHandler)
	r.HandleFunc("/chat-history", chatHistoryHandler)
	r.HandleFunc("/add-contact", addContactHandler)

	log.Printf("Server running on :%s\n", pkg.LoadFile("SERVER_PORT"))
	http.ListenAndServe(fmt.Sprintf(":%s", pkg.LoadFile("SERVER_PORT")), handler)
}
