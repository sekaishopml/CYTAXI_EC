import { Coordinates, Place, TripRequest } from "@/types";
import type { FarePayload } from "@cytaxi/events";

const API = typeof window !== "undefined"
  ? `${window.location.protocol}//${window.location.host}/api/v1`
  : "http://64.176.219.221/api/v1";

const DEFAULT_TIMEOUT = 8000;
const MAX_RETRIES = 2;

interface RetryConfig {
  timeoutMs: number;
  maxRetries: number;
}

async function req<T>(
  path: string,
  opts: RequestInit = {},
  cfg: Partial<RetryConfig> = {},
): Promise<T> {
  const timeoutMs = cfg.timeoutMs ?? DEFAULT_TIMEOUT;
  const maxRetries = cfg.maxRetries ?? MAX_RETRIES;

  let lastError: Error | null = null;

  for (let attempt = 0; attempt <= maxRetries; attempt++) {
    if (attempt > 0) {
      await new Promise(r => setTimeout(r, Math.min(1000 * Math.pow(2, attempt - 1), 4000)));
    }
    try {
      const ctrl = new AbortController();
      const timer = setTimeout(() => ctrl.abort(), timeoutMs);
      const res = await fetch(`${API}${path}`, {
        headers: { "Content-Type": "application/json" },
        signal: ctrl.signal,
        ...opts,
      });
      clearTimeout(timer);
      if (!res.ok) throw new Error(`HTTP_${res.status}`);
      return res.json();
    } catch (err: any) {
      lastError = err;
      if (err.name === "AbortError") {
        throw new Error("TIMEOUT");
      }
    }
  }
  throw lastError || new Error("REQUEST_FAILED");
}

export async function searchPlaces(query: string): Promise<Place[]> {
  if (query.length < 3) return [];
  try {
    const data: any = await req(`/geo/search?q=${encodeURIComponent(query)}`, {}, { timeoutMs: 5000 });
    const items: any[] = (Array.isArray(data) ? data : (data?.Places || data?.places || data?.results || []));
    if (!Array.isArray(items) || items.length === 0) return [];
    return items.map((p: any) => {
      if (!p || typeof p !== "object") return null;
      const name = p.Name || p.name || p.DisplayName?.split(",")[0] || "Unknown";
      const address = p.Address || p.address || p.DisplayName || p.FormattedAddress || "";
      const coords = p.Coordinates || p.coordinates || p;
      const lat = coords?.Lat ?? coords?.lat ?? null;
      const lng = coords?.Lng ?? coords?.lng ?? null;
      if (lat === null || lng === null) return null;
      return { name: String(name), address: String(address), lat: Number(lat), lng: Number(lng) } satisfies Place;
    }).filter(Boolean) as Place[];
  } catch { return []; }
}

export async function calculateRoute(
  origin: Coordinates, dest: Coordinates,
): Promise<{ distance_meters: number; duration_seconds: number; polyline: string; distance_km: number; eta_minutes: number } | null> {
  try {
    const data = await req("/geo/route", {
      method: "POST",
      body: JSON.stringify({
        origin_lat: origin.lat, origin_lng: origin.lng,
        dest_lat: dest.lat, dest_lng: dest.lng,
      }),
    });
    const route = Array.isArray(data) ? data[0] : data;
    if (!route) return null;
    const dist = route.distance || route.Distance || {};
    const dur = route.duration || route.Duration || {};
    const m = dist.meters || dist.Meters || route.distance_meters || 5000;
    const s = dur.seconds || dur.Seconds || route.duration_seconds || 600;
    const poly = route.polyline || route.Polyline || route.Geometry || "";
    return {
      distance_meters: m, duration_seconds: s,
      polyline: poly,
      distance_km: m / 1000, eta_minutes: Math.ceil(s / 60),
    };
  } catch { return null; }
}

export async function estimateFare(distanceKm: number, durationSec: number): Promise<FarePayload> {
  try {
    const data: any = await req("/pricing/estimate", {
      method: "POST",
      body: JSON.stringify({ distance_km: distanceKm, duration_sec: durationSec, region: "ecuador" }),
    });
    const f = data.fare || data;
    return {
      base: f.base || f.BaseFare || 1,
      distance: f.distance || f.DistanceFare || 0,
      time: f.time || f.TimeFare || 0,
      subtotal: f.subtotal || f.Subtotal || 0,
      total: f.total || f.Total || 0,
      currency: f.currency || f.Currency || "USD",
      distance_km: distanceKm,
      eta_minutes: Math.ceil(durationSec / 60),
      pricing_model: "standard",
    };
  } catch {
    return {
      base: 1, distance: distanceKm * 0.5, time: durationSec * 0.02,
      subtotal: 0, total: 1 + distanceKm * 0.5 + durationSec * 0.02,
      currency: "USD", distance_km: distanceKm,
      eta_minutes: Math.ceil(durationSec / 60), pricing_model: "standard",
    };
  }
}

export async function requestTrip(data: TripRequest): Promise<{ trip_id: string; status: string }> {
  return req("/trip/request", {
    method: "POST",
    body: JSON.stringify({ ...data, customer_id: `cust_${data.phone}` }),
  });
}

export async function startMatching(tripId: string, lat: number, lng: number): Promise<{ matching_id: string; candidates: any[] }> {
  return req("/matching/start", {
    method: "POST",
    body: JSON.stringify({ trip_id: tripId, pickup_lat: lat, pickup_lng: lng }),
  });
}
