package ratelimiter

import (
	"context"
	"fmt"
	helper "rl/helper"
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

func RateLimiter(username string, rdb *redis.Client) error {
	fmt.Println("Hitting RateLimiter")
	timeUsername := fmt.Sprintf("%s:time", username)
	ctx := context.Background()
	count, err := rdb.Incr(ctx, username).Result()
	if err != nil {
		panic(err)
	}

	remainingToken := maxLimit - helper.Max(5, int(count))

	currTime := time.Now()
	if count == 1 {
		value, err := rdb.Set(ctx, timeUsername, currTime, time.Duration(1000*time.Second)).Result()
		fmt.Println("value ", value)
		return err
	}

	timeStr, err := rdb.Get(ctx, timeUsername).Result()
	if err != nil {
		return err
	}

	savedTime, err := time.Parse(time.RFC3339Nano, timeStr)
	if err != nil {
		return err
	}

	diff := currTime.Sub(savedTime)
	fmt.Println("diff: ", diff)
	seconds := diff.Seconds()
	fmt.Println("seconds: ", seconds)
	extraToken := 0
	if seconds > 59 {
		extraToken = helper.Min(maxLimit, int(seconds)/renewTime)
	}

	remainingToken += extraToken

	fmt.Println("remaining token", remainingToken)
	if remainingToken > 0 {
		rdb.Set(ctx, username, remainingToken, 1000*time.Second)
		return nil
	}

	return fmt.Errorf("rate Limit Exceeded")
}
