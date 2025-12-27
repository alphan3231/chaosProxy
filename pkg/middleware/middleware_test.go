package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestChain(t *testing.T) {
	// Create a simple handler
	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Create middlewares that add headers
	middleware1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Middleware-1", "true")
			next.ServeHTTP(w, r)
		})
	}

	middleware2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Middleware-2", "true")
			next.ServeHTTP(w, r)
		})
	}

	// Chain them
	handler := Chain(finalHandler, middleware1, middleware2)

	// Test
	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Header().Get("X-Middleware-1") != "true" {
		t.Error("Middleware 1 was not applied")
	}

	if rec.Header().Get("X-Middleware-2") != "true" {
		t.Error("Middleware 2 was not applied")
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}
}

func TestLogger(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	loggedHandler := Logger(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()

	loggedHandler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}
}

func TestRecovery(t *testing.T) {
	// Handler that panics
	panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})

	recoveredHandler := Recovery(panicHandler)

	req := httptest.NewRequest("GET", "/panic", nil)
	rec := httptest.NewRecorder()

	// Should not panic
	recoveredHandler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500 after panic, got %d", rec.Code)
	}
}

func TestRateLimiter_Allow(t *testing.T) {
	limiter := NewRateLimiter(5, time.Minute)

	// First 5 requests should be allowed
	for i := 0; i < 5; i++ {
		if !limiter.Allow("192.168.1.1") {
			t.Errorf("Request %d should have been allowed", i+1)
		}
	}

	// 6th request should be blocked
	if limiter.Allow("192.168.1.1") {
		t.Error("6th request should have been blocked")
	}
}

func TestRateLimiter_DifferentIPs(t *testing.T) {
	limiter := NewRateLimiter(2, time.Minute)

	// IP 1
	limiter.Allow("192.168.1.1")
	limiter.Allow("192.168.1.1")

	// IP 1 should be blocked
	if limiter.Allow("192.168.1.1") {
		t.Error("IP 1 should be blocked after 2 requests")
	}

	// IP 2 should still be allowed
	if !limiter.Allow("192.168.1.2") {
		t.Error("IP 2 should still be allowed")
	}
}

func TestRateLimit_Middleware(t *testing.T) {
	limiter := NewRateLimiter(2, time.Minute)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	limitedHandler := RateLimit(limiter)(handler)

	// First 2 requests should pass
	for i := 0; i < 2; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "192.168.1.1:12345"
		rec := httptest.NewRecorder()

		limitedHandler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Request %d should have returned 200, got %d", i+1, rec.Code)
		}
	}

	// 3rd request should be rate limited
	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	rec := httptest.NewRecorder()

	limitedHandler.ServeHTTP(rec, req)

	if rec.Code != http.StatusTooManyRequests {
		t.Errorf("3rd request should have returned 429, got %d", rec.Code)
	}
}

func TestSanitizeBody(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Password field",
			input:    `{"username": "john", "password": "secret123"}`,
			expected: `{"username": "john", "[REDACTED]"}`,
		},
		{
			name:     "Token field",
			input:    `{"token": "abc123xyz"}`,
			expected: `{"[REDACTED]"}`,
		},
		{
			name:     "No sensitive data",
			input:    `{"name": "John", "age": 30}`,
			expected: `{"name": "John", "age": 30}`,
		},
		{
			name:     "API key field",
			input:    `{"api_key": "sk-12345"}`,
			expected: `{"[REDACTED]"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizeBody(tt.input)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestHasSensitiveHeader(t *testing.T) {
	tests := []struct {
		name     string
		headers  map[string]string
		expected bool
	}{
		{
			name:     "Authorization header",
			headers:  map[string]string{"Authorization": "Bearer token123"},
			expected: true,
		},
		{
			name:     "Cookie header",
			headers:  map[string]string{"Cookie": "session=abc"},
			expected: true,
		},
		{
			name:     "Regular header",
			headers:  map[string]string{"Content-Type": "application/json"},
			expected: false,
		},
		{
			name:     "X-API-Key header",
			headers:  map[string]string{"X-API-Key": "key123"},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)
			for k, v := range tt.headers {
				req.Header.Set(k, v)
			}

			result := hasSensitiveHeader(req)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}
