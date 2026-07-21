# ADR-007: Feature-Slice Frontend Architecture

**Status:** Accepted  
**Date:** 2026-07-19  

## Context
Frontend code was organized by technical role (components/, hooks/, services/) rather than business domain. This made it hard to find code related to a specific feature and encouraged cross-cutting dependencies.

## Decision
Restructure frontends using feature-slice architecture:

```
src/
  entities/     Domain models and types
  features/     Business features (booking, driver, profile, auth, admin)
    booking/
      ui/       State components (PickUpStep, ArrivingState, etc.)
      hooks/    Feature-specific hooks
      services/ Feature-specific API calls
  shared/       Reusable UI components and utilities
    ui/         Generic UI components (Button, Card, etc.)
    utils/      Pure utility functions
    constants/  App-wide constants
  lib/          Third-party library wrappers
  providers/    Context providers
  app/          Next.js App Router pages
```

## Consequences
- **Positive:** Co-location of related code reduces cognitive load
- **Positive:** Features can be developed, tested, and removed independently
- **Positive:** Clear boundaries prevent accidental coupling
- **Negative:** Some utilities may need to be lifted to shared/ if used by multiple features
- **Negative:** Requires discipline to avoid feature boundary violations

## Migration
Backward-compatible re-exports at original paths ensure existing imports continue to work during the transition.
