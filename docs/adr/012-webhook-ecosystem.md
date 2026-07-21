# ADR-012: Webhook Event System for Platform Ecosystem

**Status:** Accepted  
**Date:** 2026-07-19  

## Context
Enterprise integrations require outgoing webhooks for real-time event notification. Partners and internal services need to subscribe to trip, driver, and payment events with guaranteed delivery.

## Decision
Create `@cytaxi/webhooks` with:

### Event Types
- `trip.created`, `trip.started`, `trip.completed`, `trip.cancelled`
- `driver.assigned`, `driver.arrived`
- `payment.confirmed`, `payment.failed`
- `rating.submitted`

### WebhookEndpoint
- URL, method (POST/GET/PUT/PATCH), custom headers, secret per endpoint
- Per-endpoint event subscription list
- Configurable retry count and timeout (default 3 retries, 5s timeout)
- Active/inactive toggle

### Delivery
- HMAC-SHA256 signature in `X-Webhook-Signature` header
- `X-Webhook-Timestamp` and `X-Webhook-Event` headers
- Exponential backoff retry: 2s, 4s, 8s, 16s, max 30s
- Delivery status tracking: pending → retrying → delivered/failed

### WebhookDispatcher
- Register/unregister endpoints
- `dispatch(event, payload)` — sends to all matching endpoints
- Per-delivery status tracking with attempt count

### Signature Verification
```typescript
signature = sha256(`${timestamp}.${JSON.stringify(payload)}`, secret)
Header: X-Webhook-Signature: sha256=<hex>
```

## Consequences
- **Positive:** Enables partner integrations and event-driven architecture
- **Positive:** Retry with backoff handles transient failures
- **Positive:** Signature verification prevents tampering
- **Negative:** No idempotency key support yet (risk of duplicate delivery)
- **Negative:** In-memory endpoint storage (need DB for production)

## Future
- Idempotency keys (`Idempotency-Key` header)
- Webhook delivery dashboard in Super Admin
- Dead letter queue for failed deliveries
- Rate limiting per endpoint
