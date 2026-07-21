# CYDIGITAL-BLUEPRINT

- **Repo**: https://github.com/sekaishopml/cydigital-blueprint
- **Propósito**: Single Source of Truth para toda la arquitectura de CYTAXI
- **Estado**: ~32% complete (architecture designed, no production code)
- **Idioma**: Español / Inglés

## Filosofía

> ''El Blueprint define la realidad. La implementación sigue la realidad. Si la implementación difiere del Blueprint, el Blueprint gana.''

## Jerarquía de autoridad

```
PROJECT_CONSTITUTION → PRODUCT_BIBLE → ENGINEERING_CONSTITUTION → ARCHITECTURE → DOMAIN → BUSINESS_RULES → EVENTS → DATABASE → API → CODE
```

## Los 30 Capítulos

| Cap | Título |
|-----|--------|
| 1 | AI Bootstrap & Engineering Philosophy |
| 2 | Repository Architecture & Knowledge Flow |
| 3 | Knowledge Architecture & Decision System |
| 4 | Engineering Workflow & Development Lifecycle |
| 5 | Engineering Decision Framework |
| 6 | AI Collaboration Protocol |
| 7 | Domain-Driven Design (DDD) Engineering Standard |
| 8 | Clean Architecture Engineering Standard |
| 9 | Microservices Architecture Standard |
| 10 | Event-Driven Architecture Standard |
| 11 | CQRS & Intelligent Read Models Standard |
| 12 | Contract-First API Engineering Standard |
| 13 | Data Architecture & Persistence Engineering Standard |
| 14 | Caching, Performance & Distributed State Engineering Standard |
| 15 | Observability, Monitoring & Operational Excellence |
| 16 | Security, Identity & Zero Trust Architecture |
| 17 | Engineering Quality, Testing & Validation Standard |
| 18 | Knowledge Management & Documentation Engineering Standard |
| 19 | AI Coding Standards & Engineering Implementation Rules |
| 20 | DevOps, CI/CD & Release Engineering Standard |
| 21 | Artificial Intelligence Architecture & Cognitive System |
| 22 | Conversation Intelligence & WhatsApp Cognitive Engine |
| 23 | Geospatial Intelligence & Location Engine |
| 24 | Mobility Decision Engine & Smart Dispatch Standard |
| 25 | Revenue Intelligence & Fare Decision Engine |
| 26 | Customer Intelligence, Memory & Recommendation Engine |
| 27 | Driver Intelligence, Reputation & Operational Excellence |
| 28 | Trust Intelligence Platform & Risk Management Standard |
| 29 | Autonomous AI Engineering Organization & Multi-Agent Collaboration |
| 30 | CYTAXI Engineering Constitution, Evolution & Long-Term Vision |

## 6 Principios de Ingeniería

1. Documentación antes que implementación
2. Arquitectura antes que optimización
3. Negocio antes que tecnología
4. Consistencia sobre velocidad
5. Reemplazabilidad (toda dependencia externa debe ser reemplazable)
6. Pensamiento a largo plazo (hoy, 6 meses, 1 año, 3 años)

## Clean Architecture (Capítulo 8)

```
Presentation (HTTP/WS/gRPC)
    ↓
Application (Use Cases, Commands, Queries)
    ↓
Domain (Entities, Aggregates, Value Objects, Domain Events)
    ↓
Infrastructure (PostgreSQL, Redis, NATS, Google Maps, LLM)
    ↓
External Services
```

- Las dependencias siempre apuntan hacia adentro
- Domain no depende de nada externo
- Cada caso de uso = una responsabilidad
- Repository Pattern: interfaces en Domain, implementaciones en Infrastructure

## Microservicios (Capítulo 9)

13 microservicios iniciales:

1. **Gateway** — API gateway, rate limiting, routing
2. **Authentication** — login, JWT, roles, permisos
3. **Customer** — perfil, favoritos, preferencias, lealtad
4. **Driver** — conductor, vehículo, documentos, wallet
5. **Trip** — viaje, estados, timeline (corazón del sistema)
6. **Pricing** — tarifas, promociones, surge pricing (determinístico)
7. **Dispatch** — asignación de conductores, ETA ranking
8. **Maps** — geocoding, rutas, tráfico (wrapper de Google Maps)
9. **AI** — NLP, reconocimiento de intención, asistentes
10. **Notification** — WhatsApp, push, SMS, email
11. **Payment** — pagos, wallet, refunds, settlement
12. **Admin** — dashboard, operaciones, reportes
13. **Analytics** — métricas, KPIs, forecasting

## Estado de Implementación

- **Blueprint**: 32% — chapters 1-30 diseñados
- **Producción**: 0% — no hay código de producción
- **Next milestone**: completar documentación del Blueprint antes de escribir código

## Relación con CYTAXI_EC

CYTAXI_EC es el repositorio de implementación que ya tiene código avanzado (50 sprints). El Blueprint se creó después como la fuente de verdad arquitectónica. Idealmente, la implementación debe converger hacia lo que dicta el Blueprint.
