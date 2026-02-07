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
	// Security: Do not trust X-Forwarded-For or X-Real-IP headers as they can be spoofed.
	// Since this proxy is intended to be the edge sentinel, we rely on the actual RemoteAddr.
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}
