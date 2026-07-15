================================================
SPRINT REPORT
================================================

Estado: ✅ Listo para revision
Sprint: 35
Nombre: Identity & Authentication Platform

------------------------------------------------
Funcionalidades implementadas
------------------------------------------------

Plataforma de autenticacion centralizada:
1. JWT + Refresh Tokens con HMAC-SHA256
2. 2 auth providers: Email/Password + OTP (WhatsApp)
3. RBAC con 4 roles + permisos granulares
4. Session management con revocacion + blacklist
5. Auth APIs: register, login, logout, refresh, me, forgot-password
6. Auth provider registry con Adapter Pattern

------------------------------------------------
Archivos creados
------------------------------------------------

| Archivo | Descripcion |
|---------|-------------|
| infrastructure/providers/auth.go | 2 providers (EmailPassword + OTP) + Registry + AuthMethod types |
| infrastructure/auth/manager.go | Token Manager (JWT sign/validate/refresh/revoke) + RoleManager (RBAC) |
| cmd/auth_server.go | AuthServer: 7 endpoints (register/login/refresh/logout/me/forgot/health) |
| .env.auth | JWT_SECRET + TTL config |
| docs/adr/auth_adr.md | ADR-010: Centralized Authentication |

------------------------------------------------
APIs implementadas
------------------------------------------------

| Method | Path | Auth | Descripcion |
|--------|------|------|-------------|
| POST | /auth/register | No | Registro (email+password / OTP) |
| POST | /auth/login | No | Login + tokens |
| POST | /auth/refresh | Refresh | Renovar access token |
| POST | /auth/logout | Bearer | Cerrar sesion + revocar refresh |
| GET | /auth/me | Bearer | Perfil del usuario autenticado |
| POST | /auth/forgot-password | No | Solicitar reset |
| GET | /auth/health | No | Health check |

------------------------------------------------
RBAC Roles
------------------------------------------------

| Role | Permisos |
|------|----------|
| customer | trip:create, trip:read, profile:read, profile:write, payment:read, payment:create |
| driver | trip:read, trip:accept, trip:reject, trip:start, trip:finish, profile:read, vehicle:manage, availability:manage |
| operator | trip:read, trip:cancel, driver:read, payment:read, refund:create, support:read |
| admin | * (todos) |

------------------------------------------------
Token config
------------------------------------------------

| Param | Default | Descripcion |
|-------|---------|-------------|
| JWT_ACCESS_TTL | 900s (15min) | Access token TTL |
| JWT_REFRESH_TTL | 604800s (7d) | Refresh token TTL |

------------------------------------------------
Arquitectura respetada
------------------------------------------------
DDD            ✅ Trust Engine dueño del dominio
Clean Architecture ✅ domain → infrastructure/auth
Adapter Pattern ✅ Provider Registry (email+password, OTP, Google OAuth stub)
Contract First ✅ APIs estandar JWT
Zero Trust     ✅ Secretos en .env, never exposed

------------------------------------------------
Commit sugerido
------------------------------------------------
feat(auth): implement centralized identity & authentication platform

------------------------------------------------
