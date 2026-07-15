package ratelimiter

import (
	"time"

	"golang.org/x/time/rate"
)

type LimiterWrapper struct {
	limiter *rate.Limiter
	limit   Limit
}

func NewLimiterWrapper(limit Limit) *LimiterWrapper {
	// Convert Requests/Window → rate.Limit (events per second)
	r := rate.Limit(float64(limit.Requests) / limit.Window.Seconds())

	return &LimiterWrapper{
		limiter: rate.NewLimiter(r, limit.Burst),
		limit:   limit,
	}
}

func (l *LimiterWrapper) Allow() bool {
	return l.limiter.Allow()
}
func (l *LimiterWrapper) Reserve() *rate.Reservation {
	return l.limiter.Reserve()
}

func (l *LimiterWrapper) Limit() Limit {
	return l.limit
}
func (l *LimiterWrapper) RemainingTokens() int {
	res := l.limiter.Reserve()
	if !res.OK() {
		return 0
	}

	delay := res.Delay()
	res.Cancel() // very important to not consume token

	if delay > 0 {
		return 0
	}

	return l.limit.Burst
}

func (l *LimiterWrapper) RetryAfter() time.Duration {
	res := l.limiter.Reserve()
	if !res.OK() {
		return time.Second
	}

	delay := res.Delay()
	res.Cancel()

	return delay
}
