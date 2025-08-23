package ratelimiter

import (
	"context"
	"fmt"
	"rl/helper"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	maxLimit  = 10 // max tokens per user
	renewTime = 60 // 1 token per 60 seconds
)

func CreateClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	return client
}

func RateLimiter(username string, rdb *redis.Client, mutexClinet *sync.Mutex) error {
	ctx := context.Background()

	tokenKey := fmt.Sprintf("%s:tokens", username)
	timeKey := fmt.Sprintf("%s:lastRefill", username)

	mutexClinet.Lock()
	tokenStr, err := rdb.Get(ctx, tokenKey).Result()
	mutexClinet.Unlock()
	tokens := maxLimit
	if err == nil {
		fmt.Sscanf(tokenStr, "%d", &tokens)
	}

	mutexClinet.Lock()
	lastRefillStr, err := rdb.Get(ctx, timeKey).Result()
	mutexClinet.Unlock()
	lastRefill := time.Now()
	if err == nil {
		parsed, perr := time.Parse(time.RFC3339Nano, lastRefillStr)
		if perr == nil {
			lastRefill = parsed
		}
	}

	now := time.Now()
	diff := now.Sub(lastRefill).Seconds()
	addedTokens := int(diff) / renewTime

	if addedTokens > 0 {
		tokens = helper.Min(maxLimit, tokens+addedTokens)
		lastRefill = now
	}

	// fmt.Println("Available tokens before request:", tokens)

	if tokens <= 0 {
		return fmt.Errorf("rate Limit Exceeded")
	}

	tokens--

	mutexClinet.Lock()
	rdb.Set(ctx, tokenKey, tokens, 1000*time.Second)
	rdb.Set(ctx, timeKey, lastRefill.Format(time.RFC3339Nano), 1000*time.Second)
	mutexClinet.Unlock()

	fmt.Println("Remaining tokens:", tokens)
	return nil
}
