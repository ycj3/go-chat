package main

import (
	"go-chat/internal/router"
	"go-chat/internal/rpc"
	"go-chat/pkg/utils"

	"github.com/sirupsen/logrus"
)

func main() {
	utils.InitLogrus()

	//init rpc client
	if err := rpc.InitRpcClient(); err != nil {
		logrus.Fatalf("init rpc client fail:%s", err.Error())
	}
	logrus.Info("init rpc client success")

	router := router.NewRouter()
	router.Run()
}
