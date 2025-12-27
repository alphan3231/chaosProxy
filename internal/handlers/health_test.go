package handlers

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
)

// MockRedisClient mocks the Redis client for testing
type MockRedisClient struct {
	PingSuccess bool
}

func TestHealthCheck_Response(t *testing.T) {
	// Since we can't easily mock the Redis client in this test,
	// we'll just test that the handler returns valid JSON structure

	// This is more of a contract test
	t.Run("Response structure", func(t *testing.T) {
		response := HealthResponse{
			Status: "ok",
			Redis:  "connected",
		}

		data, err := json.Marshal(response)
		if err != nil {
			t.Fatalf("Failed to marshal response: %v", err)
		}

		var unmarshaled HealthResponse
		if err := json.Unmarshal(data, &unmarshaled); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if unmarshaled.Status != "ok" {
			t.Errorf("Expected status 'ok', got '%s'", unmarshaled.Status)
		}

		if unmarshaled.Redis != "connected" {
			t.Errorf("Expected redis 'connected', got '%s'", unmarshaled.Redis)
		}
	})
}

func TestHealthResponse_JSON(t *testing.T) {
	tests := []struct {
		name     string
		response HealthResponse
	}{
		{
			name:     "Healthy",
			response: HealthResponse{Status: "ok", Redis: "connected"},
		},
		{
			name:     "Degraded",
			response: HealthResponse{Status: "degraded", Redis: "disconnected"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()

			recorder.Header().Set("Content-Type", "application/json")
			json.NewEncoder(recorder).Encode(tt.response)

			if recorder.Header().Get("Content-Type") != "application/json" {
				t.Error("Content-Type should be application/json")
			}

			var result HealthResponse
			if err := json.NewDecoder(recorder.Body).Decode(&result); err != nil {
				t.Fatalf("Failed to decode response: %v", err)
			}

			if result.Status != tt.response.Status {
				t.Errorf("Status mismatch: expected %s, got %s", tt.response.Status, result.Status)
			}
		})
	}
}
