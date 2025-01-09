package connect

import (
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/ycj3/go-chat/server/pb"
	"google.golang.org/protobuf/proto"
)

type Server struct {
	Options ServerOptions
}

type ServerOptions struct {
	WriteWait       time.Duration
	PongWait        time.Duration
	PingPeriod      time.Duration
	MaxMessageSize  int64
	ReadBufferSize  int
	WriteBufferSize int
	BroadcastSize   int
}

func NewServer(options ServerOptions) *Server {
	return &Server{
		Options: options,
	}
}

func (s *Server) readPump(ss *Session) {
	defer func() {
		ss.conn.Close()
	}()
	ss.conn.SetReadLimit(s.Options.MaxMessageSize)
	ss.conn.SetReadDeadline(time.Now().Add(s.Options.PongWait))
	ss.conn.SetPongHandler(func(string) error { ss.conn.SetReadDeadline(time.Now().Add(s.Options.PongWait)); return nil })
	for {
		_, message, err := ss.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logrus.Errorf("Unexpected close error: %v", err)
			} else {
				logrus.Debugf("Connection closed: %v", err)
			}
			break
		}
		chatMessage := &pb.ChatMessage{}
		if err := proto.Unmarshal(message, chatMessage); err != nil {
			logrus.Errorf("Failed to unmarshal message: %v", err)
			continue
		}
		ss.broadcast <- chatMessage
		logrus.Debugf("Broadcasted message from user %s", chatMessage.User)
	}
}

func (s *Server) writePump(ss *Session) {
	ticker := time.NewTicker(s.Options.PingPeriod)
	defer func() {
		ticker.Stop()
		ss.conn.Close()
	}()
	for {
		select {
		case chatMessage, ok := <-ss.broadcast:
			ss.conn.SetWriteDeadline(time.Now().Add(s.Options.WriteWait))
			if !ok {
				// The hub closed the channel.
				ss.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := ss.conn.NextWriter(websocket.BinaryMessage)
			if err != nil {
				return
			}
			message, err := proto.Marshal(chatMessage)
			if err != nil {
				logrus.Debug("Error marshaling chat message:", err)
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(ss.broadcast)
			for i := 0; i < n; i++ {
				message, err := proto.Marshal(<-ss.broadcast)
				if err != nil {
					logrus.Debug("Error marshaling chat message:", err)
					return
				}
				w.Write(message)
			}

			if err := w.Close(); err != nil {
				logrus.Debug("Error closing writer:", err)
				return
			}
		case <-ticker.C:
			ss.conn.SetWriteDeadline(time.Now().Add(s.Options.WriteWait))
			if err := ss.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				logrus.Debug("Error writing ping message:", err)
				return
			}
		}
	}
}
