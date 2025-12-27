package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/elliot/chaosProxy/pkg/infrastructure/redis"
)

const (
	MaxBodySize = 1024 * 1024 // 1MB limit
)

// Sensitive headers that should not be logged
var sensitiveHeaders = map[string]bool{
	"authorization":  true,
	"cookie":         true,
	"x-api-key":      true,
	"x-auth-token":   true,
	"x-access-token": true,
	"set-cookie":     true,
}

// Patterns to mask in request/response bodies
var sensitivePatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)"password"\s*:\s*"[^"]*"`),
	regexp.MustCompile(`(?i)"token"\s*:\s*"[^"]*"`),
	regexp.MustCompile(`(?i)"secret"\s*:\s*"[^"]*"`),
	regexp.MustCompile(`(?i)"api_key"\s*:\s*"[^"]*"`),
	regexp.MustCompile(`(?i)"credit_card"\s*:\s*"[^"]*"`),
}

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

// sanitizeBody masks sensitive data in the body
func sanitizeBody(body string) string {
	for _, pattern := range sensitivePatterns {
		body = pattern.ReplaceAllString(body, `"[REDACTED]"`)
	}
	return body
}

// hasSensitiveHeader checks if request has sensitive headers
func hasSensitiveHeader(r *http.Request) bool {
	for header := range r.Header {
		if sensitiveHeaders[strings.ToLower(header)] {
			return true
		}
	}
	return false
}

func TrafficLogger(redisClient *redis.Client) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// 1. Capture Request Body with Size Limit
			var reqBody []byte
			if r.Body != nil {
				// Limit body size to prevent DoS
				limitedReader := io.LimitReader(r.Body, MaxBodySize)
				reqBody, _ = io.ReadAll(limitedReader)
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

			// 4. Async Log to Redis (with sanitization)
			duration := time.Since(start)

			go func() {
				// Skip logging if request has sensitive headers
				if hasSensitiveHeader(r) {
					log.Printf("[Traffic] Skipped logging for %s %s (sensitive headers)", r.Method, r.URL.Path)
					return
				}

				// Sanitize bodies
				sanitizedReqBody := sanitizeBody(string(reqBody))
				sanitizedRespBody := sanitizeBody(wrapper.body.String())

				// Create log entry
				entry := redis.TrafficLog{
					Timestamp:    start,
					Method:       r.Method,
					Path:         r.URL.Path,
					RequestBody:  sanitizedReqBody,
					Status:       wrapper.statusCode,
					ResponseBody: sanitizedRespBody,
					Duration:     duration.String(),
				}

				// Encode to JSON
				data, err := json.Marshal(entry)
				if err != nil {
					log.Printf("Failed to marshal log entry: %v", err)
					return
				}

				if err := redisClient.GetRawClient().Publish(context.Background(), "chaos:traffic", data).Err(); err != nil {
					log.Printf("Failed to publish traffic log: %v", err)
				}
			}()
		})
	}
}
