package ws

import (
	"log/slog"

	"github.com/kartikx04/chat/internal/models"
)

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan *models.Chat
	register   chan *Client
	unregister chan *Client
}

var HubInstance *Hub

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan *models.Chat, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func InitHub() {
	HubInstance = NewHub()
	go HubInstance.Run()
	slog.Info("ws hub started")
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			slog.Debug("ws client registered",
				"username", client.Username,
				"total", len(h.clients),
			)

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				slog.Info("ws client disconnected",
					"username", client.Username,
					"total", len(h.clients),
				)
			}

		case message := <-h.broadcast:
			slog.Debug("ws broadcasting",
				"from", message.FromId,
				"to", message.ToId,
			)
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
					slog.Warn("ws send buffer full, dropping client",
						"username", client.Username,
					)
				}
			}
		}
	}
}
