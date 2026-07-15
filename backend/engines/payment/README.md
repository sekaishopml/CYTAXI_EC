# Payment Engine

Payment lifecycle management for the CYTAXI platform.

## Purpose

The Payment Engine is the exclusive owner of all financial operations. No other Engine may handle money directly — everything routes through Payment Engine.

## Architecture

DDD ✓ Clean Architecture ✓ CQRS ✓ Event Driven ✓ Contract First ✓ Zero Trust ✓

## Domain

### Aggregates
- **Payment** — Payment lifecycle (pending→authorized→captured→refunded)
- **Settlement** — Driver payout calculation and processing
- **Refund** — Refund management (pending→processing→completed)

### Entities
- **Transaction** — Immutable transaction record (authorization, capture, refund, settlement, payout)
- **Wallet** — Customer/driver balance with Credit/Debit
- **Invoice** — Customer invoice with line items
- **Receipt** — Payment confirmation receipt
- **PaymentMethod** — Stored payment method details
- **PaymentIntent** — Pre-authorization intent with expiration

### Value Objects
PaymentID, TransactionID, WalletID, RefundID, SettlementID, InvoiceNumber, ReceiptNumber, PaymentMethodID, Money, Percentage, PaymentStatus(7), PaymentMethodType(5), SettlementStatus(4), RefundStatus(4)

## CQRS

**Commands:** CreatePayment, AuthorizePayment, CapturePayment, CancelPayment, RefundPayment, CreateSettlement, PayDriver, GenerateInvoice, GenerateReceipt

**Queries:** GetPayment, GetPayments, GetTransaction, GetWallet, GetRefund, GetSettlement, GetInvoice, GetReceipt

## PaymentGateway Adapter

```go
type PaymentGateway interface {
    Authorize(ctx, amount, method, ref) (*GatewayResult, error)
    Capture(ctx, gatewayRef, amount) (*GatewayResult, error)
    Refund(ctx, gatewayRef, amount) (*GatewayResult, error)
    Cancel(ctx, gatewayRef) error
}
```

**GatewayRegistry** maps payment methods to gateways. Designed for: Stripe, PayPhone, DataFast, Kushki, PayPal, MercadoPago.

## Events

| Event | Description |
|-------|-------------|
| `payment.created` | Payment created |
| `payment.authorized` | Payment authorized |
| `payment.captured` | Payment captured |
| `payment.refund_created` | Refund initiated |
| `payment.refund_completed` | Refund completed |
| `payment.settlement_created` | Settlement for driver created |
| `payment.driver_paid` | Driver payout completed |
| `payment.invoice_generated` | Invoice generated |
| `payment.receipt_generated` | Receipt generated |

## API

| Method | Path | Description |
|--------|------|-------------|
| GET | /health | Health check |
| GET | /payments/{id} | Get payment |
| GET | /wallets/{owner_id} | Get wallet |

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `PAYMENT_PORT` | 8091 | HTTP server port |

## Development

```bash
go run ./cmd/payment
```
