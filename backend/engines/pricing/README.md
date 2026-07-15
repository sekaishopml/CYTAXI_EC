# Pricing Engine

Fare calculation engine for the CYTAXI platform.

## Purpose

The Pricing Engine is the exclusive owner of all pricing, fare calculation, taxes, commissions, promotions, and coupons. No other Engine may calculate prices directly.

## Architecture

DDD ✓ Clean Architecture ✓ CQRS ✓ Event Driven ✓ Contract First ✓ Zero Trust ✓

## Domain

### Aggregates
- **Fare** — Complete fare with breakdown (base, distance, time, waiting, night/demand surcharge, discounts, taxes, commission, driver earnings)
- **Promotion** — Time-limited or rule-based discounts

### Entities
- **Coupon** — Single-use or limited-use discount codes
- **Tax** — Country/region specific taxes
- **Commission** — Platform commission
- **Discount** — General discounts

### Value Objects
Money, FareID, PromotionID, CouponCode, Percentage, TaxRate, FareComponents

### Fare Components (11-item breakdown)
```
BaseFare + DistanceFare + TimeFare + WaitingFare + NightSurcharge + DemandSurcharge
= Subtotal
- Promotion - Coupon
+ Tax
+ Commission
= Total → DriverEarnings
```

## CQRS

**Commands:** CalculateFare, ApplyPromotion, ApplyCoupon, RemoveCoupon, CalculateTaxes, CalculateCommission

**Queries:** GetFare, GetFareHistory, GetPromotions, GetCoupons, PreviewFare

## Events

| Event | Description |
|-------|-------------|
| `pricing.fare_calculated` | Fare calculated |
| `pricing.fare_updated` | Fare updated |
| `pricing.promotion_applied` | Promotion applied |
| `pricing.promotion_removed` | Promotion removed |
| `pricing.coupon_applied` | Coupon applied |
| `pricing.coupon_rejected` | Coupon rejected |
| `pricing.taxes_calculated` | Taxes calculated |
| `pricing.commission_calculated` | Commission calculated |

## API

| Method | Path | Description |
|--------|------|-------------|
| GET | /health | Health check |
| GET | /fares/{fare_id} | Get fare |
| GET | /trips/{trip_id}/fares | Fare history |
| GET | /promotions | Active promotions |

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `PRICING_PORT` | 8088 | HTTP server port |

## Development

```bash
go run ./cmd/pricing
```
