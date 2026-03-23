package cache

import (
	"net/url"
)

// A KeyGenerator is responsible for generating cache keys based on URLs.
type KeyGenerator interface {
	HashURL(url.URL) string
}
