package redis

import (
	"context"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
)

var dbr *redis.Client

func Connent() *redis.Client {
	if dbr != nil {
		return dbr
	}

	db, _ := strconv.Atoi(os.Getenv("REDIS_DB"))

	option := &redis.Options{
		Addr:         os.Getenv("REDIS_ADDR"),
		Password:     os.Getenv("REDIS_PASS"),
		DB:           db,
		PoolSize:     32,
		MinIdleConns: 8,
	}

	dbr = redis.NewClient(option)

	pong, err := dbr.Ping(context.Background()).Result()
	if err != nil || pong != "PONG" {
		panic(err)
	}

	return dbr
}
