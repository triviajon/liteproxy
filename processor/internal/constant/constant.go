package constant

import "time"

const (
	// DefaultCacheTTL is the default time-to-live for cached items.
	DefaultCacheTTL = 1 * time.Hour

	// DefaultReadTimeout is the default maximum duration for reading the entire request, including the body.
	DefaultReadTimeout = 15 * time.Second
	// DefaultWriteTimeout is the default maximum duration before timing out writes of the response.
	DefaultWriteTimeout = 30 * time.Second
	// DefaultIdleTimeout is the default maximum amount of time to wait for the next request when keep-alives are enabled.
	DefaultIdleTimeout = 120 * time.Second

	// Blake3DigestSize is the size of the BLAKE3 hash digest in bytes.
	Blake3DigestSize = 32
)
