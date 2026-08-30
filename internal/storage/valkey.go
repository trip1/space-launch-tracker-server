package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"ds9labs.com/space-launch-server/internal/config"
)

type ValkeyStore struct {
	client *redis.Client
}

func NewValkeyStore(cfg config.ValkeyConfig) (*ValkeyStore, error) {
	client := redis.NewClient(&redis.Options{
		Addr:            cfg.Address,
		Password:        cfg.Password,
		DB:              cfg.DB,
		DialTimeout:     cfg.DialTimeout,
		ReadTimeout:     cfg.ReadTimeout,
		WriteTimeout:    cfg.WriteTimeout,
		PoolSize:        cfg.PoolSize,
		MinIdleConns:    cfg.MinIdleConns,
		MaxRetries:      cfg.MaxRetries,
		PoolTimeout:     cfg.PoolTimeout,
		ConnMaxIdleTime: cfg.ConnMaxIdleTime,
		ConnMaxLifetime: cfg.ConnMaxLifetime,
	})

	ctx, cancel := context.WithTimeout(context.Background(), cfg.DialTimeout)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("ping valkey: %w", err)
	}

	return &ValkeyStore{client: client}, nil
}

func (s *ValkeyStore) Ping(ctx context.Context) error {
	if err := s.client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("ping valkey: %w", err)
	}

	return nil
}

func (s *ValkeyStore) Get(ctx context.Context, key string) ([]byte, error) {
	value, err := s.client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ErrCacheMiss
		}

		return nil, fmt.Errorf("get key %q: %w", key, err)
	}

	return value, nil
}

func (s *ValkeyStore) Set(ctx context.Context, key string, value []byte, expiration time.Duration) error {
	if err := s.client.Set(ctx, key, value, expiration).Err(); err != nil {
		return fmt.Errorf("set key %q: %w", key, err)
	}

	return nil
}

func (s *ValkeyStore) Delete(ctx context.Context, key string) error {
	if err := s.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("delete key %q: %w", key, err)
	}
	return nil
}

func (s *ValkeyStore) Keys(ctx context.Context, pattern string) ([]string, error) {
	var keys []string
	var cursor uint64
	for {
		page, next, err := s.client.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return nil, fmt.Errorf("scan keys %q: %w", pattern, err)
		}
		keys = append(keys, page...)
		cursor = next
		if cursor == 0 {
			return keys, nil
		}
	}
}

func (s *ValkeyStore) Close() error {
	return s.client.Close()
}
