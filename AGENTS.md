# CYTAXI — Project instructions

## Stack
- **Frontend**: Next.js 14 (App Router), TypeScript, Leaflet/Google Maps, framer-motion
- **Backend**: Go 1.22+, chi router, OSRM/Nominatim geospatial engine
- **Packages**: `@cytaxi/ride-machine` (state machine), `@cytaxi/map-engine`, `@cytaxi/design-tokens`
- **Design**: Cobalt theme (Hallmark) — `#3b82f6` accent, Space Grotesk + Inter + JetBrains Mono

## Where to find things
- `travel/` — Next.js frontend app
- `backend/engines/geospatial/` — Go geospatial engine (routing, reverse geocode)
- `backend/gateway/` — API gateway
- `packages/` — shared packages (ride-machine, map-engine, design-tokens)
- `knowledge/` — Logseq knowledge graph for persistent project documentation
- `deploy/`, `infra/` — deployment and infrastructure

## Key conventions
- Run `npm run build` from `travel/` before deploying
- Deploy: `systemctl restart cytaxi-travel` (frontend), `systemctl restart cytaxi-geospatial` (backend)
- State machine lives in `packages/ride-machine/`; transitions in `useJourneyEngine.ts`
- All UI text in Spanish
