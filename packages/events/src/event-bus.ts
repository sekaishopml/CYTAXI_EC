import type { RideEventName, EventPayloadMap, AuditLogEntry } from "./types";

type Listener<T> = (payload: T) => void;
type Unsubscribe = () => void;

let _idCounter = 0;
const uid = () => `evt_${++_idCounter}_${Date.now()}`;

export class EventBus {
  private listeners = new Map<RideEventName, Set<Listener<unknown>>>();
  private auditLog: AuditLogEntry[] = [];
  private maxAuditSize: number;
  private source?: string;

  constructor(opts?: { maxAuditSize?: number; source?: string }) {
    this.maxAuditSize = opts?.maxAuditSize ?? 1000;
    this.source = opts?.source;
  }

  on<E extends RideEventName>(
    event: E,
    listener: Listener<EventPayloadMap[E]>
  ): Unsubscribe {
    if (!this.listeners.has(event)) {
      this.listeners.set(event, new Set());
    }
    this.listeners.get(event)!.add(listener as Listener<unknown>);
    return () => {
      this.listeners.get(event)?.delete(listener as Listener<unknown>);
    };
  }

  once<E extends RideEventName>(
    event: E,
    listener: Listener<EventPayloadMap[E]>
  ): Unsubscribe {
    const wrapper: Listener<EventPayloadMap[E]> = (payload) => {
      unsub();
      listener(payload);
    };
    const unsub = this.on(event, wrapper);
    return unsub;
  }

  off<E extends RideEventName>(
    event: E,
    listener: Listener<EventPayloadMap[E]>
  ): void {
    this.listeners.get(event)?.delete(listener as Listener<unknown>);
  }

  emit<E extends RideEventName>(
    event: E,
    payload: EventPayloadMap[E]
  ): void {
    const logEntry: AuditLogEntry = {
      id: uid(),
      event,
      payload,
      timestamp: new Date().toISOString(),
      source: this.source,
    };

    this.auditLog.push(logEntry);
    if (this.auditLog.length > this.maxAuditSize) {
      this.auditLog.shift();
    }

    const subs = this.listeners.get(event);
    if (subs) {
      subs.forEach((fn) => {
        try {
          fn(payload);
        } catch (err) {
          console.error(`[EventBus] Error in listener for ${event}:`, err);
        }
      });
    }
  }

  clear(): void {
    this.listeners.clear();
  }

  getAuditLog(): readonly AuditLogEntry[] {
    return this.auditLog;
  }

  getAuditLogForEvent(event: RideEventName): AuditLogEntry[] {
    return this.auditLog.filter((e) => e.event === event);
  }

  clearAuditLog(): void {
    this.auditLog = [];
  }

  listenerCount(event: RideEventName): number {
    return this.listeners.get(event)?.size ?? 0;
  }
}

// ─── Singleton global ──────────────────────────────────────────────
let _globalBus: EventBus | null = null;

export function getGlobalBus(): EventBus {
  if (!_globalBus) {
    _globalBus = new EventBus({ source: "global" });
  }
  return _globalBus;
}

export function resetGlobalBus(): void {
  if (_globalBus) {
    _globalBus.clear();
    _globalBus.clearAuditLog();
  }
  _globalBus = null;
}
