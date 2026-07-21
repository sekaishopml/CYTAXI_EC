# ADR-002: Ride State Machine for Trip Lifecycle

**Status:** Accepted  
**Date:** 2026-07-19  

## Context
Trip flow was managed by a series of boolean flags and enums (TripState) scattered across components. This made it impossible to guarantee valid state transitions and led to inconsistent UI states.

## Decision
Replace ad-hoc boolean state management with a finite state machine in `@cytaxi/ride-machine`:

- 12 states covering the full passenger journey
- 15 typed events that trigger transitions
- Transition table as a single source of truth
- `transitionRide()` pure function: (current, event) → next state | null

## Consequences
- **Positive:** Impossible states are eliminated by construction
- **Positive:** New states (arriving, destination, payment, rating) added without changing existing transitions
- **Positive:** Audit trail via Event Bus events emitted on every transition
- **Positive:** Animation config per state lives alongside the machine definition
- **Negative:** Adding a new state requires updating both the machine and transition table

## State Map
```
pickup_select → input → confirm → searching → driver_found → arriving → arrived → in_progress → destination → payment → rating → completed
```
