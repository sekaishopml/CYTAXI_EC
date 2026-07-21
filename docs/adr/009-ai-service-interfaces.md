# ADR-009: AI Service Interfaces for Matching, ETA, Pricing & Fraud

**Status:** Accepted  
**Date:** 2026-07-19  

## Context
The platform needs intelligent matching, ETA prediction, dynamic pricing, and fraud detection. The existing matching/pricing engines use heuristic algorithms. A clean interface is needed to swap rule-based logic for ML models later.

## Decision
Define an `AIService` interface in `@cytaxi/ai` with four capabilities:

### Interface
```typescript
interface AIService {
  matchDrivers(request, candidates): MatchingResult
  predictETA(origin, dest, timestamp?): ETAPrediction
  calculateDynamicPrice(baseParams, demand, supply, hour, zone): DynamicPrice
  assessFraud(tripData): FraudAssessment
}
```

### Rule-Based Implementation (`RuleBasedAIService`)
- **Matching**: Weighted multi-factor scoring (distance 0.35, ETA 0.25, rating 0.15, accept rate 0.10, trust score 0.10, zone 0.05)
- **ETA**: Haversine distance / 25 km/h average speed with random traffic factor
- **Pricing**: Base + distance + time, demand/supply surge multipliers (0.8–3.0), night/airport surcharges
- **Fraud**: Flag-based risk scoring (short distance, high fare, cash payment, late hour)

### Provider Abstraction
- `AIConfig.provider: "rule_based" | "ml_service"`
- Future ML provider connects to external inference endpoint
- No breaking changes when switching providers

## Consequences
- **Positive:** Clear contract enables parallel development of ML models
- **Positive:** Rule-based system works immediately without ML infrastructure
- **Positive:** Fraud detection catches obvious patterns without training data
- **Negative:** Rule-based ETA ignores real-time traffic conditions
- **Negative:** Matching weights need tuning per-zone (future: per-tenant config)

## Future
- Connect `RuleBasedAIService` to the Go matching/pricing engines
- ML model endpoints for surge prediction and ETA
- Real-time traffic data integration
- Per-tenant AI model fine-tuning
