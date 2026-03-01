package controllers

import (
	"log"
	"net/http"

	redisrepo "github.com/kartikx04/chat/internal/redis-repo"
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

	log.Println("Server running on :8007")
	http.ListenAndServe(":8007", handler)
}
