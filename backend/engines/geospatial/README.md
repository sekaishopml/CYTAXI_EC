# Geospatial Engine

Geographic services abstraction for the CYTAXI platform.

## Purpose

Provides a unified interface for all map-related operations (geocoding, routing, places, zones) while keeping the rest of the platform independent of any specific map provider (Google Maps, OpenStreetMap, Mapbox, etc.).

## Architecture

```
Application (Conversation Engine, etc.)
       ↓ (GeocodeInputPort, RouteInputPort, PlaceInputPort, ZoneInputPort)
Geospatial Engine Use Cases
       ↓ (GeospatialProvider interface)
Provider Adapters (Google Maps, OSM, Mapbox)
       ↓
External APIs
```

## Provider Adapters

| Provider | Package | Status |
|----------|---------|--------|
| Google Maps | `infrastructure/googlemaps` | Stub |
| OpenStreetMap (Nominatim + OSRM) | `infrastructure/osm` | Stub |
| Mapbox | `infrastructure/mapbox` | Stub |

All adapters implement the `GeospatialProvider` interface:

```go
type GeospatialProvider interface {
    Geocoder     // Geocode, ReverseGeocode
    DirectionFinder  // FindRoute, GetDistanceMatrix
    PlaceSearcher    // SearchPlaces, Autocomplete, PlaceDetails
    Name() string
}
```

## Services

| Service | Description |
|---------|-------------|
| **Geocode** | Address → Coordinates |
| **Reverse Geocode** | Coordinates → Address |
| **Directions** | Route between points with waypoints |
| **Distance Matrix** | Multi-origin to multi-destination distances |
| **Places Search** | Text search, autocomplete, details |
| **Zone Service** | Geo-fencing (polygon/circle), zone management |

## Domain Types

| Type | Description |
|------|-------------|
| `Coordinates` | Latitude/Longitude with distance calculation |
| `Address` | Structured address with components |
| `Route` | Route with distance, duration, polyline, steps |
| `DistanceMatrix` | Multi-point distance/duration matrix |
| `Place` | Point of interest with metadata |
| `GeoFence` | Polygon or circle geographic boundary |
| `Zone` | Named geographic area with business rules |

## Cache

- `GeospatialCache` wraps any `Cache` implementation (Redis, in-memory).
- `GetOrFetch` pattern: returns cached data or fetches from provider.
- Configurable TTL per cache entry.

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `GEOSPATIAL_PORT` | 8082 | HTTP server port |
| `GEOSPATIAL_PROVIDER` | google_maps | Map provider |
| `GEOSPATIAL_API_KEY` | "" | API key for provider |
| `GEOSPATIAL_CACHE_TTL` | 300 | Cache TTL in seconds |
| `APP_ENV` | development | Environment |
| `LOG_LEVEL` | info | Log level |

## Events

| Event | Description |
|-------|-------------|
| `geospatial.geocode_requested` | Geocode operation started |
| `geospatial.geocode_completed` | Geocode completed with coordinates |
| `geospatial.route_requested` | Route calculation started |
| `geospatial.route_found` | Route calculated successfully |
| `geospatial.provider_error` | Provider API error occurred |

## Development

```bash
go run ./cmd/geospatial
```
