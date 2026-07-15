# MiniWeb OpenAPI Client

Customer MiniWeb communicates exclusively with the CYTAXI API Gateway.

Base URL: `http://localhost:8000/api/v1`

## Endpoints

| Method | Path | Module | Description |
|--------|------|--------|-------------|
| GET | /trip/trips/{id} | Trip | Get trip |
| GET | /trip/customers/{id}/trips | Trip | Trip history |
| GET | /pricing/fares/{id} | Pricing | Get fare |
| GET | /pricing/trips/{id}/fares | Pricing | Fare history |
| GET | /customer/customers/{id}/profile | Customer | Profile |
| GET | /customer/customers/{id}/preferences | Customer | Preferences |
| GET | /customer/customers/{id}/favorites | Customer | Favorite places |
| GET | /notification/notifications/{id} | Notification | Get notification |
| GET | /notification/recipients/{id}/notifications | Notification | History |
| GET | /health | Gateway | Health check |
