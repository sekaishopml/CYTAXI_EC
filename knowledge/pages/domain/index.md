# Domain Model

Basado en CYDIGITAL-BLUEPRINT y la implementación en CYTAXI_EC.

## Agregados Principales

### Trip (Viaje)
- **Domain Events**: TripRequested, TripCreated, TripStarted, PassengerBoarded, TripCompleted, TripCancelled
- **Estados**: solicitando → creado → conductor_asignado → conductor_llegando → pasajero_a_bordo → en_viaje → completado / cancelado
- **Value Objects**: Location (lat/lng/dirección), Route, Fare, Distance, Duration

### Customer (Cliente)
- Atributos: profile, favoritos, recentTrips, preferencias, emergencyContacts, loyalty
- **Domain Events**: CustomerCreated, CustomerUpdated, FavoriteLocationAdded

### Driver (Conductor)
- Atributos: profile, vehículo, documentos, status, ratings, disponibilidad, wallet
- **Domain Events**: DriverRegistered, DriverOnline, DriverOffline, DriverAcceptedTrip

### Fare (Tarifa)
- **Determinístico** — calculado por reglas de negocio, no por IA
- Factores: base, distancia, tiempo, demanda (surge), promociones
- **Value Objects**: BaseFare, DistanceFare, TimeFare, SurgeMultiplier

## State Machine (Ride Machine)

```
pickup_select → destination_select → trip_preview → requesting
    → trip_created → driver_assigned → driver_arriving
    → passenger_picked_up → in_transit → trip_complete
    → rating → payment
```

- `packages/ride-machine/` contiene la implementación
- Los eventos se disparan desde el frontend y se procesan en el backend
- Estados visibles al usuario vs estados internos del sistema

## Value Objects Compartidos

- **Location**: `{ lat, lng, address, placeId }`
- **Route**: `{ origin, destination, distance, duration, polyline }`
- **Price**: `{ amount, currency: "USD", breakdown }`
- **Rating**: `{ score, comment, timestamp }`
- **Vehicle**: `{ type, plate, color, model }`
