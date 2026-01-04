package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port          string
	TargetURL     string
	RedisAddr     string
	RedisPassword string
	AppEnv        string
	WebhookURL    string
	CanaryURL     string
	CanaryWeight  int
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func LoadConfig() *Config {
	return &Config{
		Port:          getEnv("PORT", "8080"),
		TargetURL:     getEnv("TARGET_URL", "http://httpbin.org"),
		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		AppEnv:        getEnv("APP_ENV", "development"),
		WebhookURL:    getEnv("WEBHOOK_URL", ""),
		CanaryURL:     getEnv("CANARY_URL", ""),
		CanaryWeight:  getEnvInt("CANARY_WEIGHT", 0),
	}
}

func getEnvInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		var i int
		if _, err := fmt.Sscanf(value, "%d", &i); err == nil {
			return i
		}
	}
	return fallback
}
