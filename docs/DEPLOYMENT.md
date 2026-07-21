# CYTAXI — Despliegue y Arquitectura de Subdominios (sekaishopec.com)

Fuente de verdad operativa. Actualizado 2026-07-20.

## DNS (Cloudflare, Proxy ON → 64.176.219.221)
Todos los subdominios apuntan a la misma IP; nginx enruta por `Host`.

| Subdominio | App | Puerto local | systemd unit | Puerto gateway |
|---|---|---|---|---|
| `sekaishopec.com` / `www` | landing (nginx root) | — | nginx | — |
| `travel.sekaishopec.com` | App clientes (travel) | 3000 | `cytaxi-travel` | — |
| `driver.sekaishopec.com` | Driver Web | 3002 | `cytaxi-driver-web` | — |
| `admin.sekaishopec.com` | Admin Dashboard | 3003 | `cytaxi-dashboard` | — |
| `status.sekaishopec.com` | Status Page | 3001 | `cytaxi-status` | — |
| `bot.sekaishopec.com` | Bot/WhatsApp | (pendiente) | — | — |
| `api.sekaishopec.com` | API Gateway | 8000 | `cytaxi-gateway` | — |

## Backend engines (detrás del gateway, vía `api.sekaishopec.com`)
El gateway (`/api/v1/<engine>[/s]/...`) hace reverse proxy al engine en su puerto.
Acepta forma singular y plural (`/api/v1/trip/...` y `/api/v1/trips/...`).

| Engine | Puerto | Unit | Estado |
|---|---|---|---|
| trip | 8087 | `cytaxi-trip` | ✅ (flujo completo implementado) |
| pricing | 8088 | (manual) | ✅ |
| matching | 8089 | (manual) | ✅ |
| geospatial | 8082 | (manual) | ✅ |
| customer | 8085 | `cytaxi-customer` | ✅ desplegado |
| driver | 8086 | `cytaxi-driver` | ✅ desplegado |
| payment | 8091 | `cytaxi-payment` | ✅ desplegado |
| conversation / notification / admin / analytics | — | — | ⚠️ no desplegados (código incompleto) |

## Flujo core de viaje (end-to-end verificado)
`POST /api/v1/trip/request` → `searching`
`POST /api/v1/trip/accept` → `accepted`
`POST /api/v1/trip/start` → `in_progress` (arrived→started)
`POST /api/v1/trip/location` → ok
`POST /api/v1/trip/complete` → `completed`
`GET  /api/v1/trip/trips/{id}` → detalle
`GET  /api/v1/trip/customers/{id}/trips` → historial

## nginx
Config activa en la VM: `/etc/nginx/sites-available/wildcard.sekaishopec.com.conf`
+ map de puertos en `/etc/nginx/nginx.conf` (`map $subdomain $upstream_port`).
El repo tiene `nginx/nginx.conf` como referencia fiel de esa config.
`nginx/cytaxi.conf` es LEGACY (dominio viejo cytaxi.app) — ignorar.

## Comandos útiles
```bash
# Estado de todos los servicios
systemctl status cytaxi-*

# Ver logs de un engine
journalctl -u cytaxi-trip -f

# Rebuild + restart un servicio frontend
cd /home/CYTAXI_EC/travel && npm run build && systemctl restart cytaxi-travel

# Rebuild + restart un engine Go
cd /home/CYTAXI_EC/backend/engines/trip && \
  go build -o /opt/cytaxi/bin/cytaxi-trip ./cmd/trip && systemctl restart cytaxi-trip

# Rebuild + restart gateway
cd /home/CYTAXI_EC/backend/gateway && \
  go build -o /opt/cytaxi/bin/cytaxi-gateway ./cmd && systemctl restart cytaxi-gateway
```

## Cloudflare (dev)
- **Cache Rule `CYTAXI dev: bypass cache all subdomains` ACTIVA** (ruleset `4459840af3e2478cb893ff52a8475890`, fase `http_request_cache_settings`).
  Expresión: todos los subdominios de sekaishopec.com → `cache:false`, `browser_ttl:bypass`.
  Verificado: respuestas externas muestran `cf-cache-status: DYNAMIC` (no cachea).
- Purge ejecutado 2026-07-20 (cache limpio).
- **Development Mode**: NO activado vía API (el token no tiene `Zone Settings:Edit`).
  Como la Cache Rule ya hace bypass, no es necesario. Si se quiere, activarlo manualmente
  en dashboard (Overview → Development Mode → On) o crear token con `Zone Settings:Edit`.
