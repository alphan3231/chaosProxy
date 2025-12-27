package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/elliot/chaosProxy/internal/config"
	"github.com/elliot/chaosProxy/internal/handlers"
	"github.com/elliot/chaosProxy/pkg/infrastructure/redis"
	"github.com/elliot/chaosProxy/pkg/middleware"
)

func main() {
	// Load Configuration
	cfg := config.LoadConfig()

	// Initialize Redis
	redisClient, err := redis.NewClient(cfg.RedisAddr)
	if err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}
	defer redisClient.Close()
	log.Printf("⚡ Connected to Redis at %s", cfg.RedisAddr)

	// Parse the target URL
	target, err := url.Parse(cfg.TargetURL)
	if err != nil {
		log.Fatalf("Failed to parse target URL: %v", err)
	}

	// Create the Reverse Proxy
	proxy := httputil.NewSingleHostReverseProxy(target)

	// Update the Director to set the Host header correctly
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Host = target.Host // Important: Set the Host header to the target's host
	}

	// Define the main handler with the proxy
	proxyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	})

	// Setup Router (Mux) to handle specific paths
	mux := http.NewServeMux()

	// Register Health Check
	mux.HandleFunc("/healthz", handlers.HealthCheck(redisClient))

	// Catch-all to Proxy
	mux.Handle("/", proxyHandler)

	// Wrap handler with Middleware Chain
	// Order: Recovery -> Logger -> TrafficLogger -> Mux
	trafficMiddleware := middleware.TrafficLogger(redisClient)
	handler := middleware.Chain(mux, trafficMiddleware, middleware.Logger, middleware.Recovery)

	// Start the Server
	log.Printf("👻 Chaos-Proxy Sentinel starting on :%s", cfg.Port)
	log.Printf("🎯 Forwarding to: %s", cfg.TargetURL)

	if err := http.ListenAndServe(":"+cfg.Port, handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
