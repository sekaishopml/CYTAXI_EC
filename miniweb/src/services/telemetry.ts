interface TelemetryEvent {
  event: string;
  state?: string;
  prevState?: string | null;
  duration?: number;
  error?: string;
  latency?: number;
  timestamp: number;
}

const MAX_EVENTS = 200;
const FLUSH_INTERVAL = 30000;

class TelemetryCollector {
  private events: TelemetryEvent[] = [];
  private timers: Map<string, number> = new Map();
  private flushTimer: ReturnType<typeof setInterval> | null = null;

  constructor() {
    if (typeof window !== "undefined") {
      this.flushTimer = setInterval(() => this.flush(), FLUSH_INTERVAL);
    }
  }

  track(event: string, data?: Partial<TelemetryEvent>) {
    this.events.push({
      event,
      timestamp: Date.now(),
      ...data,
    });
    if (this.events.length > MAX_EVENTS) this.events.shift();
  }

  startTimer(id: string) {
    this.timers.set(id, Date.now());
  }

  endTimer(id: string, event: string, extra?: Partial<TelemetryEvent>) {
    const start = this.timers.get(id);
    if (!start) return;
    const duration = Date.now() - start;
    this.timers.delete(id);
    this.track(event, { duration, ...extra });
  }

  private flush() {
    if (this.events.length === 0) return;
    try {
      if (typeof navigator !== "undefined" && navigator.onLine) {
        const payload = this.events.splice(0, this.events.length);
        fetch("/api/v1/telemetry", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ events: payload, source: "miniweb" }),
          keepalive: true,
        }).catch(() => {});
      }
    } catch {}
  }

  getSnapshot() {
    return {
      total: this.events.length,
      recent: this.events.slice(-20),
      activeTimers: this.timers.size,
    };
  }

  destroy() {
    if (this.flushTimer) clearInterval(this.flushTimer);
    this.flush();
  }
}

export const telemetry = typeof window !== "undefined" ? new TelemetryCollector() : null;

export function trackJourneyEvent(event: string, data?: Partial<TelemetryEvent>) {
  telemetry?.track(event, data);
}

export function trackStateDuration(from: string, to: string, durationMs: number) {
  telemetry?.track("STATE_DURATION", {
    state: to,
    prevState: from,
    duration: durationMs,
  });
}

export function trackError(code: string, message: string) {
  telemetry?.track("ERROR", { error: `${code}: ${message}` });
}

export function trackLatency(operation: string, latencyMs: number) {
  telemetry?.track("LATENCY", { event: operation, latency: latencyMs });
}
