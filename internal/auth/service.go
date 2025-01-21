package auth

import (
	"context"
	"go-chat/internal/rpc"
	"go-chat/tools"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	UserID string
}

type LoginResponse struct {
	Code      int
	AuthToken string
}

func Login(c *gin.Context) {
	user_id := c.Param("user_id")
	req := &LoginRequest{
		UserID: user_id,
	}
	msg := ""
	reply := &LoginResponse{}
	err := rpc.RpcClient.Call(context.Background(), "Login", req, reply)
	if err != nil {
		msg = err.Error()
	}
	code := reply.Code
	authToken := reply.AuthToken

	if code == tools.CodeFail || authToken == "" {
		tools.FailWithMsg(c, msg)
		return
	}
	tools.SuccessWithMsg(c, nil, gin.H{"token": authToken})
}
