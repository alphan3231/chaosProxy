package middleware

import "net/http"

// PoweredBy adds the X-Powered-By header to the response.
func PoweredBy(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Powered-By", "ChaosProxy/1.0.0")
		next.ServeHTTP(w, r)
	})
}
