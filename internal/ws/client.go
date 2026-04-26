package ws

import (
	"encoding/json"
	"log/slog"
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
	send     chan *models.Chat
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.Conn.Close()
		slog.Debug("ws read pump closed", "username", c.Username)
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
			if websocket.IsUnexpectedCloseError(err,
				websocket.CloseGoingAway,
				websocket.CloseNormalClosure,
			) {
				slog.Warn("ws unexpected close", "username", c.Username, "error", err)
			}
			return
		}

		m := &models.Message{}
		if err := json.Unmarshal(p, m); err != nil {
			slog.Warn("ws unmarshal error", "username", c.Username, "error", err)
			continue
		}

		if m.Type == "bootup" {
			c.Username = m.User
			id, err := uuid.Parse(m.UserId)
			if err != nil {
				slog.Warn("ws invalid user id", "raw", m.UserId, "error", err)
				continue
			}
			c.Id = id
			redisrepo.SetUsernameLookup(id, m.User)
			redisrepo.SetIdLookup(m.User, id)
			slog.Info("ws client mapped", "username", c.Username, "user_id", c.Id)
			continue
		}

		if m.Type == "message" {
			chat := m.Chat
			if chat.FromId == uuid.Nil || chat.ToId == uuid.Nil || chat.Message == "" {
				slog.Warn("ws invalid message payload", "username", c.Username)
				continue
			}

			now := time.Now()
			chat.CreatedAt = now
			chat.CreatedAtUnix = now.Unix()

			if err := database.DB.Create(&models.Chat{
				FromId:        chat.FromId,
				ToId:          chat.ToId,
				Message:       chat.Message,
				CreatedAt:     now,
				CreatedAtUnix: chat.CreatedAtUnix,
			}).Error; err != nil {
				slog.Error("ws db save failed",
					"from", chat.FromId,
					"to", chat.ToId,
					"error", err,
				)
				continue
			}

			id, err := redisrepo.CreateChat(&chat)
			if err != nil {
				slog.Error("ws redis save failed",
					"from", chat.FromId,
					"to", chat.ToId,
					"error", err,
				)
				continue
			}
			chat.Id = id

			slog.Debug("ws message saved and broadcasting",
				"from", chat.FromId,
				"to", chat.ToId,
				"chat_id", chat.Id,
			)

			c.hub.broadcast <- &chat
		}
	}
}

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
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.Conn.WriteJSON(message); err != nil {
				slog.Warn("ws write error", "username", c.Username, "error", err)
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				slog.Debug("ws ping failed, closing", "username", c.Username)
				return
			}
		}
	}
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
