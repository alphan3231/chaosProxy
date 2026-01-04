package middleware

import (
	"net/http"
	"time"
)

// RetryTransport wraps a RoundTripper to add retry logic.
type RetryTransport struct {
	Base       http.RoundTripper
	MaxRetries int
	Delay      time.Duration
}

// NewRetryTransport creates a new RetryTransport.
// If base is nil, http.DefaultTransport is used.
func NewRetryTransport(base http.RoundTripper, maxRetries int, delayMs int) *RetryTransport {
	if base == nil {
		base = http.DefaultTransport
	}
	return &RetryTransport{
		Base:       base,
		MaxRetries: maxRetries,
		Delay:      time.Duration(delayMs) * time.Millisecond,
	}
}

// RoundTrip executes the request and retries on failure.
func (t *RetryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Skeleton: Just pass through for now
	// Logic to be implemented in next commit
	return t.Base.RoundTrip(req)
}
