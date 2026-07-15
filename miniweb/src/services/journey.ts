const GATEWAY = process.env.NEXT_PUBLIC_GATEWAY_URL || "http://localhost:8000";

export interface JourneyRequest {
  phone: string;
  passenger_name: string;
  origin_address: string;
  origin_lat: number;
  origin_lng: number;
  dest_address: string;
  dest_lat: number;
  dest_lng: number;
}

export interface JourneyResponse {
  session_id: string;
  trip_status: string;
  fare_estimate?: {
    base: number;
    distance: number;
    subtotal: number;
    total: number;
    currency: string;
  };
  timeline: Array<{ step: string; status: string; duration: string }>;
  success: boolean;
  error?: string;
}

export interface TripStatus {
  id: string;
  status: string;
  origin: string;
  destination: string;
  fare: string;
  created_at: string;
}

export async function startConversation(phone: string, message: string): Promise<{ session_id?: string }> {
  const res = await fetch(`${GATEWAY}/api/v1/conversation/start`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ phone, message }),
  });
  if (!res.ok) throw new Error("Failed to start conversation");
  return res.json();
}

export async function requestTrip(data: Omit<JourneyRequest, "phone"> & { phone: string; customer_id: string }): Promise<{ status: string }> {
  const res = await fetch(`${GATEWAY}/api/v1/trip/request`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(data),
  });
  if (!res.ok) throw new Error("Failed to request trip");
  return res.json();
}

export async function estimateFare(tripId: string, distanceKm: number, durationSec: number): Promise<{ fare: JourneyResponse["fare_estimate"] }> {
  const res = await fetch(`${GATEWAY}/api/v1/pricing/estimate`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ trip_id: tripId, distance_km: distanceKm, duration_sec: durationSec, region: "ecuador" }),
  });
  if (!res.ok) throw new Error("Failed to estimate fare");
  return res.json();
}

export async function getTripStatus(tripId: string): Promise<TripStatus> {
  const res = await fetch(`${GATEWAY}/api/v1/trip/trips/${tripId}`);
  if (!res.ok) throw new Error("Trip not found");
  return res.json();
}

export async function executeFullJourney(data: JourneyRequest): Promise<JourneyResponse> {
  const timeline: JourneyResponse["timeline"] = [];
  const start = performance.now();

  try {
    timeline.push({ step: "start_conversation", status: "running", duration: "0ms" });
    const conv = await startConversation(data.phone, "Hola, quiero un taxi");
    timeline[0] = { ...timeline[0], status: "completed", duration: `${(performance.now() - start).toFixed(0)}ms` };

    timeline.push({ step: "create_trip", status: "running", duration: "0ms" });
    const trip = await requestTrip({ ...data, customer_id: `cust_${data.phone}` });
    timeline[1] = { ...timeline[1], status: "completed", duration: `${(performance.now() - start).toFixed(0)}ms` };

    timeline.push({ step: "estimate_fare", status: "running", duration: "0ms" });
    const fare = await estimateFare(`trip_${data.phone}`, 5.5, 900);
    timeline[2] = { ...timeline[2], status: "completed", duration: `${(performance.now() - start).toFixed(0)}ms` };

    return {
      session_id: conv.session_id || "",
      trip_status: "created",
      fare_estimate: fare.fare,
      timeline,
      success: true,
    };
  } catch (e) {
    return {
      session_id: "",
      trip_status: "failed",
      timeline,
      success: false,
      error: (e as Error).message,
    };
  }
}
