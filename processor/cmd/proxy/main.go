package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/triviajon/liteproxy/processor/internal/auth"
	"github.com/triviajon/liteproxy/processor/internal/cache"
	"github.com/triviajon/liteproxy/processor/internal/constant"
	"github.com/triviajon/liteproxy/processor/internal/proxy"
	"github.com/triviajon/liteproxy/processor/internal/rewritepipeline"
)

func main() {
	proxySecret := os.Getenv("PROXY_AUTH_TOKEN")
	cacheSalt := os.Getenv("CACHE_SALT")
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	listenPort := os.Getenv("PROCESSOR_PORT")

	// Validation: BLAKE3 requires exactly 32 bytes for keyed hashing
	if len(cacheSalt) != 32 {
		log.Fatalf("FATAL: CACHE_SALT must be exactly 32 bytes. Got %d", len(cacheSalt))
	}
	if proxySecret == "" {
		log.Fatal("FATAL: PROXY_AUTH_TOKEN is required for middleware")
	}

	// Initialize the Hashing Identity
	keyGen := cache.NewRedisKeyGeneratorFromStringKey(cacheSalt)

	// Initialize the Storage Adapter
	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)
	redisStore := cache.NewRedisCache(redisAddr, keyGen)

	// Build the Logic Pipeline
	pipeline := rewritepipeline.NewPipeline(
		&rewritepipeline.ImageStripper{},
	)

	// Initialize the Core Proxy Server
	proxySrv := &proxy.ProxyServer{
		Pipeline: *pipeline,
		Cache:    redisStore,
	}

	// Chain the Middlewares
	handler := auth.WithHeaderAuth(proxySrv, proxySecret)

	// Lifecycle Management
	server := &http.Server{
		Addr:         ":" + listenPort,
		Handler:      handler,
		ReadTimeout:  constant.DefaultReadTimeout,
		WriteTimeout: constant.DefaultWriteTimeout,
		IdleTimeout:  constant.DefaultIdleTimeout,
	}

	log.Printf("LiteProxy Processor started on port %s", listenPort)

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("FATAL: Server failed: %s", err)
	}
}
