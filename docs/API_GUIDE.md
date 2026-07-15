# CYTAXI API Guide v1.0.0-rc1

Base URL: `https://api.cytaxi.app/api/v1`

## Authentication

All endpoints require JWT Bearer token:
```
Authorization: Bearer <token>
```

## Endpoints

### Customer Engine (`/customer`)
| Method | Path | Description |
|--------|------|-------------|
| GET | /customers/{id}/profile | Get customer profile |
| GET | /customers/{id}/preferences | Get preferences |
| GET | /customers/{id}/favorites | Get favorite places |

### Driver Engine (`/driver`)
| Method | Path | Description |
|--------|------|-------------|
| GET | /drivers/{id} | Get driver profile |
| GET | /drivers/{id}/vehicles | Get vehicles |
| GET | /drivers/{id}/availability | Get availability status |
| GET | /driver/requests | Get pending trip requests |
| POST | /driver/accept | Accept trip request |
| POST | /driver/reject | Reject trip request |

### Trip Engine (`/trip`)
| Method | Path | Description |
|--------|------|-------------|
| GET | /trips/{id} | Get trip details |
| POST | /trip/request | Create trip request |
| POST | /trip/start | Start trip (driver) |
| POST | /trip/location | Update driver location |
| POST | /trip/finish | Complete trip |
| GET | /trip/ws?trip_id= | SSE tracking stream |

### Pricing Engine (`/pricing`)
| Method | Path | Description |
|--------|------|-------------|
| POST | /pricing/estimate | Get fare estimate |
| GET | /pricing/fares/{id} | Get fare details |

### Payment Engine (`/payment`)
| Method | Path | Description |
|--------|------|-------------|
| POST | /payments | Create payment |
| POST | /payments/confirm | Confirm payment |
| POST | /payments/refund | Refund payment |
| GET | /payments/{id} | Get payment |
| GET | /payments/history | Payment history |
| GET | /receipts/{id} | Get receipt |
| GET | /payments/driver/{id}/earnings | Driver earnings |

### Matching Engine (`/matching`)
| Method | Path | Description |
|--------|------|-------------|
| POST | /matching/start | Start driver search |
| GET | /matching/{id}/candidates | Get candidates |

### Notification Engine (`/notification`)
| Method | Path | Description |
|--------|------|-------------|
| GET | /notifications/{id} | Get notification |
| GET | /recipients/{id}/notifications | Notification history |

### Admin Engine (`/admin`)
| Method | Path | Description |
|--------|------|-------------|
| GET | /admin/roles | List roles |
| GET | /admin/feature-flags | Feature flags |

### Analytics Engine (`/analytics`)
| Method | Path | Description |
|--------|------|-------------|
| GET | /analytics/dashboard | Dashboard metrics |
| GET | /analytics/metrics | Business metrics |

## Common Headers

| Header | Description |
|--------|-------------|
| `X-Correlation-ID` | Track request across services |
| `X-Request-ID` | Unique request identifier |
| `Authorization` | JWT Bearer token |
| `Content-Type` | application/json |

## Error Format

```json
{
  "error": "description",
  "details": "additional info",
  "correlation_id": "corr_xxx"
}
```

## Rate Limits

- Default: 100 requests/second per IP
- Authenticated: 500 requests/second per user
- Payment endpoints: 10 requests/second per user
