package ratelimiter

import "context"

type KeyParts struct {
	UserId string
	Tier   string
	Route  string
}

type Decision struct {
	Allowed    bool
	RetryAfter int64
	Remaining  int64
	Limit      int64
}

type Limiter interface {
	Allow(ctx context.Context, k KeyParts) (Decision, error)
}
