package logic

import (
	"context"
	"errors"
	"strconv"
	"time"

	"go-chat/server/config"
	"go-chat/server/proto"
	"go-chat/server/tools"

	"github.com/sirupsen/logrus"
)

type RpcLogic struct {
}

func (rpc *RpcLogic) Login(ctx context.Context, args *proto.LoginRequest, reply *proto.LoginResponse) (err error) {
	// Initialize the response code to fail
	reply.Code = config.FailReplyCode

	// Convert UserID to integer
	userID, err := strconv.Atoi(args.UserID)
	if err != nil {
		return errors.New("invalid user ID")
	}

	// Generate session ID based on user ID
	loginSessionId := tools.GetSessionIdByUserId(userID)

	// Generate a random token and create a session ID
	randToken := tools.GetRandomToken(32)
	sessionId := tools.CreateSessionId(randToken)

	// Prepare user data to store in Redis
	userData := make(map[string]interface{})
	userData["userId"] = args.UserID

	// Check if the user is already logged in
	token, _ := RedisSessClient.Get(loginSessionId).Result()
	if token != "" {
		// Logout the already logged-in user session
		oldSession := tools.CreateSessionId(token)
		err := RedisSessClient.Del(oldSession).Err()
		if err != nil {
			return errors.New("logout user fail!token is:" + token)
		}
	}

	// Store session data in Redis
	RedisSessClient.Do("MULTI")
	RedisSessClient.HMSet(sessionId, userData)
	RedisSessClient.Expire(sessionId, 86400*time.Second)
	RedisSessClient.Set(loginSessionId, randToken, 86400*time.Second)
	err = RedisSessClient.Do("EXEC").Err()
	if err != nil {
		logrus.Infof("register set redis token fail!")
		return err
	}

	// Set success response
	reply.Code = config.SuccessReplyCode
	reply.AuthToken = randToken
	return
}
