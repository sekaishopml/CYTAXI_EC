================================================
SPRINT REPORT
================================================

Estado: ✅ Listo para revision
Sprint: 37
Nombre: Trust, Reputation & Safety Platform

------------------------------------------------
Funcionalidades implementadas
------------------------------------------------

Sistema de confianza y reputacion:
1. Ratings (1-5 estrellas) con comentarios y moderacion
2. Trust Score dinamico (0-100) con formula ponderada
3. Incidentes con 6 tipos + 4 severidades + workflow de resolucion
4. Apelaciones con revision y aprobacion/rechazo
5. Historial de cambios de Trust Score
6. 5 niveles de confianza: excellent, good, fair, poor, critical

------------------------------------------------
Archivos creados
------------------------------------------------

| Archivo | Descripcion |
|---------|-------------|
| infrastructure/trust/manager.go | Trust Manager: ratings, incidents, appeals, trust score recalculation engine |
| cmd/trust_server.go | TrustServer: 6 endpoints (rate/report/trust-score/appeal/resolve-incident/incidents) |

------------------------------------------------
APIs implementadas
------------------------------------------------

| Method | Path | Descripcion |
|--------|------|-------------|
| POST | /ratings | Calificar viaje (1-5 + comentario) |
| POST | /reports | Reportar incidente |
| GET | /trust-score/{user_id} | Consultar Trust Score |
| POST | /appeals | Apelar incidente |
| POST | /incidents/resolve | Resolver incidente |
| GET | /incidents | Listar incidentes (abiertos o por usuario) |

------------------------------------------------
Trust Score Formula
------------------------------------------------

```
Score = RatingScore*0.4 + CompletionRate*0.2 + IncidentScore*0.4

Incident deductions:
  low: -5, medium: -15, high: -30, critical: -50
```

Levels: excellent (85+), good (70+), fair (50+), poor (30+), critical (<30)

------------------------------------------------
Incident Types
------------------------------------------------

| Tipo | Descripcion |
|------|-------------|
| safety | Riesgo de seguridad |
| behavior | Mal comportamiento |
| fraud | Fraude |
| vehicle | Problema con vehiculo |
| payment | Problema de pago |
| other | Otro |

------------------------------------------------
Arquitectura respetada
------------------------------------------------
DDD            ✅ Trust Engine dueño unico de reputacion
Clean Architecture ✅ domain → infrastructure/trust
Event Driven   ✅ RatingSubmitted, ReportCreated, TrustScoreUpdated, etc.
Zero Trust     ✅ Moderacion + apelaciones + audit trail

------------------------------------------------
Commit sugerido
------------------------------------------------
feat(trust): implement trust, reputation & safety platform

------------------------------------------------
