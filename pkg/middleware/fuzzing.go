package middleware

import (
	"log"
	"net/http"

	"github.com/elliot/chaosProxy/pkg/security"
)

// SecurityFuzzer returns a middleware that passively fuzzes requests.
func SecurityFuzzer(enabled bool) Middleware {
	if !enabled {
		return func(next http.Handler) http.Handler {
			return next
		}
	}

	log.Println("🛡️ Security Fuzzing Enabled! Be careful, this sends malicious requests to your backend.")

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 1. Pass request to next handler immediately (non-blocking)
			next.ServeHTTP(w, r)

			// 2. Clone request for fuzzing (very basic clone)
			// Cloning http.Request correctly for async use is tricky because the Context and Body might be invalid.
			// Since our Fuzzer currently only looks at the URL for GET params, we just clone the URL.
			// Ideally, we'd do a deep clone if we wanted to fuzz POST bodies later.

			// We only fuzz if traffic was successful (200-299) to avoid fuzzing already broken endpoints
			// But we don't have access to the status code here without wrapping the writer again.
			// For simplicity in this MVP, we fuzz everything.

			go security.Fuzz(r)
		})
	}
}
