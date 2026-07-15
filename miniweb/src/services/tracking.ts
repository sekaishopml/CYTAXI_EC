import { TrackingUpdate } from "@/types";

const API = typeof window !== "undefined"
  ? `${window.location.protocol}//${window.location.host}/api/v1`
  : "http://64.176.219.221/api/v1";

export function subscribeToTrip(tripId: string, onUpdate: (data: TrackingUpdate) => void): () => void {
  const es = new EventSource(`${API}/trip/ws?trip_id=${tripId}`);

  es.onmessage = (event) => {
    try {
      const data: TrackingUpdate = JSON.parse(event.data);
      onUpdate(data);
      if (data.type === "trip_completed") es.close();
    } catch (e) { console.warn("SSE parse error", e); }
  };
  es.onerror = () => { setTimeout(() => es.close(), 2000); };

  return () => es.close();
}

console.info(
  "%c📡 CYTAXI Tracking %cSSE vs WebSocket",
  "font-weight:bold;color:#16a34a",
  "color:#6b7280;font-size:12px"
);
console.info(
  "%cSSE ✓%c Unidirectional (server→client) — perfect for tracking\n" +
  "%cSSE ✓%c Native EventSource API, auto-reconnect, HTTP/2 compatible\n" +
  "%cSSE ✓%c Works through Nginx reverse proxy without extra config\n" +
  "%cSSE ✓%c Simpler implementation, no upgrade handshake needed\n" +
  "%cWebSocket ✗%c Bidirectional overhead for one-way tracking\n" +
  "%cWebSocket ✗%c Requires proxy upgrade config, sticky sessions",
  "color:#16a34a;font-weight:bold", "color:#374151",
  "color:#16a34a;font-weight:bold", "color:#374151",
  "color:#16a34a;font-weight:bold", "color:#374151",
  "color:#16a34a;font-weight:bold", "color:#374151",
  "color:#dc2626;font-weight:bold", "color:#374151",
  "color:#dc2626;font-weight:bold", "color:#374151"
);
