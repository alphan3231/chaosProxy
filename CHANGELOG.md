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
- Added `brain/learner.py` implementing "Exact Match" learning strategy (caching 200 OK responses to Redis).
- Implemented **Ghost Mode** recovery logic in Sentinel Proxy (`proxy.ErrorHandler`).
- Added `verify_ghost.sh` for semantic verification of recovery mode.
- **Phase 5:** Initialized `dashboard` (Next.js 14) for real-time monitoring.
- Added Redis Stats tracking (`request_count`, `ghost_count`) to Brain and Proxy.

### 🔒 Security Hardening
- Added Redis password authentication support.
- Implemented request body size limit (1MB) to prevent DoS attacks.
- Added sensitive header filtering (`Authorization`, `Cookie`, etc.).
- Implemented body content sanitization (password, token, secret masking).
- Added Rate Limiting middleware (100 req/min per IP).
- Implemented Dashboard Basic Authentication.
- Added CORS configuration and security headers.
- Created `.env.example` for secure configuration reference.
