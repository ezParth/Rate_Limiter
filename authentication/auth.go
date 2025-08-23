package authentication

import (
	"fmt"
	"net/http"
	rateLimiter "rl/ratelimiter"
	"strings"

	"github.com/redis/go-redis/v9"
)

// format Authorization: Bearer YOUR_TOKEN_HERE

var demo_token = "YOUR_TOKEN_HERE"

func IsAuthenticated(rdb *redis.Client, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("hitting auth")
		authHeader := r.Header.Get("Authorization")
		fmt.Println("auth header", authHeader, " url: ", r.URL)

		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		token := parts[1]
		if token != demo_token {
			http.Error(w, "Token is incorrect", http.StatusBadRequest)
		}

		bodyHeader := r.Header.Get("Username")
		fmt.Println(bodyHeader)

		if strings.TrimSpace(bodyHeader) == "" {
			http.Error(w, "Invalid Username", http.StatusBadRequest)
			return
		}

		err := rateLimiter.RateLimiter(bodyHeader, rdb)
		if err != nil {
			http.Error(w, fmt.Sprintf("%s", err), http.StatusBadRequest)
			return
		}

		fmt.Println("Authorized")
		next(w, r)
	}
}
