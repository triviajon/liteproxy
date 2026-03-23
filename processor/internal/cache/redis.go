package cache

import (
	"context"
	"net/url"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client       *redis.Client
	keyGenerator KeyGenerator
}

func NewRedisCache(addr string, keyGen KeyGenerator) *RedisCache {
	return &RedisCache{
		client:       redis.NewClient(&redis.Options{Addr: addr}),
		keyGenerator: keyGen,
	}
}

func (r *RedisCache) Get(ctx context.Context, url url.URL) ([]byte, error) {
	key := r.keyGenerator.HashURL(url)
	return r.client.Get(ctx, key).Bytes()
}

func (r *RedisCache) Set(ctx context.Context, url url.URL, val []byte, ttl time.Duration) error {
	key := r.keyGenerator.HashURL(url)
	return r.client.Set(ctx, key, val, ttl).Err()
}
