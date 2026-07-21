# Backend

## Stack

- **Lenguaje**: Go 1.22+
- **Router**: chi
- **Base de datos**: PostgreSQL
- **Geospatial**: OSRM (routing) + Nominatim (geocoding)
- **Infra**: Docker, Terraform (GCP)

## Servicios Actuales

### Gateway (`backend/gateway/`)
- API Gateway
- Routing de peticiones
- Punto de entrada único

### Geospatial Engine (`backend/engines/geospatial/`)
- Geocoding directo: `/v1/geocode`
- Geocoding inverso: `/v1/reverse`
- Cálculo de rutas: `/v1/routes`
- ETA: `/v1/eta`

## Estructura de Directorios (Objetivo Clean Architecture)

```
backend/
  gateway/
    cmd/
    internal/
    configs/
  engines/
    geospatial/
      cmd/
      internal/
        domain/           → entities, value objects (Location, Route, Distance)
        application/      → use cases (GeocodeLocation, CalculateRoute, ReverseGeocode)
        infrastructure/   → OSRM client, Nominatim client, PostgreSQL repos
        presentation/     → HTTP handlers, DTOs
      configs/
      tests/
```

## Pendiente (Blueprint)

- Implementar NATS JetStream para event bus
- Separar en microservicios por dominio
- Agregar CQRS
- Implementar Zero Trust
- Servicio de IA dedicado

## Base de Datos

- PostgreSQL
- Migraciones en cada servicio
- Cada servicio dueño de su esquema
- Sin consultas directas entre servicios
