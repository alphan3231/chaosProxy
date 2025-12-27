# Chaos-Proxy: Yol Haritası (Roadmap)

## 🎯 Proje Vizyonu
**Chaos-Proxy**, API servisleri için bir "ölümsüzlük" katmanıdır. Backend servisleri ve istemciler (client) arasına girer, normal zamanda trafiği izleyip öğrenir, backend çöktüğünde ise yapay zeka destekli "Ghost Mode" ile trafiği simüle ederek kesintisiz hizmet algısı yaratır.

---

## 🏗 Mimari Genel Bakış
Sistem 3 ana bileşenden oluşur:
1.  **The Sentinel (Go):** Yüksek performanslı Reverse Proxy. Trafiği karşılar, Redis'e loglar ve gerekirse Ghost response döner.
2.  **The Memory (Redis):** Canlı istek/cevap verilerinin ve AI modellerinin (veya kurallarının) tutulduğu hızlı önbellek.
3.  **The Brain (Python):** Arka planda çalışır. Redis'teki veriyi analiz eder, pattern'leri öğrenir ve "Ghost Response" modellerini oluşturur.

---

## 🚀 Geliştirme Fazları

### Faz 1: The Sentinel (Temel Proxy ve İzleme)
*Hedef: Trafiği aktaran ve kaydeden çalışan bir Go Proxy.*

- [ ] **Proje Kurulumu:** Go modül yapısı ve temel dizinler.
- [ ] **Reverse Proxy Çekirdeği:** `net/http/httputil` kullanarak temel proxy mantığı.
- [ ] **Middleware Yapısı:** Request/Response body'sini yakalamak için hook noktaları.
- [ ] **Redis Entegrasyonu:** Her işlem (Method, Path, Body, Response) için loglama yapısı.
- [ ] **Health Check:** Backend'in hayatta olup olmadığını sürekli kontrol eden mekanizma.

### Faz 2: The Brain (Öğrenme Motoru)
*Hedef: Normal trafikten anlamlı veri çıkaran Python servisi.*

- [ ] **Veri Tüketici (Python):** Redis'ten logları okuyan worker.
- [ ] **Basit Öğrenme (Heuristic):** Static pathler için (örn: `/api/v1/users`) son başarılı 200 OK cevaplarını saklama.
- [ ] **Dinamik Parametre Analizi:** URL query parametreleri veya JSON body'ye göre değişen cevapları gruplama (cluster).
- [ ] **Model Eğitimi (MVP):** Basit bir "Nearest Neighbor" veya kural tabanlı eşleştirme. "Bu request'e en çok şu response benziyor".

### Faz 3: The Ghost (Ölümsüzlük Modu)
*Hedef: Backend çöktüğünde devreye giren simülasyon.*

- [ ] **Circuit Breaker:** Go tarafında backend %X hata verdiğinde veya timeout olduğunda trafiği kesme.
- [ ] **Ghost Handover:** Proxy'nin trafiği Python servisine (veya Redis'teki ön-hazırlanmış verilere) yönlendirmesi.
- [ ] **Semantic Matching:** Gelen isteği analiz edip, en mantıklı "sahte" cevabı üretme.
- [ ] **Chaos Testing:** Bilerek backend'i kapatıp sistemin davranışını test etme.

### Faz 4: Dashboard & SaaS (Ürünleştirme)
*Hedef: Kullanıcıya görünürlük sağlama.*

- [ ] **Web UI (React/Next.js):** Canlı trafik akışı.
- [ ] **Health Monitor:** Backend uptime ve Ghost Mode devreye girme sayıları.
- [ ] **Traffic Replay:** Geçmiş trafiği tekrar oynatma özelliği.
- [ ] **Anomaly Detection:** "API'niz normalden yavaş" veya "Garip requestler geliyor" uyarıları.

---

## 🛠 Teknoloji Yığını (Tech Stack)

| Bileşen | Teknoloji | Neden? |
| --- | --- | --- |
| **Proxy Core** | **Go (Golang)** | Yüksek concurrency, düşük latency, `goroutines` ile non-blocking IO. |
| **Cache/Bus** | **Redis** | Çok hızlı yazma/okuma, Pub/Sub yeteneği (Go -> Python haberleşmesi). |
| **AI/ML** | **Python (FastAPI + Scikit-learn/PyTorch)** | Zengin ML kütüphaneleri, hızlı prototipleme. |
| **Database** | **PostgreSQL / TimescaleDB** | Kalıcı log saklama ve zaman serisi (analytics) için. |
| **Frontend** | **Next.js + Tailwind** | Modern, hızlı dashboard geliştirme. |

## 📅 İlk Adım (MVP)
Öncelikle **Faz 1**'i tamamlayıp, basit bir Go Proxy'yi ayağa kaldıracağız. Bu proxy, gelen isteği "gerçek" sunucuya iletecek ve dönen cevabı Redis'e yazacak.
