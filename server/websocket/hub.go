package websocket

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
	userConnections map[string]*Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:       make(chan []byte),
		register:        make(chan *Client),
		unregister:      make(chan *Client),
		clients:         make(map[*Client]bool),
		userConnections: make(map[string]*Client),
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
func (h *Hub) GetClientByUser(user string) (*Client, bool) {
	client, ok := h.userConnections[user]
	return client, ok
}

// GetOnlineCount returns the number of online users.
func (h *Hub) GetOnlineCount() int {
	return len(h.userConnections)
}

// GetOnlineMembers returns the list of online users.
func (h *Hub) GetOnlineMembers() []string {
	members := make([]string, 0, len(h.userConnections))
	for user := range h.userConnections {
		members = append(members, user)
	}
	return members
}
