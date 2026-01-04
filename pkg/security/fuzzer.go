package security

import (
	"log"
	"net/http"
	"net/url"
)

// Common payloads for basic fuzzing
var payloads = []string{
	"' OR 1=1--",                // SQL Injection
	"<script>alert(1)</script>", // XSS
	"../../etc/passwd",          // Path Traversal
	"WAITFOR DELAY '0:0:5'",     // Time-based SQLi
}

// Fuzz takes a request and asynchronously replays it with injected payloads.
// It logs any suspicious responses (e.g. 500 errors).
func Fuzz(originalReq *http.Request) {
	// We only fuzz GET parameters for safety/simplicity in this basic version
	// Fuzzing POST bodies requires deeper parsing which we can add later.

	targetURL := originalReq.URL
	params := targetURL.Query()

	if len(params) == 0 {
		return // Nothing to fuzz
	}

	for param, values := range params {
		for _, payload := range payloads {
			// Clone the values to avoid modifying the original map for other iterations
			newValues := make([]string, len(values))
			copy(newValues, values)

			// Inject payload into the first value of the parameter
			if len(newValues) > 0 {
				newValues[0] = payload // Replace with payload
			} else {
				newValues = []string{payload}
			}

			// Create new query string
			newParams := url.Values{}
			for p, v := range params {
				if p == param {
					newParams[p] = newValues
				} else {
					newParams[p] = v
				}
			}

			// Construct new URL
			fuzzURL := *targetURL
			fuzzURL.RawQuery = newParams.Encode()

			// Execute Fuzz Request
			go executeFuzzRequest(fuzzURL.String(), originalReq.Header)
		}
	}
}

func executeFuzzRequest(fuzzURL string, headers http.Header) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fuzzURL, nil)
	if err != nil {
		return
	}

	// Copy headers (excluding sensitive ones potentially, but strictly we want to mimic the user)
	for k, v := range headers {
		req.Header[k] = v
	}
	req.Header.Set("X-Chaos-Fuzzer", "true") // Tag our traffic

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// Analyze Result
	if resp.StatusCode >= 500 {
		log.Printf("🚨 [SECURITY] Vulnerability Detected? Payload caused %d at %s", resp.StatusCode, fuzzURL)
	} else if resp.StatusCode != 404 && resp.StatusCode != 400 {
		// Log successes for analysis too, if needed
		// log.Printf("⚠️ [SECURITY] Payload accepted (%d) at %s", resp.StatusCode, fuzzURL)
	}
}
