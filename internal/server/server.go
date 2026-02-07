package server

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/elliot/chaosProxy/internal/config"
	"github.com/elliot/chaosProxy/internal/handlers"
	"github.com/elliot/chaosProxy/pkg/infrastructure/redis"
	"github.com/elliot/chaosProxy/pkg/middleware"
)

type Server struct {
	cfg         *config.Config
	redisClient *redis.Client
}

func NewServer(cfg *config.Config, redisClient *redis.Client) *Server {
	return &Server{
		cfg:         cfg,
		redisClient: redisClient,
	}
}

func (s *Server) Start() error {
	// Parse the target URL
	target, err := url.Parse(s.cfg.TargetURL)
	if err != nil {
		return err
	}

	// Create the Reverse Proxy
	proxy := httputil.NewSingleHostReverseProxy(target)

	// Set Custom Transport for Retries
	if s.cfg.RetryMax > 0 {
		proxy.Transport = middleware.NewRetryTransport(http.DefaultTransport, s.cfg.RetryMax, s.cfg.RetryDelay)
	}

	s.setupProxy(proxy, target)

	// Setup Router (Mux)
	mux := http.NewServeMux()

	// Register Health Check
	mux.HandleFunc("/healthz", handlers.HealthCheck(s.redisClient))
	// Register Blocked IPs API
	mux.HandleFunc("/api/blocked-ips", handlers.GetBlockedIPs(s.redisClient))

	// Calculate Canary Proxy
	canaryMiddleware := middleware.NewCanary(s.cfg.CanaryURL, s.cfg.CanaryWeight)
	finalProxy := canaryMiddleware(proxy)

	// Catch-all to Proxy
	mux.Handle("/", finalProxy)

	// Setup Middleware Chain
	handler := s.setupMiddleware(mux)

	// Start the Server
	srv := &http.Server{
		Addr:    ":" + s.cfg.Port,
		Handler: handler,
	}

	go func() {
		log.Printf("👻 Chaos-Proxy Sentinel starting on :%s", s.cfg.Port)
		log.Printf("🎯 Forwarding to: %s", s.cfg.TargetURL)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("👋 Server exited gracefully")
	return nil
}

func (s *Server) setupProxy(proxy *httputil.ReverseProxy, target *url.URL) {
	// Custom Error Handler for Ghost Mode
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("🔥 Backend failed: %v. Attempting Ghost Mode...", err)

		// Try to get ghost response
		ghost, ghostErr := s.redisClient.GetGhostResponse(r.Context(), r.Method, r.URL.Path)
		if ghostErr == nil && ghost != nil {
			log.Printf("👻 Ghost Mode Activated for: %s %s", r.Method, r.URL.Path)

			// Increment Stats
			s.redisClient.GetRawClient().Incr(r.Context(), "chaos:stats:ghost_count")

			// Set Ghost Headers
			w.Header().Set("X-Chaos-Ghost", "true")

			// Replay original headers if available
			for k, values := range ghost.ResponseHeaders {
				// Skip sensitive or problematic headers (like Content-Length which will be recalculated)
				if k == "Content-Length" || k == "Connection" || k == "Date" {
					continue
				}
				for _, v := range values {
					w.Header().Add(k, v)
				}
			}

			// If Content-Type wasn't in the recorded headers, try to detect it
			if w.Header().Get("Content-Type") == "" {
				if http.DetectContentType([]byte(ghost.ResponseBody)) != "application/octet-stream" {
					w.Header().Set("Content-Type", http.DetectContentType([]byte(ghost.ResponseBody)))
				} else {
					w.Header().Set("Content-Type", "application/json")
				}
			}

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
}

func (s *Server) setupMiddleware(handler http.Handler) http.Handler {
	// Order: Recovery -> RateLimit -> IPFilter -> Chaos -> Logger -> TrafficLogger -> Mux
	rateLimiter := middleware.NewRateLimiter(100, time.Minute) // 100 requests per minute per IP
	chaosMiddleware := middleware.NewChaosMiddleware(s.redisClient, s.cfg.SimulateRegion)
	trafficMiddleware := middleware.TrafficLogger(s.redisClient)

	return middleware.Chain(
		handler,
		middleware.Recovery, // Recovery should be first (outermost) to catch panics in any middleware
		middleware.RequestID,
		middleware.PoweredBy,
		trafficMiddleware,
		middleware.Logger,
		middleware.SecurityFuzzer(s.cfg.SecurityFuzzingEnabled),
		chaosMiddleware.Chaos,
		middleware.IPFilter(s.redisClient),
		middleware.RateLimit(rateLimiter),
	)
}
