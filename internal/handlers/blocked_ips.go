package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/elliot/chaosProxy/pkg/infrastructure/redis"
)

type BlockedIPsResponse struct {
	BlockedIPs []string `json:"blocked_ips"`
	Count      int      `json:"count"`
}

func GetBlockedIPs(redisClient *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Fetch blocked IPs set from Redis
		ips, err := redisClient.GetRawClient().SMembers(ctx, "chaos:settings:blocked_ips").Result()
		if err != nil {
			http.Error(w, "Failed to fetch blocked IPs", http.StatusInternalServerError)
			return
		}

		response := BlockedIPsResponse{
			BlockedIPs: ips,
			Count:      len(ips),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
