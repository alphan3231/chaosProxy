package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/elliot/chaosProxy/internal/config"
	"github.com/elliot/chaosProxy/internal/handlers"
	"github.com/elliot/chaosProxy/pkg/infrastructure/redis"
	"github.com/elliot/chaosProxy/pkg/middleware"
)

func main() {
	// Load Configuration
	cfg := config.LoadConfig()

	// Initialize Redis
	redisClient, err := redis.NewClient(cfg.RedisAddr, cfg.RedisPassword)
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

	// Custom Error Handler for Ghost Mode
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("🔥 Backend failed: %v. Attempting Ghost Mode...", err)

		// Try to get ghost response
		ghost, ghostErr := redisClient.GetGhostResponse(r.Context(), r.Method, r.URL.Path)
		if ghostErr == nil && ghost != nil {
			log.Printf("👻 Ghost Mode Activated for: %s %s", r.Method, r.URL.Path)

			// Increment Stats
			redisClient.GetRawClient().Incr(r.Context(), "chaos:stats:ghost_count")

			// Set Ghost Headers
			w.Header().Set("X-Chaos-Ghost", "true")
			w.Header().Set("X-Chaos-Original-Status", fmt.Sprintf("%d", ghost.Status))
			w.Header().Set("Content-Type", "application/json") // Assuming JSON for now

			w.WriteHeader(ghost.Status)
			w.Write([]byte(ghost.ResponseBody))
			return
		}

		log.Printf("💀 Ghost Mode failed (no data found). Returning 502.")
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte(`{"error": "Service temporarily unavailable"}`))
	}

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
	// Order: Recovery -> RateLimit -> Chaos -> Logger -> TrafficLogger -> Mux
	rateLimiter := middleware.NewRateLimiter(100, time.Minute) // 100 requests per minute per IP
	chaosMiddleware := middleware.NewChaosMiddleware(redisClient)
	trafficMiddleware := middleware.TrafficLogger(redisClient)
	handler := middleware.Chain(mux, trafficMiddleware, middleware.Logger, chaosMiddleware.Chaos, middleware.RateLimit(rateLimiter), middleware.Recovery)

	// Start the Server
	log.Printf("👻 Chaos-Proxy Sentinel starting on :%s", cfg.Port)
	log.Printf("🎯 Forwarding to: %s", cfg.TargetURL)

	if err := http.ListenAndServe(":"+cfg.Port, handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
