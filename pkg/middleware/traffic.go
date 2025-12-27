package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/elliot/chaosProxy/pkg/infrastructure/redis"
)

// responseWriterWrapper wraps http.ResponseWriter to capture the status and body
type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

func (rw *responseWriterWrapper) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriterWrapper) Write(b []byte) (int, error) {
	rw.body.Write(b)
	return rw.ResponseWriter.Write(b)
}

func TrafficLogger(redisClient *redis.Client) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// 1. Capture Request Body
			var reqBody []byte
			if r.Body != nil {
				reqBody, _ = io.ReadAll(r.Body)
				r.Body = io.NopCloser(bytes.NewBuffer(reqBody)) // Restore body
			}

			// 2. Wrap Response Writer
			wrapper := &responseWriterWrapper{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
				body:           &bytes.Buffer{},
			}

			// 3. Process Request
			next.ServeHTTP(wrapper, r)

			// 4. Async Log to Redis
			duration := time.Since(start)

			go func() {
				// Create log entry
				entry := redis.TrafficLog{
					Timestamp:    start,
					Method:       r.Method,
					Path:         r.URL.Path,
					RequestBody:  string(reqBody),
					Status:       wrapper.statusCode,
					ResponseBody: wrapper.body.String(),
					Duration:     duration.String(),
				}

				// Encode to JSON
				data, err := json.Marshal(entry)
				if err != nil {
					log.Printf("Failed to marshal log entry: %v", err)
					return
				}

				// We can't pass the original entry directly because go-redis Publish expects string/byte/int etc
				// So we marshal it to JSON bytes (or string) first, or pass it as is if it implemented BinaryMarshaler
				// But we just manually marshalled it above.

				if err := redisClient.GetRawClient().Publish(context.Background(), "chaos:traffic", data).Err(); err != nil {
					log.Printf("Failed to publish traffic log: %v", err)
				}
			}()
		})
	}
}
