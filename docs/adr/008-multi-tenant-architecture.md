# ADR-008: Multi-Tenant Architecture

**Status:** Accepted  
**Date:** 2026-07-19  

## Context
CYTAXI must support multiple cooperatives and enterprises (tenants) from a single infrastructure. Each tenant needs isolated configuration, branding, plans, and permissions.

## Decision
Adopt **logical isolation** (tenant_id column on all data) rather than database-per-tenant:

### Backend (`backend/tenant/`)
- `Tenant` model with ID, name, slug, plan, branding, features, limits
- `Repository` interface with `InMemoryRepository` seed (2 default tenants)
- `context.Context` propagation via `WithTenant`/`FromContext`/`MustFromContext`
- `domain` (subdomain) and `X-Tenant-ID` header resolution

### Gateway middleware (`TenantResolver`)
- Resolves tenant from `X-Tenant-ID` header → subdomain → default
- Sets `X-Tenant-ID` response header
- Rejects requests for inactive tenants (HTTP 403)
- Admin CRUD endpoints at `/admin/tenants/*`

### Frontend (`@cytaxi/multi-tenant`)
- `TenantContext` + `useTenant()` hook for React apps
- `useFeature()` gate feature access by tenant plan
- `useCanManageTenants()` for super admin UI
- `planLabels` and `planLimits` define per-plan constraints

### Tenant Resolution Order
1. `X-Tenant-ID` HTTP header
2. `tenant_id` query parameter
3. Subdomain (first DNS label)
4. Default tenant fallback

## Consequences
- **Positive:** Single deployment serves all tenants
- **Positive:** Shared infrastructure reduces ops cost
- **Positive:** No data leaks via tenant_id scoping
- **Negative:** Requires all queries to include tenant_id (must audit)
- **Negative:** In-memory repo is for MVP; production needs PostgreSQL

## Future
- Migrate `InMemoryRepository` to PostgreSQL with tenant_id index
- Add tenant-scoped API keys
- Per-tenant rate limiting and feature flags
