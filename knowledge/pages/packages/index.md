# Paquetes Compartidos

Monorepo con 3 paquetes publicados como `@cytaxi/*`.

## `@cytaxi/ride-machine`

State machine para el ciclo de vida del viaje.

```
pickup_select → destination_select → trip_preview → requesting
    → trip_created → driver_assigned → driver_arriving
    → passenger_picked_up → in_transit → trip_complete
    → rating → payment
```

- Define estados, transiciones, y side-effects
- Usado por el frontend (miniweb) para manejar el flujo del viaje
- Integrado via `useJourneyEngine` hook

## `@cytaxi/map-engine`

Lógica de mapas independiente del proveedor.

- Leaflet (actual) como proveedor default
- Abstraction layer para cambiar de proveedor sin cambiar lógica de negocio
- Geocoding, markers, rutas, bounds, eventos de mapa

## `@cytaxi/design-tokens`

Tokens de diseño del tema Cobalt Hallmark.

```json
{
  "colors": {
    "accent": "#3b82f6",
    "paper": "#f5f6f8",
    "ink": "#1e293b"
  },
  "fonts": {
    "display": "Space Grotesk",
    "body": "Inter",
    "mono": "JetBrains Mono"
  },
  "radii": {
    "sm": "6px",
    "md": "10px"
  },
  "easing": "cubic-bezier(0.16, 1, 0.3, 1)"
}
```
