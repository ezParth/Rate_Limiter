package main

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func redisOptions(addr string, pass string, db int) *redis.Options {
	return &redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       db,
	}
}

func main() {
	rdb := redis.NewClient(redisOptions("localhost:6379", "", 0))

	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		fmt.Println("Error: ", err)
	}

	fmt.Println("PONG: ", pong)
}
