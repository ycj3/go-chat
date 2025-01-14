package connect

import (
	"encoding/json"
	"time"

	"go-chat/proto"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type Server struct {
	Options  ServerOptions
	operator Operator
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

func NewServer(o Operator, options ServerOptions) *Server {
	return &Server{
		Options:  options,
		operator: o,
	}
}

func (s *Server) readPump(ss *Session, c *Connect) {
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
		var connReq *proto.ConnectRequest
		logrus.Infof("get a message :%s", message)
		if err := json.Unmarshal([]byte(message), &connReq); err != nil {
			logrus.Errorf("message struct %+v", connReq)
		}
		if connReq == nil || connReq.AuthToken == "" {
			logrus.Errorf("s.operator.Connect no authToken")
			return
		}
		connReq.ServerId = c.ServerId
		userId, err := s.operator.Connect(connReq)
		if err != nil {
			logrus.Errorf("s.operator.Connect error %s", err.Error())
			return
		}
		if userId == 0 {
			logrus.Error("Invalid AuthToken ,userId empty")
			return
		}
		logrus.Infof("websocket rpc call return userId:%d,RoomId:%d", userId, connReq.RoomId)
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
			message, err := json.Marshal(chatMessage)
			if err != nil {
				logrus.Debug("Error marshaling chat message:", err)
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(ss.broadcast)
			for i := 0; i < n; i++ {
				message, err := json.Marshal(<-ss.broadcast)
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
