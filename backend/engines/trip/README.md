# Trip Engine

Core operational domain for the CYTAXI platform.

## Purpose

The Trip Engine administers the complete lifecycle of a trip — from request through completion or cancellation. It does NOT decide which driver gets assigned (that's the Mobility Decision Engine), calculate prices, process payments, or interact with map providers.

## Architecture

DDD ✓ Clean Architecture ✓ CQRS ✓ Event Driven ✓ Contract First ✓ Zero Trust ✓

## Domain

### Aggregate
- **Trip** — Root aggregate managing the full trip lifecycle (11 statuses)

### Status Lifecycle
```
requested → created → searching → driver_assigned → accepted → arrived → started → (paused ↔ resumed) → completed
                                                                                                    → cancelled
```

### Entities
- **Passenger** — Trip passenger identity
- **Stop** — Pickup, waypoints
- **Destination** — Final destination
- **TripAssignment** — Driver assignment record
- **TripTimeline** — Audit trail of trip events

### Value Objects
TripID, CustomerID, DriverID, Coordinates, Distance, ETA, Money, TripStatus, StopID

## Events

15 domain events: TripRequested, TripCreated, DriverAssigned, DriverUnassigned, TripAccepted, TripRejected, DriverArrived, TripStarted, TripPaused, TripResumed, TripCompleted, TripCancelled, DestinationChanged, StopAdded, StopRemoved

## CQRS

**Commands:** CreateTrip, AssignDriver, UnassignDriver, AcceptTrip, RejectTrip, StartTrip, PauseTrip, ResumeTrip, CompleteTrip, CancelTrip, AddStop, RemoveStop, ChangeDestination

**Queries:** GetTrip, GetTripHistory, GetTripTimeline, GetDriverTrips, GetCustomerTrips, GetActiveTrips

## API

| Method | Path | Description |
|--------|------|-------------|
| GET | /health | Health check |
| GET | /trips/{trip_id} | Get trip |
| GET | /customers/{customer_id}/trips | Customer trip history |
| GET | /trips/active | Active trips |

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `TRIP_PORT` | 8087 | HTTP server port |

## Development

```bash
go run ./cmd/trip
```
