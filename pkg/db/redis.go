package db

import (
	"os"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient() *redis.Client {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "redis:6379"
	}

	return redis.NewClient(&redis.Options{
		Addr: addr,
	})
}
