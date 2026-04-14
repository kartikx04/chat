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

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("ok")
	})

	r.HandleFunc("/google-sso", GoogleSignOn)
	r.HandleFunc("/callback", Callback)
	r.HandleFunc("/", RenderPage)
	r.HandleFunc("/home", Home)

	handler := cors.Default().Handler(r)

	log.Printf("Server running on :%s\n", pkg.LoadFile("SERVER_PORT"))
	http.ListenAndServe(fmt.Sprintf(":%s", pkg.LoadFile("SERVER_PORT")), handler)
}
