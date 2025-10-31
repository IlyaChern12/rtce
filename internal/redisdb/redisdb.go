package redisdb

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)


func RedisConnect(addr string) (*redis.Client, error) {
	redisDB := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	if err := redisDB.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis is not available: %w", err)
	}

	return redisDB, nil
}