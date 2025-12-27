package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/elliot/chaosProxy/pkg/infrastructure/redis"
)

type HealthResponse struct {
	Status string `json:"status"`
	Redis  string `json:"redis"`
}

func HealthCheck(redisClient *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := "ok"
		redisStatus := "connected"

		// Check Redis
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		if err := redisClient.GetRawClient().Ping(ctx).Err(); err != nil {
			redisStatus = "disconnected"
			status = "degraded"
		}

		response := HealthResponse{
			Status: status,
			Redis:  redisStatus,
		}

		w.Header().Set("Content-Type", "application/json")
		if status != "ok" {
			w.WriteHeader(http.StatusServiceUnavailable)
		} else {
			w.WriteHeader(http.StatusOK)
		}

		json.NewEncoder(w).Encode(response)
	}
}
