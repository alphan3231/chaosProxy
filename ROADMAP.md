# Chaos-Proxy: Roadmap

## 🎯 Project Vision
**Chaos-Proxy** is an "immortality" layer for API services. It sits between backend services and clients, monitoring and learning from traffic during normal operation. When the backend fails, it activates AI-powered "Ghost Mode" to simulate traffic, creating a perception of uninterrupted service.

---

## 🏗 Architecture Overview
The system consists of 3 main components:
1.  **The Sentinel (Go):** High-performance Reverse Proxy. Handles traffic, logs to Redis, and serves Ghost responses if needed.
2.  **The Memory (Redis):** Fast cache holding live request/response data and AI models (or rules).
3.  **The Brain (Python):** Runs in the background. Analyzes data in Redis, learns patterns, and creates "Ghost Response" models.

---

## 🚀 Development Phases

### Phase 1: The Sentinel (Proxy Core & Monitoring)
*Goal: A working Go Proxy that forwards and logs traffic.*

- [x] **Project Setup:** Go module structure and basic directories.
- [ ] **Reverse Proxy Core:** Basic proxy logic using `net/http/httputil`.
- [ ] **Middleware Structure:** Hook points to capture Request/Response bodies.
- [ ] **Redis Integration:** Logging structure for every operation (Method, Path, Body, Response).
- [ ] **Health Check:** Mechanism to continuously check if the "real" backend is alive.

### Phase 2: The Brain (Learning Engine)
*Goal: Python service extracting meaningful data from normal traffic.*

- [ ] **Data Consumer (Python):** Worker reading logs from Redis.
- [ ] **Basic Learning (Heuristic):** Storing last successful 200 OK responses for static paths (e.g., `/api/v1/users`).
- [ ] **Dynamic Parameter Analysis:** Clustering responses based on URL query parameters or JSON body.
- [ ] **Model Training (MVP):** Simple "Nearest Neighbor" or rule-based matching. "This request looks most like this response".

### Phase 3: The Ghost (Immortality Mode)
*Goal: Simulation activated when Backend fails.*

- [ ] **Circuit Breaker:** Cut traffic on Go side when backend gives %X errors or times out.
- [ ] **Ghost Handover:** Proxy directs traffic to Python service (or pre-prepared data in Redis).
- [ ] **Semantic Matching:** Analyzing incoming request to generate the most logical "fake" response.
- [ ] **Chaos Testing:** Intentionally shutting down backend to test system behavior.

### Phase 4: Dashboard & SaaS (Productization)
*Goal: Providing visibility to the user.*

- [ ] **Web UI (React/Next.js):** Live traffic flow.
- [ ] **Health Monitor:** Backend uptime and Ghost Mode activation counts.
- [ ] **Traffic Replay:** Feature to replay past traffic.
- [ ] **Anomaly Detection:** Warnings like "Your API is slower than usual" or "Strange requests incoming".

---

## 🛠 Tech Stack

| Component | Technology | Why? |
| --- | --- | --- |
| **Proxy Core** | **Go (Golang)** | High concurrency, low latency, non-blocking IO with `goroutines`. |
| **Cache/Bus** | **Redis** | Very fast R/W, Pub/Sub capability (Go -> Python communication). |
| **AI/ML** | **Python (FastAPI + Scikit-learn/PyTorch)** | Rich ML libraries, fast prototyping. |
| **Database** | **PostgreSQL / TimescaleDB** | Persistent log storage and time-series (analytics). |
| **Frontend** | **Next.js + Tailwind** | Modern, fast dashboard development. |

## 📅 First Step (MVP)
We will first complete **Phase 1** and get a simple Go Proxy up and running. This proxy will forward incoming requests to the "real" server and write the returned response to Redis.
