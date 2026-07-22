package middleware

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/go-api/internal/shared/ratelimiter"
)

type RateLimitMiddleware struct {
	limiter ratelimiter.RateLimiter
}

// create a new middleware
func NewRateLimitMiddleware(limiter ratelimiter.RateLimiter) *RateLimitMiddleware {
	return &RateLimitMiddleware{
		limiter: limiter,
	}
}

// return the gin middleware
func (m *RateLimitMiddleware) Handler(group string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Extract auth info from context (based on your existing middleware)
		isAuthenticated := false
		userID := ""

		// Adjust these keys to match the auth middleware
		if uid, exists := c.Get("user_id"); exists {
			if id, ok := uid.(string); ok && id != "" {
				isAuthenticated = true
				userID = id
			}
		}

		// 2. Check rate limit
		allowed, retryAfter, limit := m.limiter.Allow(c, isAuthenticated, userID, group)

		// 3. Set standard headers
		c.Header("X-RateLimit-Limit", strconv.Itoa(limit.Requests))

		if !allowed {
			// Retry-After in seconds
			seconds := int(retryAfter.Seconds())
			if seconds <= 0 {
				seconds = 1
			}

			c.Header("Retry-After", strconv.Itoa(seconds))

			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":       "rate_limit_exceeded",
				"message":     "Too many requests. Please try again later.",
				"retry_after": seconds,
			})
			return
		}

		// 4. Continue request
		c.Next()
	}
}
