package connect

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

var DefaultServer *Server

type Connect struct {
	ServerId string
}

func NewConnect() *Connect {
	return &Connect{
		ServerId: fmt.Sprintf("%s-%s", "ws", uuid.New().String()),
	}
}

func (c *Connect) Run() {
	DefaultServer = NewServer(
		ServerOptions{
			WriteWait:       10 * time.Second,
			PongWait:        60 * time.Second,
			PingPeriod:      (60 * time.Second * 9) / 10,
			MaxMessageSize:  512,
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			BroadcastSize:   512,
		})

	c.InitWebSocket()
}
