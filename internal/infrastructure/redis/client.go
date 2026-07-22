package redis

import (
	"context"
	"fmt"

	goredis "github.com/redis/go-redis/v9"

	"github.com/go-api/internal/shared/config"
)

type Client struct {
	client *goredis.Client
}

func New(cfg config.RedisConfig) (*Client, error) {
	client := goredis.NewClient(NewOptions(cfg))

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrConnectionFailed, err)
	}

	return &Client{
		client: client,
	}, nil
}

func (c *Client) Client() *goredis.Client {
	return c.client
}

func (c *Client) Ping(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}

func (c *Client) Close() error {
	return c.client.Close()
}
