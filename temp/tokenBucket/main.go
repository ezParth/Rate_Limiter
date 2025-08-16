package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	capacity   = 5
	refillTime = 10 // seconds for 1 token
)

func main() {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer rdb.Close()

	userID := "user:123"

	for i := 0; i < 15; i++ {
		accepted, remaining := tokenBucket(ctx, rdb, userID)
		if accepted {
			fmt.Println("✅ Accepted, Remaining:", remaining)
		} else {
			if remaining == -1 {
				fmt.Println("Currently working on one request")
			} else {
				fmt.Println("❌ Rejected, Remaining:", remaining)
			}
		}
		time.Sleep(2 * time.Second)
	}
}

func tokenBucket(ctx context.Context, rdb *redis.Client, userID string) (bool, int) {
	// Keys in Redis
	tokensKey := fmt.Sprintf("%s:tokens", userID)
	timeKey := fmt.Sprintf("%s:last_refill", userID)
	working := fmt.Sprintf("%s:working", userID)

	// Get current values
	workingStr, _ := rdb.Get(ctx, working).Result()
	tokensStr, _ := rdb.Get(ctx, tokensKey).Result()
	lastRefillStr, _ := rdb.Get(ctx, timeKey).Result()

	var tokens int
	var lastRefill time.Time

	if workingStr == "1" {
		return false, -1
	} else {
		rdb.Set(ctx, working, "1", 0)
	}

	if tokensStr == "" { // first request
		tokens = capacity
		lastRefill = time.Now()
	} else {
		tokens, _ = strconv.Atoi(tokensStr)
		lastRefillInt, _ := strconv.ParseInt(lastRefillStr, 10, 64)
		lastRefill = time.Unix(lastRefillInt, 0)
	}

	// Refill logic
	elapsed := time.Since(lastRefill).Seconds()
	if elapsed >= refillTime {
		addTokens := int(elapsed) / refillTime
		tokens = min(capacity, tokens+addTokens)
		lastRefill = time.Now()
	}

	// Decide
	if tokens > 0 {
		tokens--
		// Save state
		rdb.Set(ctx, tokensKey, tokens, 0)
		rdb.Set(ctx, timeKey, lastRefill.Unix(), 0)
		rdb.Set(ctx, working, "0", 0)
		return true, tokens
	}

	// Save state even if denied
	rdb.Set(ctx, tokensKey, tokens, 0)
	rdb.Set(ctx, timeKey, lastRefill.Unix(), 0)
	rdb.Set(ctx, working, "0", 0)
	return false, tokens
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
