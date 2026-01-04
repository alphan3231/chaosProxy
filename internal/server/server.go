package server

import (
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
	log.Printf("👻 Chaos-Proxy Sentinel starting on :%s", s.cfg.Port)
	log.Printf("🎯 Forwarding to: %s", s.cfg.TargetURL)

	return http.ListenAndServe(":"+s.cfg.Port, handler)
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
			// w.Header().Set("X-Chaos-Original-Status", fmt.Sprintf("%d", ghost.Status)) // Original status might be useful
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
}

func (s *Server) setupMiddleware(handler http.Handler) http.Handler {
	// Order: Recovery -> RateLimit -> IPFilter -> Chaos -> Logger -> TrafficLogger -> Mux
	rateLimiter := middleware.NewRateLimiter(100, time.Minute) // 100 requests per minute per IP
	chaosMiddleware := middleware.NewChaosMiddleware(s.redisClient)
	trafficMiddleware := middleware.TrafficLogger(s.redisClient)

	return middleware.Chain(
		handler,
		middleware.RequestID,
		middleware.PoweredBy,
		trafficMiddleware,
		middleware.Logger,
		middleware.SecurityFuzzer(s.cfg.SecurityFuzzingEnabled),
		chaosMiddleware.Chaos,
		middleware.IPFilter(s.redisClient),
		middleware.RateLimit(rateLimiter),
		middleware.Recovery,
	)
}
