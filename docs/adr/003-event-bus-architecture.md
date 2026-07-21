# ADR-003: Typed Event Bus for Cross-Cutting Communication

**Status:** Accepted  
**Date:** 2026-07-19  

## Context
Components and services needed a decoupled way to communicate across the platform. Prop drilling and callback chains were making the code brittle. Trip lifecycle events needed to trigger analytics, UI updates, sounds, and haptics simultaneously.

## Decision
Create `@cytaxi/events` with:
- Singleton `EventBus` class with typed emit/on/off
- 25+ typed event names with per-event payload interfaces
- Built-in audit log (last 1000 entries, in-memory)
- `useEvent` and `useEmitEvent` React hooks
- `getGlobalBus()` singleton for cross-module access

## Consequences
- **Positive:** Complete decoupling of event producers and consumers
- **Positive:** Audit log enables debugging and replay
- **Positive:** TypeScript ensures type-safe payloads per event
- **Negative:** Overuse can make data flow hard to trace
- **Negative:** Events are fire-and-forget (no return value)

## Event Categories
- **Journey:** LOCATION_DETECTED, ROUTE_CALCULATED, FARE_ESTIMATED, ...
- **Trip:** TRIP_STARTED, TRIP_COMPLETED, TRIP_CANCELLED, ...
- **Payment:** PAYMENT_CONFIRMED, ...
- **Analytics:** SCREEN_VIEW, BUTTON_CLICK, ...
- **Error:** ERROR_OCCURRED
