const GATEWAY = process.env.NEXT_PUBLIC_GATEWAY_URL || "http://localhost:8000";

export interface DriverRequest {
  id: string; trip_id: string; pickup: string; destination: string;
  fare: string; eta_seconds: number; status: string; created_at: string; expires_at: string;
}

export async function getDriverRequests(): Promise<{ requests: DriverRequest[]; status: string }> {
  const res = await fetch(`${GATEWAY}/api/v1/driver/requests`);
  if (!res.ok) throw new Error("Failed to fetch requests");
  return res.json();
}

export async function acceptRequest(requestId: string): Promise<{ status: string; driver: any }> {
  const res = await fetch(`${GATEWAY}/api/v1/driver/accept`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ request_id: requestId }),
  });
  if (!res.ok) throw new Error("Failed to accept");
  return res.json();
}

export async function rejectRequest(requestId: string, reason: string): Promise<{ status: string }> {
  const res = await fetch(`${GATEWAY}/api/v1/driver/reject`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ request_id: requestId, reason }),
  });
  if (!res.ok) throw new Error("Failed to reject");
  return res.json();
}

export async function getDriverStatus(): Promise<{ status: string; name: string; rating: number }> {
  const res = await fetch(`${GATEWAY}/api/v1/driver/drivers/status`);
  if (!res.ok) throw new Error("Failed to get status");
  return res.json();
}

export async function startMatching(tripId: string, lat: number, lng: number): Promise<{ candidates: any[]; matching_id: string }> {
  const res = await fetch(`${GATEWAY}/api/v1/matching/start`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ trip_id: tripId, pickup_lat: lat, pickup_lng: lng }),
  });
  if (!res.ok) throw new Error("Failed to start matching");
  return res.json();
}
