# Arquitectura Objetivo (Blueprint Vision)

Basada en CYDIGITAL-BLUEPRINT capítulos 7-13.

## Principios Arquitectónicos

1. **Domain-Driven Design** — cada microservicio posee un dominio de negocio
2. **Clean Architecture** — dependencias hacia adentro, dominio aislado
3. **Event-Driven Architecture** — eventos de negocio vía NATS JetStream
4. **CQRS** — comandos vs queries, read models separados
5. **Contract-First** — API diseñada antes de implementar
6. **Zero Trust** — toda comunicación autenticada y autorizada

## Capas por Microservicio

```
cmd/                       → entrypoint
internal/
  domain/                  → entities, aggregates, value objects, domain events
  application/             → use cases, commands, queries, application services
  infrastructure/          → repos, external APIs, persistence
  presentation/            → HTTP handlers, DTOs, middleware
configs/                   → configuración por entorno
tests/                     → tests unitarios, de integración, e2e
docs/                      → documentación del servicio
```

## Comunicación entre Servicios

- **Síncrona** (HTTP/gRPC): solo cuando se necesita respuesta inmediata
  - Autenticación
  - Creación de viaje
  - Cálculo de precio
- **Asíncrona** (NATS JetStream): para consistencia eventual
  - Notificaciones
  - Analítica
  - Recomendaciones
  - Historial
  - Lealtad

## Data Ownership

Cada servicio es dueño de su propio esquema de base de datos. Nadie consulta tablas de otro servicio directamente. El intercambio de datos ocurre solo via APIs, eventos, o read models.

## Reglas de Negocio Permanentes

- El pricing es determinístico (la IA nunca decide precios)
- El dispatch es determinístico (la IA nunca asigna conductores)
- La IA nunca aprueba pagos
- La IA nunca modifica estado de negocio directamente
- Las reglas de negocio siempre ganan

## Aislamiento de Falla

La falla de un servicio no debe detener la plataforma:
- Notification falla → Trip continúa
- Analytics falla → Payments continúa  
- AI falla → Customer puede solicitar viaje sin IA

## Seguridad

- Zero Trust Architecture
- JWT + refresh tokens
- Rate limiting por servicio
- Auditoría de todas las operaciones
- Cifrado en tránsito y en reposo
