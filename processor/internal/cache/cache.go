package cache

import (
	"context"
	"net/url"
	"time"
)

// Cache defines the interface for caching operations.
type Cache interface {
	// Get retrieves the value associated with the given url.
	Get(ctx context.Context, url url.URL) ([]byte, error)
	// Set stores the value with the specified url and time-to-live (TTL).
	Set(ctx context.Context, url url.URL, value []byte, ttl time.Duration) error
}
