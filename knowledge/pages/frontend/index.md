# Frontend: Travel (formerly Miniweb)

## Stack

- **Framework**: Next.js 14 (App Router)
- **Lenguaje**: TypeScript
- **Mapas**: Leaflet (OpenStreetMap)
- **Geocoding**: OpenStreetMap Nominatim
- **Animaciones**: framer-motion
- **State Machine**: `@cytaxi/ride-machine`
- **Diseño**: Cobalt Hallmark

## Diseño (Cobalt Hallmark)

- **Accent**: `#3b82f6`
- **Paper**: `#f5f6f8`
- **Fonts**: Space Grotesk (display), Inter (body), JetBrains Mono (labels/mono)
- **Radii**: 6px / 10px
- **Easing**: `cubic-bezier(0.16,1,0.3,1)`
- **Borders**: hairline (sin glassmorphism)
- **Breakpoints**: mobile-first

## Estructura

```
travel/
  src/
    app/              → App Router pages (Next.js 14)
    components/       → UI components
    hooks/            → Custom hooks
    lib/              → Utilities, API clients
    state/            → State machine integration
    styles/           → Global styles, tokens
    types/            → TypeScript types
  public/             → Static assets
```

## State Machine (useJourneyEngine)

El viaje se maneja como una máquina de estados:

1. `pickup_select` — selección de origen en mapa
2. `destination_select` — búsqueda y selección de destino
3. `trip_preview` — vista previa del viaje (ruta, precio estimado)
4. `requesting` — solicitando conductor
5. `trip_created` — viaje creado
6. `driver_assigned` — conductor asignado
7. `driver_arriving` — conductor llegando
8. `passenger_picked_up` — pasajero a bordo
9. `in_transit` — en viaje
10. `trip_complete` — viaje completado
11. `rating` — calificación
12. `payment` — pago

## UI Text

- **Idioma**: Español (latinoamericano)
- Todas las cadenas de UI en español
- Formato de moneda: USD ($)
- Formato de distancia: km
- Formato de hora: 12h (AM/PM) o 24h según región

## Paquetes

- `@cytaxi/ride-machine` — state machine viaje
- `@cytaxi/map-engine` — lógica de mapas Leaflet
- `@cytaxi/design-tokens` — tokens de diseño cobalt
