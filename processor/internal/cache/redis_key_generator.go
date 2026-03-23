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

func NewRedisKeyGenerator(secretKey []byte) KeyGenerator {
	if len(secretKey) != constant.Blake3DigestSize {
		panic(fmt.Sprintf("secretKey must be exactly %d bytes", constant.Blake3DigestSize))
	}
	return &RedisKeyGenerator{secretKey: secretKey}
}

func NewRedisKeyGeneratorFromStringKey(secretKey string) KeyGenerator {
	if len(secretKey) != constant.Blake3DigestSize {
		panic(fmt.Sprintf("secretKey must be exactly %d bytes", constant.Blake3DigestSize))
	}
	secretKeyBytes := []byte(secretKey)
	return &RedisKeyGenerator{secretKey: secretKeyBytes}
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
