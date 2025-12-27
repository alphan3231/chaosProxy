# 👻 Chaos-Proxy

**API'ler İçin "Ölümsüzlük" Katmanı**

Chaos-Proxy, microservice ve API mimarilerinde backend servislerinin çökmesi durumunda bile istemcilere (clients) kesintisiz hizmet vermeyi sağlayan akıllı bir Reverse Proxy'dir.

Normal çalışma zamanında trafiği izler ve öğrenir (The Sentinel & The Brain). Backend çöktüğünde ise "Ghost Mode" devreye girer ve öğrenilen verilerle gerçekçi cevaplar üretmeye devam eder.

## 🚀 Özellikler

- **🪟 Sentinel Proxy (Go):** Yüksek performanslı, düşük gecikmeli reverse proxy.
- **🧠 The Brain (Python + AI):** Trafik analizi ve response modelleme.
- **👻 Ghost Mode:** Backend kesintilerinde otomatik devreye giren akıllı simülasyon.
- **⚡ Redis Backed:** Hızlı veri erişimi ve cache yönetimi.

## 🛠 Kurulum

### Gereksinimler
- Go 1.21+
- Python 3.10+
- Docker & Docker Compose

### Hızlı Başlangıç

1. **Repoyu klonlayın:**
   ```bash
   git clone https://github.com/elliot/chaosProxy.git
   cd chaosProxy
   ```

2. **Altyapıyı ayağa kaldırın (Redis):**
   ```bash
   docker-compose up -d
   ```

3. **Proxy'yi çalıştırın:**
   ```bash
   go run cmd/sentinel/main.go
   ```

## 🗺 Yol Haritası

Detaylı gelişim planı için [ROADMAP.md](ROADMAP.md) dosyasına göz atabilirsiniz.

## 📄 Lisans

MIT
