package connect

import (
	"fmt"
	"go-chat/config"
	"runtime"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
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

	connectConfig := config.Conf.Connect

	runtime.GOMAXPROCS(connectConfig.ConnectBucket.CpuNum)

	if err := c.InitLogicRpcClient(); err != nil {
		logrus.Panicf("InitLogicRpcClient err:%s", err.Error())
	}
	operator := new(DefaultOperator)
	DefaultServer = NewServer(operator, ServerOptions{
		WriteWait:       10 * time.Second,
		PongWait:        60 * time.Second,
		PingPeriod:      54 * time.Second,
		MaxMessageSize:  512,
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		BroadcastSize:   512,
	})
	c.ServerId = fmt.Sprintf("%s-%s", "ws", uuid.New().String())
	if err := c.InitConnectWebsocketRpcServer(); err != nil {
		logrus.Panicf("InitConnectWebsocketRpcServer Fatal error: %s \n", err.Error())
	}

	c.InitWebSocket()
}
