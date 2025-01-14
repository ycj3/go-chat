package tools

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis"
)

var RedisClientMap = map[string]*redis.Client{}
var syncLock sync.Mutex

type RedisOption struct {
	Address  string
	Password string
	Db       int
}

func GetRedisInstance(redisOpt RedisOption) *redis.Client {
	address := redisOpt.Address
	db := redisOpt.Db
	password := redisOpt.Password
	addr := fmt.Sprintf("%s", address)
	syncLock.Lock()
	if redisCli, ok := RedisClientMap[addr]; ok {
		return redisCli
	}
	client := redis.NewClient(&redis.Options{
		Addr:       addr,
		Password:   password,
		DB:         db,
		MaxConnAge: 20 * time.Second,
	})
	RedisClientMap[addr] = client
	syncLock.Unlock()
	return RedisClientMap[addr]
}
