package storage

import (
	"context"
	"errors"
	"time"
)

var ErrCacheMiss = errors.New("cache miss")

type Store interface {
	Ping(ctx context.Context) error
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
	Keys(ctx context.Context, pattern string) ([]string, error)
	Close() error
}
