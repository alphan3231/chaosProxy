package config

import (
	"os"
	"testing"
)

func TestLoadConfig_Defaults(t *testing.T) {
	// Clear any existing env vars
	os.Unsetenv("PORT")
	os.Unsetenv("TARGET_URL")
	os.Unsetenv("REDIS_ADDR")
	os.Unsetenv("REDIS_PASSWORD")
	os.Unsetenv("APP_ENV")

	cfg := LoadConfig()

	if cfg.Port != "8080" {
		t.Errorf("Expected Port to be '8080', got '%s'", cfg.Port)
	}

	if cfg.TargetURL != "http://httpbin.org" {
		t.Errorf("Expected TargetURL to be 'http://httpbin.org', got '%s'", cfg.TargetURL)
	}

	if cfg.RedisAddr != "localhost:6379" {
		t.Errorf("Expected RedisAddr to be 'localhost:6379', got '%s'", cfg.RedisAddr)
	}

	if cfg.RedisPassword != "" {
		t.Errorf("Expected RedisPassword to be empty, got '%s'", cfg.RedisPassword)
	}

	if cfg.AppEnv != "development" {
		t.Errorf("Expected AppEnv to be 'development', got '%s'", cfg.AppEnv)
	}
}

func TestLoadConfig_FromEnv(t *testing.T) {
	// Set custom env vars
	os.Setenv("PORT", "9090")
	os.Setenv("TARGET_URL", "http://my-backend:8000")
	os.Setenv("REDIS_ADDR", "redis.example.com:6380")
	os.Setenv("REDIS_PASSWORD", "supersecret")
	os.Setenv("APP_ENV", "production")

	defer func() {
		os.Unsetenv("PORT")
		os.Unsetenv("TARGET_URL")
		os.Unsetenv("REDIS_ADDR")
		os.Unsetenv("REDIS_PASSWORD")
		os.Unsetenv("APP_ENV")
	}()

	cfg := LoadConfig()

	if cfg.Port != "9090" {
		t.Errorf("Expected Port to be '9090', got '%s'", cfg.Port)
	}

	if cfg.TargetURL != "http://my-backend:8000" {
		t.Errorf("Expected TargetURL to be 'http://my-backend:8000', got '%s'", cfg.TargetURL)
	}

	if cfg.RedisAddr != "redis.example.com:6380" {
		t.Errorf("Expected RedisAddr to be 'redis.example.com:6380', got '%s'", cfg.RedisAddr)
	}

	if cfg.RedisPassword != "supersecret" {
		t.Errorf("Expected RedisPassword to be 'supersecret', got '%s'", cfg.RedisPassword)
	}

	if cfg.AppEnv != "production" {
		t.Errorf("Expected AppEnv to be 'production', got '%s'", cfg.AppEnv)
	}
}

func TestGetEnv_Fallback(t *testing.T) {
	os.Unsetenv("TEST_VAR")

	result := getEnv("TEST_VAR", "default_value")

	if result != "default_value" {
		t.Errorf("Expected 'default_value', got '%s'", result)
	}
}

func TestGetEnv_ExistingVar(t *testing.T) {
	os.Setenv("TEST_VAR", "actual_value")
	defer os.Unsetenv("TEST_VAR")

	result := getEnv("TEST_VAR", "default_value")

	if result != "actual_value" {
		t.Errorf("Expected 'actual_value', got '%s'", result)
	}
}
