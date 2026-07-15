package ratelimiter

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

type KeyGenerator struct{}

func NewKeyGenerator() *KeyGenerator {
	return &KeyGenerator{}
}
func (k *KeyGenerator) GenerateKey(c *gin.Context, isAuthenticated bool, userID string) string {
	method := c.Request.Method
	path := normalizePath(c.FullPath())

	// Route-specific key (highest priority)
	routeKey := fmt.Sprintf("%s:%s", method, path)

	// If authenticated → use user ID
	if isAuthenticated && userID != "" {
		return fmt.Sprintf("%s:user:%s", routeKey, userID)
	}

	// Otherwise → fallback to IP
	ip := getClientIP(c)
	return fmt.Sprintf("%s:ip:%s", routeKey, ip)
}

func (k *KeyGenerator) GenerateGroupKey(c *gin.Context, group string, isAuthenticated bool, userID string) string {
	if isAuthenticated && userID != "" {
		return fmt.Sprintf("%s:user:%s", group, userID)
	}

	ip := getClientIP(c)
	return fmt.Sprintf("%s:ip:%s", group, ip)
}

func normalizePath(path string) string {
	if path == "" {
		return "unknown"
	}
	path = strings.TrimSuffix(path, "/")
	return path
}

func getClientIP(c *gin.Context) string {
	// If behind proxy (important for production)
	if ip := c.GetHeader("X-Forwarded-For"); ip != "" {
		parts := strings.Split(ip, ",")
		return strings.TrimSpace(parts[0])
	}

	if ip := c.GetHeader("X-Real-IP"); ip != "" {
		return ip
	}

	return c.ClientIP()
}
