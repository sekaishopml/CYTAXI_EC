export type LogLevel = "debug" | "info" | "warn" | "error";

export interface LogEntry {
  level: LogLevel;
  message: string;
  context?: Record<string, unknown>;
  timestamp: string;
  source?: string;
  traceId?: string;
}

export interface MetricPoint {
  name: string;
  value: number;
  tags?: Record<string, string>;
  timestamp: string;
}

export interface HealthCheckResult {
  service: string;
  status: "healthy" | "degraded" | "unhealthy";
  latencyMs: number;
  lastCheck: string;
  error?: string;
}

class Logger {
  private log(level: LogLevel, message: string, context?: Record<string, unknown>): void {
    const entry: LogEntry = { level, message, context, timestamp: new Date().toISOString(), source: "web" };
    if (level === "error") console.error(entry);
    else if (level === "warn") console.warn(entry);
    else console.log(entry);
  }

  debug(msg: string, ctx?: Record<string, unknown>): void { this.log("debug", msg, ctx); }
  info(msg: string, ctx?: Record<string, unknown>): void { this.log("info", msg, ctx); }
  warn(msg: string, ctx?: Record<string, unknown>): void { this.log("warn", msg, ctx); }
  error(msg: string, ctx?: Record<string, unknown>): void { this.log("error", msg, ctx); }
}

class MetricsCollector {
  private metrics: MetricPoint[] = [];
  private maxPoints = 1000;

  record(name: string, value: number, tags?: Record<string, string>): void {
    this.metrics.push({ name, value, tags, timestamp: new Date().toISOString() });
    if (this.metrics.length > this.maxPoints) this.metrics.shift();
  }

  flush(): MetricPoint[] {
    const snapshot = [...this.metrics];
    this.metrics = [];
    return snapshot;
  }

  gauge(name: string, value: number, tags?: Record<string, string>): void {
    this.record(name, value, tags);
  }

  increment(name: string, tags?: Record<string, string>): void {
    this.record(name, 1, tags);
  }

  timing(name: string, durationMs: number, tags?: Record<string, string>): void {
    this.record(name, durationMs, { ...tags, unit: "ms" });
  }
}

class Tracer {
  private traceId: string;

  constructor() {
    this.traceId = `trace_${Date.now()}_${Math.random().toString(36).slice(2, 8)}`;
  }

  getTraceId(): string { return this.traceId; }

  span(name: string, fn: () => Promise<unknown>): Promise<unknown>;
  span<T>(name: string, fn: () => T): T;
  span<T>(name: string, fn: (() => T) | (() => Promise<unknown>)): T | Promise<unknown> {
    const start = performance.now();
    const result = fn();
    if (result instanceof Promise) {
      return result.finally(() => {
        const duration = performance.now() - start;
        metrics.timing(`span.${name}`, duration, { traceId: this.traceId });
      });
    }
    const duration = performance.now() - start;
    metrics.timing(`span.${name}`, duration, { traceId: this.traceId });
    return result;
  }
}

export const logger = new Logger();
export const metrics = new MetricsCollector();
export const tracer = new Tracer();

export function observe<T>(name: string, fn: () => Promise<T>): Promise<T> {
  const start = performance.now();
  return fn().finally(() => {
    const duration = performance.now() - start;
    metrics.timing(name, duration);
    logger.debug(`Observed ${name}`, { durationMs: Math.round(duration) });
  });
}

export function reportWebVitals(): void {
  if (typeof window === "undefined" || !("performance" in window)) return;

  try {
    const observer = new PerformanceObserver((list) => {
      list.getEntries().forEach((entry) => {
        if (entry.entryType === "largest-contentful-paint") {
          metrics.gauge("web_vitals.lcp", entry.startTime);
          logger.info("LCP", { value: Math.round(entry.startTime) });
        }
        if (entry.entryType === "first-input") {
          metrics.gauge("web_vitals.inp", entry.duration);
          logger.info("INP", { value: Math.round(entry.duration) });
        }
        if (entry.entryType === "layout-shift") {
          const cls = (entry as any).value || 0;
          metrics.gauge("web_vitals.cls", cls);
        }
      });
    });

    observer.observe({ type: "largest-contentful-paint", buffered: true });
    observer.observe({ type: "first-input", buffered: true });
    observer.observe({ type: "layout-shift", buffered: true });
  } catch {}
}
