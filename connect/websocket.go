package connect

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (c *Connect) handleConnections(server *Server, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logrus.Debug("Error upgrading connection:", err)
		return
	}

	ss := NewSession(nil, conn, nil)
	go server.writePump(ss)
	go server.readPump(ss, c)
}

func (c *Connect) InitWebSocket() {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		c.handleConnections(DefaultServer, w, r)
	})
	log.Println("WebSocket server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
