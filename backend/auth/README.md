# Auth

Authentication and authorization for CYTAXI backend.

## Usage

```go
import "github.com/sekaishopml/cytaxi/backend/auth"

// Authenticate
authenticator := auth.NewAuthenticator("my-secret")
principal, err := authenticator.Authenticate("token...")

// Generate tokens
token, err := auth.GenerateToken("my-secret", principal, 24*time.Hour)

// Authorize
authorizer := auth.NewAuthorizer(map[string][]string{
    "admin": {"*"},
    "driver": {"trip:read", "trip:write"},
})
allowed := authorizer.Authorize(principal, "trip", "read")
```

## How it works

- HMAC-SHA256 JWT implementation (stdlib, zero external deps).
- Claims: `sub` (user ID), `role`, `scope`, `exp`, `iat`.
- Role-based authorization with resource:action rules.
