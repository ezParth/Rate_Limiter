package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer rdb.Close()

	players := []string{"alice", "bob", "charlie", "dave", "eve"}
	rand.Seed(time.Now().UnixNano())

	// Start pipeline
	pipe := rdb.Pipeline()

	// Simulate batch score updates
	for _, player := range players {
		score := rand.Intn(1000) // random score
		pipe.ZAdd(ctx, "game_leaderboard", redis.Z{
			Score:  float64(score),
			Member: player,
		})
	}

	// Execute all updates in one go
	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Fatal("Pipeline execution failed:", err)
	}

	// Get top 3 players
	topPlayers, err := rdb.ZRevRangeWithScores(ctx, "game_leaderboard", 0, 6).Result()
	if err != nil {
		log.Fatal("Failed to fetch leaderboard:", err)
	}

	fmt.Println("🏆 Top Players:")
	for _, p := range topPlayers {
		fmt.Printf("%s: %.0f\n", p.Member, p.Score)
	}
}
