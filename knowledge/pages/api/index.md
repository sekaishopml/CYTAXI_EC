# API Reference

## Gateway API (`backend/gateway/`)

Endpoints manejados por el API Gateway (Go + chi).

```
GET    /health          → Health check
POST   /v1/geocode      → Geocoding directo (dirección → coordenadas)
POST   /v1/reverse      → Geocoding inverso (coordenadas → dirección)
GET    /v1/routes       → Cálculo de rutas entre dos puntos
GET    /v1/eta          → Tiempo estimado de llegada
```

## Geospatial Engine (`backend/engines/geospatial/`)

Interno, no expuesto directamente. El gateway enruta las peticiones.

### `POST /v1/geocode`
- Input: `{ query: string, limit?: number }`
- Output: `[{ lat, lng, displayName, placeId }]`
- Backend: Nominatim

### `POST /v1/reverse`
- Input: `{ lat: number, lng: number }`
- Output: `{ lat, lng, displayName, address, placeId }`
- Backend: Nominatim

### `GET /v1/routes`
- Query: `origin={lat,lng}&destination={lat,lng}`
- Output: `{ distance, duration, polyline, steps }`
- Backend: OSRM

### `GET /v1/eta`
- Query: `origin={lat,lng}&destination={lat,lng}`
- Output: `{ duration, distance, traffic }`
- Backend: OSRM

## API Futura (Blueprint)

| Servicio | Endpoints |
|----------|-----------|
| Auth | `POST /v1/auth/login`, `POST /v1/auth/register`, `POST /v1/auth/refresh` |
| Customer | `GET/PUT /v1/customers/:id`, `GET /v1/customers/:id/favorites`, `POST /v1/customers/:id/favorites` |
| Driver | `POST /v1/drivers/register`, `PUT /v1/drivers/:id/status`, `GET /v1/drivers/:id` |
| Trip | `POST /v1/trips`, `GET /v1/trips/:id`, `PATCH /v1/trips/:id/status` |
| Pricing | `POST /v1/pricing/estimate`, `GET /v1/pricing/fare/:tripId` |
| Dispatch | `POST /v1/dispatch/assign`, `GET /v1/dispatch/nearby` |
| Payment | `POST /v1/payments/authorize`, `POST /v1/payments/confirm`, `POST /v1/payments/refund` |
| Notification | `POST /v1/notifications/send`, `GET /v1/notifications/:id/status` |
| AI | `POST /v1/ai/interpret`, `POST /v1/ai/recommend` |

## Estándares API

- RESTful con versionado (`/v1/`)
- Request/Response en JSON
- Errores estandarizados: `{ error: { code, message, details } }`
- Paginación: `{ data, pagination: { page, limit, total } }`
- Autenticación: Bearer JWT
- Rate limiting por endpoint
- Health check en cada servicio: `GET /health`
