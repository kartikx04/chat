package controllers

import (
	"fmt"
	"log"
	"net/http"

	redisrepo "github.com/kartikx04/chat/internal/redis-repo"
	"github.com/kartikx04/chat/pkg"
	"github.com/rs/cors"
)

func StartHTTPServer() {
	// initialise redis
	redisClient := redisrepo.InitRedis()
	defer redisClient.Close()

	// create indexes
	// redisClient.CreateFetchChatBetweenIndex()

	r := http.NewServeMux()

	// test route health

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("ok")
	})

	r.HandleFunc("/google-sso", GoogleSignOn)
	r.HandleFunc("/callback", Callback)
	r.HandleFunc("/", RenderPage)
	r.HandleFunc("/home", Home)

	// Use default options
	handler := cors.Default().Handler(r)

	log.Printf("Server running on :%s\n", pkg.LoadFile("SERVER_PORT"))
	http.ListenAndServe(fmt.Sprintf(":%s", pkg.LoadFile("SERVER_PORT")), handler)
}
