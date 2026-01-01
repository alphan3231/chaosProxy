# 👻 Chaos-Proxy

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Python Version](https://img.shields.io/badge/Python-3.10+-3776AB?style=flat&logo=python)](https://python.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](CONTRIBUTING.md)

**Immortality Layer for APIs** — Your backend can crash, but your users won't notice.

<p align="center">
  <img src="https://img.shields.io/badge/Status-MVP%20Ready-success" alt="Status">
</p>

---


---

## 🏗 Architecture Overview

```mermaid
graph TD
    Client[Client / User] -->|HTTP Requests| Sentinel[Sentinel (Go Proxy)]
    Sentinel -->|Log Traffic| Redis[(Redis Cache)]
    Sentinel -->|Forward| Backend[Real Backend]
    
    Backend -->|Response| Sentinel
    Sentinel -->|Response| Client
    
    subgraph "Ghost Mode"
    Backend -.->|Failure| Sentinel
    Sentinel -.->|Fetch Ghost Response| Redis
    Redis -.->|Ghost Response| Sentinel
    end
    
    Redis -->|Traffic Data| Brain[Brain (Python AI)]
    Brain -->|Learned Patterns| Redis
```

The system consists of 3 main components:

```mermaid
graph TD
    Client[Client] -->|Traffic| Sentinel[Chaos-Proxy Sentinel]
    Sentinel -->|Forward| Backend[Backend Service]
    
    subgraph "The Immortality Layer"
        Sentinel -->|Log Traffic| Redis[(Redis)]
        Sentinel -->|Check Blocklist| Redis
        Sentinel -->|Get Chaos Settings| Redis
        Sentinel -->|Get Ghost Response| Redis
        
        Brain[The Brain (Python ML)] -->|Learn Patterns| Redis
        Brain -->|Update Ghost Cache| Redis
        
        Dashboard[Dashboard (Next.js)] -->|Monitor & Control| Redis
    end
    
    Backend -.->|Failure (50x)| Sentinel
    Sentinel -.->|Ghost Mode Activated| Client
```

## ✨ Features

- **🛡️ Sentinel Proxy (Go)** — High-performance reverse proxy with middleware support
- **🧠 The Brain (Python)** — Traffic analysis and pattern learning
- **👻 Ghost Mode** — Automatic failover with cached responses
- **😈 Chaos Mode** — Inject artificial latency and failures to test resilience
- **🚫 IP Blocking** — Real-time IP filtering and access control
- **📊 Real-time Dashboard** — Monitor traffic, ghost activations, and system health
- **🔒 Security First** — Rate limiting, auth, sensitive data filtering
- **⚡ Redis-powered** — Fast caching and pub/sub messaging

## 🚀 Quick Start

### Prerequisites

- Go 1.21+
- Python 3.10+
- Docker & Docker Compose
- Node.js 18+ (for Dashboard)

### 1. Clone & Setup

```bash
git clone https://github.com/elliot/chaosProxy.git
cd chaosProxy

# Copy environment template
cp .env.example .env
# Edit .env with your configuration
```

### 2. Start Infrastructure

```bash
docker-compose up -d  # Starts Redis
```

### 🐳 Docker (Recommended)

Run the entire stack with one command:

```bash
# Clone the repo
git clone https://github.com/elliot/chaosProxy.git
cd chaosProxy

# Configure
cp .env.example .env
# Edit .env with your settings

# Start everything
docker-compose up -d

# Check status
docker-compose ps

# View logs
docker-compose logs -f
```

Services will be available at:
- **Proxy:** http://localhost:8080
- **Dashboard:** http://localhost:3000
- **Redis:** localhost:6379

### 3. Run the Proxy

```bash
go run cmd/sentinel/main.go
```

### 4. Start the Brain (Learning Module)

```bash
cd brain
pip install -r requirements.txt  # Use virtualenv recommended
python main.py
```

### 5. Launch Dashboard

```bash
cd dashboard
npm install
npm run dev
# Open http://localhost:3000
```

### 6. Generate Traffic

```bash
./dev_test.sh  # Sends sample requests through the proxy
```

## 🛠 CLI Reference

Managing Chaos-Proxy is easy with the included CLI tool:

```bash
# View logs
chaos-proxy logs -f

# Block an IP
chaos-proxy block 192.168.1.5

# List blocked IPs
chaos-proxy ls-blocked

# Unblock an IP
chaos-proxy unblock 192.168.1.5
```

## 🧪 Testing Ghost Mode

```bash
# This script simulates a backend failure and verifies ghost mode activation
./verify_ghost.sh
```

## 📁 Project Structure

```
chaosProxy/
├── cmd/
│   └── sentinel/         # Main proxy entry point
├── internal/
│   ├── config/           # Configuration management
│   └── handlers/         # HTTP handlers (health check)
├── pkg/
│   ├── infrastructure/
│   │   └── redis/        # Redis client wrapper
│   └── middleware/       # Proxy middlewares (logging, rate-limit, traffic)
├── brain/                # Python learning module
│   ├── main.py           # Redis consumer
│   └── learner.py        # Learning logic
├── dashboard/            # Next.js monitoring UI
├── docker-compose.yml    # Redis infrastructure
└── .env.example          # Environment template
```

## ⚙️ Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Proxy listen port | `8080` |
| `TARGET_URL` | Backend service URL | `http://httpbin.org` |
| `REDIS_ADDR` | Redis connection address | `localhost:6379` |
| `REDIS_PASSWORD` | Redis password | _(empty)_ |
| `APP_ENV` | Environment mode | `development` |
| `DASHBOARD_USER` | Dashboard auth username | `admin` |
| `DASHBOARD_PASSWORD` | Dashboard auth password | `chaos123` |

## 🔒 Security Features

- ✅ Redis password authentication
- ✅ Request body size limiting (DoS protection)
- ✅ Sensitive header filtering
- ✅ Body content sanitization
- ✅ Rate limiting (100 req/min per IP)
- ✅ Dashboard Basic Authentication
- ✅ CORS & security headers

## 🗺️ Roadmap

See [ROADMAP.md](ROADMAP.md) for the detailed development plan.

- [x] Phase 1: Project Setup
- [x] Phase 2: Sentinel Proxy Core
- [x] Phase 3: The Brain (Learning)
- [x] Phase 4: Ghost Mode
- [x] Phase 5: Dashboard
- [x] Security Hardening
- [x] Phase 6: Chaos Mode (Fault Injection)
- [x] Phase 7: Polish & Refinement (CLI Logs, Refactoring)
- [x] Phase 8: Security Access Control (IP Filtering)
- [ ] Phase 9: Cloud Deployment (AWS/GCP)

## 🤝 Contributing

Contributions are welcome! Please read [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 💡 Use Cases

- **Staging/Demo Environments** — Never show errors to stakeholders
- **Chaos Engineering** — Test your frontend's resilience
- **API Mocking** — Generate realistic mock responses from real traffic
- **Graceful Degradation** — Serve cached data when services are down

---

<p align="center">
  Made with 👻 by <a href="https://github.com/alphan3231">elliot</a>
</p>
