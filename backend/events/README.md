# Events

Event bus for domain events in CYTAXI backend.

## Usage

```go
import "github.com/sekaishopml/cytaxi/backend/events"

bus := events.NewMemoryBus()

bus.Subscribe("trip.created", func(evt events.DomainEvent) error {
    fmt.Println("trip created:", evt.Payload)
    return nil
})

bus.Publish(events.DomainEvent{
    Type:    "trip.created",
    Version: 1,
    Payload: tripData,
})
```

## Available implementations

- `NewMemoryBus()` — In-memory, synchronous. Ideal for testing and single-process deployments.
- For production, implement `events.Bus` with NATS or RabbitMQ.
