# ADR-007: Payment Provider Adapter Pattern

**Status:** Accepted
**Date:** Sprint 33

**Context:** Need real payment processing (Stripe, Kushki, PayPhone, PayPal) without coupling the Payment Engine to any specific provider.

**Decision:** All payment providers implement the `Provider` interface:
```go
type Provider interface {
    CreatePaymentIntent(ctx, req) (*Response, error)
    Authorize(ctx, req) (*Response, error)
    Capture(ctx, req) (*Response, error)
    Refund(ctx, req) (*Response, error)
    IsAvailable(ctx) bool
}
```
The `ProviderRegistry` selects the active provider at runtime based on the `PAYMENT_PROVIDER` env var. The Payment Engine never knows which concrete provider is being used.

**Consequences:**
- Providers can be swapped without code changes
- Each provider has its own adapter implementing the same interface
- Mock provider enables end-to-end testing without real payments
- New providers (MercadoPago, PayU) can be added as new adapters

---

# ADR-008: Webhook Idempotency

**Status:** Accepted
**Date:** Sprint 33

**Context:** Payment providers send webhooks that may arrive more than once. Duplicate processing would cause double charges or refunds.

**Decision:** The Webhook Receiver maintains an in-memory idempotency cache (processed event IDs). Events with the same ID are rejected on retry. In production, this cache will be backed by Redis with a 24-hour TTL.

**Consequences:**
- Guaranteed exactly-once processing
- Memory-based cache may grow large (mitigated by periodic cleanup)
- Redis backup ensures idempotency survives restarts
- HMAC signature verification ensures webhook authenticity
