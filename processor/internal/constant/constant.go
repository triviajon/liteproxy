package constant

import "time"

const (
	// DefaultCacheTTL is the default time-to-live for cached items.
	DefaultCacheTTL = 1 * time.Hour

	// Blake3DigestSize is the size of the BLAKE3 hash digest in bytes.
	Blake3DigestSize = 32
)
