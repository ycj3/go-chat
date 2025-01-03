package websocket

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/ycj3/go-chat/server/models"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	// The user associated with the connection.
	user *models.User
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
		logrus.Debug("Client unregistered and connection closed for user:", c.user.UserID)
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logrus.Errorf("Unexpected close error: %v", err)
			} else {
				logrus.Debugf("Connection closed: %v", err)
			}
			break
		}
		logrus.Debugf("Message received from user %s: %s", c.user.UserID, string(message))
		// Handle heartbeat messages
		var msg map[string]interface{}
		if err := json.Unmarshal(message, &msg); err == nil {
			if msg["type"] == "heartbeat" {
				c.user.LastActive = time.Now().Unix()
				c.user.IsOnline = c.user.IsCurrentlyOnline()
				// Update the user's last_active timestamp in the cache or database
				if err := c.hub.db.Save(c.user).Error; err != nil {
					logrus.Errorf("Error updating last_active for user %s: %v", c.user.UserID, err)
				} else {
					logrus.Debugf("Updated last_active for user %s", c.user.UserID)
				}
				continue
			}
		}
		c.hub.broadcast <- message
		logrus.Debugf("Broadcasted message from user %s", c.user.UserID)
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The application
// ensures that there is at most one writer to a connection by executing all
// writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.BinaryMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				logrus.Debug("Error closing writer:", err)
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				logrus.Debug("Error writing ping message:", err)
				return
			}
		}
	}
}

// ServeWs handles websocket requests from the peer.
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	logrus.Debug("ServeWs called with URL:", r.URL)
	userID := r.URL.Query().Get("user_id")
	logrus.Debug("User ID from query parameter:", userID)
	user, err := models.GetUserByID(hub.db, userID)
	if err != nil {
		logrus.Debug("Error getting user by ID:", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logrus.Debug("Error upgrading connection:", err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256), user: user}
	client.hub.register <- client
	logrus.Debug("Client registered:", client.user.UserID)
	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
	logrus.Debug("Client writePump and readPump started for user:", client.user.UserID)
}
