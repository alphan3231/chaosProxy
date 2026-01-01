package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
)

type contextKey string

const RequestIDKey contextKey = "requestID"

// RequestID adds a unique ID to the request context and headers.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get("X-Request-ID")
		if reqID == "" {
			reqID = generateID()
		}

		ctx := context.WithValue(r.Context(), RequestIDKey, reqID)
		w.Header().Set("X-Request-ID", reqID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func generateID() string {
	b := make([]byte, 8) // 16 characters hex
	if _, err := rand.Read(b); err != nil {
		// Fallback if random fails (unlikely)
		return fmt.Sprintf("req-%d", 0)
	}
	return hex.EncodeToString(b)
}
