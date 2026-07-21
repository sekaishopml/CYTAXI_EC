# Pilot Deployment Checklist — Journey Engine

## 1. Environment & Config
- [ ] `.env.production` has correct `API_URL`, `WS_URL`, `MAPS_API_KEY`
- [ ] All feature flags (`FEATURE_LLM`, `FEATURE_SSE`) match production intent
- [ ] Next.js `output` in `next.config.js` set to `standalone`
- [ ] Node.js version pinned (`.nvmrc` / `engines.node`)

## 2. Build Verification
- [ ] `npm run build` succeeds with 0 errors
- [ ] `npx tsc --noEmit` passes with 0 errors
- [ ] No `any` type violations in new FASE 6/7 code
- [ ] `next.config.js` `experimental` flags reviewed for stability
- [ ] Bundle size checked: `next build` output shows each page's size

## 3. State Management
- [ ] Journey Engine FSM transitions verified for all 12 states
- [ ] Session persistence: reload at `pickup_select`, `input`, `confirm` restores state
- [ ] Session TTL (24h) evicts stale sessions
- [ ] Active trips (`searching` → `destination`) NOT persisted (intentional)
- [ ] Offline queue stores pending actions when `navigator.onLine === false`

## 4. Real-Time Tracking (SSE)
- [ ] `subscribeToTrip` reconnects with exponential backoff (1s–30s)
- [ ] Max 10 retry attempts before giving up
- [ ] Connection lost mid-trip shows user-facing error banner
- [ ] `trip_completed` event tears down EventSource cleanly
- [ ] SSE endpoint (`/api/v1/trip/ws`) reachable behind Cloudflare/Nginx

## 5. Error Handling
- [ ] `<ErrorBoundary>` wraps entire page — catches render errors
- [ ] `getDerivedStateFromError` shows "Recargar página" button
- [ ] Errors are reported to `/api/v1/telemetry`
- [ ] API client retries on failure (up to 2 retries, exponential backoff)
- [ ] API client times out after 8s by default
- [ ] Network errors show banner with dismiss action
- [ ] `trackError()` recorded for all caught exceptions

## 6. Telemetry
- [ ] State transitions recorded with duration
- [ ] API latency tracked (`search_places`, `calculate_route`, `estimate_fare`, `request_trip`)
- [ ] SSE errors and reconnection events recorded
- [ ] Telemetry flush every 30s to `/api/v1/telemetry`
- [ ] Max 200 buffered events

## 7. Performance
- [ ] All callback handlers in `useJourneyEngine` use `useCallback`
- [ ] `MemoEngine` marker/layer updates batched (no per-frame re-render)
- [ ] BottomSheet content uses `React.memo` where applicable
- [ ] Search debounce at 400ms
- [ ] No unnecessary `setState` on every SSE tick (only on meaningful change)

## 8. Accessibility
- [ ] All interactive elements have `aria-label`
- [ ] Error/offline banners have `role="alert"` and `aria-live="assertive"`
- [ ] Map pin and GPS button labelled for screen readers
- [ ] Keyboard navigation: tab order covers all interactive elements
- [ ] Touch targets ≥ 44px per WCAG 2.5.8
- [ ] Color contrast ratios meet WCAG AA (4.5:1 normal, 3:1 large)

## 9. Map (MapEngine / MapController)
- [ ] Map renders at `pickup_select`, `input`, `confirm`, `searching`, `driver_found`, `arriving`, `arrived`, `in_progress`, `destination`
- [ ] Pin overlay shows at `pickup_select`, `input`
- [ ] Map is interactive only at `pickup_select`, `input`
- [ ] GPS button works on mobile (high accuracy, 8s timeout)
- [ ] Reverse geocoding validates street access (river/sea check)
- [ ] Map click in `input` mode selects destination

## 10. Mobile
- [ ] `visualViewport` resize handler correctly detects keyboard open/close
- [ ] Bottom sheet scrolls to focused input when keyboard opens
- [ ] `safe-area-inset-bottom` applied to navbar
- [ ] Touch actions: `touchAction: "none"` on container, `"auto"` on map area
- [ ] Bottom sheet snap points (40%, 70%, 90%) work on iOS Safari
- [ ] No horizontal scroll on 375px viewport

## 11. API Endpoints (Production)
- [ ] `GET /api/v1/geo/search?q=...` — places search
- [ ] `POST /api/v1/geo/route` — route calculation
- [ ] `POST /api/v1/pricing/estimate` — fare estimation
- [ ] `POST /api/v1/trip/request` — trip creation
- [ ] `GET /api/v1/trip/ws?trip_id=...` — SSE tracking endpoint
- [ ] `POST /api/v1/matching/start` — driver matching
- [ ] `POST /api/v1/telemetry` — telemetry ingestion

## 12. Post-Deploy Smoke Test
- [ ] Open app → map loads → GPS detects location → address displays
- [ ] Enter destination → route + fare appear → confirm shows details
- [ ] Request trip → "Buscando conductor" animation → driver assigned
- [ ] Driver accepted → SSE connects → position updates arrive
- [ ] Complete flow: arrive → start → destination → payment → rating → completed
- [ ] Kill network mid-trip → offline banner → restore → actions replayed
- [ ] Hard refresh mid-flow → state restored from localStorage

## 13. Rollback Plan
- [ ] Previous build zip available at `https://sekaishopec.com/cydigital-v1.0.0.zip`
- [ ] Nginx `default` configured to swap symlink between releases
- [ ] Cloudflare health check configured (expects 200 on `/`)
- [ ] `git tag` applied at current HEAD for rollback reference
