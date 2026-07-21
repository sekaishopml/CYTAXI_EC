import type { RideState, RideEvent } from "@cytaxi/ride-machine";
import type { DriverPayload, FarePayload, RoutePayload } from "@cytaxi/events";

export type { RideState, RideEvent };
export type TripState = RideState;

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

export type DriverInfo = DriverPayload;

export type FareBreakdown = FarePayload;

export interface TrackingUpdate {
  type: string;
  trip_id: string;
  status: string;
  driver?: DriverPayload;
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

export type { RoutePayload };
