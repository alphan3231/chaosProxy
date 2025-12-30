package middleware

import (
	"context"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/elliot/chaosProxy/pkg/infrastructure/redis"
)

func IPFilter(redisClient *redis.Client) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := getRealIP(r)

			blocked, err := redisClient.IsIPBlocked(context.Background(), ip)
			if err != nil {
				// Fail open on Redis error, log it
				log.Printf("Failed to check blocklist for %s: %v", ip, err)
			} else if blocked {
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
	// Check X-Forwarded-For
	xfwd := r.Header.Get("X-Forwarded-For")
	if xfwd != "" {
		// Can look like "client_ip, proxy1, proxy2"
		ips := strings.Split(xfwd, ",")
		return strings.TrimSpace(ips[0])
	}

	// Fallback to RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}
