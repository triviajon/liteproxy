package cache

import (
	"encoding/hex"
	"fmt"
	"net/url"
	"strings"

	"github.com/triviajon/liteproxy/processor/internal/constant"
	"lukechampine.com/blake3"
)

type RedisKeyGenerator struct {
	secretKey []byte // Must be constant.Blake3DigestSize bytes
}

// NewRedisKeyGenerator creates a new RedisKeyGenerator.
// Requires that secretKey is exactly constant.Blake3DigestSize bytes.
// Returns a KeyGenerator, otherwise an error describing the constraint violation.
func NewRedisKeyGenerator(secretKey []byte) (KeyGenerator, error) {
	if len(secretKey) != constant.Blake3DigestSize {
		return nil, fmt.Errorf("secretKey must be exactly %d bytes, got %d", constant.Blake3DigestSize, len(secretKey))
	}
	return &RedisKeyGenerator{secretKey: secretKey}, nil
}

// NewRedisKeyGeneratorFromStringKey creates a new RedisKeyGenerator from a string key.
// Requires that secretKey is exactly constant.Blake3DigestSize bytes.
// Returns a KeyGenerator, otherwise an error describing the constraint violation.
func NewRedisKeyGeneratorFromStringKey(secretKey string) (KeyGenerator, error) {
	if len(secretKey) != constant.Blake3DigestSize {
		return nil, fmt.Errorf("secretKey must be exactly %d bytes, got %d", constant.Blake3DigestSize, len(secretKey))
	}
	secretKeyBytes := []byte(secretKey)
	return &RedisKeyGenerator{secretKey: secretKeyBytes}, nil
}

func (kg *RedisKeyGenerator) HashURL(url url.URL) string {
	// Normalize the URL
	url.Fragment = ""                        // Strip the fragment (e.g., #anchor)
	url.Scheme = strings.ToLower(url.Scheme) // Lowercase the scheme
	url.Host = strings.ToLower(url.Host)     // Lowercase the host

	hasher := blake3.New(constant.Blake3DigestSize, kg.secretKey)
	hasher.Write([]byte(url.String()))
	return hex.EncodeToString(hasher.Sum(nil))
}
