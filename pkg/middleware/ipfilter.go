package middleware

import (
	"context"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/elliot/chaosProxy/pkg/infrastructure/redis"
)

type CachedIPFilter struct {
	redisClient *redis.Client
	blockedIPs  map[string]bool
	lastFetch   time.Time
	mu          sync.RWMutex
}

func NewCachedIPFilter(redisClient *redis.Client) *CachedIPFilter {
	return &CachedIPFilter{
		redisClient: redisClient,
		blockedIPs:  make(map[string]bool),
	}
}

func (f *CachedIPFilter) refresh() {
	f.mu.RLock()
	if time.Since(f.lastFetch) < 5*time.Second {
		f.mu.RUnlock()
		return
	}
	f.mu.RUnlock()

	// Fetch all blocked IPs
	// Note: We need to implement GetBlockedIPs in redis client or use SMembers
	ips, err := f.redisClient.GetBlockedIPs(context.Background())
	if err != nil {
		log.Printf("Failed to refresh blocked IPs: %v", err)
		return
	}

	newMap := make(map[string]bool)
	for _, ip := range ips {
		newMap[ip] = true
	}

	f.mu.Lock()
	f.blockedIPs = newMap
	f.lastFetch = time.Now()
	f.mu.Unlock()
}

func (f *CachedIPFilter) IsBlocked(ip string) bool {
	f.refresh()
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.blockedIPs[ip]
}

func IPFilter(redisClient *redis.Client) Middleware {
	filter := NewCachedIPFilter(redisClient)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := getRealIP(r)

			// Use cached check
			if filter.IsBlocked(ip) {
				log.Printf("🚫 Blocked request from %s", ip)
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte(`{"error": "Access Denied"}`))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func getRealIP(r *http.Request) string {
	// Security: Do not trust X-Forwarded-For or X-Real-IP headers as they can be spoofed.
	// Since this proxy is intended to be the edge sentinel, we rely on the actual RemoteAddr.
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}
