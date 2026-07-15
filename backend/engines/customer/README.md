# Customer Engine

Customer domain management for the CYTAXI platform.

## Purpose

The Customer Engine is the exclusive owner of customer data. No other Engine may directly modify customer information; all interactions happen through contracts and events.

## Architecture

```
Conversation Engine / Other Engines
       ↓ (contracts + events)
Customer Engine
  ├── Profile Service
  ├── Preference Service
  ├── Favorite Places Service
  └── Customer Context Service
       ↓
  Repository interfaces (PostgreSQL, Redis)
       ↓
  Event Publisher (NATS, RabbitMQ)
```

## Domain Entities

| Entity | Description |
|--------|-------------|
| `Customer` | Core customer aggregate (id, phone, name, status) |
| `Profile` | Extended profile (name, email, avatar, language, timezone) |
| `Preferences` | Travel preferences (vehicle type, baby seat, wheelchair, payment) |
| `FavoritePlace` | Saved places (home, work, gym, etc.) with coordinates |
| `CustomerContext` | Lightweight context for real-time decisions (preferences, recent trips) |

## Services

| Service | Description |
|---------|-------------|
| **Profile Service** | Get/update customer profile |
| **Preference Service** | Get/update travel preferences with partial updates |
| **Favorite Place Service** | CRUD for favorite places |
| **Customer Context Service** | Lightweight context for decision engines |

## Events

| Event | Description |
|-------|-------------|
| `customer.created` | New customer registered |
| `customer.updated` | Customer data changed |
| `customer.blocked` | Customer blocked |
| `customer.profile_updated` | Profile updated |
| `customer.preferences_updated` | Preferences changed |
| `customer.favorite_place_added` | New favorite place saved |
| `customer.favorite_place_removed` | Favorite place removed |
| `customer.context_changed` | Customer context changed |

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | /health | Health check |
| GET | /customers/{customer_id}/profile | Get profile |
| GET | /customers/{customer_id}/preferences | Get preferences |
| GET | /customers/{customer_id}/favorites | Get favorite places |
| GET | /customers/{customer_id}/context | Get customer context |

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `CUSTOMER_PORT` | 8085 | HTTP server port |

## Development

```bash
go run ./cmd/customer
```
