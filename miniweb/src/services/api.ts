const GATEWAY_URL = process.env.NEXT_PUBLIC_GATEWAY_URL || "http://localhost:8000";

interface RequestOptions {
  method?: "GET" | "POST" | "PUT" | "DELETE";
  body?: unknown;
  headers?: Record<string, string>;
}

async function request<T>(path: string, options: RequestOptions = {}): Promise<T> {
  const { method = "GET", body, headers = {} } = options;
  const url = `${GATEWAY_URL}/api/v1${path}`;

  const res = await fetch(url, {
    method,
    headers: {
      "Content-Type": "application/json",
      ...headers,
    },
    body: body ? JSON.stringify(body) : undefined,
  });

  if (!res.ok) {
    const error = await res.json().catch(() => ({ error: "Request failed" }));
    throw new Error(error.error || `HTTP ${res.status}`);
  }

  return res.json();
}

export const api = {
  trips: {
    get: (id: string) => request(`/trip/trips/${id}`),
    history: (customerId: string) => request(`/trip/customers/${customerId}/trips`),
  },
  pricing: {
    get: (id: string) => request(`/pricing/fares/${id}`),
    history: (tripId: string) => request(`/pricing/trips/${tripId}/fares`),
  },
  customer: {
    profile: (id: string) => request(`/customer/customers/${id}/profile`),
    preferences: (id: string) => request(`/customer/customers/${id}/preferences`),
    favorites: (id: string) => request(`/customer/customers/${id}/favorites`),
  },
  notifications: {
    get: (id: string) => request(`/notification/notifications/${id}`),
    history: (recipientId: string) => request(`/notification/recipients/${recipientId}/notifications`),
  },
  health: {
    check: () => request("/health"),
  },
};

export type { RequestOptions };
