package cache

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client       *redis.Client
	keyGenerator KeyGenerator
}

// NewRedisCache creates a new RedisCache.
// Requires that addr is not empty and keyGen is not nil.
// Returns a RedisCache, otherwise an error describing which constraint was violated.
func NewRedisCache(addr string, keyGen KeyGenerator) (*RedisCache, error) {
	if addr == "" {
		return nil, fmt.Errorf("addr must not be empty")
	}
	if keyGen == nil {
		return nil, fmt.Errorf("keyGen must not be nil")
	}
	return &RedisCache{
		client:       redis.NewClient(&redis.Options{Addr: addr}),
		keyGenerator: keyGen,
	}, nil
}

func (r *RedisCache) Get(ctx context.Context, url url.URL) ([]byte, error) {
	key := r.keyGenerator.HashURL(url)
	return r.client.Get(ctx, key).Bytes()
}

func (r *RedisCache) Set(ctx context.Context, url url.URL, val []byte, ttl time.Duration) error {
	key := r.keyGenerator.HashURL(url)
	return r.client.Set(ctx, key, val, ttl).Err()
}
