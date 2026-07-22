package redis

import (
	"fmt"

	"github.com/go-api/internal/shared/config"

	goredis "github.com/redis/go-redis/v9"
)

func NewOptions(cfg config.RedisConfig) *goredis.Options {
	return &goredis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	}
}
