package middleware

import (
	"log"
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
	var resp *http.Response
	var err error

	// Try at least once + retries
	for i := 0; i <= t.MaxRetries; i++ {
		if i > 0 {
			log.Printf("🔄 [Retry] Attempt %d/%d for %s", i, t.MaxRetries, req.URL.Path)
			time.Sleep(t.Delay)
		}

		resp, err = t.Base.RoundTrip(req)

		// Check for success or terminal error
		if err == nil && resp.StatusCode < 500 {
			return resp, nil
		}

		// If it's a 5xx error, we close the body so we can retry cleanly
		// (though RoundTrip handles connection state, we should be careful with bodies)
		if resp != nil {
			resp.Body.Close()
		}
	}

	// Return last error or response
	return resp, err
}
