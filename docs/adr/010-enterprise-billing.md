# ADR-010: Enterprise Billing & Subscription System

**Status:** Accepted  
**Date:** 2026-07-19  

## Context
Multi-tenant operation requires per-tenant billing. Each tenant needs a subscription plan with driver-based pricing, trip commissions, coupons, and tax support.

## Decision
Create `@cytaxi/billing` with:

### Plans & Pricing
- 4 tiers: Free ($0/3 drivers), Starter ($99/50 drivers), Business ($299/200 drivers), Enterprise ($999/unlimited)
- Monthly and yearly billing (yearly saves ~15%)
- Per-additional-driver pricing above included count
- Per-trip commission percentage

### Subscription Model
- `Subscription` tracks status (active/trialing/past_due/canceled/incomplete)
- Period tracking with start/end dates
- Trial support with `trialEnd` date
- Driver count monitoring for overage billing

### Invoicing
- `Invoice` with line items (subscription, additional drivers, trip commission, adjustments, coupons)
- Status lifecycle: draft → open → paid → void
- Tax rate support per country/region

### Calculation
```typescript
subtotal = planAmount + max(0, driverCount - included) * perDriverPrice
discount = coupon value (percentage or fixed)
tax = (subtotal - discount) * taxRate
total = subtotal - discount + tax
```

## Consequences
- **Positive:** Self-service billing enables tenant onboarding without sales
- **Positive:** Driver-based pricing scales with tenant growth
- **Positive:** Coupon/tax support handles promotions and compliance
- **Negative:** No payment gateway integration yet (Stripe/PayPal needed)
- **Negative:** No dunning (automated retry for failed payments)

## Future
- Stripe Connect integration for payment processing
- Automated invoice generation and email delivery
- Dunning workflow for past-due subscriptions
- Usage-based billing (per-trip commission automation)
