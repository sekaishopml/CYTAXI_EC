================================================
SPRINT REPORT
================================================

Estado: ✅ Listo para revision
Sprint: 31
Nombre: Public IP Beta Deployment

------------------------------------------------
IP publica
------------------------------------------------
64.176.219.221

------------------------------------------------
Servicios publicados (via Nginx port 80)
------------------------------------------------

| URL | Servicio | Puerto interno |
|-----|----------|----------------|
| http://64.176.219.221/ | Customer MiniWeb | :3000 |
| http://64.176.219.221/driver | Driver Web Portal | :3001 |
| http://64.176.219.221/admin | Admin Dashboard | :3002 |
| http://64.176.219.221/api/v1 | API Gateway | :8000 |
| http://64.176.219.221/docs | Swagger UI | :8081 |
| http://64.176.219.221/health | Health check | :8000 |

------------------------------------------------
Servicios internos (Docker Network, no expuestos)
------------------------------------------------

| Servicio | Puerto | Red |
|----------|--------|-----|
| Trip Engine | 8087 | cytaxi_net |
| Pricing Engine | 8088 | cytaxi_net |
| Payment Engine | 8091 | cytaxi_net |
| Customer Engine | 8085 | cytaxi_net |
| Driver Engine | 8086 | cytaxi_net |
| Notification Engine | 8090 | cytaxi_net |
| Matching Engine | 8089 | cytaxi_net |
| Admin Engine | 8094 | cytaxi_net |
| Analytics Engine | 8093 | cytaxi_net |
| Trust Engine | 8092 | cytaxi_net |
| PostgreSQL | 5432 | cytaxi_net |
| Redis | 6379 | cytaxi_net |

------------------------------------------------
Archivos creados
------------------------------------------------

| Archivo | Descripcion |
|---------|-------------|
| nginx/nginx.conf | Reverse proxy path-based routing: / → miniweb, /driver → driver-web, /admin → dashboard, /api → gateway, /docs → swagger |
| nginx/swagger/index.html | Swagger UI page |
| nginx/swagger/spec.json | OpenAPI 3.0 spec con 20+ endpoints |
| nginx/swagger/Dockerfile | Swagger UI Docker image |
| docker-compose.beta.yml | 18 servicios: nginx + gateway + 3 frontends + swagger + 10 engines + postgres + redis |
| .env.beta | Variables de entorno para IP publica |
| dashboard/package.json | Dashboard React/Next.js |
| dashboard/tsconfig.json | TypeScript config |
| dashboard/next.config.js | Next.js config |
| dashboard/src/pages/index.tsx | Admin dashboard: metrics, service status table, refresh |
| dashboard/Dockerfile | Dashboard Docker image |
| miniweb/Dockerfile | MiniWeb Docker image |
| driver-web/Dockerfile | Driver Web Portal Docker image |
| scripts/deploy_beta.sh | Script de deploy automatizado |

------------------------------------------------
Arquitectura de routing
------------------------------------------------

```
Internet → 64.176.219.221:80 (Nginx)
                │
    ┌───────────┼───────────┬──────────┬──────────┐
    ▼           ▼           ▼          ▼          ▼
   /           /driver     /admin     /api        /docs
 MiniWeb    Driver Web   Dashboard  Gateway    Swagger
 :3000       :3001        :3002      :8000      :8081
                │
    ┌───────────┼──────────────────┐
    ▼           ▼                  ▼
 Trip/Price  Customer/Driver  Payment/Match
```

------------------------------------------------
Arquitectura respetada
------------------------------------------------
DDD            ✅ Engines internos sin exponer
Clean Architecture ✅ Gateway unico punto de entrada
CQRS           ✅ Contratos intactos
Event Driven   ✅ SSE via /api
Zero Trust     ✅ Solo puerto 80 expuesto; engines internos

------------------------------------------------
Riesgos
------------------------------------------------

| Riesgo | Impacto | Mitigacion |
|--------|---------|------------|
| HTTP sin SSL | Alto | Beta controlada; HTTPS en siguiente sprint |
| Sin autenticacion en endpoints | Medio | JWT preparado pero no requerido en beta |
| Todos los servicios en una IP | Bajo | Suficiente para beta con pocos usuarios |

------------------------------------------------
Mejoras futuras (Sprint 32)
------------------------------------------------
- Comprar dominio (cytaxi.ec)
- Configurar HTTPS + Let's Encrypt
- Configurar Cloudflare
- Cambiar nginx.conf de IP a dominio

------------------------------------------------
Commit sugerido
------------------------------------------------
chore(deployment): deploy beta over public IP

------------------------------------------------
