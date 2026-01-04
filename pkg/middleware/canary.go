package middleware

import (
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// NewCanary returns a middleware that routes traffic to a canary URL based on weight.
func NewCanary(canaryURL string, weight int) Middleware {
	// 1. Validation: If no URL or 0 weight, return a pass-through
	if canaryURL == "" || weight <= 0 {
		return func(next http.Handler) http.Handler {
			return next
		}
	}

	// 2. Parse Canary URL
	target, err := url.Parse(canaryURL)
	if err != nil {
		log.Printf("⚠️ Invalid Canary URL received: %v. Canary routing disabled.", err)
		return func(next http.Handler) http.Handler {
			return next
		}
	}

	// 3. Create Reverse Proxy for Canary
	canaryProxy := httputil.NewSingleHostReverseProxy(target)

	// Copy Director logic to ensure Host header is set correctly
	originalDirector := canaryProxy.Director
	canaryProxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Host = target.Host
	}

	log.Printf("🐤 Canary Routing Enabled: targeting %s with %d%% weight", canaryURL, weight)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 4. Random dice roll
			if rand.Intn(100) < weight {
				// 🐤 Canary Hit!
				log.Printf("🐤 Routing to Canary: %s", r.URL.Path)
				canaryProxy.ServeHTTP(w, r)
				return
			}

			// 5. Normal Path
			next.ServeHTTP(w, r)
		})
	}
}
