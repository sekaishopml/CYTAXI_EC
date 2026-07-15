# API Gateway

Single HTTP entry point for the CYTAXI platform.

## Purpose

The API Gateway is the only HTTP entry point. No frontend (MiniWeb, Dashboard, Driver App) communicates directly with any Engine. All requests go through the Gateway.

## Architecture

```
MiniWeb / Dashboard / Driver App / WhatsApp API
                     │
              API Gateway (port 8000)
                     │
        ┌────────────┼────────────┐
        ▼            ▼            ▼
   /api/v1/trip  /api/v1/pricing  /api/v1/customer ...
```

## Routes

| Path | Backend | Purpose |
|------|---------|---------|
| `/api/v1/customer/*` | Customer Engine | Customer operations |
| `/api/v1/driver/*` | Driver Engine | Driver operations |
| `/api/v1/trip/*` | Trip Engine | Trip lifecycle |
| `/api/v1/pricing/*` | Pricing Engine | Fare and pricing |
| `/api/v1/payment/*` | Payment Engine | Payment operations |
| `/api/v1/notification/*` | Notification Engine | Notifications |
| `/api/v1/admin/*` | Administration Engine | Admin operations |
| `/api/v1/analytics/*` | Analytics Engine | Business intelligence |
| `/api/v1/matching/*` | Matching Engine | Driver matching |
| `/health` | Gateway | Health check |

## Middleware Chain

| Order | Middleware | Description |
|-------|-----------|-------------|
| 1 | Recovery | Panic recovery |
| 2 | Correlation ID | X-Correlation-ID header |
| 3 | CORS | Cross-Origin |
| 4 | Request Logger | Structured request logging |
| 5 | Rate Limiter | Token bucket rate limiting |
| 6 | JWT Auth | Authorization header passthrough |

## OpenAPI

OpenAPI 3.0 base spec available with all engine tags, security schemes (bearerAuth, apiKey), and error schemas.

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `GATEWAY_PORT` | 8000 | Gateway HTTP port |
| `GATEWAY_RATE_LIMIT_RPS` | 100 | Rate limit requests/sec |
| `BACKEND_TRIP` | localhost:8087 | Trip Engine URL |
| `BACKEND_PRICING` | localhost:8088 | Pricing Engine URL |
| `BACKEND_PAYMENT` | localhost:8091 | Payment Engine URL |
| `BACKEND_CUSTOMER` | localhost:8085 | Customer Engine URL |
| `BACKEND_DRIVER` | localhost:8086 | Driver Engine URL |
| `BACKEND_NOTIFICATION` | localhost:8090 | Notification Engine URL |
| `BACKEND_ADMIN` | localhost:8094 | Admin Engine URL |
| `BACKEND_ANALYTICS` | localhost:8093 | Analytics Engine URL |
| `BACKEND_MATCHING` | localhost:8089 | Matching Engine URL |

## Development

```bash
go run ./cmd
```
