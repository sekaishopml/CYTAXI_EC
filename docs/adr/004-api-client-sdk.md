# ADR-004: Unified API Client SDK

**Status:** Accepted  
**Date:** 2026-07-19  

## Context
Each frontend had its own API call logic with inconsistent error handling, no retry logic, and no token management. This led to duplicated auth logic and brittle error handling.

## Decision
Create `@cytaxi/api-client` with:
- `ApiClient` class wrapping `fetch` with interceptors
- Automatic JWT token injection and 401 refresh
- Exponential backoff retry (configurable codes: 408, 429, 5xx)
- In-memory GET response cache with TTL
- Request ID and trace ID per request
- `AbortController` timeout (default 15s)
- Typed `ApiRequest<T>` and `ApiResponse<T>` interfaces

## Consequences
- **Positive:** Consistent error handling across all frontends
- **Positive:** Token refresh handled automatically at the transport layer
- **Positive:** Request tracing enables debugging distributed flows
- **Negative:** Adds ~4KB to bundle (acceptable for the capability)
- **Negative:** Retry logic can mask transient failures in dev
