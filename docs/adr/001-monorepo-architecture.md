# ADR-001: Monorepo Architecture with npm Workspaces

**Status:** Accepted  
**Date:** 2026-07-19  
**Decision Makers:** Platform Engineering Team

## Context
CYTAXI has three frontends (travel, driver-web, dashboard), shared packages, and a Go backend. Each frontend previously managed its own copy of UI components, types, and utilities, leading to duplication and drift.

## Decision
Use npm workspaces monorepo with:
- `packages/*` for shared libraries (design-tokens, ui, events, api-client, etc.)
- Each app (travel, driver-web, dashboard) imports shared packages via workspace references
- Single `node_modules` at root

## Consequences
- **Positive:** Single source of truth for types, components, and business logic
- **Positive:** Shared build tooling, ESLint, and TypeScript configs
- **Positive:** Atomic commits across packages and apps
- **Negative:** Requires discipline to keep packages focused and versioned

## Package Map
```
packages/
  api-client/     SDK HTTP con auth, retry, cache
  design-tokens/  Colores, spacing, tipografía, shadows
  events/         Event Bus tipado con auditoría
  fonts/          Inter + JetBrains Mono
  i18n/           Internacionalización (es/en/pt)
  map-engine/     Google Maps utilities
  observability/  Logging, métricas, tracing, Web Vitals
  offline/        Cola de acciones IndexedDB
  realtime/       WebSocket con reconnect y heartbeat
  ride-machine/   Máquina de estados del viaje (12 estados)
  security/       JWT, RBAC, validaciones
  sounds/         Sonidos sintéticos + hápticos
  trust-score/    Sistema de reputación (5 tiers)
  ui/             Componentes React compartidos
```
