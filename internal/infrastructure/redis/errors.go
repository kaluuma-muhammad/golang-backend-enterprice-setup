package redis

import "errors"

var (
	ErrConnectionFailed = errors.New("failed to connect to redis")
)
