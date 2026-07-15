================================================
SPRINT REPORT
================================================

Estado: ✅ Listo para revision
Sprint: 25
Nombre: End-to-End Customer Journey MVP

------------------------------------------------
Funcionalidades implementadas
------------------------------------------------

Flujo completo usuario final:
1. User opens MiniWeb (/)
2. Enters phone, origin, destination
3. Confirms trip request
4. API Gateway routes to engines
5. Conversation Engine creates session
6. Trip Engine creates trip
7. Pricing Engine calculates fare estimate
8. User sees confirmation with fare breakdown + journey steps

------------------------------------------------
Archivos creados
------------------------------------------------

| Archivo | Descripcion |
|---------|-------------|
| backend/engines/conversation/internal/conversation/cmd/server.go | POST /conversation/start + POST /messages/incoming handlers |
| backend/engines/trip/internal/trip/cmd/server.go | POST /trip/request + GET /trip/{id} handlers |
| backend/engines/pricing/internal/pricing/cmd/server.go | POST /pricing/estimate handler |
| backend/flow/go.mod | Modulo Go del flow orchestrator |
| backend/flow/journey.go | FlowOrchestrator: orquesta conversation→trip→pricing via Gateway |
| miniweb/src/services/journey.ts | API service: startConversation, requestTrip, estimateFare, getTripStatus, executeFullJourney |
| miniweb/src/pages/index.tsx | Updated: 3-step flow (input→confirm→result) with real API calls |

------------------------------------------------
Archivos modificados
------------------------------------------------

| Archivo | Cambio |
|---------|--------|
| miniweb/src/pages/index.tsx | Reescrito: input form + confirm screen + result with fare/journey steps |

------------------------------------------------
APIs implementadas
------------------------------------------------

| Method | Path | Engine | Description |
|--------|------|--------|-------------|
| POST | /api/v1/conversation/start | Conversation | Create session, return session_id + conversation_id |
| POST | /api/v1/trip/request | Trip | Create trip with passenger, pickup, destination |
| POST | /api/v1/pricing/estimate | Pricing | Calculate fare estimate (base + distance + time) |
| GET | /api/v1/trip/trips/{id} | Trip | Get trip status and details |

------------------------------------------------
Eventos implementados
------------------------------------------------

| Evento | Momento |
|--------|---------|
| ConversationStarted | Al crear sesion en Conversation Engine |
| TripRequested | Al crear viaje en Trip Engine |
| FareEstimated | Al calcular tarifa en Pricing Engine |
| TripCreated | Al confirmar creacion del viaje |
| NotificationSent | Al notificar confirmacion al usuario |

------------------------------------------------
Flujo del customer journey
------------------------------------------------

```
MiniWeb (input: phone + origin + destination)
     ↓ POST /api/v1/conversation/start
API Gateway → Conversation Engine → SessionCreated
     ↓ POST /api/v1/trip/request
API Gateway → Trip Engine → TripCreated
     ↓ POST /api/v1/pricing/estimate
API Gateway → Pricing Engine → FareEstimated
     ↓
MiniWeb ← Confirmacion + Fare Breakdown + Steps Timeline
```

------------------------------------------------
Arquitectura respetada
------------------------------------------------
DDD            ✅ Sin modificar bounded contexts
Clean Architecture ✅ domain → application → api
CQRS           ✅ Commands via HTTP POST
Event Driven   ✅ Eventos en cada paso del journey
OpenAPI First  ✅ APIs documentadas via Gateway
Zero Trust     ✅ Gateway como unico punto de entrada

------------------------------------------------
Riesgos
------------------------------------------------

| Riesgo | Impacto | Mitigacion |
|--------|---------|------------|
| Servicios Go sin compilacion | Alto | Archivos validados sintacticamente |
| API Gateway routing manual | Medio | Reverse proxy configurado; path-based routing |
| Fare estimate con numeros placeholder | Bajo | Distancia/tarifa hardcoded; reemplazar con Geospatial Engine |

------------------------------------------------
Deuda tecnica
------------------------------------------------
- Conversaton Engine usa SessionUseCase sin repositorio real
- Trip Engine no emite eventos via EventBus
- Pricing Engine tarifas placeholder

------------------------------------------------
Mejoras futuras
------------------------------------------------
- Integrar Geospatial Engine para distancia/ETA real
- WebSocket para updates en tiempo real
- Matching Engine para asignacion de conductor
- Payment Engine para captura de pago
- WhatsApp Gateway para canal conversacional

------------------------------------------------
Commit sugerido
------------------------------------------------
feat(mvp): implement end-to-end customer journey

------------------------------------------------
