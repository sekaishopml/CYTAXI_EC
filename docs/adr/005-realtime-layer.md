# ADR-005: WebSocket Realtime Layer

**Status:** Accepted  
**Date:** 2026-07-19  

## Context
Tracking driver location, trip status updates, and notifications required real-time communication. Multiple ad-hoc EventSource connections were being created per trip.

## Decision
Create `@cytaxi/realtime` with:
- `RealtimeClient` class wrapping a single WebSocket connection
- Automatic reconnection with exponential backoff (1s → 30s max)
- Heartbeat/ping-pong every 30s with 5s timeout
- Typed `subscribe(type, handler)` and `send(type, payload)`
- Wildcard `subscribeToAll()` for debugging
- Connection state tracking: disconnected → connecting → connected → reconnecting

## Consequences
- **Positive:** Single persistent connection, no connection per trip
- **Positive:** Heartbeat detects stale connections within 35s
- **Positive:** Reconnection preserves subscription state
- **Negative:** Requires WebSocket proxy config in Nginx (already done)
- **Negative:** Falls back to polling if WebSocket not available
