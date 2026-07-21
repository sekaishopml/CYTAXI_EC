import type { Place, RoutePayload, FareBreakdown, DriverInfo } from "@/types";
import type { RideState } from "@cytaxi/ride-machine";

export interface ValidationResult {
  valid: boolean;
  errors: string[];
  warnings: string[];
}

const STATE_FLOW: RideState[] = [
  "travel_home", "pickup_select", "input", "confirm", "searching", "driver_found",
  "arriving", "arrived", "in_progress", "destination", "payment", "rating", "completed",
];

export function validateJourneyState(curr: RideState, prev: RideState | null): ValidationResult {
  const errors: string[] = [];
  const warnings: string[] = [];

  if (!STATE_FLOW.includes(curr)) {
    errors.push(`Estado desconocido: ${curr}`);
  }

  if (prev && STATE_FLOW.includes(prev)) {
    const currIdx = STATE_FLOW.indexOf(curr);
    const prevIdx = STATE_FLOW.indexOf(prev);
    if (currIdx < prevIdx - 1 && curr !== "pickup_select") {
      warnings.push(`Salto inesperado: ${prev} → ${curr}`);
    }
  }

  return { valid: errors.length === 0, errors, warnings };
}

export function validatePickupData(data: {
  address: string;
  coords: { lat: number; lng: number } | null;
}): ValidationResult {
  const errors: string[] = [];
  const warnings: string[] = [];

  if (!data.address || data.address.trim().length === 0) {
    errors.push("Dirección de recogida requerida");
  }
  if (!data.coords || data.coords.lat === 0) {
    errors.push("Coordenadas de recogida inválidas");
  }
  if (data.address.startsWith("⚠️")) {
    warnings.push("Dirección de recogida sin acceso vial confirmado");
  }

  return { valid: errors.length === 0, errors, warnings };
}

export function validateDestinationData(data: {
  dest: Place | null;
  destQuery: string;
}): ValidationResult {
  const errors: string[] = [];
  const warnings: string[] = [];

  if (!data.dest) {
    errors.push("Destino no seleccionado");
  } else {
    if (!data.dest.lat || !data.dest.lng) {
      errors.push("Coordenadas de destino inválidas");
    }
    if (data.dest.address === data.destQuery && data.dest.name === "Destino seleccionado") {
      warnings.push("Destino seleccionado desde el mapa sin dirección estructurada");
    }
  }

  return { valid: errors.length === 0, errors, warnings };
}

export function validateRouteData(data: {
  route: RoutePayload | null;
  pickupCoords: { lat: number; lng: number } | null;
  dest: Place | null;
}): ValidationResult {
  const errors: string[] = [];
  const warnings: string[] = [];

  if (!data.route) {
    errors.push("Ruta no calculada");
  } else {
    if (data.route.distance_km <= 0) {
      errors.push("Distancia de ruta inválida");
    }
    if (data.route.duration_seconds <= 0) {
      errors.push("Duración de ruta inválida");
    }
    if (!data.route.polyline && data.route.distance_km > 0.1) {
      warnings.push("Ruta sin polyline — posible problema de renderizado en mapa");
    }
  }

  return { valid: errors.length === 0, errors, warnings };
}

export function validateFareData(data: {
  fare: FareBreakdown | null;
  route: RoutePayload | null;
}): ValidationResult {
  const errors: string[] = [];
  const warnings: string[] = [];

  if (!data.fare) {
    errors.push("Tarifa no calculada");
  } else {
    if (data.fare.total <= 0) {
      errors.push("Tarifa total inválida");
    }
    if (data.fare.currency !== "USD") {
      warnings.push(`Moneda no estándar: ${data.fare.currency}`);
    }
    if (data.route && data.fare.distance_km !== data.route.distance_km) {
      warnings.push("Discrepancia entre distancia de ruta y tarifa");
    }
  }

  return { valid: errors.length === 0, errors, warnings };
}

export function validateDriverData(data: {
  driver: DriverInfo | null;
  state: RideState;
}): ValidationResult {
  const errors: string[] = [];
  const warnings: string[] = [];

  const needsDriver: RideState[] = ["driver_found", "arriving", "arrived", "in_progress", "destination"];

  if (needsDriver.includes(data.state) && !data.driver) {
    errors.push(`Conductor requerido en estado ${data.state}`);
  }

  if (data.driver) {
    if (!data.driver.id) errors.push("ID de conductor inválido");
    if (!data.driver.name) warnings.push("Conductor sin nombre");
    if (data.driver.rating < 0 || data.driver.rating > 5) {
      warnings.push(`Rating de conductor fuera de rango: ${data.driver.rating}`);
    }
  }

  return { valid: errors.length === 0, errors, warnings };
}

export function validateTrackingData(data: {
  tripId: string;
  state: RideState;
}): ValidationResult {
  const errors: string[] = [];

  const needsTripId: RideState[] = [
    "searching", "driver_found", "arriving", "arrived",
    "in_progress", "destination", "payment", "rating", "completed",
  ];

  if (needsTripId.includes(data.state) && !data.tripId) {
    errors.push(`Trip ID requerido en estado ${data.state}`);
  }

  return { valid: errors.length === 0, errors, warnings: [] };
}

export function validateFullJourney(snapshot: {
  state: RideState;
  prevState: RideState | null;
  pickupAddress: string;
  pickupCoords: { lat: number; lng: number } | null;
  dest: Place | null;
  destQuery: string;
  route: RoutePayload | null;
  fare: FareBreakdown | null;
  driver: DriverInfo | null;
  tripId: string;
}): ValidationResult {
  const all: ValidationResult = { valid: true, errors: [], warnings: [] };

  const results = [
    validateJourneyState(snapshot.state, snapshot.prevState),
    validatePickupData({ address: snapshot.pickupAddress, coords: snapshot.pickupCoords }),
    validateDestinationData({ dest: snapshot.dest, destQuery: snapshot.destQuery }),
    validateRouteData({ route: snapshot.route, pickupCoords: snapshot.pickupCoords, dest: snapshot.dest }),
    validateFareData({ fare: snapshot.fare, route: snapshot.route }),
    validateDriverData({ driver: snapshot.driver, state: snapshot.state }),
    validateTrackingData({ tripId: snapshot.tripId, state: snapshot.state }),
  ];

  for (const r of results) {
    all.errors.push(...r.errors);
    all.warnings.push(...r.warnings);
  }

  all.valid = all.errors.length === 0;
  return all;
}
