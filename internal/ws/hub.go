package ws

import (
	"github.com/kartikx04/chat/internal/models"
)

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan *models.Chat
	register   chan *Client
	unregister chan *Client
}

var HubInstance *Hub

func InitHub() {
	HubInstance = NewHub()
	go HubInstance.Run()
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan *models.Chat, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run is the ONLY goroutine that touches hub.clients — no mutex needed.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send) // signals writePump to exit
			}

		case message := <-h.broadcast:
			for client := range h.clients {
				if client.Id != message.FromId && client.Id != message.ToId {
					continue
				}
				msgCopy := *message
				msgCopy.IsSelf = client.Id == message.FromId
				select {
				case client.send <- &msgCopy:
				default:
					delete(h.clients, client)
					close(client.send)
				}
			}
		}
	}
}
