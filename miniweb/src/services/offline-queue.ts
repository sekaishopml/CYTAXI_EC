const QUEUE_KEY = "cytaxi_offline_queue";

export interface QueuedAction {
  id: string;
  type: "TRIP_REQUEST" | "LOCATION_UPDATE" | "STATUS_CHANGE" | "PAYMENT_CONFIRM" | "RATING_SUBMIT";
  payload: unknown;
  timestamp: number;
  retries: number;
  maxRetries: number;
}

function getQueue(): QueuedAction[] {
  try {
    const raw = localStorage.getItem(QUEUE_KEY);
    return raw ? JSON.parse(raw) : [];
  } catch { return []; }
}

function saveQueue(queue: QueuedAction[]): void {
  try { localStorage.setItem(QUEUE_KEY, JSON.stringify(queue)); } catch {}
}

export function enqueueAction(type: QueuedAction["type"], payload: unknown): void {
  const queue = getQueue();
  queue.push({
    id: `${type}_${Date.now()}_${Math.random().toString(36).slice(2, 8)}`,
    type,
    payload,
    timestamp: Date.now(),
    retries: 0,
    maxRetries: 3,
  });
  saveQueue(queue);
}

export function dequeueAction(id: string): void {
  const queue = getQueue().filter(a => a.id !== id);
  saveQueue(queue);
}

export function getPendingActions(): QueuedAction[] {
  return getQueue().filter(a => a.retries < a.maxRetries);
}

export function incrementRetry(id: string): void {
  const queue = getQueue();
  const action = queue.find(a => a.id === id);
  if (action) {
    action.retries++;
    saveQueue(queue);
  }
}

export function clearQueue(): void {
  try { localStorage.removeItem(QUEUE_KEY); } catch {}
}

export function isOnline(): boolean {
  return typeof navigator !== "undefined" ? navigator.onLine : true;
}

export function onOnline(callback: () => void): () => void {
  window.addEventListener("online", callback);
  return () => window.removeEventListener("online", callback);
}

export function onOffline(callback: () => void): () => void {
  window.addEventListener("offline", callback);
  return () => window.removeEventListener("offline", callback);
}
