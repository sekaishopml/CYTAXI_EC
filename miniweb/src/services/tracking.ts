import { TrackingUpdate } from "@/types";

const API = typeof window !== "undefined"
  ? `${window.location.protocol}//${window.location.host}/api/v1`
  : "http://64.176.219.221/api/v1";

const MAX_RECONNECT_DELAY = 30000;
const INITIAL_RECONNECT_DELAY = 1000;
const MAX_RETRIES = 10;

interface ReconnectState {
  attempt: number;
  timer: ReturnType<typeof setTimeout> | null;
  closed: boolean;
}

export function subscribeToTrip(
  tripId: string,
  onUpdate: (data: TrackingUpdate) => void,
  onError?: (err: string) => void,
): () => void {
  const state: ReconnectState = { attempt: 0, timer: null, closed: false };

  function connect() {
    if (state.closed) return;
    const es = new EventSource(`${API}/trip/ws?trip_id=${tripId}`);

    es.onmessage = (event) => {
      try {
        const data: TrackingUpdate = JSON.parse(event.data);
        state.attempt = 0;
        onUpdate(data);
        if (data.type === "trip_completed") {
          state.closed = true;
          es.close();
        }
      } catch (e) {
        onError?.("parse_error");
      }
    };

    es.onerror = () => {
      es.close();
      if (state.closed) return;
      state.attempt++;
      if (state.attempt > MAX_RETRIES) {
        onError?.("max_retries");
        return;
      }
      const delay = Math.min(
        INITIAL_RECONNECT_DELAY * Math.pow(2, state.attempt - 1) + Math.random() * 1000,
        MAX_RECONNECT_DELAY,
      );
      state.timer = setTimeout(connect, delay);
    };
  }

  connect();

  return () => {
    state.closed = true;
    if (state.timer) clearTimeout(state.timer);
  };
}

export function createTrackingUrl(tripId: string): string {
  return `${API}/trip/ws?trip_id=${tripId}`;
}
