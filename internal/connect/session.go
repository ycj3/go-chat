package connect

import (
	"go-chat/models"
	"go-chat/pb"

	"github.com/gorilla/websocket"
)

type Session struct {
	Channel   *Channel
	Next      *Channel
	Prev      *Channel
	broadcast chan *pb.ChatMessage
	user      *models.User
	conn      *websocket.Conn
}

func NewSession(user *models.User, conn *websocket.Conn, channel *Channel) *Session {
	return &Session{
		Channel:   channel,
		broadcast: make(chan *pb.ChatMessage),
		user:      user,
		conn:      conn,
	}
}
