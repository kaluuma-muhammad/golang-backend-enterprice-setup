package ratelimiter

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-api/internal/application/common"
)

type RedisService struct {
	config       *Config
	cache        common.Cache
	keyGenerator *KeyGenerator
}

func NewRedisService(config *Config, cache common.Cache) *RedisService {
	return &RedisService{
		config:       config,
		cache:        cache,
		keyGenerator: NewKeyGenerator(),
	}
}

var _ RateLimiter = (*RedisService)(nil)

func (s *RedisService) Allow(c *gin.Context, isAuthenticated bool, userID string, group string) (bool, time.Duration, Limit) {

	limit := s.resolveLimit(c, group)

	key := s.keyGenerator.GenerateKey(
		c,
		isAuthenticated,
		userID,
	)

	count, err := s.cache.Increment(
		c.Request.Context(),
		key,
		limit.Window,
	)

	if err != nil {
		return false, time.Second, limit
	}

	if count > int64(limit.Requests) {
		return false, limit.Window, limit
	}

	return true, 0, limit
}

func (s *RedisService) resolveLimit(c *gin.Context, group string) Limit {
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
