package redisManager

import (
	"github.com/garyburd/redigo/redis"

	"os"
)

const POOL_SIZE = 100
const REDIS_CLOUD = "localhost:6379"

var (
	redisPool    *redis.Pool
	redisAddress = os.Getenv("REDIS_URL")
)

func GetClient() redis.Conn {
	return redisPool.Get()
}

func ClosePool() {
	if redisPool != nil {
		redisPool.Close()
	}
}

func init() {
	if redisAddress == "" {
		redisAddress = REDIS_CLOUD
	}
	redisPool = redis.NewPool(func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", redisAddress)
		if err != nil {
			return nil, err
		}
		return c, err
	}, POOL_SIZE)
}
