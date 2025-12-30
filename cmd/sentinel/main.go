package main

import (
	"log"

	"github.com/elliot/chaosProxy/internal/config"
	"github.com/elliot/chaosProxy/internal/server"
	"github.com/elliot/chaosProxy/pkg/infrastructure/redis"
)

func main() {
	// Load Configuration
	cfg := config.LoadConfig()

	// Initialize Redis
	redisClient, err := redis.NewClient(cfg.RedisAddr, cfg.RedisPassword)
	if err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}
	defer redisClient.Close()
	log.Printf("⚡ Connected to Redis at %s", cfg.RedisAddr)

	// Initialize and Start Server
	srv := server.NewServer(cfg, redisClient)
	if err := srv.Start(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
