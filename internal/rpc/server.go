package rpc

import (
	"fmt"
	"go-chat/config"
	"go-chat/internal/redis"
	"go-chat/tools"
	"runtime"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rcrowley/go-metrics"
	"github.com/rpcxio/rpcx-etcd/serverplugin"
	"github.com/sirupsen/logrus"
	"github.com/smallnest/rpcx/server"
)

type RpcServer struct {
	ServerId string
	Rss      []RegisterService
}

type RegisterService struct {
	Name string
	Rcvr interface{}
}

func NewRpcServer(rss []RegisterService) *RpcServer {
	s := new(RpcServer)
	s.Rss = rss
	return s
}

func (srv *RpcServer) Run() {
	//read config
	logicConfig := config.Conf.Logic

	runtime.GOMAXPROCS(logicConfig.LogicBase.CpuNum)
	srv.ServerId = fmt.Sprintf("logic-%s", uuid.New().String())
	//init publish redis
	if err := redis.InitRedisClient(); err != nil {
		logrus.Panicf("logic init publishRedisClient fail,err:%s", err.Error())
	}

	//init rpc server
	if err := srv.InitRpcServer(); err != nil {
		logrus.Panicf("logic init rpc server fail")
	}
	select {}
}
func (srv *RpcServer) InitRpcServer() (err error) {
	logrus.Info("Initializing RPC Server")
	var network, addr string
	// a host multi port case
	rpcAddressList := strings.Split(config.Conf.Logic.LogicBase.RpcAddress, ",")
	for _, bind := range rpcAddressList {
		if network, addr, err = tools.ParseNetwork(bind); err != nil {
			logrus.Panicf("InitLogicRpc ParseNetwork error : %s", err.Error())
		}
		logrus.Infof("logic start run at-->%s:%s", network, addr)
		for _, rs := range srv.Rss {
			go srv.createRpcServer(rs.Name, rs.Rcvr, network, addr)
		}
	}
	logrus.Info("RPC Server initialized successfully")
	return
}

func (srv *RpcServer) createRpcServer(name string, rcvr interface{}, network string, addr string) {
	logrus.Infof("Creating RPC Server at %s:%s", network, addr)
	s := server.NewServer()
	srv.addRegistryPlugin(s, network, addr)

	if err := s.RegisterName(name, rcvr, fmt.Sprintf("%s", srv.ServerId)); err != nil {
		logrus.Errorf("register error:%s", err.Error())
	} else {
		logrus.Infof("RPC Server registered with ID: %s", srv.ServerId)
	}

	s.RegisterOnShutdown(func(s *server.Server) {
		logrus.Info("Shutting down RPC Server")
		s.UnregisterAll()
	})
	s.Serve(network, addr)
}

func (srv *RpcServer) addRegistryPlugin(s *server.Server, network string, addr string) {
	logrus.Infof("Adding registry plugin for %s:%s", network, addr)
	r := &serverplugin.EtcdV3RegisterPlugin{
		ServiceAddress: network + "@" + addr,
		EtcdServers:    []string{config.Conf.Common.CommonEtcd.Host},
		BasePath:       config.Conf.Common.CommonEtcd.BasePath,
		Metrics:        metrics.NewRegistry(),
		UpdateInterval: time.Minute,
	}
	err := r.Start()
	if err != nil {
		logrus.Fatal(err)
	} else {
		logrus.Info("Registry plugin added successfully")
	}
	s.Plugins.Add(r)
}
