import { colors } from "@cytaxi/design-tokens";

export type RideState =
  | "travel_home"
  | "pickup_select"
  | "input"
  | "confirm"
  | "searching"
  | "driver_found"
  | "arriving"
  | "arrived"
  | "in_progress"
  | "destination"
  | "payment"
  | "rating"
  | "completed";

export type RideEvent =
  | "START_TRIP"
  | "SELECT_PICKUP"
  | "SELECT_DEST"
  | "CONFIRM"
  | "REQUEST"
  | "DRIVER_FOUND"
  | "DRIVER_ACCEPTED"
  | "DRIVER_ARRIVING"
  | "DRIVER_ARRIVED"
  | "TRIP_START"
  | "TRIP_ARRIVED"
  | "TRIP_COMPLETE"
  | "PAYMENT_DONE"
  | "RATING_DONE"
  | "CANCEL"
  | "REJECT"
  | "RESET"
  | "SELECT_VEHICLE"
  | "APPLY_COUPON"
  | "SET_NOTE"
  | "SCHEDULE_TRIP";

export interface RideStateConfig {
  id: RideState;
  label: string;
  labelShort: string;
  icon: string;
  color: string;
  bgColor: string;
  description: string;
  showMap: boolean;
  showPin: boolean;
  mapInteractive: boolean;
  showNavbar: boolean;
}

export const RIDE_STATES: Record<RideState, RideStateConfig> = {
  travel_home: {
    id: "travel_home",
    label: "Viajá con CYTAXI",
    labelShort: "Inicio",
    icon: "🏠",
    color: colors.brand.green,
    bgColor: colors.brand.greenBg,
    description: "Elegí cómo moverte hoy",
    showMap: false,
    showPin: false,
    mapInteractive: false,
    showNavbar: true,
  },
  pickup_select: {
    id: "pickup_select",
    label: "Selecciona tu ubicación",
    labelShort: "Origen",
    icon: "📍",
    color: colors.brand.green,
    bgColor: colors.brand.greenBg,
    description: "Arrastra el mapa para elegir tu punto de partida",
    showMap: true,
    showPin: true,
    mapInteractive: true,
    showNavbar: true,
  },
  input: {
    id: "input",
    label: "¿A dónde vamos?",
    labelShort: "Destino",
    icon: "🎯",
    color: colors.status.info,
    bgColor: "rgba(68,138,255,0.08)",
    description: "Ingresa tu destino o selecciónalo en el mapa",
    showMap: true,
    showPin: true,
    mapInteractive: true,
    showNavbar: true,
  },
  confirm: {
    id: "confirm",
    label: "Confirma tu viaje",
    labelShort: "Confirmar",
    icon: "📋",
    color: colors.brand.green,
    bgColor: colors.brand.greenBg,
    description: "Revisa la ruta y el precio antes de solicitar",
    showMap: true,
    showPin: false,
    mapInteractive: false,
    showNavbar: true,
  },
  searching: {
    id: "searching",
    label: "Buscando conductor",
    labelShort: "Buscando",
    icon: "🔍",
    color: colors.status.info,
    bgColor: "rgba(68,138,255,0.08)",
    description: "Conectando con los mejores conductores cercanos",
    showMap: true,
    showPin: false,
    mapInteractive: false,
    showNavbar: false,
  },
  driver_found: {
    id: "driver_found",
    label: "Conductor encontrado",
    labelShort: "Conductor",
    icon: "🚗",
    color: colors.brand.green,
    bgColor: colors.brand.greenBg,
    description: "Tu conductor está en camino",
    showMap: true,
    showPin: false,
    mapInteractive: false,
    showNavbar: false,
  },
  arriving: {
    id: "arriving",
    label: "Conductor llegando",
    labelShort: "Llegando",
    icon: "🔄",
    color: colors.status.warning,
    bgColor: "rgba(255,193,7,0.08)",
    description: "El conductor está cerca de tu ubicación",
    showMap: true,
    showPin: false,
    mapInteractive: false,
    showNavbar: false,
  },
  arrived: {
    id: "arrived",
    label: "El conductor ha llegado",
    labelShort: "Llegó",
    icon: "✅",
    color: colors.status.success,
    bgColor: "rgba(0,161,82,0.08)",
    description: "Tu conductor te está esperando",
    showMap: true,
    showPin: false,
    mapInteractive: false,
    showNavbar: false,
  },
  in_progress: {
    id: "in_progress",
    label: "En viaje",
    labelShort: "Viajando",
    icon: "🏁",
    color: colors.brand.green,
    bgColor: colors.brand.greenBg,
    description: "Rumbo a tu destino",
    showMap: true,
    showPin: false,
    mapInteractive: false,
    showNavbar: false,
  },
  destination: {
    id: "destination",
    label: "Llegando al destino",
    labelShort: "Casi llega",
    icon: "📍",
    color: colors.status.info,
    bgColor: "rgba(68,138,255,0.08)",
    description: "Estás por llegar a tu destino",
    showMap: true,
    showPin: false,
    mapInteractive: false,
    showNavbar: false,
  },
  payment: {
    id: "payment",
    label: "Procesando pago",
    labelShort: "Pago",
    icon: "💳",
    color: colors.brand.green,
    bgColor: colors.brand.greenBg,
    description: "Procesando tu pago",
    showMap: false,
    showPin: false,
    mapInteractive: false,
    showNavbar: true,
  },
  rating: {
    id: "rating",
    label: "Califica tu viaje",
    labelShort: "Calificar",
    icon: "⭐",
    color: colors.status.warning,
    bgColor: "rgba(255,193,7,0.08)",
    description: "¿Cómo fue tu experiencia?",
    showMap: false,
    showPin: false,
    mapInteractive: false,
    showNavbar: true,
  },
  completed: {
    id: "completed",
    label: "Viaje completado",
    labelShort: "Completado",
    icon: "🎉",
    color: colors.status.success,
    bgColor: "rgba(0,161,82,0.08)",
    description: "Gracias por viajar con nosotros",
    showMap: false,
    showPin: false,
    mapInteractive: false,
    showNavbar: true,
  },
};

