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

4. **Start The Brain (Learning Module):**
   ```bash
   # In a new terminal
   python3 brain/main.py
   ```

## 👻 How It Works

1.  **Learning:** Send requests to the proxy (`Standard Mode`). The usage data is silently logged to Redis.
2.  **Simulation:** If the backend (e.g., httpbin.org) fails, the Proxy enters `Ghost Mode`.
3.  **Immortality:** The Proxy checks Redis for a previously "learned" response for that specific method/path and returns it instantly.

## 🧪 Testing

We have included scripts to help you verify the "Immortality":

- `./dev_test.sh`: Generates normal traffic to train the brain.
- `./verify_ghost.sh`: Simulates a backend failure and checks if `Ghost Mode` activates successfully.

## 📄 License

MIT
