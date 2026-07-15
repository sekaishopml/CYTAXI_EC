# Notification Engine

Multi-channel notification delivery for the CYTAXI platform.

## Purpose

The Notification Engine is the exclusive owner of all outbound communications. No other Engine may send messages directly — everything routes through Notification Engine.

## Architecture

DDD ✓ Clean Architecture ✓ CQRS ✓ Event Driven ✓ Contract First ✓ Zero Trust ✓

## Supported Channels

| Channel | Kind | Description |
|---------|------|-------------|
| WhatsApp | `whatsapp` | WhatsApp outgoing messages |
| Push | `push` | iOS/Android push notifications |
| Email | `email` | Email notifications |
| SMS | `sms` | Text messages |
| WebSocket | `websocket` | Real-time in-app delivery |
| In-App | `in_app` | In-app notification center |

## ChannelProvider Adapter

Each channel has a `ChannelProvider` adapter implementing:

```go
type ChannelProvider interface {
    Name() string
    Kind() ChannelType
    Send(ctx context.Context, to string, body string) (*SendResult, error)
    IsAvailable(ctx context.Context) bool
}
```

## Domain

### Aggregates
- **Notification** — Core notification with 7-state lifecycle
- **Delivery** — Per-channel delivery tracking with retries
- **NotificationTemplate** — Templated messages with variable substitution

### Status Lifecycle
```
pending → queued → sending → sent → delivered
                                ↓
                              failed → (retry) → sending
                                ↓
                              cancelled
```

### Entities
- **Recipient** — User with device tokens
- **DeliveryAttempt** — Individual send attempt record
- **NotificationPreference** — Per-channel user preferences

### Value Objects
NotificationID, RecipientID, TemplateID, ChannelType (6 types), Priority (4 levels), DeliveryStatus (6 states), Locale, AttemptID

## CQRS

**Commands:** CreateNotification, QueueNotification, SendNotification, RetryNotification, CancelNotification, UpdateDeliveryStatus

**Queries:** GetNotification, GetNotificationHistory, GetDeliveryStatus, GetPendingNotifications, GetTemplates

## Events

| Event | Description |
|-------|-------------|
| `notification.created` | Notification created |
| `notification.queued` | Added to delivery queue |
| `notification.sent` | Sent to provider |
| `notification.delivered` | Confirmed delivered |
| `notification.failed` | Delivery failed |
| `notification.retried` | Retry attempted |
| `notification.cancelled` | Cancelled before send |
| `notification.template_updated` | Template modified |

## API

| Method | Path | Description |
|--------|------|-------------|
| GET | /health | Health check |
| GET | /notifications/{id} | Get notification |
| GET | /recipients/{id}/notifications | History |
| GET | /notifications/templates | Templates |

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `NOTIFICATION_PORT` | 8090 | HTTP server port |

## Development

```bash
go run ./cmd/notification
```
