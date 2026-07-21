# UX / Diseño

## Identidad Visual

- **Tema**: Cobalt Hallmark
- **Color principal**: `#3b82f6` (azul cobalt)
- **Fondo**: `#f5f6f8` (paper)
- **Texto**: `#1e293b` (ink)

## Tipografía

| Uso | Fuente | Peso |
|-----|--------|------|
| Display/Headings | Space Grotesk | 500-700 |
| Body | Inter | 400-600 |
| Labels/Mono | JetBrains Mono | 400 |

## Espaciado y Formas

- **Border radius**: 6px (small), 10px (medium)
- **Bordes**: hairline (1px, `rgba(0,0,0,0.08)`)
- **Sin glassmorphism**
- **Easing**: `cubic-bezier(0.16, 1, 0.3, 1)`

## Principios de UI

1. **Mobile-first** — todo diseñado para móvil primero
2. **Conversacional** — la interfaz imita una conversación natural
3. **Minimalista** — sin elementos decorativos innecesarios
4. **Feedback inmediato** — cada acción tiene respuesta visual
5. **Estado visible** — el state machine muestra el progreso del viaje

## Flujo de Usuario (Miniweb)

1. Pantalla de carga → mapa centrado en ubicación del usuario
2. Pin de pickup en ubicación actual (ajustable)
3. Barra de búsqueda para destino
4. Confirmación de ruta con precio estimado
5. "Solicitar viaje" → búsqueda de conductor
6. Tracking en tiempo real del conductor
7. Check-in al iniciar viaje
8. Checkout al llegar a destino
9. Pantalla de pago y calificación

## Texto de UI

- **Idioma**: Español (latinoamericano)
- **Moneda**: USD ($)
- **Distancia**: km
- **Tiempo**: formato local (12h o 24h)
