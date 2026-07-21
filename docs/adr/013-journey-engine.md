# ADR-013: Journey Engine y Continuidad del Viaje

**Estado:** Implementado
**Fecha:** 2026-07-19

## Contexto
El flujo del pasajero sufría pérdidas de estado, transiciones bruscas y comportamientos
inconsistentes del frontend. La máquina de estados (`ride-machine`) existía pero no
controlaba la interfaz. No había un motor de mapa desacoplado ni persistencia robusta.

## Decisión
Implementar un **Journey Engine** como única fuente de verdad, con:

1. **JourneyEngine** (`@cytaxi/ride-machine`) — Clase con estado interno, métodos `send(event)`,
   `goTo(state)`, `onTransition(listener)`, y auto-detección de dirección forward/back.

2. **MapEngine** (`@cytaxi/map-engine`) — Controlador único que orquesta marcadores A/B,
   polylines, animación de ruta, cámara, zoom y marcador del conductor. Reemplaza la
   lógica de dibujo inline que estaba en `MapView.tsx`.

3. **useJourneyEngine** — Hook que envuelve `JourneyEngine` + API calls + persistencia
   (localStorage) + tracking SSE. Reemplaza `useTripFlow` + `useAnimatedState`.

4. **BottomSheet** — Componente reutilizable con animaciones por estado (Framer Motion)
   que se mantiene visible sin desmontar.

5. **TripTimeline** — Línea de progreso visual con 8 pasos (searching → completed),
   estados completed/active/pending, animaciones de aparición.

6. **MapController** — Wrapper que instancia `MapEngine` y reacciona a cambios de
   `pickupCoords`, `destCoords`, `route`, y posición del conductor.

7. **ConfirmState extendido** — Selector de tipo de vehículo (Standard/XL/Premium),
   cupón de descuento, observaciones, programación de viaje, edición de origen/destino.

## Cambios concretos

### ride-machine (packages/ride-machine/src/index.ts)
- Nuevos eventos: `SELECT_VEHICLE`, `APPLY_COUPON`, `SET_NOTE`, `SCHEDULE_TRIP`
- Nueva clase `JourneyEngine` con `send()`, `goTo()`, `onTransition()`, `snapshot`
- Función `createJourneyEngine()` factory
- Auto-detección de dirección forward/back basada en tipo de evento

### map-engine (packages/map-engine/src/)
- Nuevo archivo `engine.ts`: clase `MapEngine` con métodos:
  - `drawOrigin()`, `drawDestination()` con animaciones
  - `drawRoute()`, `updateRoute()` con animación progresiva
  - `fitToMarkers()` con padding configurable
  - `showDriver()` con marcador del conductor
  - `setInteractive()`, `setCenter()`, `panTo()`
  - `clearAll()`, `destroy()`
- `routes.ts`: `decodePolyline` ahora es exportado (era privado)
- `markers.ts`: Fix tipos `heading` (number no undefined), fix `rotation` → `labelOrigin`

### travel
- **Nuevo hook** `useJourneyEngine` — 376 líneas, state + FSM + API + persistencia
- **Nuevo componente** `BottomSheet` — sheet persistente con variantes de animación
- **Nuevo componente** `TripTimeline` — 8 pasos visuales con dots + líneas + labels
- **Nuevo componente** `MapController` — integra MapEngine, maneja marcadores/ruta/driver
- **ConfirmState mejorado** — selector vehículo, cupón, notas, programar, editar
- **page.tsx** — Reescrito para usar JourneyEngine + BottomSheet + MapController + Timeline
- **types** — Consistentes entre entities/trip.ts, types.ts, eventos
- **features/booking** — Exporta todos los nuevos componentes
- **shared/index** — Exporta JourneyEngine, MapEngine

## Consecuencias
- **Positivo**: Estado inmutable, transiciones predecibles, mapa desacoplado, persistencia
  robusta, panel de confirmación completo, build exitoso (0 errores TS, 99.8kB page)
- **Negativo**: Dependencia adicional de `@types/google.maps`, los componentes legacy
  (`useTripFlow`, `useAnimatedState`) aún existen pero no se usan
- **Riesgo**: `MapEngine` depende de `google.maps.Map` — cualquier cambio en la API
  de Google Maps requiere actualizar el adapter
