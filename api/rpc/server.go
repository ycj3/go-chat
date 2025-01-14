package rpc

import (
	"context"
	"log"
	"net"

	"go-chat/proto"

	"github.com/smallnest/rpcx/server"
)

type LogicService struct{}

func (s *LogicService) Login(ctx context.Context, req *proto.LoginRequest, reply *proto.LoginResponse) error {
	// Implement your login logic here
	// This is a placeholder implementation
	if req.UserID == "test" {
		reply.Code = 0
		reply.AuthToken = "test-token"
	} else {
		reply.Code = 1
		reply.AuthToken = ""
	}
	return nil
}

func StartRpcServer() {
	s := server.NewServer()
	s.RegisterName("LogicService", new(LogicService), "")
	listener, err := net.Listen("tcp", ":8972")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("RPC server started on :8972")
	s.ServeListener("tcp", listener)
}
