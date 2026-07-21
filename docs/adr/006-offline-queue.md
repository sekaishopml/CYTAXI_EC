# ADR-006: Offline Action Queue with IndexedDB

**Status:** Accepted  
**Date:** 2026-07-19  

## Context
Network connectivity in ride-hailing scenarios is unreliable. Users in tunnels, basements, or areas with poor coverage should be able to continue using the app without data loss.

## Decision
Create `@cytaxi/offline` with:
- `OfflineQueue` class backed by IndexedDB
- Action queue: enqueue → persist → process → complete/fail
- Configurable max retries (default 3) with status tracking
- Network status watcher (online/offline events)
- Auto-processing when connection restores
- Failed item inspection and manual retry

## Consequences
- **Positive:** Actions survive page refresh and tab closure
- **Positive:** Transparent retry with no data loss
- **Positive:** Works with any async action via `onSync` callback
- **Negative:** IndexedDB quota limits (typically 50MB–1GB)
- **Negative:** Queue items can become stale if not processed promptly
