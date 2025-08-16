package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Redis server
	})

	defer rdb.Close()

	for i := 0; i < 5; i++ {
		msg := fmt.Sprintf("This is a message from news number %d", i+1)
		err := rdb.Publish(ctx, "news", msg).Err()
		if err != nil {
			log.Fatal("Error is publishing the message -> ", err)
		}
		fmt.Println("Message published: ", msg)
		time.Sleep(time.Second * 2)
	}
}
