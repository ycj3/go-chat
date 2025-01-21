package main

import (
	"go-chat/config"
	"go-chat/internal/auth"
	"go-chat/internal/rpc"
	"go-chat/pkg/utils"

	"github.com/sirupsen/logrus"
)

func main() {
	utils.InitLogrus()

	rpc.NewRpcServer([]rpc.RegisterService{
		// {config.Conf.Common.CommonEtcd.ServerPathRpc, new(rpc.RpcServer)},
		{config.Conf.Common.CommonEtcd.ServerPathRpc, new(auth.RpcHandler)},
	}).Run()
	logrus.Info("init rpc server success")
}
