const GATEWAY = process.env.NEXT_PUBLIC_GATEWAY_URL || "http://localhost:8000";

export interface TrackingUpdate {
  type: string; trip_id: string; status: string;
  driver?: { id: string; name: string; vehicle: string; plate: string; lat: number; lng: number; rating: number };
  eta_seconds?: number; timestamp: string;
}

type Listener = (update: TrackingUpdate) => void;

export function subscribeToTrip(tripId: string, onUpdate: Listener): () => void {
  const url = `${GATEWAY.replace("http", "ws") || "ws://localhost:8000"}/api/v1/trip/ws?trip_id=${tripId}`;

  // Use SSE fallback since native WebSocket support varies
  const evtSource = new EventSource(`${GATEWAY}/api/v1/trip/ws?trip_id=${tripId}`);

  evtSource.onmessage = (event) => {
    try {
      const data: TrackingUpdate = JSON.parse(event.data);
      if (data.type === "trip_completed") {
        evtSource.close();
      }
      onUpdate(data);
    } catch (e) {
      console.error("SSE parse error", e);
    }
  };

  evtSource.onerror = () => {
    evtSource.close();
  };

  return () => evtSource.close();
}

export async function startTrip(tripId: string, driverId: string): Promise<any> {
  const res = await fetch(`${GATEWAY}/api/v1/trip/start`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ trip_id: tripId, driver_id: driverId }),
  });
  return res.json();
}

export async function updateLocation(tripId: string, driverId: string, lat: number, lng: number): Promise<any> {
  const res = await fetch(`${GATEWAY}/api/v1/trip/location`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ trip_id: tripId, driver_id: driverId, lat, lng }),
  });
  return res.json();
}

export async function finishTrip(tripId: string): Promise<any> {
  const res = await fetch(`${GATEWAY}/api/v1/trip/finish`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ trip_id: tripId }),
  });
  return res.json();
}
