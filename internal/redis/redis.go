package redis

import (
	"go-chat/config"
	"go-chat/tools"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

var RedisClient *redis.Client

func InitRedisClient() error {
	logrus.Info("Initializing Publish Redis Client")
	redisOpt := tools.RedisOption{
		Address:  config.Conf.Common.CommonRedis.RedisAddress,
		Password: config.Conf.Common.CommonRedis.RedisPassword,
		Db:       config.Conf.Common.CommonRedis.Db,
	}
	client := tools.GetRedisInstance(redisOpt)
	if pong, err := client.Ping().Result(); err != nil {
		logrus.Infof("RedisCli Ping Result pong: %s,  err: %s", pong, err)
		return err
	} else {
		logrus.Info("Publish Redis Client initialized successfully")
	}
	RedisClient = client
	return nil
}
