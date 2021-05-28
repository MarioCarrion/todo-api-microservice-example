package internal

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-redis/redis/v8"

	"github.com/MarioCarrion/todo-api/internal/envvar"
)

// NewRedis instantiates the Redis client using configuration defined in environment variables.
func NewRedis(conf *envvar.Configuration) (*redis.Client, error) {
	host, err := conf.Get("REDIS_HOST")
	if err != nil {
		return nil, fmt.Errorf("conf.Get REDIS_HOST %w", err)
	}

	db, err := conf.Get("REDIS_DB")
	if err != nil {
		return nil, fmt.Errorf("conf.Get REDIS_DB %w", err)
	}

	dbi, _ := strconv.Atoi(db)

	rdb := redis.NewClient(&redis.Options{
		Addr: host,
		DB:   dbi,
	})

	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		return nil, fmt.Errorf("rdb.Ping %w", err)
	}

	return rdb, nil
}
