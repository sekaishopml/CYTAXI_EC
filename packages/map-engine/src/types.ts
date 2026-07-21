export interface LatLng {
  lat: number;
  lng: number;
}

export interface DriverMarkerData {
  id: string;
  lat: number;
  lng: number;
  heading?: number;
  name?: string;
  vehicle?: string;
  plate?: string;
  eta_seconds?: number;
  status?: "idle" | "searching" | "arriving" | "waiting" | "driving";
  pulse?: boolean;
}

export interface RouteData {
  polyline: string;
  distance_km: number;
  duration_seconds: number;
}

export interface MapConfig {
  center: LatLng;
  zoom: number;
  styles?: unknown[];
  gestureHandling?: "greedy" | "cooperative" | "none";
}

export type MapEngineCallback = (map: google.maps.Map) => void;

export interface MarkerAnimationConfig {
  bounce?: boolean;
  pulse?: boolean;
  pulseColor?: string;
  pulseDuration?: number;
  rotate?: boolean;
  label?: string;
}
