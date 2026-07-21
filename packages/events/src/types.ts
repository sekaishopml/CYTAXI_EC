export type JourneyEventName =
  | "LOCATION_DETECTED"
  | "ROUTE_CALCULATED"
  | "FARE_ESTIMATED"
  | "SEARCHING_DRIVERS"
  | "DRIVER_FOUND"
  | "DRIVER_ASSIGNED"
  | "DRIVER_ACCEPTED"
  | "DRIVER_REJECTED"
  | "DRIVER_ARRIVING"
  | "DRIVER_ARRIVED"
  | "PASSENGER_BOARDED"
  | "TRIP_STARTED"
  | "DESTINATION_ARRIVED"
  | "TRIP_COMPLETED"
  | "TRIP_CANCELLED"
  | "PAYMENT_CONFIRMED"
  | "RATING_SUBMITTED"
  | "JOURNEY_RESET"
  | "ERROR_OCCURRED";

export type AnalyticsEventName =
  | "SCREEN_VIEW"
  | "BUTTON_CLICK"
  | "SEARCH_QUERY"
  | "MAP_INTERACTION"
  | "TRIP_REQUESTED"
  | "TRIP_ACCEPTED"
  | "TRIP_REJECTED"
  | "TRIP_CANCELLED"
  | "PAYMENT_INITIATED"
  | "PAYMENT_COMPLETED"
  | "PAYMENT_FAILED"
  | "RATING_SUBMITTED"
  | "ERROR_LOGGED";

export type RideEventName = JourneyEventName | AnalyticsEventName;

export interface LocationPayload {
  lat: number;
  lng: number;
  address: string;
  source: "gps" | "map_drag" | "search";
}

export interface RoutePayload {
  distance_km: number;
  distance_meters: number;
  duration_seconds: number;
  eta_minutes: number;
  polyline: string;
  pickup: { lat: number; lng: number };
  dest: { lat: number; lng: number };
}

export interface FarePayload {
  base: number;
  distance: number;
  time: number;
  subtotal: number;
  total: number;
  currency: string;
  distance_km: number;
  eta_minutes: number;
  pricing_model: "standard" | "surge" | "premium";
  surge_multiplier?: number;
}

export interface DriverPayload {
  id: string;
  name: string;
  vehicle: string;
  plate: string;
  rating: number;
  photo: string;
  lat: number;
  lng: number;
  eta_seconds: number;
  tier?: "silver" | "gold" | "platinum" | "elite";
  trust_score?: number;
}

export interface TripPayload {
  trip_id: string;
  status: string;
  passenger_id?: string;
  driver_id?: string;
  origin: { lat: number; lng: number; address: string };
  destination: { lat: number; lng: number; address: string };
  fare?: FarePayload;
  created_at: string;
  cancelled_by?: "passenger" | "driver" | "system";
  cancel_reason?: string;
}

export interface PaymentPayload {
  trip_id: string;
  amount: number;
  currency: string;
  method: "cash" | "card" | "wallet";
  status: "pending" | "completed" | "failed";
  transaction_id?: string;
  paid_at?: string;
}

export interface RatingPayload {
  trip_id: string;
  from: "passenger" | "driver";
  to: string;
  score: number;
  comment?: string;
  categories?: { punctuality: number; service: number; comfort: number };
}

export interface ErrorPayload {
  code: string;
  message: string;
  context?: Record<string, unknown>;
  timestamp: string;
}

export interface AnalyticsPayload {
  event: AnalyticsEventName;
  properties?: Record<string, unknown>;
  timestamp: string;
}

export type EventPayloadMap = {
  LOCATION_DETECTED: LocationPayload;
  ROUTE_CALCULATED: RoutePayload;
  FARE_ESTIMATED: FarePayload;
  SEARCHING_DRIVERS: { count?: number };
  DRIVER_FOUND: DriverPayload;
  DRIVER_ASSIGNED: DriverPayload;
  DRIVER_ACCEPTED: DriverPayload;
  DRIVER_REJECTED: { driver_id: string };
  DRIVER_ARRIVING: DriverPayload;
  DRIVER_ARRIVED: DriverPayload;
  PASSENGER_BOARDED: { driver_id: string; trip_id: string };
  TRIP_STARTED: TripPayload;
  DESTINATION_ARRIVED: TripPayload;
  TRIP_COMPLETED: TripPayload;
  TRIP_CANCELLED: TripPayload;
  PAYMENT_CONFIRMED: PaymentPayload;
  RATING_SUBMITTED: RatingPayload;
  JOURNEY_RESET: Record<string, never>;
  ERROR_OCCURRED: ErrorPayload;
  SCREEN_VIEW: AnalyticsPayload;
  BUTTON_CLICK: AnalyticsPayload;
  SEARCH_QUERY: AnalyticsPayload;
  MAP_INTERACTION: AnalyticsPayload;
  TRIP_REQUESTED: AnalyticsPayload;
  TRIP_ACCEPTED: AnalyticsPayload;
  TRIP_REJECTED: AnalyticsPayload;
  TRIP_CANCELLED_A: AnalyticsPayload;
  PAYMENT_INITIATED: AnalyticsPayload;
  PAYMENT_COMPLETED: AnalyticsPayload;
  PAYMENT_FAILED: AnalyticsPayload;
  RATING_SUBMITTED_A: AnalyticsPayload;
  ERROR_LOGGED: AnalyticsPayload;
};

export interface AuditLogEntry {
  id: string;
  event: RideEventName;
  payload: unknown;
  timestamp: string;
  source?: string;
  session_id?: string;
}
