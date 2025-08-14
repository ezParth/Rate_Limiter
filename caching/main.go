package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type WeatherData struct {
	Temp     string `json:"temperature"`
	Wind     string `json:"wind"`
	Forecast []struct {
		Day  string `json:"day"`
		Temp string `json:"temperature"`
	} `json:"forecast"`
}

func fetchWeather(city string) (WeatherData, error) {
	// Example API: GoWeather
	url := fmt.Sprintf("https://goweather.herokuapp.com/weather/%s", city)

	resp, err := http.Get(url)
	if err != nil {
		return WeatherData{}, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var data WeatherData
	err = json.Unmarshal(body, &data)
	if err != nil {
		return WeatherData{}, err
	}

	return data, nil
}

func main() {
	// Connect to Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	city := "london"

	// 1. Try to get from cache
	cached, err := rdb.Get(ctx, city).Result()
	if err == nil {
		fmt.Println("✅ From Cache:", cached)
		return
	}

	// 2. Cache miss → fetch from API
	fmt.Println("📡 Fetching from API...")
	weather, err := fetchWeather(city)
	if err != nil {
		panic(err)
	}

	// Convert to JSON for storing in Redis
	jsonData, _ := json.Marshal(weather)

	// Store in Redis with 30s TTL
	err = rdb.Set(ctx, city, jsonData, 30*time.Second).Err()
	if err != nil {
		panic(err)
	}

	fmt.Println("🌤 From API:", string(jsonData))
}
