package redis

import (
	"context"
	"time"

	"github.com/go-api/internal/application/common"
)

var _ common.Cache = (*Client)(nil)

func (c *Client) Get(ctx context.Context, key string) ([]byte, error) {
	return c.client.Get(ctx, key).Bytes()
}

func (c *Client) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return c.client.Set(ctx, key, value, ttl).Err()
}

func (c *Client) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

func (c *Client) Exists(ctx context.Context, key string) (bool, error) {
	count, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (c *Client) Increment(ctx context.Context, key string, ttl time.Duration) (int64, error) {
	pipe := c.client.TxPipeline()

	incr := pipe.Incr(ctx, key)

	pipe.Expire(ctx, key, ttl)

	_, err := pipe.Exec(ctx)

	if err != nil {
		return 0, err
	}

	return incr.Val(), nil
}
