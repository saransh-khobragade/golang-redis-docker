package cache

import (
	"time"

	"github.com/go-redis/redis"
)

type redisCache struct {
	host string
	db   int
	exp  time.Duration
}

func NewRedisCache(host string, db int, exp time.Duration) *redisCache {
	return &redisCache{
		host: host,
		db:   db,
		exp:  exp,
	}
}

func (cache redisCache) getClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cache.host,
		Password: "",
		DB:       cache.db,
	})
}
