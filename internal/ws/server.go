package ws

import (
	"fmt"
	"log"
	"net/http"
)

func StartWebsocketServer() {
	fmt.Println("Starting WebSocket server on :8081")

	hub := NewHub()
	go hub.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ServeWs(hub, w, r)
	})

	log.Fatal(http.ListenAndServe(":8081", nil))
}
