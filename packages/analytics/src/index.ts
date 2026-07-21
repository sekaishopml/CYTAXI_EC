export type AnalyticsEvent =
  | { name: "trip_created"; properties: { trip_id: string; origin_lat: number; origin_lng: number; dest_lat: number; dest_lng: number; distance_km: number } }
  | { name: "trip_accepted"; properties: { trip_id: string; driver_id: string; eta_seconds: number } }
  | { name: "trip_started"; properties: { trip_id: string; driver_id: string; passenger_id: string } }
  | { name: "trip_completed"; properties: { trip_id: string; fare: number; duration_seconds: number; distance_km: number } }
  | { name: "trip_cancelled"; properties: { trip_id: string; cancelled_by: "passenger" | "driver" | "system"; reason?: string } }
  | { name: "payment_completed"; properties: { trip_id: string; amount: number; method: "cash" | "card" | "wallet"; status: string } }
  | { name: "screen_view"; properties: { screen: string; referrer?: string } }
  | { name: "search_query"; properties: { query: string; results_count: number } }
  | { name: "rating_submitted"; properties: { trip_id: string; score: number; from: "passenger" | "driver" } }
  | { name: "driver_online"; properties: { driver_id: string } }
  | { name: "driver_offline"; properties: { driver_id: string } }
  | { name: "error_occurred"; properties: { code: string; message: string; context?: Record<string, unknown> } };

interface AnalyticsConfig {
  endpoint?: string;
  flushIntervalMs?: number;
  maxBatchSize?: number;
  onEvent?: (event: AnalyticsEvent) => void;
}

export class Analytics {
  private queue: AnalyticsEvent[] = [];
  private config: Required<AnalyticsConfig>;
  private timer: ReturnType<typeof setInterval> | null = null;
  private sessionId: string;

  constructor(config: AnalyticsConfig = {}) {
    this.config = {
      endpoint: config.endpoint ?? "/api/v1/analytics/events",
      flushIntervalMs: config.flushIntervalMs ?? 5000,
      maxBatchSize: config.maxBatchSize ?? 20,
      onEvent: config.onEvent ?? (() => {}),
    };
    this.sessionId = `sess_${Date.now()}_${Math.random().toString(36).slice(2, 8)}`;
    if (typeof window !== "undefined") {
      this.timer = setInterval(() => this.flush(), this.config.flushIntervalMs);
      window.addEventListener("beforeunload", () => this.flush());
    }
  }

  track(event: AnalyticsEvent): void {
    this.queue.push(event);
    this.config.onEvent(event);
    if (this.queue.length >= this.config.maxBatchSize) this.flush();
  }

  async flush(): Promise<void> {
    if (this.queue.length === 0) return;
    const batch = this.queue.splice(0, this.config.maxBatchSize);
    try {
      const payload = { session_id: this.sessionId, events: batch, timestamp: new Date().toISOString() };
      if (typeof navigator !== "undefined" && "sendBeacon" in navigator) {
        navigator.sendBeacon(this.config.endpoint, JSON.stringify(payload));
      } else {
        await fetch(this.config.endpoint, { method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify(payload), keepalive: true });
      }
    } catch {}
  }

  getSessionId(): string { return this.sessionId; }

  destroy(): void {
    this.flush();
    if (this.timer) clearInterval(this.timer);
  }
}

export function createFunnel(steps: string[]): { trackStep: (step: string, properties?: Record<string, unknown>) => void } {
  const funnelId = `funnel_${Date.now()}`;
  let currentStep = 0;

  return {
    trackStep: (step: string, properties?: Record<string, unknown>) => {
      if (currentStep < steps.length && step === steps[currentStep]) {
        const event: AnalyticsEvent = {
          name: "screen_view",
          properties: { screen: `funnel:${funnelId}:${step}`, ...properties } as any,
        };
        currentStep++;
      }
    },
  };
}
