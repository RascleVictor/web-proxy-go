# 🚀 Go Reverse Proxy

A lightweight and high-performance reverse proxy written in **Go**, built for speed, observability, and reliability.

---

## 📦 Features
- ⚡ Optimized HTTP transport (connection pooling, timeouts)
- 🧠 Structured logging with [Zap](https://github.com/uber-go/zap)
- 🧱 Middlewares: Logging, Recovery, CORS, Metrics
- 📊 Prometheus metrics at `/metrics`
- 🔒 Error handling and safe panic recovery

---

## 🛠️ Requirements
- Go **1.21+**
- A valid `config.yaml` file in `./config/`

Example `config/config.yaml`:
```yaml
server:
  port: 8080
  target: "https://example.com"

logging:
  level: "info"

git clone https://github.com/yourusername/proxy-web-go.git
cd proxy-web-go
go mod tidy
go run ./cmd/proxy
