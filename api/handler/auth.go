package handler

import (
	"go-chat/api/rpc"
	"go-chat/proto"
	"go-chat/tools"

	"github.com/gin-gonic/gin"
)

func Auth(c *gin.Context) {
	user_id := c.Param("user_id")
	req := &proto.LoginRequest{
		UserID: user_id,
	}
	code, authToken, msg := rpc.RpcLogicObj.Login(req)
	if code == tools.CodeFail || authToken == "" {
		tools.FailWithMsg(c, msg)
		return
	}
	tools.SuccessWithMsg(c, nil, gin.H{"token": authToken})
}
