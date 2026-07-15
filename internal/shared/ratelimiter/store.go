package ratelimiter

import "time"

type Entry struct {
	Limiter  *LimiterWrapper
	LastSeen time.Time
}

type Store interface {
	Get(key string) (*Entry, bool)
	Set(key string, entry *Entry)
	Delete(key string)
	Keys() []string
}
