package middleware

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/elliot/chaosProxy/pkg/infrastructure/redis"
)

type ChaosMiddleware struct {
	redisClient *redis.Client
	settings    *redis.ChaosSettings
	mu          sync.RWMutex
	lastFetch   time.Time
}

func NewChaosMiddleware(redisClient *redis.Client) *ChaosMiddleware {
	return &ChaosMiddleware{
		redisClient: redisClient,
		settings:    &redis.ChaosSettings{}, // Default empty
	}
}

func (c *ChaosMiddleware) refreshSettings() {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Cache for 1 second to avoid hitting Redis on every request
	if time.Since(c.lastFetch) < 1*time.Second {
		return
	}

	settings, err := c.redisClient.GetChaosSettings(context.Background())
	if err != nil {
		// Log error but continue with old settings
		// log.Printf("Failed to fetch chaos settings: %v", err)
		return
	}

	c.settings = settings
	c.lastFetch = time.Now()
}

func (c *ChaosMiddleware) Chaos(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Update settings async or sync? Sync is safer for immediate effect but adds latency.
		// Given we cache for 1s, sync is fine.
		c.refreshSettings()

		c.mu.RLock()
		settings := *c.settings // Copy
		c.mu.RUnlock()

		// 1. Latency Injection
		if settings.LatencyEnabled {
			min := settings.LatencyMin
			max := settings.LatencyMax
			if max > min {
				// Sleep random duration
				delay := time.Duration(rand.Intn(max-min)+min) * time.Millisecond
				time.Sleep(delay)
			}
		}

		// 2. Failure Injection
		if settings.FailureEnabled {
			// Random number between 0-99
			if rand.Intn(100) < settings.FailureRate {
				log.Printf("💀 CHAOS: Injecting failure for %s", r.URL.Path)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error": "Chaos Monkey Struck!"}`))
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
