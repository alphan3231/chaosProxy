package config

import (
	"os"
)

type Config struct {
	Port      string
	TargetURL string
	RedisAddr string
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func LoadConfig() *Config {
	return &Config{
		Port:      getEnv("PORT", "8080"),
		TargetURL: getEnv("TARGET_URL", "http://httpbin.org"),
		RedisAddr: getEnv("REDIS_ADDR", "localhost:6379"),
	}
}
