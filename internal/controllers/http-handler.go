package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/kartikx04/chat/internal/ws"
	"github.com/kartikx04/chat/pkg"
	"github.com/rs/cors"
)

func StartHTTPServer() {
	r := http.NewServeMux()

	c := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:3000",
			"https://chat-0rnj.onrender.com",
			"https://banterrr.vercel.app",
		},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("ok")
	})

	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWs(ws.HubInstance, w, r)
	})

	r.HandleFunc("/google-sso", GoogleSignOn)
	r.HandleFunc("/auth/google/callback", Callback)

	r.HandleFunc("/contacts", contactListHandler)
	r.HandleFunc("/chat-history", chatHistoryHandler)
	r.HandleFunc("/add-contact", addContactHandler)
	r.HandleFunc("/verify-contact", verifyContactHandler)

	log.Printf("Server running on :%s\n", pkg.LoadFile("SERVER_PORT"))
	http.ListenAndServe(fmt.Sprintf(":%s", pkg.LoadFile("SERVER_PORT")),
		http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if req.URL.Path == "/ws" {
				r.ServeHTTP(w, req) // 🚀 bypass CORS
				return
			}
			c.Handler(r).ServeHTTP(w, req)
		}),
	)
}
