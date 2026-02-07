# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Chaos-Proxy is an "Immortality Layer" for APIs. It acts as a reverse proxy that learns from successful traffic and serves cached "Ghost" responses when the real backend fails. It also supports Chaos Engineering features like fault injection.

### Architecture
- **Sentinel (Go):** High-performance reverse proxy. Handles request forwarding, middleware (auth, rate-limiting, chaos injection), and Ghost Mode failover.
  - Entry point: `cmd/sentinel/main.go`
  - Core logic: `pkg/`
- **Brain (Python):** Machine Learning component. Consumes traffic logs from Redis, learns response patterns, and populates the Ghost cache.
  - Location: `brain/`
- **Dashboard (Next.js):** Real-time monitoring and control interface.
  - Location: `dashboard/`
- **CLI (TypeScript):** Command-line tool for managing the proxy.
  - Location: `packages/cli/`
- **Redis:** The central message bus and storage. Used for:
  - Traffic logs (Sentinel -> Brain)
  - Ghost responses (Brain -> Sentinel)
  - Configuration/Control (Dashboard/CLI -> Sentinel)
  - Pub/Sub events

## Development Environment Setup

### Prerequisites
- Go 1.21+
- Python 3.10+
- Node.js 18+
- Docker & Docker Compose (for Redis)

### Infrastructure
Start Redis (required for all components):
```bash
docker-compose up -d
```

## Build & Run Commands

### Sentinel (Proxy)
- **Run:** `go run cmd/sentinel/main.go`
- **Test:** `go test ./...`
- **Format:** `go fmt ./...`

### Brain (ML Engine)
- **Setup:** `cd brain && pip install -r requirements.txt`
- **Run:** `cd brain && python main.py`
- **Test:** `cd brain && python -m unittest test_learner.py`
- **Format:** Use `black`

### Dashboard (UI)
- **Setup:** `cd dashboard && npm install`
- **Run:** `cd dashboard && npm run dev`
- **Lint:** `cd dashboard && npm run lint`

### CLI Tool
- **Setup:** `cd packages/cli && npm install`
- **Build:** `cd packages/cli && npm run build`
- **Run:** `cd packages/cli && npm start` or `node dist/index.js`

### Integration Tests
- **Run Traffic Simulation:** `./tests/run_tests.sh`
- **Verify Ghost Mode:** `./verify_ghost.sh`

## Code Style & Standards

- **Go:** Follow standard Go idioms. Always run `go fmt` before committing.
- **Python:** Follow PEP 8. Use `black` for formatting.
- **TypeScript:** Follow the ESLint configuration.
- **Commits:** Follow [Conventional Commits](https://www.conventionalcommits.org/) (e.g., `feat:`, `fix:`, `docs:`).

## Configuration
- Environment variables are loaded from `.env` (copy from `.env.example`).
- Key variables:
  - `TARGET_URL`: The backend service being proxied.
  - `REDIS_ADDR`: Redis connection string.
  - `APP_ENV`: `development` or `production`.
