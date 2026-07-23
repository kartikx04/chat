package ws

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/kartikx04/chat/internal/models"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 4096
)

type Client struct {
	hub      *Hub
	Conn     *websocket.Conn
	Username string
	Id       uuid.UUID
	send     chan *models.Chat
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func (c *Client) readPump() {

}

func (c *Client) writePump() {

}

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("ws upgrade failed", "error", err, "remote", r.RemoteAddr)
		return
	}

	slog.Debug("ws connection upgraded", "remote", r.RemoteAddr)

	client := &Client{
		hub:  hub,
		Conn: conn,
		send: make(chan *models.Chat, 256),
	}

	hub.register <- client
	go client.writePump()
	go client.readPump()
}
