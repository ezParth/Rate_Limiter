package auth

import (
	"fmt"
	"net/http"
	"strings"
)

// format Authorization: Bearer YOUR_TOKEN_HERE

var demo_token = "YOUR_TOKEN_HERE"

func IsAuthenticated(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		fmt.Println(authHeader)

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

		fmt.Println("Authorized")
		return
	}
}
