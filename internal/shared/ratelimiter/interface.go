package ratelimiter

import (
	"time"

	"github.com/gin-gonic/gin"
)

type RateLimiter interface {
	Allow(
		c *gin.Context,
		isAuthenticated bool,
		userID string,
		group string,
	) (bool, time.Duration, Limit)
}
