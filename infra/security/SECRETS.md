# CYTAXI Secrets Management

## Environment Separation

| Environment | File | Scope |
|-------------|------|-------|
| Development | .env | Local, debug=true |
| Beta | .env.beta | Public IP testing |
| Production | .env.prod | Live services |

## Secrets Inventory

| Secret | Location | Rotation | Owner |
|--------|----------|----------|-------|
| JWT_SECRET | .env.prod | Every 90 days | Security |
| DB_PASSWORD | .env.prod | Every 90 days | DevOps |
| REDIS_PASSWORD | .env.prod | On change | DevOps |
| STRIPE_API_KEY | .env.payment | On provider rotation | Finance |
| PAYPHONE_API_KEY | .env.payment | On provider rotation | Finance |
| GOOGLE_MAPS_API_KEY | .env.geo | On provider rotation | Engineering |
| WHATSAPP_API_TOKEN | .env.whatsapp | On Meta rotation | Engineering |
| WEBHOOK_SECRET | .env.payment | Every 90 days | Security |

## Rules

1. **NEVER commit .env files** to version control
2. **NEVER hardcode secrets** in source code
3. **Use .env.example** templates with empty values
4. **Rotate on compromise** immediately
5. **Audit access** quarterly
6. **All secrets at least 256 bits** of entropy
