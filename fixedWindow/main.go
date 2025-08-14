package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func main() {
	// Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	userID := "user_123"
	limit := 5                 // max requests
	window := 10 * time.Second // time window

	for i := 1; i <= 15; i++ {
		allowed, remaining := allowRequest(rdb, userID, limit, window)
		if allowed {
			fmt.Printf("✅ Request %d allowed (%d remaining)\n", i, remaining)
		} else {
			fmt.Printf("❌ Request %d blocked — Rate limit exceeded\n", i)
		}
		time.Sleep(1 * time.Second)
	}
}

func allowRequest(rdb *redis.Client, userID string, limit int, window time.Duration) (bool, int) {
	key := fmt.Sprintf("rate_limit:%s", userID)

	// Increment request count
	count, err := rdb.Incr(ctx, key).Result()
	if err != nil {
		log.Fatal(err)
	}

	// If it's the first request, set the expiration
	if count == 1 {
		rdb.Expire(ctx, key, window)
	}

	// Allow if under limit
	if count <= int64(limit) {
		return true, limit - int(count)
	}
	return false, 0
}
