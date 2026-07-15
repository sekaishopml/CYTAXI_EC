# Integration Layer

Platform-wide integration infrastructure connecting all CYTAXI Engines through events, sagas, and observability.

## Purpose

The Integration Layer connects all Engines without breaking DDD boundaries. It provides the transport and coordination patterns that make the platform behave as a unified system.

## Components

| Component | Description |
|-----------|-------------|
| **EventBus** | Abstract event bus (MemoryBus implemented; NATS/Kafka/RabbitMQ interfaces ready) |
| **SagaCoordinator** | Distributed transaction coordinator with compensation |
| **OutboxPublisher** | Guaranteed event publication pattern |
| **InboxProcessor** | Idempotent event consumption |
| **RetryManager** | Exponential backoff retry with configurable parameters |
| **DeadLetterQueue** | Failed event isolation and inspection |
| **CorrelationManager** | Correlation ID + Trace ID context propagation |
| **TraceManager** | Distributed tracing span management |
| **Observer** | Integration metrics and structured logging |

## Architecture

```
Engine A ─→ OutboxPublisher ─→ EventBus ─→ InboxProcessor ─→ Engine B
              ↑                                  ↑
        RetryManager                        CorrelationID
              ↓                                  ↓
       DeadLetterQueue                    TraceManager
              ↓
       SagaCoordinator (multi-step orchestration)
```

## Integration Flow Example

```
Conversation → Trip Engine → Pricing → Matching → Notification → Payment → Analytics
     ↓             ↓            ↓           ↓           ↓           ↓            ↓
  Outbox        In/Outbox    In/Outbox  In/Outbox   In/Outbox   In/Outbox    Inbox
     └──────────────┴───────────┴─────────┴──────────┴──────────┴──────────────┘
                                    EventBus
```

## EventBus

```go
type Bus interface {
    Publish(ctx, envelope) error
    Subscribe(eventType, handler)
}
```

`MemoryBus` included for testing and single-process deployment.

## Saga Pattern

```go
saga.Register(SagaDefinition{
    Steps: []saga.Step{
        {Name: "create_trip", Execute: createTrip, Compensate: cancelTrip},
        {Name: "charge_customer", Execute: chargeCustomer, Compensate: refundCustomer},
        {Name: "notify_driver", Execute: notifyDriver},
    },
})
saga.Execute(ctx, "trip_request_saga", tripData)
```

## Outbox Pattern

Ensures events are published even if the transaction fails between DB write and event emission.

## Inbox Pattern

Idempotent event consumption — same event processed only once.

## Correlation & Tracing

Every event carries `correlation_id` and `trace_id` through the entire flow.

## Retry Policy

| Parameter | Default |
|-----------|---------|
| Max Retries | 3 |
| Initial Delay | 500ms |
| Max Delay | 30s |
| Backoff Multiplier | 2x |

## Future Broker Adapters

NATS, RabbitMQ, Kafka, Redis Streams — all behind the `BrokerProvider` interface.
