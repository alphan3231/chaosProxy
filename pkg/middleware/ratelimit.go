package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"
)

type RateLimiter struct {
	requests map[string][]time.Time
	mu       sync.RWMutex
	limit    int
	window   time.Duration
	stopChan chan struct{}
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
		stopChan: make(chan struct{}),
	}

	// Cleanup goroutine
	go func() {
		ticker := time.NewTicker(window)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				rl.cleanup()
			case <-rl.stopChan:
				return
			}
		}
	}()

	return rl
}

func (rl *RateLimiter) Stop() {
	close(rl.stopChan)
}

func (rl *RateLimiter) cleanup() {
	// First pass: just read to collect keys (holding Read Lock)
	rl.mu.RLock()
	var keys []string
	for k := range rl.requests {
		keys = append(keys, k)
	}
	rl.mu.RUnlock()

	// Second pass: cleanup individual entries (holding Write Lock briefly per entry)
	now := time.Now()
	for _, ip := range keys {
		rl.mu.Lock()
		timestamps, exists := rl.requests[ip]
		if !exists {
			rl.mu.Unlock()
			continue
		}

		var valid []time.Time
		for _, t := range timestamps {
			if now.Sub(t) < rl.window {
				valid = append(valid, t)
			}
		}

		if len(valid) == 0 {
			delete(rl.requests, ip)
		} else {
			rl.requests[ip] = valid
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	timestamps := rl.requests[ip]

	// Filter old timestamps
	var valid []time.Time
	for _, t := range timestamps {
		if now.Sub(t) < rl.window {
			valid = append(valid, t)
		}
	}

	if len(valid) >= rl.limit {
		return false
	}

	valid = append(valid, now)
	rl.requests[ip] = valid
	return true
}

func getClientIP(r *http.Request) string {
	// Security: Do not trust headers. Use RemoteAddr.
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

func RateLimit(limiter *RateLimiter) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := getClientIP(r)

			if !limiter.Allow(ip) {
				w.Header().Set("Retry-After", "60")
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte(`{"error": "Too many requests. Please try again later."}`))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
