export type ConnectionState = "disconnected" | "connecting" | "connected" | "reconnecting";

export interface RealtimeConfig {
  url: string;
  getToken?: () => Promise<string | null>;
  reconnectDelayMs?: number;
  maxReconnectDelayMs?: number;
  heartbeatIntervalMs?: number;
  heartbeatTimeoutMs?: number;
  onStateChange?: (state: ConnectionState) => void;
}

export type MessageHandler = (data: unknown) => void;
export type Unsubscribe = () => void;

export class RealtimeClient {
  private ws: WebSocket | null = null;
  private state: ConnectionState = "disconnected";
  private listeners = new Map<string, Set<MessageHandler>>();
  private config: Required<RealtimeConfig>;
  private reconnectAttempts = 0;
  private heartbeatTimer: ReturnType<typeof setInterval> | null = null;
  private heartbeatTimeout: ReturnType<typeof setTimeout> | null = null;
  private intentionalClose = false;

  constructor(config: RealtimeConfig) {
    this.config = {
      url: config.url,
      getToken: config.getToken ?? (async () => null),
      reconnectDelayMs: config.reconnectDelayMs ?? 1000,
      maxReconnectDelayMs: config.maxReconnectDelayMs ?? 30000,
      heartbeatIntervalMs: config.heartbeatIntervalMs ?? 30000,
      heartbeatTimeoutMs: config.heartbeatTimeoutMs ?? 5000,
      onStateChange: config.onStateChange ?? (() => {}),
    };
  }

  get connectionState(): ConnectionState { return this.state; }

  async connect(): Promise<void> {
    if (this.ws && (this.state === "connected" || this.state === "connecting")) return;
    this.intentionalClose = false;
    this.setState("connecting");

    const url = new URL(this.config.url);
    url.searchParams.set("_t", Date.now().toString());

    const token = await this.config.getToken();
    if (token) url.searchParams.set("token", token);

    this.ws = new WebSocket(url.toString());

    this.ws.onopen = () => {
      this.reconnectAttempts = 0;
      this.setState("connected");
      this.startHeartbeat();
    };

    this.ws.onmessage = (event) => {
      try {
        const msg = JSON.parse(event.data) as { type: string; payload: unknown };
        if (msg.type === "pong") {
          this.clearHeartbeatTimeout();
          return;
        }
        const handlers = this.listeners.get(msg.type);
        if (handlers) {
          handlers.forEach((h) => {
            try { h(msg.payload); } catch {}
          });
        }
        const wildcard = this.listeners.get("*");
        if (wildcard) {
          wildcard.forEach((h) => {
            try { h(msg); } catch {}
          });
        }
      } catch {}
    };

    this.ws.onclose = () => {
      this.stopHeartbeat();
      if (!this.intentionalClose) {
        this.setState("reconnecting");
        this.scheduleReconnect();
      } else {
        this.setState("disconnected");
      }
    };

    this.ws.onerror = () => {};
  }

  disconnect(): void {
    this.intentionalClose = true;
    this.stopHeartbeat();
    if (this.ws) {
      this.ws.close(1000, "Client disconnect");
      this.ws = null;
    }
    this.setState("disconnected");
  }

  send(type: string, payload?: unknown): void {
    if (this.ws?.readyState !== WebSocket.OPEN) return;
    this.ws.send(JSON.stringify({ type, payload, timestamp: new Date().toISOString() }));
  }

  subscribe(type: string, handler: MessageHandler): Unsubscribe {
    if (!this.listeners.has(type)) this.listeners.set(type, new Set());
    this.listeners.get(type)!.add(handler);
    return () => { this.listeners.get(type)?.delete(handler); };
  }

  subscribeToAll(handler: (msg: { type: string; payload: unknown }) => void): Unsubscribe {
    return this.subscribe("*", handler as MessageHandler);
  }

  private setState(state: ConnectionState): void {
    this.state = state;
    this.config.onStateChange(state);
  }

  private scheduleReconnect(): void {
    const delay = Math.min(
      this.config.reconnectDelayMs * Math.pow(2, this.reconnectAttempts),
      this.config.maxReconnectDelayMs
    );
    this.reconnectAttempts++;
    setTimeout(() => this.connect(), delay);
  }

  private startHeartbeat(): void {
    this.stopHeartbeat();
    this.heartbeatTimer = setInterval(() => {
      this.send("ping");
      this.heartbeatTimeout = setTimeout(() => {
        this.intentionalClose = false;
        this.ws?.close(4000, "Heartbeat timeout");
      }, this.config.heartbeatTimeoutMs);
    }, this.config.heartbeatIntervalMs);
  }

  private stopHeartbeat(): void {
    if (this.heartbeatTimer) { clearInterval(this.heartbeatTimer); this.heartbeatTimer = null; }
    this.clearHeartbeatTimeout();
  }

  private clearHeartbeatTimeout(): void {
    if (this.heartbeatTimeout) { clearTimeout(this.heartbeatTimeout); this.heartbeatTimeout = null; }
  }
}
