# AUDITORÍA FRONT-END CYTAXI — Control de calidad total

> Generado: julio 2026 · Cobertura: 100% archivos en `miniweb/src/`
> Organizado por: archivo → botones → animaciones → responsive → colores → issues

---

## ÍNDICE

1. [app/page.tsx — Orquestador principal](#1-apppagetsx)
2. [app/layout.tsx — Root layout](#2-applayouttsx)
3. [app/globals.css — CSS global](#3-appglobalscss)
4. [app/history/page.tsx — Historial](#4-historypagetsx)
5. [app/profile/page.tsx — Perfil](#5-profilepagetsx)
6. [hooks/useJourneyEngine.ts — Hook central](#6-hooksusejourneyenginets)
7. [components/BottomSheet.tsx](#7-bottomsheettsx)
8. [components/MapController.tsx](#8-mapcontrollertsx)
9. [components/TripTimeline.tsx](#9-triptimelinetsx)
10. [components/ErrorBoundary.tsx](#10-errorboundarytsx)
11. [components/states/PickUpStep.tsx](#11-pickupsteptsx)
12. [components/states/FormState.tsx](#12-formstatetsx)
13. [components/states/ConfirmState.tsx](#13-confirmstatetsx)
14. [components/states/TrackingState.tsx](#14-trackingstatetsx)
15. [components/states/ArrivingState.tsx](#15-arrivingstatetsx)
16. [components/states/PaymentState.tsx](#16-paymentstatetsx)
17. [components/states/DestinationState.tsx](#17-destinationstatetsx)
18. [components/states/RatingState.tsx](#18-ratingstatetsx)
19. [components/states/CompletedState.tsx](#19-completedstatetsx)
20. [components/layout/header.tsx, footer.tsx, layout.tsx](#20-layout-components)
21. [services/api.ts](#21-apits)
22. [services/tracking.ts](#22-trackingts)
23. [services/demo.ts](#23-demots)
24. [services/offline-queue.ts](#24-offline-queuets)
25. [services/payments.ts](#25-paymentsts)
26. [services/state-recovery.ts](#26-state-recoveryts)
27. [services/telemetry.ts](#27-telemetryts)
28. [services/journey-validator.ts](#28-journey-validatorts)
29. [shared/geo.ts](#29-sharedgeots)
30. [styles/design.ts](#30-stylesdesignts)
31. [types.ts / entities/](#31-typests--entitiests)
32. [features/](#32-features)
33. [Resumen consolidado](#33-resumen-consolidado)

---

## 1) app/page.tsx

### Botones

| Label | Tipo | Línea | Estados | issue |
|-------|------|-------|---------|-------|
| GPS "Mi ubicación" | `<button>` | 451 | idle, hover | Sin `type="button"`, sin feedback de error si geolocation denegada |
| Dismiss error "✕" | `<button>` | 488 | idle, hover | Sin `type="button"` |
| Navbar 5 tabs | `<a>` | 529 | active/inactive, hover | Tabs 2 y 4 usan `href="#"` + `e.preventDefault()` — sucio, deberían ser `<button>` o `role="tab"` |

### Animaciones

| Elemento | Tipo | Detalle |
|----------|------|---------|
| Map pin SVG | framer-motion spring | `stiffness: 480, damping: 22, mass: 0.6` — escala 1.15 durante drag |
| Pin label | framer-motion spring | `stiffness: 300, damping: 25, delay: 0.1` |
| Pin entrada | `initial: { opacity:0, scale:0.3, y:10 }` | Salida simétrica |
| Navbar icon stroke | CSS transition | `transition: "stroke 0.25s"` — solo en svg inline |

### Responsive

| Aspecto | Implementación | issue |
|---------|----------------|-------|
| Full viewport | `100vw` / `100dvh` | ✅ |
| Keyboard iOS | `visualViewport` resize → `keyboardH` | ✅ |
| Safe area navbar | `max(12px, env(safe-area-inset-bottom))` | ✅ |
| Touch action | `touchAction: "none"` en contenedor | ✅ |
| Body overflow | `overflow: hidden` en layout | ✅ |
| **Sin breakpoints** | No hay media queries | ⚠️ No adapta para tablets/landscape |

### Colores hardcodeados

| Valor | Línea | Debería usar |
|-------|-------|--------------|
| `#ea580c` (offline banner) | 471 | `colors.warning` o token |
| `#dc2626` (error banner) | 483 | `colors.danger` o `red` |
| `#121212` (pin label text) | 416 | `colors.textPrimary` |
| `#006c49` (pin A badge) | 429 | `colors.green` |
| `#448aff` (pin B badge) | 439 | `colors.blue` |
| `rgba(255,255,255,0.92)` (label bg) | 408 | `glass.surface` |
| `rgba(0,0,0,0.1)` (label shadow) | 413 | `shadows.card` |
| `"0 2px 8px rgba(0,0,0,0.12)"` (GPS btn) | 456 | `shadows.button` |
| `rgba(255,255,255,0.95)` (GPS bg) | 455 | `colors.surfaceWhite` |
| `#9ea5a0` (navbar inactive) | 534, 538 | `colors.textMuted` |
| `#82b1ff→#448aff→#1565c0` (pin B grad) | 380 | usar tokens blue* |
| `#00e676→#00a152→#006c49` (pin A grad) | 384 | usar tokens green* |
| `colors.brand.green` | 458, 517, 534, 535, 538 | Error: local `design.ts` exporta plano `colors.green`, no `colors.brand.green` |

### Issues por archivo

1. **`colors.brand.green`** en líneas 458, 517, 534, 535, 538 — la `design.ts` local exporta `colors.green` plano, no anidado en `colors.brand`. Esto tira runtime error (undefined).
2. **`colors.brand.greenBg`** en línea 535 — no existe en `design.ts` local.
3. **`zIndex`** en línea 517 — importado de `@cytaxi/design-tokens` pero la `design.ts` local no exporta `zIndex`.
4. **`glass.surface`** y `glass.blur` en línea 516 — sí existen en `design.ts` local ✅.
5. **Pin SVG gradientes duplicados**: `pinGradA`/`pinGradB` como IDs estáticos — si ambos estados coexisten, hay conflicto de IDs SVG.
6. **`getCenter()`** en línea 196 usa `c.lat()` (Leaflet function syntax), pero en línea 91 el handler usa `c.lat` (Google Maps property syntax) — error en tiempo de ejecución si se usa Leaflet.
7. **Navbar tabs 2 y 4** (`href="#"`) — no deberían renderizarse como links sin href real.
8. **No loading spinner** para geolocalización inicial — solo texto "Detectando ubicación...".
9. **`validateGoogle` timeout 3s** — puede dejar la UI colgada si Google Maps lento.

---

## 2) app/layout.tsx

### Botones — ninguno

### Animaciones — ninguna

### Responsive

| Aspecto | Implementación | issue |
|---------|----------------|-------|
| Viewport | `width=device-width, initialScale=1, maximumScale=1, userScalable=false, viewportFit=cover` | ✅ |
| Body | `overflow:hidden, position:fixed, width:100%, height:100%, touchAction:none, overscrollBehavior:none` | ✅ |
| **Sin safe area** en body | falta `env(safe-area-inset-*)` | ⚠️ Potencial en notched devices |

### Colores hardcodeados — ninguno

### Issues

1. `font-family` de body usa `class="antialiased"` de Tailwind — mezcla estilos inline con clases.
2. Google Maps key en línea 19: fallback a key hardcodeada — cualquier persona que clone el repo tiene acceso a la key.

---

## 3) app/globals.css

### Botones (CSS classes)

| Clase | Propósito | ¿Usada? |
|-------|-----------|---------|
| `.btn` | Base button | ❌ Nunca — todos los componentes usan inline styles |
| `.btn-primary` | Primary green | ❌ |
| `.btn-secondary` | Outline secondary | ❌ |
| `.btn-ghost` | Ghost button | ❌ |
| `.btn:disabled` | Disabled state | ❌ |
| `button:active:not(:disabled)` | Press effect | ✅ Vía selector global |

### Animaciones (keyframes)

| Nombre | Propósito | ¿Usada? |
|--------|-----------|---------|
| `spin` | Loading spinner | ✅ En PickUpStep, ConfirmState |
| `pulse` | Scale pulse | ❌ No referenciada |
| `radarPulse` | Radar rings | ✅ En TrackingState "searching" |
| `dotPulse` | Blinking dot | ✅ En TripTimeline, TrackingState |
| `fadeInUp` | Fade + translate | ❌ No referenciada |
| `pinEntrance` | Pin bounce | ❌ page.tsx usa framer-motion en vez de esto |
| `labelFadeIn` | Label fade | ❌ |
| `shimmer` | Skeleton loading | ✅ En TrackingState |
| `fadeScale` | Fade + scale | ✅ Clases `.animate-fade-scale` |
| `limeGlow` | Box-shadow pulse | ✅ Clase `.animate-lime-glow` |

**10 keyframes definidos → solo 5 están referenciados efectivamente.**

### Responsive

| Media Query | Propósito | issue |
|-------------|-----------|-------|
| `@media (hover: none)` | Desactivar hover en touch | ✅ |
| **Sin breakpoints** | No hay `@media (min-width...)` | ⚠️ Sin adaptación responsive |

### Colores hardcodeados (fuera de CSS vars)

| Valor | Línea | Debería usar |
|-------|-------|--------------|
| `#e8ebef` | 43 | Algo razonable para GMaps override |
| `#d4d9df` | 61 | `--cy-border-strong` |
| `#f4f6f8` | 87, 232, 280 | `--cy-surface` |
| `#eef1f4` | 91 | `--cy-surface-alt` |
| `#ffffff` | 104, 132 | `--cy-card` |
| `#fafbfc` | 137 | `--cy-surface-alt` |
| `#1565c0` | 169 | `--cy-blue-dark` (no existe como var) |
| `#b45309` | 170 | `--cy-warning` |
| `#448aff` | 237, 284 | `--cy-blue` |
| `#f4f6f8` | 280 | `--cy-surface` |

### Issues

1. **10 keyframes no utilizados** se envían al bundle CSS.
2. **CSS classes inaplicadas** (`.btn`, `.card`, `.sheet`, `.input-field`, `.vehicle-card`, `.pill-*`, `.driver-avatar`, `.float-el`, etc.) — 100% de los componentes usan inline styles, estas clases son código muerto.
3. **`--cy-glass-blur`**: `blur(20px)` en CSS vs `blur(24px) saturate(180%)` en `design.ts` — discrepancia.
4. **`--cy-success` y `--cy-primary`**: ambos `#006c49` — redundante.
5. **`text-blue-600`**: `#448aff` hardcodeado, debería ser `var(--cy-blue)`.
6. **`.card-lime`** nombrado como legacy "lime" (antes `#00b36b`) — debería renombrarse a `.card-green` o `.card-primary`.

---

## 4) history/page.tsx

### Botones

| Label | Tipo | Línea | Estados | issue |
|-------|------|-------|---------|-------|
| "← Volver" | `<Link>` | 54 | idle/hover | ✅ |
| "Solicita uno" | `<Link>` | 65 | idle/hover | ✅ |
| Trip items | `<li>` | 70 | — | No son clickables (no hay detalle) |

### Estados

| Estado | Implementación |
|--------|----------------|
| Loading | ✅ "Cargando…" |
| Error | ✅ texto rojo `#ba1a1a` |
| Empty | ✅ "Aún no tienes viajes" |
| Success | ✅ lista de viajes |

### Animaciones — ninguna

### Responsive — ninguno (pixel fijo, sin breakpoints)

### Colores hardcodeados

| Valor | Línea | Debería usar |
|-------|-------|--------------|
| `#ba1a1a` (error) | 60, 91 | `colors.red` |
| `var(--uk-primary)` | 65, 90, 92 | `var(--cy-primary)` — usa prefijo `uk-` que no existe |
| `var(--uk-secondary)` | 92 | No definido |
| `var(--uk-on-surface-variant)` | 93 | No definido |

### Issues

1. **Variables CSS `--uk-*` inexistentes**: el CSS global solo define `--cy-*`. Las referencias `var(--uk-primary)`, `var(--uk-secondary)`, `var(--uk-background)` (línea 52), `var(--uk-on-surface-variant)` (línea 93) **nunca están definidas** — se resuelven como `undefined` y el color cae a negro.
2. **`background: var(--uk-background)`** en línea 52 — no existe. Tampoco hay fallback.
3. **No paginación** — si hay muchos viajes, la lista crece sin límite.

---

## 5) profile/page.tsx

### Botones

| Label | Tipo | Línea | Estados | issue |
|-------|------|-------|---------|-------|
| "← Volver" | `<Link>` | 65 | idle/hover | ✅ |

### Estados

| Estado | Implementación |
|--------|----------------|
| Loading | ✅ "Cargando…" |
| Error | ✅ texto rojo |
| Empty (no trips) | ✅ "Aún no tienes viajes" |
| Success | ✅ profile + recent trips |

### Animaciones — ninguna

### Responsive — ninguno

### Colores hardcodeados

| Valor | Línea | Debería usar |
|-------|-------|--------------|
| `#ba1a1a` (error) | 71 | `colors.red` |
| `var(--uk-primary)` | 117, 118, 119, 120 | Mismo problema que history |
| `var(--uk-background)` | 63 | No existe |
| `var(--uk-on-surface-variant)` | 120 | No existe |

### Issues

1. **Mismas variables `--uk-*` rotas** que history/page.tsx.
2. **`colors.green` en línea 78** — importado de `@cytaxi/design-tokens`, no de `design.ts` local, pero funciona porque el paquete externo las exporta.
3. **Sin logout/edición** — página de solo lectura.

---

## 6) hooks/useJourneyEngine.ts

### Botones — ningún render (hook de lógica)

### Estados manejados (28 variables)

| Variable | Tipo | ¿Se expone? | issue |
|----------|------|-------------|-------|
| `state` | `RideState` | ✅ |
| `prevState` | `RideState \| null` | ✅ |
| `direction` | `"forward" \| "back"` | ✅ |
| `isTransitioning` | `boolean` | ✅ |
| `online` | `boolean` | ✅ |
| `pickupAddress` | `string` | ✅ |
| `pickupCoords` | `{lat,lng} \| null` | ✅ |
| `dest` | `Place \| null` | ✅ |
| `destQuery` | `string` | ✅ |
| `destSuggestions` | `Place[]` | ✅ |
| `route` | `RoutePayload \| null` | ✅ |
| `fare` | `FareBreakdown \| null` | ✅ |
| `driver` | `DriverInfo \| null` | ✅ |
| `tripId` | `string` | ✅ |
| `tracking` | `TrackingUpdate \| null` | ✅ — pero nunca se usa en UI |
| `eta` | `number` | ✅ |
| `loading` | `boolean` | ✅ |
| `paymentMethod` | `"cash" \| "card"` | ✅ |
| `vehicleType` | `string` | ✅ |
| `note` | `string` | ✅ |
| `coupon` | `string` | ✅ |
| `scheduledAt` | `string \| null` | ✅ |
| `error` | `string \| null` | ✅ |
| `noDrivers` | `boolean` | ✅ |
| `sheetRef` | `RefObject` | ✅ |

### Callbacks expuestos (22)

| Callback | issue |
|----------|-------|
| `send(event)` | ✅ |
| `goTo(state)` | ✅ |
| `handleCenterChange` | ✅ |
| `handleMapDestChange` | ✅ |
| `doSearchDest` | ✅ — debounce 400ms |
| `selectDest` | ✅ |
| `handleClearDest` | ✅ |
| `handleConfirmPickup` | ✅ |
| `handleBackToPickup` | ✅ |
| `handleConfirmDest` | ⚠️ no captura error de route/fare API |
| `handleRequestTrip` | ✅ — SSE temprano + noDrivers |
| `handleAcceptDriver` | ⚠️ sin feedback visual |
| `handleRetrySearch` | ✅ |
| `handleRejectDriver` | ✅ |
| `handleArriveAtPickup` | ✅ |
| `handleTripStart` | ✅ |
| `handleDestinationArriving` | ✅ |
| `handleTripComplete` | ✅ |
| `handlePaymentDone` | ✅ |
| `handleRatingDone` | ✅ |
| `startTracking` | ✅ |
| `reset` | ⚠️ duplicado parcial con `handleCancelTrip` |
| `handleCancelTrip` | ⚠️ duplicado parcial con `reset` |
| `dismissError` | ✅ |

### Side effects

| Effect | Propósito | issue |
|--------|-----------|-------|
| `online/offline` listeners | Detectar conectividad | ✅ |
| Sesión recovery on mount | Restaurar journey | ✅ — solo pre-trip |
| Session auto-save | Persistir estado | ✅ |
| Offline queue flush | Reintentar encolados | ✅ |
| `tracking` cleanup | Cleanup SSE | ✅ |
| State duration tracking | Telemetría | ✅ |

### Issues

1. **`handleConfirmDest` envía `SELECT_DEST` incluso si API falla** — el `try/catch` de línea 278 solo captura `requestTrip`, no `calculateRoute`/`estimateFare`.
2. **`reset` y `handleCancelTrip`**: ~80% lógica duplicada — refactorizar a helper común.
3. **`searchSubRef` asignación dinámica**: si se llama `handleRequestTrip` dos veces rápido, `searchSubRef.current?.()` libera la sub anterior — pero `= subscribeToTrip(...)` se ejecuta antes de que la anterior se cierre → posible fuga de EventSource.
4. **`tracking` estado seteado pero nunca leído en UI** — variable muerta.
5. **`(bus as any).emit(...)`**: 5+ ocurrencias que evaden type safety.
6. **`handleAcceptDriver` no dispara transición de máquina de estados** — llama `engineRef.current.send("DRIVER_ACCEPTED")` solo si `tripId && driver && state === "driver_found"`, pero no hay setDriver intermedio.
7. **Demo mode bypass**: cuando `DEMO_CONFIG.enabled=true`, el SSE real se suscribe pero nunca se usa porque `setTimeout` del demo pisa la respuesta.

---

## 7) BottomSheet.tsx

### Botones — ninguno (solo drag handle decorativo)

### Animaciones

| Elemento | Tipo | Detalle |
|----------|------|---------|
| Sheet height | framer-motion spring | `stiffness: 380, damping: 28, mass: 0.75` |
| Content enter/exit | tween | `ease: [0.4, 0, 0.2, 1]` — direction-aware |
| Forward enter | `opacity:0, y:22, scale:0.99` | ✅ |
| Back enter | `opacity:0, y:-12, scale:0.99` | ✅ |
| Forward exit | `opacity:0, y:-12, scale:0.99` | ✅ |
| Back exit | `opacity:0, y:22, scale:0.99` | ✅ |

### Responsive

| Aspecto | Implementación | issue |
|---------|----------------|-------|
| Posición | `position:fixed, left:0, right:0, bottom:0` | ✅ |
| maxHeight dinámico | `100dvh - navbar - keyboard` | ✅ |
| ResizeObserver | content height measurement | ✅ |
| iOS scroll | `WebkitOverflowScrolling: "touch"` | ✅ |
| **Sin media queries** | No breakpoints | ⚠️ |

### Colores hardcodeados

| Valor | Línea | Debería usar |
|-------|-------|--------------|
| `#ffffff` | 76 | `colors.surfaceWhite` |
| `rgba(0,0,0,0.12)` | 76 | token |
| `"0 -4px 24px rgba(0,0,0,0.08)"` | 77 | `shadows.float` |
| `border-radius: 22` | 77 | `radius.xl=16` — discrepancia |
| `willChange: "height"` | 73 | ✅ pero precaución performance |

### Issues

1. **Drag handle `pointerEvents: "none"`** — no es interactivo, no se puede hacer drag-to-dismiss.
2. **`border-radius: 22`** hardcodeado, no usa `radius.xl` (que es 16).
3. **`willChange: "height"`** en el motion.div — previene optimizaciones del browser si se usa en todos los renders.

---

## 8) MapController.tsx

### Botones — ninguno (solo mapa)

### Animaciones

| Elemento | Tipo | Detalle |
|----------|------|---------|
| Route polyline draw | `requestAnimationFrame` | 900ms cubic ease-out (`1 - (1-t)^3`) |
| Driver marker update | JS directo | Sin animación CSS |

### Responsive

| Aspecto | Implementación | issue |
|---------|----------------|-------|
| Contenedor | `width:100%, height:100%` | ✅ |
| fitBounds padding | `[40, 120]` / `[40, 320]` | ✅ |
| **Sin responsive** | No breakpoints | ⚠️ |

### Colores hardcodeados

| Valor | Línea | Debería usar |
|-------|-------|--------------|
| `#1c1c1e` | inline | `colors.mapDark` |
| `#006c49` (pickup marker) | pickupHtml | `colors.green` |
| `#276ef1` (dest Leaflet) | destHtml | `colors.blue` |
| `#fff` (dots) | pickupHtml/destHtml | token |
| `"0 2px 6px rgba(0,0,0,.35)"` | marker shadow | `shadows.button` |

### Issues

1. **`(window as any).__cymap`**: 5 referencias a global — rompe encapsulamiento.
2. **`map.getCenter()` es llamado en click** en vez de usar coordenadas del click — línea 93-96.
3. **Marker cleanup**: Leaflet driver marker no se limpia en unmount.
4. **`pickupHtml`/`driverHtml`**: strings HTML inline — riesgo XSS si label/color es controlado por usuario.
5. **`DEFAULT_CENTER`**: Guayaquil hardcodeado.
6. **Google Maps polling**: 6×300ms (1.8s) — si no carga, fallback a Leaflet, pero no hay indicación al usuario.

---

## 9) TripTimeline.tsx

### Botones — ninguno

### Animaciones

| Elemento | Tipo | Detalle |
|----------|------|---------|
| Step circles | framer-motion spring | `stiffness: 350, damping: 20`, delay escalonado `i * 0.04` |
| Step labels | CSS transition | `transition: "color 0.3s"` |
| Connector lines | framer-motion | `initial→animate` background, `duration: 0.5, ease: [0.4,0,0.2,1]` |

### Responsive

| Aspecto | Implementación | issue |
|---------|----------------|-------|
| Posición | `position:fixed, top:0, left:0, right:0` | ✅ |
| zIndex | 600 | ✅ |
| **Sin responsive** | Tamaños fijos (26px, 8px, 56px) | ⚠️ Puede desbordar en pantallas <320px |

### Colores hardcodeados

| Valor | Línea | Debería usar |
|-------|-------|--------------|
| `rgba(255,255,255,0.95)` | 25 | `glass.surface` |
| `rgba(0,0,0,0.04)` | 27 | `glass.border` |
| `rgba(0,0,0,0.06)` | 65, 118, 122 | `colors.borderLight` |
| `rgba(0,0,0,0.1)` | 67 | token |
| `rgba(0,108,73,0.15)` | 80 | token |
| `rgba(0,0,0,0.15)` | 91 | token |

### Issues

1. **Duplicación de lógica de timeline** con `TrackingState.tsx` — ambos componentes renderizan pasos similares con estilos diferentes. Si se añade un paso, hay que actualizar ambos.
2. **`maxWidth: 56`** en labels — texto truncado sin tooltip.
3. **Sin `top` offset** para safe area en notched devices — el contenido queda detrás del notch.

---

## 10) ErrorBoundary.tsx

### Botones

| Label | Tipo | Línea | Estados | issue |
|-------|------|-------|---------|-------|
| "Recargar página" | `<button>` | 57 | idle | Sin `type="button"` |

### Estados

| Estado | Implementación |
|--------|----------------|
| Sin error | Render children |
| Con error | Fallback UI con reloj ⚠️ |
| Error + fallback prop | Render fallback personalizado |

### Animaciones — ninguna

### Colores hardcodeados

| Valor | Línea | Debería usar |
|-------|-------|--------------|
| `colors.text.primary` | 50 | Error — `colors` es plano, no `colors.text` |
| `colors.text.muted` | 53, 67 | Mismo error |
| `colors.brand.green` | 60 | Error — `colors.green` plano |

### Issues

1. **`colors.text.primary`** y **`colors.text.muted`** — `colors` en `@cytaxi/design-tokens` tiene estructura anidada (`.text.primary`), pero en `design.ts` local es plano (`.textPrimary`). El import viene del paquete externo, no del local — depende de qué bundle se resuelva.
2. **`colors.brand.green`** — mismo problema, debe ser `colors.green`.
3. **Sin `type="button"`** en el botón de recarga.

---

## 11) PickUpStep.tsx

### Botones

| Label | Tipo | Línea | Estados |
|-------|------|-------|---------|
| "Confirmar ubicación" | `<button>` | 43 | idle, loading ("Detectando..."), disabled |

### Estados

| Estado | Implementación |
|--------|----------------|
| Loading | Spinner + "Detectando..." + botón deshabilitado (gris `#999`, opacidad 0.6) |
| Empty (address) | Muestra "Ubicación actual" |
| Success | Dirección + botón verde |

### Animaciones

| Elemento | Tipo | Detalle |
|----------|------|---------|
| Spinner | `spin 0.7s linear infinite` | CSS keyframe |
| Botón | `transition: "all 0.2s cubic-bezier(0.4,0,0.2,1)"` | ✅ |

### Responsive — ninguno

### Colores hardcodeados

| Valor | Línea | Debería usar |
|-------|-------|--------------|
| `#999` (loading) | 46 | `colors.textMuted` |
| `rgba(0,108,73,0.12)` (A shadow) | 31 | token |
| `rgba(0,108,73,0.25)` (btn shadow) | 49 | `shadows.buttonGreen` |
| `rgba(255,255,255,0.88)` (card) | 23 | `glass.surface` |
| `"blur(24px) saturate(180%)"` | 23 | `glass.blur` |
| `rgba(0,0,0,0.04)` (border) | 24 | `glass.border` |

### Issues

1. **Importa de `@cytaxi/design-tokens`** — el paquete externo. Debe decidir si usar el local o el externo, no mezclar.
2. **Sin `type="button"`** en el botón.

---

## 12) FormState.tsx

### Botones

| Label | Tipo | Línea | Estados |
|-------|------|-------|---------|
| ← Back arrow | `<div role="button">` | ~30 | idle, hover (via onMouseEnter) |
| X Clear dest | `<button>` | ~50 | visible solo si `destQuery` tiene texto |
| Suggestion items | `<button>` individual | ~70 | idle, hover, hasta 4 items |
| "Calculando..." / "Buscar viaje" / "Selecciona tu destino" | `<button>` | ~90 | loading, success, empty, disabled |

### Estados

| Estado | Implementación |
|--------|----------------|
| Loading (confirm) | Spinner + "Calculando..." + botón gris `#999` |
| Empty (no query) | Placeholder "Buscar dirección o lugar" |
| Empty (no results) | ❌ **No manejado** — espacio en blanco |
| Success (dest selected) | Highlight azul + confirm button verde |
| Suggestions visible | Hasta 4 resultados |
| Greeting dinámico | "Buenos días/tardes/noches" según hora |

### Animaciones

| Elemento | Tipo | Detalle |
|----------|------|---------|
| Loading spinner | `spin 0.7s linear infinite` | CSS keyframe |
| Suggestions bg | `transition: "background 0.12s"` | ✅ |
| Main button | `transition: "all 0.25s cubic-bezier(0.4,0,0.2,1)"` | ✅ |
| Dest row bg | `transition: "background 0.2s"` | ✅ |

### Responsive — ninguno

### Colores hardcodeados

| Valor | Línea | Debería usar |
|-------|-------|--------------|
| `#448aff` (B marker) | múltiple | `colors.blue` |
| `rgba(68,138,255,0.04)` (sel bg) | — | `colors.blue` con opacidad |
| `rgba(68,138,255,0.12)` (shadow) | — | token |
| `#9a9a9a` (muted) | — | `colors.textMuted` |
| `#999` (loading) | — | `colors.textMuted` |
| `rgba(0,108,73,0.25)` (btn shadow) | — | `shadows.buttonGreen` |

### Issues

1. **Error state not handled** — si `onSearch` o `onConfirm` fallan, no se muestra error.
2. **Empty results not handled** — si hay query pero 0 resultados, se ve espacio vacío.
3. **Pickup row es clickeable vía `onBack`** — confuso, el usuario espera editar pickup, no ir atrás.
4. **Sin `type="button"`** en ningún botón.
5. **Pickup row** usa `role="button"` + `tabIndex={0}` + `onKeyDown` — accesible pero frágil.

---

## 13) ConfirmState.tsx

### Botones

| Label | Tipo | Línea | Estados |
|-------|------|-------|---------|
| "Editar origen" | hover inline | ~215 | idle, hover |
| "Editar destino" | hover inline | ~215 | idle, hover |
| 🚗 Standard | `<button>` con aria-pressed | ~260 | idle, selected (green border + scale 1.02) |
| 🚐 XL | ídem | ~260 | idle, selected |
| 🚙 Premium | ídem | ~260 | idle, selected |
| 💵 Efectivo | `<button>` con aria-pressed | ~300 | idle, selected |
| 💳 Tarjeta | ídem | ~300 | idle, selected |
| "+ Agregar cupón de descuento" | `<button>` | ~320 | collapsed/expanded |
| "Aplicar" | `<button>` | ~340 | idle |
| "+ Agregar observación" | `<button>` | ~360 | collapsed/expanded |
| "🕐 Programar viaje" | `<button>` | ~380 | idle/set |
| "Quitar" (schedule) | `<button>` | ~400 | visible solo si scheduled |
| "Atrás" | `<button>` | ~420 | idle, disabled durante loading |
| "Solicitar viaje" / "Procesando..." | `<button>` | ~430 | idle, loading (spinner + disabled) |

### Estados

| Estado | Implementación |
|--------|----------------|
| Loading | Spinner + "Procesando..." + botones deshabilitados |
| Vehicle selection | 3 opciones, active = green border + scale(1.02) |
| Payment toggle | Cash/Card con green border activo |
| Coupon toggle | Colapsable, input + "Aplicar" |
| Notes toggle | Colapsable, visible si hay nota |
| Schedule toggle | Set 1h desde ahora, "Quitar" para limpiar |
| Fare breakdown | Base × multiplier + distance × multiplier + time × multiplier + total |
| Route cards | distance_km + eta_minutes |

### Animaciones

| Elemento | Tipo | Detalle |
|----------|------|---------|
| Vehicle cards | `transition: "all 0.2s cubic-bezier(0.4,0,0.2,1)"` | ✅ |
| Payment cards | `transition: "all 0.2s"` | ✅ |
| Loading spinner | `spin 0.7s linear infinite` | ✅ |

### Responsive — ninguno

### Colores hardcodeados

| Valor | Línea | Debería usar |
|-------|-------|--------------|
| `#448aff` (dest marker) | ~50 | `colors.blue` |
| `#999` (loading bg) | — | token |
| `rgba(0,108,73,0.25)` (shadow) | — | `shadows.buttonGreen` |
| `rgba(0,108,73,0.06)` / `0.08` (selected) | — | token |

### Issues

1. **Coupon "Aplicar" llama `onCouponChange?.("")`** — no aplica, solo limpia.
2. **Schedule siempre 1h desde ahora** — no hay date picker real.
3. **Sin error state** — si route/fare falla, sección oculta pero no hay feedback.
4. **`onBack` es mismo callback para editar origen y destino** — no diferencia.
5. **Sin `type="button"`** en botones.

---

## 14) TrackingState.tsx

### Botones

| Label | Tipo | Línea | Estados |
|-------|------|-------|---------|
| "Reintentar" | `<button>` | 34 | visible solo si `noDrivers && onRetry` |
| "Volver" | `<button>` | 47 | visible solo si `noDrivers && onCancel` |
| "Cancelar" (searching) | `<button>` | 67 | hover (bg rojo suave) |
| "📞 Llamar" | `<button>` | ~210 | hover (cambia bg) |
| "💬 Mensaje" | `<button>` | ~220 | hover (cambia bg) |
| "✕" (reject driver) | `<button>` | ~230 | hover (cambia bg) |

### Estados

| Estado | Implementación |
|--------|----------------|
| No drivers | ⚠️ icon + "No encontramos conductores" + Reintentar/Volver |
| Searching (loading) | Radar rings + shimmer skeleton + "Buscando conductor" |
| Driver found | ETA header, driver card, timeline 5 pasos, payment badge, action buttons |
| Timeline steps | `completed` / `active` / `pending` |

### Animaciones

| Elemento | Tipo | Detalle |
|----------|------|---------|
| Radar rings | `radarPulse 2s ease-out infinite` | CSS keyframe, 3 círculos concéntricos con delay 0.5s |
| Shimmer skeleton | `shimmer 1.5s ease-in-out infinite` | CSS keyframe, delays escalonados 0.1s/0.15s/0.2s |
| Dot pulse | `dotPulse 1.5s infinite` | Timeline indicator |
| Timeline circles | `transition: "all 0.3s"` | ✅ |

### Responsive — ninguno (tamaños fijos)

### Colores hardcodeados

| Valor | Línea | Debería usar |
|-------|-------|--------------|
| `rgba(186,26,26,0.06)` / `0.04` | — | `colors.red` con opacidad |
| `#eee` / `#f5f5f5` (shimmer) | — | `colors.surfaceSkeleton` |
| `linear-gradient(90deg, #eee ...)` | — | token |
| `rgba(0,0,0,0.1)` / `0.06` / `0.04` / `0.08` | — | tokens border |
| `"linear-gradient(135deg, #e8f5e9, #c8e6c9)"` (avatar) | — | token |
| `rgba(0,108,73,0.08)` / `0.15` / `0.25` | — | token |
| `"0 4px 20px rgba(0,108,73,0.25)"` | — | `shadows.buttonGreen` |

### Issues

1. **`pickup` y `dest` props aceptadas pero nunca renderizadas** — props muertas.
2. **Sin feedback de loading en `handleRejectDriver`** — el botón ✕ no muestra spinner.
3. **Timeline duplicado** con `TripTimeline.tsx` — riesgo de desincronización.
4. **Fallback hardcodeado `"0.5"` km** cuando `route` es null.
5. **Sin error state** — todo error se maneja en page.tsx via banner.

---

## 15) ArrivingState.tsx

### Botones

| Label | Tipo | Línea | Estados |
|-------|------|-------|---------|
| "El conductor llegó" | `<button>` | ~30 | idle, hover (scale 1.01) |
| "Iniciar viaje" | `<button>` | ~40 | idle, hover (scale 1.01) |
| "Cancelar" | `<button>` | ~50 | visible solo si `onCancel` |

### Estados

| Estado | Implementación |
|--------|----------------|
| Not arrived | Progress ring + elapsed timer (1s interval) + "Conductor en camino" |
| Arrived | Checkmark + "El conductor ha llegado" + "Te está esperando" |
| Empty (driver null) | Driver card oculto |

### Animaciones

| Elemento | Tipo | Detalle |
|----------|------|---------|
| Progress ring | SVG `stroke-dashoffset` | `transition: 1s cubic-bezier(0.4, 0, 0.2, 1)` — actualiza cada 1s |

### Responsive — ninguno

### Colores hardcodeados

| Valor | Línea | Debería usar |
|-------|-------|--------------|
| `rgba(0,0,0,0.06)` (ring bg) | — | token |
| `rgba(0,161,82,0.08)` (arrived circle) | — | token |
| `rgba(0,108,73,0.06)` (searching circle) | — | token |
| `rgba(0,108,73,0.25)` (btn shadow) | — | `shadows.buttonGreen` |
| `"linear-gradient(135deg, #e8f5e9, #c8e6c9)"` | — | token |
| `rgba(186,26,26,0.04)` (cancel hover) | — | token |

### Issues

1. **Sin error state** — si driver data es stale, no hay feedback.
2. **Fallback `"0.5"` km** cuando `route` es null.
3. **Sin `type="button"`** en los botones.

---

## 16) PaymentState.tsx

### Botones

| Label | Tipo | Línea | Estados |
|-------|------|-------|---------|
| "Calificar viaje" | `<button>` | ~60 | idle, hover (scale 1.01) |

### Estados

| Estado | Implementación |
|--------|----------------|
| Processing | Progress ring 0→100% en ~2s (interval 40ms) + "Procesando pago" |
| Complete | Checkmark + "Pago confirmado" + fare summary + "Calificar viaje" |
| Empty (fare null) | Fare amount oculto |

### Animaciones

| Elemento | Tipo | Detalle |
|----------|------|---------|
| Progress ring | SVG `stroke-dashoffset` | `transition: 0.2s cubic-bezier(0.4, 0, 0.2, 1)` — actualiza cada 40ms |
| Botón hover | `transition: "all 0.2s"` | scale(1.01) |

### Responsive — ninguno

### Colores hardcodeados

| Valor | Línea | Debería usar |
|-------|-------|--------------|
| `rgba(0,0,0,0.06)` (ring bg) | — | token |
| `rgba(0,161,82,0.08)` (check bg) | — | token |
| `rgba(255,255,255,0.7)` (fare card) | — | `colors.surfaceCard` |
| `rgba(0,0,0,0.04)` (border) | — | token |
| `"0 4px 20px rgba(0,108,73,0.25)"` | — | `shadows.buttonGreen` |

### Issues

1. **No hay error state** — el pago siempre "exitoso" (simulado 2s, sin API real).
2. **Processing puramente cosmético** — no hay llamada a `payments.ts`.
3. **Sin `type="button"`**.

---

## 17) DestinationState.tsx

### Botones

| Label | Tipo | Línea | Estados |
|-------|------|-------|---------|
| "Finalizar viaje" | `<button>` | 61 | idle, hover (scale 1.01) |

### Estados

| Estado | Implementación |
|--------|----------------|
| Success | Icono mapa + "Llegando a tu destino" + destino card |
| Empty (dest null) | "Prepárate para bajar" |
| Driver info | Nombre + "Tu conductor" |

### Animaciones

| Elemento | Tipo |
|----------|------|
| Botón hover | `transition: "all 0.2s"`, scale(1.01) |

### Responsive — ninguno

### Colores hardcodeados

| Valor | Línea | Debería usar |
|-------|-------|--------------|
| `rgba(0,108,73,0.06)` (icon bg) | 16 | token |
| `rgba(255,255,255,0.88)` (card) | 34 | `glass.surface` |
| `rgba(0,0,0,0.04)` (border) | 35 | `glass.border` |
| `"0 4px 20px rgba(0,108,73,0.25)"` (shadow) | 67 | `shadows.buttonGreen` |

### Issues

1. **Sin `type="button"`**.
2. **`driver` prop aceptada pero solo muestra nombre** — foto y rating ignorados.

---

## 18) RatingState.tsx

### Botones

| Label | Tipo | Línea | Estados |
|-------|------|-------|---------|
| 5 star buttons | `<button>` | 74 | idle, hover (scale 1.15 + full color), selected |
| "Enviar N estrellas" / "Selecciona una calificación" | `<button>` | 103 | disabled (score=0) → gris, enabled → verde |

### Estados

| Estado | Implementación |
|--------|----------------|
| Unrated | Stars grises (grayscale(1) + opacity 0.35) |
| Hovered | Preview de estrella seleccionada |
| Rated (before submit) | Stars a color + input comment visible |
| Submitted | 🎉 + "¡Gracias por tu calificación!" + delay 600ms → callback |
| Empty (driver null) | Driver section oculto |

### Animaciones

| Elemento | Tipo | Detalle |
|----------|------|---------|
| Star hover | `transition: "all 0.2s cubic-bezier(0.34,1.56,0.64,1)"` | Spring-like |
| Star scale | `transform: active ? "scale(1.15)" : "scale(1)"` | ✅ |
| Grayscale filter | `filter: active ? "none" : "grayscale(1) opacity(0.35)"` | ✅ |

### Responsive — ninguno

### Colores hardcodeados

| Valor | Línea | Debería usar |
|-------|-------|--------------|
| `rgba(0,108,73,0.06)` (icon bg) | 28 | token |
| `"linear-gradient(135deg, #e8f5e9, #c8e6c9)"` (avatar) | 57 | token |
| `rgba(0,0,0,0.08)` (border) | 96 | token |
| `rgba(255,255,255,0.8)` (input bg) | 98 | token |
| `"0 4px 20px rgba(0,108,73,0.25)"` (shadow) | 112 | `shadows.buttonGreen` |

### Issues

1. **Sin `type="button"`** en stars ni submit.
2. **Comentario se recolecta pero nunca se envía** — `onDone(score)` solo envía score.
3. **Emoji 👤 en avatar por defecto** — el audit anterior pidió remover emojis pero este sobrevive.

---

## 19) CompletedState.tsx

### Botones

| Label | Tipo | Línea | Estados |
|-------|------|-------|---------|
| "Nuevo viaje" | `<button>` | ~60 | idle, hover (scale 1.01) |

### Estados

| Estado | Implementación |
|--------|----------------|
| Success | Checkmark + "Viaje completado" + "Gracias por viajar con CYTAXI" |
| Empty (no pickup/dest) | Route summary oculto |
| Empty (no fare) | Fare detail oculto |

### Animaciones

| Elemento | Tipo |
|----------|------|
| Botón hover | `transition: "all 0.2s"`, scale(1.01) |

### Responsive — ninguno

### Colores hardcodeados

| Valor | Línea | Debería usar |
|-------|-------|--------------|
| `rgba(0,161,82,0.08)` (check bg) | — | token |
| `rgba(255,255,255,0.88)` (card) | — | `glass.surface` |
| `rgba(255,255,255,0.7)` (surface) | — | `colors.surfaceCard` |
| `#448aff` (dest dot) | — | `colors.blue` |
| `"0 4px 20px rgba(0,108,73,0.25)"` (shadow) | — | `shadows.buttonGreen` |

### Issues

1. **`driver` prop muerta** — aceptada pero nunca renderizada.
2. **Sin `type="button"`**.
3. **Route summary sin línea de ruta** — solo dots pickup/dest.

---

## 20) Layout components (header.tsx, footer.tsx, layout.tsx)

### Botones

| Label | Tipo | Archivo | issue |
|-------|------|---------|-------|
| CYTAXI logo link | `<Link>` | header.tsx | ✅ |
| Trip link | `<Link>` | header.tsx | ✅ |
| Profile link | `<Link>` | header.tsx | ✅ |
| Notifications link | `<Link>` | header.tsx | ✅ |

### Estados — solo render estático

### Animaciones — ninguna

### Responsive

| Aspecto | Implementación | issue |
|---------|----------------|-------|
| Max width | `max-w-3xl` | ✅ pero no usa viewport units |

### Issues

1. **Estos componentes NO se usan en el app real** — `layout.tsx` no los importa ni renderiza. El layout actual es inline en `app/layout.tsx` con body directo. Estos son **código muerto** de una versión anterior.
2. Usan Tailwind classes (`container`, `border-border`, `bg-background`, etc.) que dependen de `tailwind.config.js` — si no existe, estas clases no generan CSS.

---

## 21) api.ts

### Funciones

| Función | Timeout | Retries | issue |
|---------|---------|---------|-------|
| `searchPlaces` | 8s | 2 | ✅ mapeo robusto |
| `calculateRoute` | 8s | 2 | ⚠️ retorna null en error |
| `estimateFare` | 8s | 2 | ✅ fallback formula |
| `requestTrip` | 8s | 2 | ⚠️ sin fallback |
| `startMatching` | 8s | 2 | ❌ nunca llamada |

### Estados

| Estado | Implementación |
|--------|----------------|
| Loading | Implícito (async) |
| Error (timeout) | ✅ "TIMEOUT" |
| Error (HTTP) | ✅ "HTTP_{status}" |
| Error (network) | ✅ retry con backoff 1s, 2s (cap 4s) |
| Empty (search) | ✅ retorna [] |
| Edge case (null coords) | ✅ filtrado |
| Fallback fare | ✅ si API falla |

### Issues

1. **`startMatching` definida pero nunca llamada** — código muerto.
2. **`calculateRoute` y `searchPlaces` silencian errores** — retornan null/[] sin distinción de tipo de error.
3. **Server IP hardcodeada** `64.176.219.221` — debería ser env var.
4. **Fallback fare formula** muy básica: `1 + distancia * 0.5 + duración * 0.02`.

---

## 22) tracking.ts

### Funciones

| Función | Propósito | issue |
|---------|-----------|-------|
| `subscribeToTrip` | SSE subscription with reconnection | ✅ |
| `createTrackingUrl` | URL builder | ❌ nunca llamada |

### Estados

| Estado | Implementación |
|--------|----------------|
| Connected | onmessage → parse → onUpdate |
| Connection error | Reconnect con backoff 1s→2s→4s→...→30s |
| Max retries (10) | onError("max_retries") |
| Trip completed | Auto-close |
| Parse error | onError("parse_error") |
| Cleanup | Unsubscribe function |

### Issues

1. **`createTrackingUrl` código muerto** — nunca se usa fuera del archivo.
2. **Sin heartbeat/ping** — si conexión se corta sin `onerror`, no reconecta.
3. **Sin `withCredentials`** — no envía cookies si API las requiere.

---

## 23) demo.ts

### Config

| Campo | Valor | Propósito |
|-------|-------|-----------|
| `enabled` | `true` | Demo mode switch |
| `matchingDelay` | `3000` | Demo matching delay |
| `searchTimeout` | `30000` | Real search timeout |
| `simulateNoDrivers` | `false` | Test no-drivers flow |
| `driver` | demo driver | Placeholder data |
| `passenger` | demo passenger | Placeholder data |

### Issues

1. **Solo cubre searching/driver assignment** — falta demo para payment, rating, completion.
2. **`passenger.phone = "0000000000"`** — fake, user real phone nunca se recolecta.
3. **Driver position offset** `pickupCoords.lat - 0.01` — siempre mismo delta.

---

## 24) offline-queue.ts

### Funciones y estados

| Función | Propósito | issue |
|---------|-----------|-------|
| `enqueueAction` | Add to queue | ✅ |
| `dequeueAction` | Remove from queue | ✅ |
| `getPendingActions` | Get actions with retries < maxRetries | ✅ |
| `incrementRetry` | Increment retry count | ✅ |
| `clearQueue` | Clear all | ✅ |
| `isOnline` | navigator.onLine wrapper | ✅ |
| `onOnline` | Listener | ✅ |
| `onOffline` | Listener | ✅ |

### Issues — ninguno. Archivo sólido.

---

## 25) payments.ts

### Funciones

| Función | issue |
|---------|-------|
| `createPayment` | ⚠️ sin manejo de error HTTP |
| `confirmPayment` | ⚠️ sin manejo de error HTTP |
| `getPayment` | ⚠️ sin manejo de error HTTP |
| `getReceipt` | ⚠️ sin manejo de error HTTP |
| `getPaymentHistory` | ⚠️ sin manejo de error HTTP |
| `getDriverEarnings` | ⚠️ sin manejo de error HTTP |
| `refundPayment` | ⚠️ sin manejo de error HTTP |

### Issues

1. **Ninguna función maneja errores HTTP** — si la API retorna 4xx/5xx, `res.json()` parsea el error como éxito.
2. **Ninguna función es llamada desde la app** — `PaymentState.tsx` usa setTimeouts, no `payments.ts`. **Código 100% muerto.**

---

## 26) state-recovery.ts

### Funciones

| Función | Propósito | issue |
|---------|-----------|-------|
| `saveSession` | Guardar journey | ✅ |
| `loadSession` | Cargar journey | ✅ |
| `clearSession` | Limpiar session | ✅ |
| `isSessionValid` | Validar session | ✅ |

### Issues

1. **Expiración 24h hardcodeada** (`86400000`). OK para MVP.
2. **`isSessionValid` acepta post-trip states con `tripId`** — pero `useJourneyEngine` solo restaura pre-trip. Inconsistencia.

---

## 27) telemetry.ts

### Funciones

| Función | Propósito | issue |
|---------|-----------|-------|
| `trackJourneyEvent` | Track event | ✅ |
| `trackStateDuration` | Track state duration | ✅ |
| `trackError` | Track error | ✅ |
| `trackLatency` | Track API latency | ✅ |

### Issues

1. **Sin rate limiting** — si hay muchos eventos, el buffer puede crecer (max 200).
2. **Flush cada 30s independientemente de estado** — puede enviar payload grande con muchos eventos acumulados.
3. **No batch optimization** — envía todos los eventos acumulados en una sola petición POST.

---

## 28) journey-validator.ts

### Funciones

| Función | Propósito | issue |
|---------|-----------|-------|
| `validateJourneyState` | Validar transición | ✅ |
| `validatePickupData` | Validar pickup | ✅ |
| `validateDestinationData` | Validar destino | ✅ |
| `validateRouteData` | Validar ruta | ⚠️ polyline check |
| `validateFareData` | Validar tarifa | ✅ |
| `validateDriverData` | Validar conductor | ✅ |
| `validateTrackingData` | Validar tracking | ✅ |
| `validateFullJourney` | Validar todo | ✅ |

### Issues

1. **Solo se importa, no se llama** — `journey-validator.ts` se importa en `useJourneyEngine.ts` pero ninguna función es invocada. **Código 100% muerto.**
2. Tests unitarios existen (`__tests__/journey-validator.test.ts`) ✅.

---

## 29) shared/geo.ts

### Función

| Función | issue |
|---------|-------|
| `validateGoogleCoords` | ✅ |

### Issues

1. **Timeout 3s** — si Google Maps lento, usuario espera.
2. **Sin caché** — cada llamada geocodifica de nuevo.
3. **`globalThis` cast** — `(globalThis as any).google` evita type safety.

---

## 30) styles/design.ts

### Tokens exportados

| Grupo | Items | issue |
|-------|-------|-------|
| `colors` | 14 colores | ⚠️ faltan `danger`, `info`, `warning`, `success` |
| `spacing` | 6 valores | ✅ |
| `radius` | 6 valores | ✅ |
| `shadows` | 5 sombras | ✅ |
| `typography` | 7 estilos | ✅ |
| `glass` | 3 valores | ✅ |
| `transitions` | 3 valores | ✅ |

### Issues

1. **Faltan aliases** como `danger`, `info`, `warning`, `success` en `colors` — componentes usan `colors.danger`, `colors.info` desde el paquete externo `@cytaxi/design-tokens`.
2. **`colors` es `as const`** — no permite extensión.
3. **`surfaceDark` (="#121212") nunca se usa** en ningún componente.
4. **Discrepancia con CSS vars**: `glass.blur` = `"blur(24px) saturate(180%)"` vs `--cy-glass-blur: blur(20px)`.

---

## 31) types.ts / entities/trip.ts

### Tipos duplicados

| Tipo | types.ts | entities/trip.ts | issue |
|------|----------|------------------|-------|
| `RideState` | ✅ re-export | ✅ re-export | ❌ duplicado |
| `RideEvent` | ✅ re-export | ✅ re-export | ❌ duplicado |
| `TripState` | ✅ alias | ✅ alias | ❌ duplicado |
| `Coordinates` | ✅ definido | ✅ definido | ❌ duplicado — idéntico |
| `Place` | ✅ definido | ✅ definido | ❌ duplicado — idéntico |
| `DriverInfo` | ✅ alias | ✅ alias | ❌ duplicado |
| `FareBreakdown` | ✅ alias | ✅ alias | ❌ duplicado |
| `TrackingUpdate` | ✅ interfaz | ✅ interfaz | ❌ duplicado — idéntico |
| `TripRequest` | ✅ interfaz | ✅ interfaz | ❌ duplicado — idéntico |
| `RoutePayload` | ✅ alias | ✅ alias | ❌ duplicado |

**`types.ts` y `entities/trip.ts` son archivos idénticos en propósito y contenido.** Uno debería importar del otro.

---

## 32) features/

### booking/index.ts

| Export | Origen | issue |
|--------|--------|-------|
| 13 exports | hooks, components, states | Barrel file — re-exporta todo lo que ya se puede importar directamente |

### trip/index.ts

| Export | Origen | issue |
|--------|--------|-------|
| 13 exports | hooks, components, states | **Idéntico a booking/index.ts** — barrel duplicado exacto |

### booking/services/index.ts

| Export | Origen | issue |
|--------|--------|-------|
| 5 API functions | services/api.ts, tracking.ts | ✅ barrel útil |

### Issues

1. **`features/booking/index.ts` y `features/trip/index.ts` son idénticos** — barrel duplicado.
2. Ninguno de estos barrels es importado por `page.tsx` o `useJourneyEngine.ts` — **todo código muerto**.

---

## 33) Resumen consolidado

### Código muerto identificado

| Archivo/Sección | Tipo | Impacto |
|-----------------|------|---------|
| `globals.css` clases `.btn`, `.card`, `.sheet`, etc. | ~30 clases CSS | Bundle inflado, ~5kb |
| `globals.css` keyframes `pulse`, `fadeInUp`, `pinEntrance`, `labelFadeIn`, `limeGlow` | 5 keyframes | Bundle inflado |
| `layout/header.tsx`, `footer.tsx`, `layout.tsx` | 3 componentes | No renderizados |
| `services/payments.ts` completo | 7 funciones | Nunca llamadas |
| `services/journey-validator.ts` completo | 8 funciones | Nunca llamadas (solo tests) |
| `services/api.ts` → `startMatching` | 1 función | Nunca llamada |
| `services/tracking.ts` → `createTrackingUrl` | 1 función | Nunca llamada |
| `features/booking/index.ts` y `features/trip/index.ts` | 2 barrels | Nunca importados |
| `types.ts` | archivo completo | Duplicado de `entities/trip.ts` |
| `CompletedState.tsx` prop `driver` | 1 prop | Nunca renderizada |
| `TrackingState.tsx` props `pickup`, `dest` | 2 props | Nunca renderizadas |
| `useJourneyEngine.ts` → `tracking` state | 1 variable | Nunca leída en UI |

### Botones — Checklist general

| Aspecto | Estado |
|---------|--------|
| `type="button"` en todos los `<button>` | ❌ **Ningún botón tiene `type="button"`** en toda la app |
| `aria-label` en botones | ✅ Mayoría tienen |
| `disabled` state | ✅ En botones de confirm/submit |
| `loading` state | ✅ En PickUpStep, ConfirmState, RatingState |
| `hover` state | ✅ Via `onMouseEnter`/`onMouseLeave` (inline) |
| `:active` scale(0.97) | ✅ Global CSS |
| Hover en touch devices | ✅ Suprimido vía `@media (hover: none)` |

### Animaciones — Resumen

| Tipo | Cantidad | Archivos |
|------|----------|----------|
| CSS keyframes usados | 5 | `spin`, `radarPulse`, `dotPulse`, `shimmer`, `fadeScale` |
| CSS keyframes no usados | 5 | `pulse`, `fadeInUp`, `pinEntrance`, `labelFadeIn`, `limeGlow` |
| framer-motion springs | ~8 instancias | page.tsx (pin, label), BottomSheet (sheet, content), TripTimeline (steps) |
| framer-motion tweens | ~2 instancias | BottomSheet (content entry/exit), TripTimeline (connectors) |
| SVG stroke-dashoffset | 2 instancias | ArrivingState, PaymentState (progress rings) |
| JS requestAnimationFrame | 1 instancia | MapController (route draw) |

### Responsive — Resumen

| Aspecto | Implementación | issue |
|---------|----------------|-------|
| Viewport units | `100vw`, `100dvh` | ✅ |
| Keyboard avoidance | `visualViewport` resize → `keyboardH` | ✅ Solo page.tsx |
| Safe area navbar | `env(safe-area-inset-bottom)` | ✅ |
| Touch action | `touchAction: "none"` | ✅ |
| Overflow | `overflow: hidden` | ✅ |
| Overscroll | `overscroll-behavior: none` | ✅ |
| **Media queries** | ❌ **Ningún breakpoint** | ⚠️ No hay adaptación tablet/desktop/landscape |
| **Font sizes** | **100% px fijos** | ⚠️ Sin `rem`/`em` |
| **Sizes** | **100% px fijos** | ⚠️ Sin `%`/`vw` excepto `width: "100%"` |
| **Print styles** | ❌ | No hay `@media print` |

### Colores hardcodeados — Por archivo

| Archivo | Count | Ejemplos |
|---------|-------|----------|
| page.tsx | ~15 | `#ea580c`, `#dc2626`, `#121212`, `#9ea5a0`, `colors.brand.green` (error) |
| TrackingState.tsx | ~15 | `#eee`, `#f5f5f5`, múltiples rgba |
| globals.css | ~10 | `#e8ebef`, `#d4d9df`, `#f4f6f8`, `#448aff`, `#1565c0` |
| FormState.tsx | ~7 | `#448aff`, `#9a9a9a`, `#999` |
| ConfirmState.tsx | ~6 | `#448aff`, `#999`, múltiples rgba |
| ArrivingState.tsx | ~6 | múltiples rgba |
| PaymentState.tsx | ~5 | múltiples rgba |
| BottomSheet.tsx | ~3 | `#ffffff`, rgba |
| PickUpStep.tsx | ~5 | `#999`, rgba |
| RatingState.tsx | ~5 | rgba, gradient |
| CompletedState.tsx | ~5 | `#448aff`, rgba |
| ErrorBoundary.tsx | ~3 | `colors.text.primary` (error) |
| history/profile | ~6 | `var(--uk-*)` (inexistentes) |
| MapController.tsx | ~6 | `#1c1c1e`, `#276ef1` |
| DestinationState.tsx | ~4 | rgba |
| TripTimeline.tsx | ~6 | rgba |

**Total estimado: ~110 valores hardcodeados que deberían usar tokens.**

### Issues críticos (deben arreglarse)

1. **`colors.brand.green` en page.tsx** — crash en runtime porque `design.ts` exporta `colors.green` plano.
2. **`var(--uk-*)` en history/profile** — colores undefined porque esas variables CSS no existen.
3. **`ErrorBoundary.tsx` usa `colors.text.primary`** — la estructura de `@cytaxi/design-tokens` es anidada (`.text`), la local es plana (`.textPrimary`). Depende de qué paquete se resuelva.
4. **Ningún botón tiene `type="button"`** — riesgo de submit implícito si se envuelven en `<form>`.
5. **`types.ts` y `entities/trip.ts` duplicados** — confusión de imports.
6. **`features/booking/index.ts` y `features/trip/index.ts` duplicados** — barrel idéntico.
7. **`services/payments.ts` y `services/journey-validator.ts` completamente muertos** — 15 funciones que nadie llama.
8. **3 componentes layout sin usar** (`header.tsx`, `footer.tsx`, `layout.tsx`).

### Issues de calidad (mejorable)

1. **~110 colores hardcodeados** vs tokens de `design.ts`.
2. **Cero breakpoints responsive** — la app no adapta en landscape/tablet/desktop.
3. **~30 clases CSS inaplicadas** en `globals.css`.
4. **5 keyframes CSS no usados**.
5. **Dualidad de imports**: algunos archivos importan de `@cytaxi/design-tokens` (paquete externo), otros de `@/styles/design` (local) — pueden tener valores diferentes.
6. **Timeline duplicado** entre `TripTimeline.tsx` y `TrackingState.tsx`.
7. **`reset` y `handleCancelTrip`** en hook con lógica duplicada ~80%.
8. **Sin error state en `FormState`** (no results), `ConfirmState` (route/fare fail), `TrackingState` (driver reject fail).
9. **Driver marker no se limpia** en Leaflet en `MapController`.
10. **SSE subscription race** en `useJourneyEngine` si se llama `handleRequestTrip` múltiples veces.
