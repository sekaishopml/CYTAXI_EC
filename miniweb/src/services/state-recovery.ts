const SESSION_KEY = "cytaxi_journey";
const SESSION_VER = 5;

export interface SessionData {
  v: number;
  state?: string;
  pickupAddress?: string;
  pickupCoords?: { lat: number; lng: number } | null;
  dest?: { name: string; address: string; lat: number; lng: number } | null;
  destQuery?: string;
  route?: unknown | null;
  fare?: unknown | null;
  vehicleType?: string;
  note?: string;
  paymentMethod?: "cash" | "card";
  scheduledAt?: string | null;
  tripId?: string;
  driver?: unknown | null;
  savedAt: number;
}

export function saveSession(data: Partial<SessionData>): void {
  try {
    const existing = loadSession();
    const payload: SessionData = {
      v: SESSION_VER,
      savedAt: Date.now(),
      ...existing,
      ...data,
    };
    localStorage.setItem(SESSION_KEY, JSON.stringify(payload));
  } catch {}
}

export function loadSession(): SessionData | null {
  try {
    const raw = localStorage.getItem(SESSION_KEY);
    if (!raw) return null;
    const data: SessionData = JSON.parse(raw);
    if (data.v !== SESSION_VER) {
      localStorage.removeItem(SESSION_KEY);
      return null;
    }
    const age = Date.now() - (data.savedAt || 0);
    if (age > 86400000) {
      localStorage.removeItem(SESSION_KEY);
      return null;
    }
    return data;
  } catch {
    localStorage.removeItem(SESSION_KEY);
    return null;
  }
}

export function clearSession(): void {
  try { localStorage.removeItem(SESSION_KEY); } catch {}
}

export function isSessionValid(data: SessionData | null): data is SessionData {
  if (!data) return false;
  const preTripStates = ["pickup_select", "input", "confirm"];
  if (data.state && preTripStates.includes(data.state)) return true;
  if (data.state && data.tripId) return true;
  return false;
}
