package ratelimiter

import (
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type Service struct {
	config       *Config
	store        Store
	keyGenerator *KeyGenerator
	mu           sync.Mutex
}

func NewService(config *Config, store Store) *Service {
	return &Service{
		config:       config,
		store:        store,
		keyGenerator: NewKeyGenerator(),
	}
}

var _ RateLimiter = (*Service)(nil)

func (s *Service) Allow(c *gin.Context, isAuthenticated bool, userID string, group string) (bool, time.Duration, Limit) {
	// 1. Determine applicable limit
	limit := s.resolveLimit(c, group)

	// 2. Generate key
	key := s.keyGenerator.GenerateKey(c, isAuthenticated, userID)

	// 3. Get or create limiter
	entry := s.getOrCreateLimiter(key, limit)

	// 4. Update last seen
	entry.LastSeen = time.Now()

	// 5. Check allowance
	if entry.Limiter.Allow() {
		return true, 0, limit
	}

	// 6. Calculate retry-after
	retry := entry.Limiter.RetryAfter()

	return false, retry, limit
}

func (s *Service) resolveLimit(c *gin.Context, group string) Limit {
	method := c.Request.Method
	path := normalizePath(c.FullPath())

	routeKey := method + ":" + path

	// 1. Route-specific limit (highest priority)
	if limit, ok := s.config.Routes[routeKey]; ok {
		return limit
	}

	// 2. Group-based limit
	switch strings.ToLower(group) {
	case "public":
		return s.config.Public
	case "protected":
		return s.config.Protected
	case "verified":
		return s.config.Verified
	}

	// 3. Default fallback
	return s.config.Default
}

func (s *Service) getOrCreateLimiter(key string, limit Limit) *Entry {
	// Try fast path (read)
	if entry, ok := s.store.Get(key); ok {
		return entry
	}

	// Slow path (create)
	s.mu.Lock()
	defer s.mu.Unlock()

	// Double check after lock
	if entry, ok := s.store.Get(key); ok {
		return entry
	}

	limiter := NewLimiterWrapper(limit)

	entry := &Entry{
		Limiter:  limiter,
		LastSeen: time.Now(),
	}

	s.store.Set(key, entry)

	return entry
}
