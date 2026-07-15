const GATEWAY = process.env.NEXT_PUBLIC_GATEWAY_URL || "http://localhost:8000";

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
