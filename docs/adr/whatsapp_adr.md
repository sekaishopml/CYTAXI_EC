# ADR-009: WhatsApp Adapter Pattern

**Status:** Accepted
**Date:** Sprint 34

**Context:** CYTAXI's primary channel is WhatsApp. Need to integrate with Meta's WhatsApp Business Cloud API without coupling the domain to Meta.

**Decision:** All WhatsApp communication goes through the `Provider` interface:
```go
type Provider interface {
    SendMessage(ctx, msg) (*Response, error)
    SendTemplate(ctx, to, template) (*Response, error)
    MarkAsRead(ctx, messageID) error
    IsAvailable(ctx) bool
}
```
- `MetaCloudAdapter` handles real API calls (when token configured)
- `MockProvider` enables development/testing without Meta
- Webhook `Receiver` handles incoming messages with signature validation + idempotency

**Consequences:**
- Zero coupling to Meta in domain code
- Mock provider enables full dev/testing without WhatsApp access
- Switching from Mock to Meta requires only env var change
- Message types: text, location, image, audio, video, document, interactive buttons, lists, templates
