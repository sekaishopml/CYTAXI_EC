# Arquitectura Actual (CYTAXI_EC Implementation)

A diferencia del Blueprint (que es pura especificación), CYTAXI_EC ya tiene 50 sprints de implementación con código funcional.

## Stack Tecnológico Real

### Frontend (travel/)
- **Framework**: Next.js 14 (App Router)
- **Lenguaje**: TypeScript
- **Mapas**: Leaflet (OpenStreetMap) + OpenStreetMap Nominatim
- **Animaciones**: framer-motion
- **Estados**: `packages/ride-machine/` (state machine)
- **Tema**: Cobalt Hallmark — `#3b82f6` accent, Space Grotesk/Inter/JetBrains Mono
- **UI**: Glassmorphism temprano → reemplazado por hairline borders

### Backend (backend/)
- **Gateway**: Go + chi router
- **Geospatial**: Go — OSRM + Nominatim (en `backend/engines/geospatial/`)
- **Base de datos**: PostgreSQL

### Paquetes Compartidos (packages/)
- `@cytaxi/ride-machine` — state machine para viajes
- `@cytaxi/map-engine` — lógica de mapas
- `@cytaxi/design-tokens` — tokens de diseño

### Infraestructura
- **Terraform**: `infra/terraform/` — GCP, Cloud Run, CloudSQL, VPC
- **Docker**: `deploy/docker/` — docker-compose para dev
- **CI/CD**: GitHub Actions en `.github/workflows/`

## Diferencias con el Blueprint

| Aspecto | Blueprint (ideal) | Implementación actual |
|---------|-------------------|----------------------|
| **# Servicios** | 13 microservicios | 2 (gateway + geospatial) |
| **Arquitectura** | Clean Architecture completa | Capas básicas |
| **Event Bus** | NATS JetStream | No implementado |
| **CQRS** | Separación reads/writes | No implementado |
| **DDD** | Domain events, aggregates | Incipiente |
| **IA** | Servicio de IA dedicado | No implementado |
| **WhatsApp** | Canal primario | No implementado |
| **Pagos** | Payment service | No implementado |
| **Auth** | Servicio dedicado + Zero Trust | Básica |

## Estado por Servicio (Blueprint vs Realidad)

| Servicio | Blueprint | CYTAXI_EC |
|----------|-----------|-----------|
| Gateway | ✅ Especificado | ✅ Implementado (chi router) |
| Maps/Geospatial | ✅ Especificado | ✅ Implementado (OSRM + Nominatim) |
| Trip | ✅ Especificado | ⚠️ Parcial (state machine en travel) |
| Pricing | ✅ Especificado | ❌ No implementado |
| Dispatch | ✅ Especificado | ❌ No implementado |
| Auth | ✅ Especificado | ❌ No implementado |
| Customer | ✅ Especificado | ❌ No implementado |
| Driver | ✅ Especificado | ❌ No implementado |
| AI | ✅ Especificado | ❌ No implementado |
| Notification | ✅ Especificado | ❌ No implementado |
| Payment | ✅ Especificado | ❌ No implementado |
| Admin | ✅ Especificado | ❌ No implementado |
| Analytics | ✅ Especificado | ❌ No implementado |

## Data Flow Actual (Miniweb)

1. Usuario abre travel → state machine inicia en `pickup_select`
2. Mapa detecta GPS → reverse geocode → selección de pickup
3. Usuario confirma pickup → busca destino → selecciona → confirma
4. Trip request → driver matching → tracking → payment → rating

## Próximos Pasos (Convergencia)

La implementación actual debe converger hacia la arquitectura del Blueprint:
1. Adoptar Clean Architecture en los servicios existentes
2. Implementar event bus (NATS)
3. Separar en microservicios por dominio
4. Implementar CQRS
5. Agregar autenticación y autorización
6. Integrar WhatsApp como canal primario
