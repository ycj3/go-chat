package connect

import (
	"context"
	"go-chat/internal/redis"
	"go-chat/tools"
	"strconv"

	"github.com/sirupsen/logrus"
)

func (c *Connect) Connect(ctx context.Context, args *ConnectRequest, reply *ConnectReply) (err error) {
	if args == nil {
		logrus.Errorf("logic,connect args empty")
		return
	}
	logrus.Infof("logic,authToken is:%s", args.AuthToken)
	key := tools.GetSessionName(args.AuthToken)
	userInfo, err := redis.RedisClient.HGetAll(key).Result()
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
