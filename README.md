# 👻 Chaos-Proxy

**Immortality Layer for APIs**

Chaos-Proxy is a smart Reverse Proxy designed for microservice and API architectures. It ensures continuous service availability to clients even when backend services fail.

During normal operation, it monitors and learns from traffic (The Sentinel & The Brain). When the backend crashes, "Ghost Mode" activates, generating realistic responses based on learned data.

## 🚀 Features

- **🪟 Sentinel Proxy (Go):** High-performance, low-latency reverse proxy.
- **🧠 The Brain (Python + AI):** Traffic analysis and response modeling.
- **👻 Ghost Mode:** Smart simulation that automatically activates during backend outages.
- **⚡ Redis Backed:** Fast data access and cache management.

## 🛠 Installation

### Requirements
- Go 1.21+
- Python 3.10+
- Docker & Docker Compose

### Quick Start

1. **Clone the repo:**
   ```bash
   git clone https://github.com/elliot/chaosProxy.git
   cd chaosProxy
   ```

2. **Start the infrastructure (Redis):**
   ```bash
   docker-compose up -d
   ```

3. **Run the Proxy:**
   ```bash
   go run cmd/sentinel/main.go
   ```

## 🗺 Roadmap

See [ROADMAP.md](ROADMAP.md) for the detailed development plan.

## 📄 License

MIT
