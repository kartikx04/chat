package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/kartikx04/chat/internal/database"
	"github.com/kartikx04/chat/internal/models"
	redisrepo "github.com/kartikx04/chat/internal/redis-repo"
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
	send     chan *models.Chat // buffered — writePump is the ONLY writer to Conn
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// readPump — one goroutine, reads from WebSocket, writes to hub.broadcast.
// It is the ONLY goroutine that calls Conn.ReadMessage.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, p, err := c.Conn.ReadMessage()
		if err != nil {
			// Both expected (1000/1001) and unexpected closes exit cleanly here
			if websocket.IsUnexpectedCloseError(err,
				websocket.CloseGoingAway,
				websocket.CloseNormalClosure,
			) {
				log.Println("unexpected close:", err)
			}
			return // always return on any read error
		}

		m := &models.Message{}
		if err := json.Unmarshal(p, m); err != nil {
			log.Println("unmarshal error:", err)
			continue
		}

		if m.Type == "bootup" {
			c.Username = m.User
			id, err := uuid.Parse(m.UserId)
			if err != nil {
				log.Println("invalid user id:", err)
				continue
			}
			c.Id = id
			fmt.Println("client mapped:", c.Username, c.Id)
			continue
		}

		if m.Type == "message" {
			chat := m.Chat
			if chat.FromId == uuid.Nil || chat.ToId == uuid.Nil || chat.Message == "" {
				log.Println("invalid message payload")
				continue
			}

			now := time.Now()
			chat.CreatedAt = now
			chat.CreatedAtUnix = now.Unix()

			// Save to Postgres
			if err := database.DB.Create(&models.Chat{
				FromId:        chat.FromId,
				ToId:          chat.ToId,
				Message:       chat.Message,
				CreatedAt:     now,
				CreatedAtUnix: chat.CreatedAtUnix,
			}).Error; err != nil {
				log.Println("DB save error:", err)
				continue
			}

			// Save to Redis
			id, err := redisrepo.CreateChat(&chat)
			if err != nil {
				log.Println("Redis save error:", err)
				continue
			}
			chat.Id = id

			// Send to hub — non-blocking because hub.broadcast is buffered (256)
			c.hub.broadcast <- &chat
		}
	}
}

// writePump — one goroutine, the ONLY writer to Conn for this client.
// Gorilla WebSocket connections are NOT safe for concurrent writes.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// hub closed the channel — send close frame and exit
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.Conn.WriteJSON(message); err != nil {
				log.Println("write error:", err)
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade error:", err)
		return
	}

	client := &Client{
		hub:  hub,
		Conn: conn,
		send: make(chan *models.Chat, 256), // buffered — never blocks hub.Run()
	}

	hub.register <- client

	// Each client gets exactly two goroutines — one reads, one writes.
	go client.writePump()
	go client.readPump()
}
