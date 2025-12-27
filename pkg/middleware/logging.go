package middleware

import (
	"log"
	"net/http"
	"time"
)

// Logger logs the request details and execution time.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap ResponseWriter to capture status code could be added here later

		next.ServeHTTP(w, r)

		log.Printf("[REQ] %s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

// Recovery handles panics in the application.
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("[PANIC] %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
