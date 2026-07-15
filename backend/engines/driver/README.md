# Driver Engine

Driver domain management for the CYTAXI platform.

## Purpose

The Driver Engine is the exclusive owner of all driver-related data. No other Engine may directly modify driver information; all interactions happen through contracts and events.

## Domain Entities

| Entity | Description |
|--------|-------------|
| `Driver` | Core driver aggregate (id, phone, name, status, rating) |
| `Vehicle` | Vehicles registered per driver (plate, type, capacity, baby seat, wheelchair) |
| `License` | Driving licenses with expiry dates |
| `Availability` | Real-time driver availability with location tracking |
| `Document` | Uploaded documents (license, registration, insurance, ID, background check) |
| `Preferences` | Driver preferences (max distance, min fare, auto-accept) |
| `Capabilities` | Driver capabilities (baby seat, wheelchair, XL, premium, cash/card) |

## Driver Status Lifecycle

```
pending → approved → online ↔ offline
pending → rejected
approved → suspended
```

## Architecture

DDD ✓ Clean Architecture ✓ CQRS ✓ Event Driven ✓ Contract First ✓ Zero Trust ✓

## Events

| Event | Description |
|-------|-------------|
| `driver.created` | New driver registered |
| `driver.approved` | Driver approved |
| `driver.rejected` | Driver rejected |
| `driver.online` | Driver went online |
| `driver.offline` | Driver went offline |
| `driver.suspended` | Driver suspended |
| `driver.vehicle_updated` | Vehicle details changed |
| `driver.license_updated` | License updated |
| `driver.availability_changed` | Availability status changed |
| `driver.capabilities_changed` | Driver capabilities updated |
| `driver.preferences_updated` | Preferences changed |
| `driver.document_uploaded ` | New document uploaded |

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | /health | Health check |
| GET | /drivers/{driver_id} | Get driver |
| GET | /drivers/{driver_id}/vehicles | Get vehicles |
| GET | /drivers/{driver_id}/licenses | Get licenses |
| GET | /drivers/{driver_id}/availability | Get availability |

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `DRIVER_PORT` | 8086 | HTTP server port |

## Development

```bash
go run ./cmd/driver
```
