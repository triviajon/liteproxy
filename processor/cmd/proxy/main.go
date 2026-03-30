package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/triviajon/liteproxy/processor/internal/auth"
	"github.com/triviajon/liteproxy/processor/internal/cache"
	"github.com/triviajon/liteproxy/processor/internal/constant"
	"github.com/triviajon/liteproxy/processor/internal/logging"
	"github.com/triviajon/liteproxy/processor/internal/proxy"
	"github.com/triviajon/liteproxy/processor/internal/rewritepipeline"
)

func main() {
	if err := logging.Init(); err != nil {
		panic(err)
	}

	logging.ConfigureFromEnv()
	logging.Infof("LiteProxy Processor initializing...")

	proxySecret := os.Getenv("PROXY_AUTH_TOKEN")
	cacheSalt := os.Getenv("CACHE_SALT")
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	listenPort := os.Getenv("PROCESSOR_PORT")

	logging.Infof("Loading configuration - redis_host=%s redis_port=%s listen_port=%s", redisHost, redisPort, listenPort)

	// Precondition validation
	if proxySecret == "" {
		logging.Fatalf("PROXY_AUTH_TOKEN is required for middleware")
	}
	logging.Infof("PROXY_AUTH_TOKEN configured")

	// Initialize the Hashing Identity
	logging.Infof("Creating Redis key generator...")
	keyGen, err := cache.NewRedisKeyGeneratorFromStringKey(cacheSalt)
	if err != nil {
		logging.Fatalf("Failed to create Redis key generator: %v", err)
	}
	logging.Infof("Redis key generator created")

	// Initialize the Storage Adapter
	logging.Infof("Creating Redis cache adapter...")
	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)
	redisStore, err := cache.NewRedisCache(redisAddr, keyGen)
	if err != nil {
		logging.Fatalf("Failed to create Redis cache: %v", err)
	}

	// Build the Logic Pipeline
	logging.Infof("Building request processing pipeline...")
	pipeline, err := rewritepipeline.NewPipeline(
		&rewritepipeline.ImageStripper{},
	)
	if err != nil {
		logging.Fatalf("Failed to create pipeline: %v", err)
	}

	// Initialize the Core Proxy Server
	logging.Infof("Initializing proxy server...")
	proxySrv := &proxy.ProxyServer{
		Pipeline: *pipeline,
		Cache:    redisStore,
	}

	// Chain the Middlewares
	logging.Infof("Setting up authentication middleware...")
	handler, err := auth.WithHeaderAuth(proxySrv, proxySecret)
	if err != nil {
		logging.Fatalf("Failed to create auth middleware: %v", err)
	}

	// Lifecycle Management
	logging.Infof("Creating HTTP server - addr=:%s read_timeout=%v write_timeout=%v idle_timeout=%v",
		listenPort, constant.DefaultReadTimeout, constant.DefaultWriteTimeout, constant.DefaultIdleTimeout)
	server := &http.Server{
		Addr:         ":" + listenPort,
		Handler:      handler,
		ReadTimeout:  constant.DefaultReadTimeout,
		WriteTimeout: constant.DefaultWriteTimeout,
		IdleTimeout:  constant.DefaultIdleTimeout,
	}

	logging.Infof("LiteProxy Processor started and listening on port %s", listenPort)

	if err := server.ListenAndServe(); err != nil {
		logging.Fatalf("Server failed: %s", err)
	}
}
