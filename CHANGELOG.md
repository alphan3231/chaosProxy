# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

### 🚀 Added
- Created project skeleton (`cmd`, `internal`, `pkg`, `configs`).
- Initialized Go module: `github.com/elliot/chaosProxy`.
- Added `docker-compose.yml` for Redis infrastructure.
- Documented project vision and plan in `ROADMAP.md`.
- Implemented basic Reverse Proxy core using `httputil`.
- Added `internal/config` for environment variable management.
- Implemented Middleware pattern with `Logger` and `Recovery`.
- Integrated `go-redis/v9` client in `pkg/infrastructure/redis`.
- Implemented `TrafficLogger` middleware to capture Request/Response pairs.
- Added `PublishTraffic` method to Redis client for async logging to `chaos:traffic` channel.
- Added `/healthz` endpoint for checking service and Redis status.
- Key Phase 3: Setup Python environment structure (`brain/`).
- Implemented basic Redis consumer in `brain/main.py`.
