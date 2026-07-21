const GATEWAY_URL = process.env.NEXT_PUBLIC_GATEWAY_URL || "http://localhost:8000";

async function request<T>(path: string, options: { method?: string; body?: unknown } = {}): Promise<T> {
  const { method = "GET", body } = options;
  const res = await fetch(`${GATEWAY_URL}/api/v1${path}`, {
    method,
    headers: { "Content-Type": "application/json" },
    body: body ? JSON.stringify(body) : undefined,
  });
  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: "Request failed" }));
    throw new Error(err.error || `HTTP ${res.status}`);
  }
  return res.json();
}

export const api = {
  driver: {
    get: (id: string) => request(`/driver/drivers/${id}`),
    vehicles: (id: string) => request(`/driver/drivers/${id}/vehicles`),
    availability: (id: string) => request(`/driver/drivers/${id}/availability`),
    licenses: (id: string) => request(`/driver/drivers/${id}/licenses`),
  },
  trip: {
    get: (id: string) => request(`/trip/trips/${id}`),
    history: (driverId: string) => request(`/trip/drivers/${driverId}/trips`),
  },
  matching: {
    candidates: (id: string) => request(`/matching/matching/${id}/candidates`),
  },
  notifications: {
    history: (recipientId: string) => request(`/notification/recipients/${recipientId}/notifications`),
  },
  payment: {
    wallet: (ownerId: string) => request(`/payment/wallets/${ownerId}`),
  },
};