export interface Transition {
  from: RideState[];
  event: RideEvent;
  to: RideState;
}

export const RIDE_TRANSITIONS: Transition[] = [
  { from: ["travel_home"], event: "START_TRIP", to: "pickup_select" },
  { from: ["pickup_select"], event: "SELECT_PICKUP", to: "input" },
  { from: ["input"], event: "SELECT_DEST", to: "confirm" },
  { from: ["input"], event: "CANCEL", to: "pickup_select" },
  { from: ["confirm"], event: "CONFIRM", to: "input" },
  { from: ["confirm"], event: "REQUEST", to: "searching" },
  { from: ["confirm"], event: "SELECT_VEHICLE", to: "confirm" },
  { from: ["confirm"], event: "APPLY_COUPON", to: "confirm" },
  { from: ["confirm"], event: "SET_NOTE", to: "confirm" },
  { from: ["confirm"], event: "SCHEDULE_TRIP", to: "confirm" },
  { from: ["searching"], event: "DRIVER_FOUND", to: "driver_found" },
  { from: ["searching"], event: "CANCEL", to: "pickup_select" },
  { from: ["driver_found"], event: "DRIVER_ACCEPTED", to: "arriving" },
  { from: ["driver_found"], event: "REJECT", to: "input" },
  { from: ["arriving"], event: "DRIVER_ARRIVING", to: "arrived" },
  { from: ["arriving"], event: "CANCEL", to: "pickup_select" },
  { from: ["arrived"], event: "TRIP_START", to: "in_progress" },
  { from: ["in_progress"], event: "TRIP_ARRIVED", to: "destination" },
  { from: ["destination"], event: "TRIP_COMPLETE", to: "payment" },
  { from: ["payment"], event: "PAYMENT_DONE", to: "rating" },
  { from: ["rating"], event: "RATING_DONE", to: "completed" },
  { from: ["completed"], event: "RESET", to: "travel_home" },
  { from: ["driver_found", "arriving", "arrived", "in_progress"], event: "CANCEL", to: "travel_home" },
];

export function transitionRide(current: RideState, event: RideEvent): RideState | null {
  for (const t of RIDE_TRANSITIONS) {
    if (t.event === event && t.from.includes(current)) {
      return t.to;
    }
  }
  return null;
}

export function getStateConfig(state: RideState): RideStateConfig {
  return RIDE_STATES[state];
}

