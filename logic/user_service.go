package logic

import (
	"context"
	"errors"
	"strconv"
	"time"

	"go-chat/config"
	"go-chat/proto"
	"go-chat/tools"

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

func (rpc *RpcLogic) Connect(ctx context.Context, args *proto.ConnectRequest, reply *proto.ConnectReply) (err error) {
	if args == nil {
		logrus.Errorf("logic,connect args empty")
		return
	}
	logrus.Infof("logic,authToken is:%s", args.AuthToken)
	key := tools.GetSessionName(args.AuthToken)
	userInfo, err := RedisClient.HGetAll(key).Result()
	if err != nil {
		logrus.Infof("RedisCli HGetAll key :%s , err:%s", key, err.Error())
		return err
	}
	if len(userInfo) == 0 {
		reply.UserId = 0
		return
	}
	reply.UserId, _ = strconv.Atoi(userInfo["userId"])
	logrus.Infof("logic rpc userId:%d", reply.UserId)
	return
}
