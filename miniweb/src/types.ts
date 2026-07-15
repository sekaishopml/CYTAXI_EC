export type TripState = "idle" | "confirm" | "searching" | "driver_found" | "in_progress" | "completed";

export interface Coordinates {
  lat: number;
  lng: number;
}

export interface Place {
  name: string;
  address: string;
  lat: number;
  lng: number;
}

export interface DriverInfo {
  id: string;
  name: string;
  vehicle: string;
  plate: string;
  rating: number;
  photo: string;
  eta_seconds: number;
  lat?: number;
  lng?: number;
}

export interface FareBreakdown {
  base: number;
  distance: number;
  time: number;
  subtotal: number;
  total: number;
  currency: string;
  distance_km: number;
  eta_minutes: number;
}

export interface TrackingUpdate {
  type: string;
  trip_id: string;
  status: string;
  driver?: {
    id: string;
    name: string;
    vehicle: string;
    plate: string;
    lat: number;
    lng: number;
    rating: number;
  };
  eta_seconds?: number;
  timestamp: string;
}

export interface TripRequest {
  phone: string;
  passenger_name: string;
  origin_address: string;
  origin_lat: number;
  origin_lng: number;
  dest_address: string;
  dest_lat: number;
  dest_lng: number;
}
