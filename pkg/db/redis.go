package db

import (
	"github.com/redis/go-redis/v9"
)

func NewRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "redis:6379", // container name in docker-compose
	})
}
