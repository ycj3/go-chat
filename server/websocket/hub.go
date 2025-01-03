package websocket

import (
	"github.com/sirupsen/logrus"
	"github.com/ycj3/go-chat/server/models"
	"gorm.io/gorm"
)

// Hub maintains the set of active clients and broadcasts messages to the clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool
	// Inbound messages from the clients.
	broadcast chan []byte
	// Register requests from the clients.
	register chan *Client
	// Unregister requests from clients.
	unregister chan *Client
	// User connections cache.
	userConnections map[*models.User]*Client
	// Reference to the database
	db *gorm.DB
}

func NewHub(db *gorm.DB) *Hub {
	return &Hub{
		broadcast:       make(chan []byte),
		register:        make(chan *Client),
		unregister:      make(chan *Client),
		clients:         make(map[*Client]bool),
		userConnections: make(map[*models.User]*Client),
		db:              db,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			h.userConnections[client.user] = client
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				delete(h.userConnections, client.user)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
					delete(h.userConnections, client.user)
				}
			}
		}
	}
}

// GetClientByUser returns the client associated with the given user.
func (h *Hub) GetClientByUser(user *models.User) (*Client, bool) {
	client, ok := h.userConnections[user]
	return client, ok
}

// GetOnlineCount returns the number of online users.
func (h *Hub) GetOnlineCount() int {
	return len(h.GetOnlineMembers())
}

// GetOnlineMembers returns the list of online users.
func (h *Hub) GetOnlineMembers() []models.User {
	onlineUsers, err := models.GetOnlineUsers(h.db)
	if err != nil {
		logrus.Errorf("Error getting online users: %v", err)
		return []models.User{}
	}
	return onlineUsers
}
