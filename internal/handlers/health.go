package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/elliot/chaosProxy/pkg/infrastructure/redis"
	"github.com/elliot/chaosProxy/pkg/response"
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

		resp := HealthResponse{
			Status: status,
			Redis:  redisStatus,
		}

		httpStatus := http.StatusOK
		if status != "ok" {
			httpStatus = http.StatusServiceUnavailable
		}

		response.JSON(w, httpStatus, resp)
	}
}
