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

	// Precondition validation
	if proxySecret == "" {
		log.Fatal("FATAL: PROXY_AUTH_TOKEN is required for middleware")
	}

	// Initialize the Hashing Identity
	keyGen, err := cache.NewRedisKeyGeneratorFromStringKey(cacheSalt)
	if err != nil {
		log.Fatalf("FATAL: Failed to create key generator: %v", err)
	}

	// Initialize the Storage Adapter
	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)
	redisStore, err := cache.NewRedisCache(redisAddr, keyGen)
	if err != nil {
		log.Fatalf("FATAL: Failed to create Redis cache: %v", err)
	}

	// Build the Logic Pipeline
	pipeline, err := rewritepipeline.NewPipeline(
		&rewritepipeline.ImageStripper{},
	)
	if err != nil {
		log.Fatalf("FATAL: Failed to create pipeline: %v", err)
	}

	// Initialize the Core Proxy Server
	proxySrv := &proxy.ProxyServer{
		Pipeline: *pipeline,
		Cache:    redisStore,
	}

	// Chain the Middlewares
	handler, err := auth.WithHeaderAuth(proxySrv, proxySecret)
	if err != nil {
		log.Fatalf("FATAL: Failed to create auth middleware: %v", err)
	}

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
