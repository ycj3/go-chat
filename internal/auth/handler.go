package auth

import (
	"context"
	"errors"
	"go-chat/config"
	"go-chat/internal/redis"

	"go-chat/tools"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

type RpcHandler struct {
}

func (h *RpcHandler) Login(ctx context.Context, args *LoginRequest, reply *LoginResponse) (err error) {
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
	token, _ := redis.RedisClient.Get(loginSessionId).Result()
	if token != "" {
		// Logout the already logged-in user session
		oldSession := tools.CreateSessionId(token)
		err := redis.RedisClient.Del(oldSession).Err()
		if err != nil {
			return errors.New("logout user fail!token is:" + token)
		}
	}

	// Store session data in Redis
	redis.RedisClient.Do("MULTI")
	redis.RedisClient.HMSet(sessionId, userData)
	redis.RedisClient.Expire(sessionId, 86400*time.Second)
	redis.RedisClient.Set(loginSessionId, randToken, 86400*time.Second)
	err = redis.RedisClient.Do("EXEC").Err()
	if err != nil {
		logrus.Infof("register set redis token fail!")
		return err
	}

	// Set success response
	reply.Code = config.SuccessReplyCode
	reply.AuthToken = randToken
	return
}
