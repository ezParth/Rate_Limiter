package auth

func IsAuthenticated(token string) bool {
	if token != "" {
		return true
	}
	return false
}
