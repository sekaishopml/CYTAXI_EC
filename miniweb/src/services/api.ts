import { Coordinates, Place, FareBreakdown, TripRequest } from "@/types";

const API = typeof window !== "undefined"
  ? `${window.location.protocol}//${window.location.host}/api/v1`
  : "http://64.176.219.221/api/v1";

async function req<T>(path: string, opts: RequestInit = {}): Promise<T> {
  const res = await fetch(`${API}${path}`, {
    headers: { "Content-Type": "application/json" },
    ...opts,
  });
  if (!res.ok) throw new Error(`${res.status}`);
  return res.json();
}

export async function searchPlaces(query: string): Promise<Place[]> {
  if (query.length < 3) return [];
  try {
    const data: any = await req(`/geo/search?q=${encodeURIComponent(query)}`);
    return (data.Places || data.places || data || []).map((p: any) => ({
      name: p.Name || p.Name || p.DisplayName?.split(",")[0] || "Unknown",
      address: p.Address || p.DisplayName || p.FormattedAddress || "",
      lat: p.Coordinates?.Lat || p.Lat || p.lat || 0,
      lng: p.Coordinates?.Lng || p.Lng || p.lng || 0,
    })) as Place[];
  } catch { return []; }
}

export async function calculateRoute(origin: Coordinates, dest: Coordinates): Promise<{ distance_meters: number; duration_seconds: number; polyline: string; distance_km: number; eta_minutes: number } | null> {
  try {
    const data = await req("/geo/route", {
      method: "POST",
      body: JSON.stringify({ origin_lat: origin.lat, origin_lng: origin.lng, dest_lat: dest.lat, dest_lng: dest.lng }),
    });
    const route = Array.isArray(data) ? data[0] : data;
    if (!route) return null;
    const m = route.Distance?.Meters || route.Distance || 5000;
    const s = route.Duration?.Seconds || route.Duration || 600;
    return { distance_meters: m, duration_seconds: s, polyline: route.Polyline || route.Geometry || "", distance_km: m / 1000, eta_minutes: Math.ceil(s / 60) };
  } catch { return null; }
}

export async function estimateFare(distanceKm: number, durationSec: number): Promise<FareBreakdown> {
  try {
    const data: any = await req("/pricing/estimate", {
      method: "POST",
      body: JSON.stringify({ distance_km: distanceKm, duration_sec: durationSec, region: "ecuador" }),
    });
    const f = data.fare || data;
    return { base: f.base || f.BaseFare || 1, distance: f.distance || f.DistanceFare || 0, time: f.time || f.TimeFare || 0, subtotal: f.subtotal || f.Subtotal || 0, total: f.total || f.Total || 0, currency: f.currency || f.Currency || "USD", distance_km: distanceKm, eta_minutes: Math.ceil(durationSec / 60) };
  } catch { return { base: 1, distance: distanceKm * 0.5, time: durationSec * 0.02, subtotal: 0, total: 1 + distanceKm * 0.5 + durationSec * 0.02, currency: "USD", distance_km: distanceKm, eta_minutes: Math.ceil(durationSec / 60) }; }
}

export async function requestTrip(data: TripRequest): Promise<{ trip_id: string; status: string }> {
  return req("/trip/request", { method: "POST", body: JSON.stringify({ ...data, customer_id: `cust_${data.phone}` }) });
}

export async function startMatching(tripId: string, lat: number, lng: number): Promise<{ matching_id: string; candidates: any[] }> {
  return req("/matching/start", { method: "POST", body: JSON.stringify({ trip_id: tripId, pickup_lat: lat, pickup_lng: lng }) });
}
