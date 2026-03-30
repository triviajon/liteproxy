package cache

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/triviajon/liteproxy/processor/internal/logging"
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
	logging.Infof("Connecting to Redis at %s", addr)
	rc := &RedisCache{
		client:       redis.NewClient(&redis.Options{Addr: addr}),
		keyGenerator: keyGen,
	}
	logging.Debugf("Redis connection established")
	return rc, nil
}

func (r *RedisCache) Get(ctx context.Context, url url.URL) ([]byte, error) {
	key := r.keyGenerator.HashURL(url)
	result, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			logging.Debugf("Cache miss - key=%s url=%s", key, url.String())
		} else {
			logging.Errorf("Cache GET error - key=%s error=%v", key, err)
		}
		return nil, err
	}
	logging.Debugf("Cache hit - key=%s url=%s bytes=%d", key, url.String(), len(result))
	return result, nil
}

func (r *RedisCache) Set(ctx context.Context, url url.URL, val []byte, ttl time.Duration) error {
	key := r.keyGenerator.HashURL(url)
	err := r.client.Set(ctx, key, val, ttl).Err()
	if err != nil {
		logging.Errorf("Cache SET error - key=%s url=%s ttl=%v error=%v", key, url.String(), ttl, err)
	} else {
		logging.Debugf("Cache SET success - key=%s url=%s ttl=%v bytes=%d", key, url.String(), ttl, len(val))
	}
	return err
}
