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

	newsChx := rdb.Subscribe(ctx, "news")
	defer newsChx.Close()

	_, err := newsChx.Receive(ctx)
	if err != nil {
		log.Fatal(err)
	}

	ch := newsChx.Channel()
	fmt.Println("Subscribed to news channel")

	timeout := time.After(12 * time.Second)

	for {
		select {
		case msg := <-ch:
			if msg == nil {
				fmt.Println("Channel closed")
				return
			}
			fmt.Println("news:", msg.Payload)

		case <-timeout:
			fmt.Println("Timeout reached, closing...")
			return
		}
	}
}