export function getTimelineSteps(currentState: RideState) {
  const allStates: RideState[] = [
    "searching", "driver_found", "arriving", "arrived",
    "in_progress", "destination", "payment", "completed",
  ];
  const currentIdx = allStates.indexOf(currentState);

  return allStates.map((s, i) => ({
    id: s,
    label: RIDE_STATES[s].labelShort,
    description: RIDE_STATES[s].description,
    icon: RIDE_STATES[s].icon,
    status: i < currentIdx ? "completed" as const : i === currentIdx ? "active" as const : "pending" as const,
    color: RIDE_STATES[s].color,
  }));
}

export const stateAnimations: Record<RideState, { enter: string; exit: string; duration: number }> = {
  travel_home: { enter: "fadeIn 0.4s ease-out", exit: "fadeOut 0.15s ease-in", duration: 400 },
  pickup_select: { enter: "fadeIn 0.35s ease-out", exit: "fadeOut 0.18s ease-in", duration: 350 },
  input: { enter: "slideUp 0.35s ease-out", exit: "slideDown 0.18s ease-in", duration: 350 },
  confirm: { enter: "scaleIn 0.4s ease-out", exit: "scaleOut 0.18s ease-in", duration: 400 },
  searching: { enter: "fadeIn 0.3s ease-out", exit: "fadeOut 0.15s ease-in", duration: 300 },
  driver_found: { enter: "slideUp 0.5s ease-out", exit: "fadeOut 0.2s ease-in", duration: 500 },
  arriving: { enter: "fadeIn 0.3s ease-out", exit: "fadeOut 0.15s ease-in", duration: 300 },
  arrived: { enter: "scaleIn 0.4s ease-out", exit: "fadeOut 0.15s ease-in", duration: 400 },
  in_progress: { enter: "fadeIn 0.35s ease-out", exit: "fadeOut 0.15s ease-in", duration: 350 },
  destination: { enter: "fadeIn 0.3s ease-out", exit: "fadeOut 0.15s ease-in", duration: 300 },
  payment: { enter: "slideUp 0.35s ease-out", exit: "fadeOut 0.15s ease-in", duration: 350 },
  rating: { enter: "scaleIn 0.5s ease-out", exit: "fadeOut 0.2s ease-in", duration: 500 },
  completed: { enter: "scaleIn 0.6s ease-out", exit: "fadeOut 0.2s ease-in", duration: 600 },
};

export interface JourneyEngineState {
  state: RideState;
  prevState: RideState | null;
  direction: "forward" | "back";
  isTransitioning: boolean;
}

export type JourneyListener = (newState: RideState, prevState: RideState | null, event: RideEvent) => void;

export class JourneyEngine {
  private _state: RideState;
  private _prevState: RideState | null = null;
  private _direction: "forward" | "back" = "forward";
  private _isTransitioning = false;
  private _listeners: Set<JourneyListener> = new Set();

  constructor(initial: RideState = "travel_home") {
    this._state = initial;
  }

  get snapshot(): JourneyEngineState {
    return {
      state: this._state,
      prevState: this._prevState,
      direction: this._direction,
      isTransitioning: this._isTransitioning,
    };
  }

  get state() { return this._state; }

  send(event: RideEvent): RideState | null {
    const next = transitionRide(this._state, event);
    if (!next || next === this._state) return null;

    this._prevState = this._state;
    this._direction = this._isForward(event) ? "forward" : "back";
    this._isTransitioning = true;
    this._state = next;
    this._emit(event);
    return next;
  }

  goTo(state: RideState): boolean {
    if (state === this._state) return false;
    this._prevState = this._state;
    this._direction = "forward";
    this._isTransitioning = true;
    this._state = state;
    this._emit("RESET");
    return true;
  }

  endTransition() {
    this._isTransitioning = false;
  }

  onTransition(listener: JourneyListener): () => void {
    this._listeners.add(listener);
    return () => this._listeners.delete(listener);
  }

  private _emit(event: RideEvent) {
    this._listeners.forEach(fn => fn(this._state, this._prevState, event));
  }

  private _isForward(event: RideEvent): boolean {
    return !["CANCEL", "REJECT", "CONFIRM", "RESET"].includes(event);
  }

  reset() {
    this._state = "travel_home";
    this._prevState = null;
    this._direction = "forward";
    this._isTransitioning = false;
    this._listeners.clear();
  }
}

export function createJourneyEngine(initial?: RideState): JourneyEngine {
  return new JourneyEngine(initial);
}
