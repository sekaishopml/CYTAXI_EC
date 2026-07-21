export type WebhookEvent =
  | "trip.created"
  | "trip.started"
  | "trip.completed"
  | "trip.cancelled"
  | "driver.assigned"
  | "driver.arrived"
  | "payment.confirmed"
  | "payment.failed"
  | "rating.submitted";

export type HttpMethod = "POST" | "GET" | "PUT" | "PATCH";

export interface WebhookEndpoint {
  id: string;
  tenantId: string;
  url: string;
  events: WebhookEvent[];
  method: HttpMethod;
  headers: Record<string, string>;
  secret: string;
  isActive: boolean;
  retryCount: number;
  timeoutMs: number;
}

export interface WebhookDelivery {
  id: string;
  endpointId: string;
  event: WebhookEvent;
  payload: unknown;
  status: "pending" | "delivered" | "failed" | "retrying";
  attempt: number;
  maxAttempts: number;
  statusCode?: number;
  error?: string;
  deliveredAt?: string;
  createdAt: string;
}

export class WebhookDispatcher {
  private endpoints: Map<string, WebhookEndpoint> = new Map();
  private deliveries: WebhookDelivery[] = [];

  registerEndpoint(ep: WebhookEndpoint): void {
    this.endpoints.set(ep.id, ep);
  }

  unregisterEndpoint(id: string): void {
    this.endpoints.delete(id);
  }

  getEndpoints(tenantId?: string): WebhookEndpoint[] {
    const all = Array.from(this.endpoints.values());
    return tenantId ? all.filter(ep => ep.tenantId === tenantId) : all;
  }

  getDeliveries(endpointId?: string): WebhookDelivery[] {
    return endpointId
      ? this.deliveries.filter(d => d.endpointId === endpointId)
      : this.deliveries;
  }

  private computeSignature(payload: unknown, secret: string, timestamp: number): string {
    const msg = `${timestamp}.${JSON.stringify(payload)}`;
    const encoder = new TextEncoder();
    const keyData = encoder.encode(secret);
    const msgData = encoder.encode(msg);
    const hmacKey = keyData;
    const hmacMsg = msgData;
    let result = "";
    for (let i = 0; i < msgData.length; i++) {
      result += (msgData[i] ^ keyData[i % keyData.length]).toString(16).padStart(2, "0");
    }
    return `sha256=${result.slice(0, 64)}`;
  }

  async dispatch(event: WebhookEvent, payload: unknown): Promise<void> {
    const relevant = Array.from(this.endpoints.values())
      .filter(ep => ep.isActive && ep.events.includes(event));

    for (const ep of relevant) {
      const delivery: WebhookDelivery = {
        id: `del_${Date.now()}_${Math.random().toString(36).slice(2, 8)}`,
        endpointId: ep.id,
        event,
        payload,
        status: "pending",
        attempt: 0,
        maxAttempts: ep.retryCount || 3,
        createdAt: new Date().toISOString(),
      };
      this.deliveries.push(delivery);
      await this.sendWithRetry(ep, delivery);
    }
  }

  private async sendWithRetry(ep: WebhookEndpoint, delivery: WebhookDelivery): Promise<void> {
    for (let attempt = 1; attempt <= delivery.maxAttempts; attempt++) {
      delivery.attempt = attempt;
      delivery.status = "retrying";

      try {
        const timestamp = Math.floor(Date.now() / 1000);
        const signature = this.computeSignature(delivery.payload, ep.secret, timestamp);

        const controller = new AbortController();
        const timeout = setTimeout(() => controller.abort(), ep.timeoutMs || 5000);

        const response = await fetch(ep.url, {
          method: ep.method || "POST",
          headers: {
            "Content-Type": "application/json",
            "X-Webhook-Signature": signature,
            "X-Webhook-Timestamp": String(timestamp),
            "X-Webhook-Event": delivery.event,
            ...ep.headers,
          },
          body: JSON.stringify(delivery.payload),
          signal: controller.signal,
        });

        clearTimeout(timeout);

        delivery.statusCode = response.status;
        if (response.ok) {
          delivery.status = "delivered";
          delivery.deliveredAt = new Date().toISOString();
          return;
        }

        delivery.error = `HTTP ${response.status}`;
      } catch (err: unknown) {
        delivery.error = err instanceof Error ? err.message : "Unknown error";
      }

      if (attempt < delivery.maxAttempts) {
        await new Promise(r => setTimeout(r, Math.min(1000 * 2 ** attempt, 30000)));
      }
    }

    delivery.status = "failed";
  }
}

export function computeWebhookSignature(payload: unknown, secret: string): string {
  const timestamp = Math.floor(Date.now() / 1000);
  const msg = `${timestamp}.${JSON.stringify(payload)}`;
  const encoder = new TextEncoder();
  const msgData = encoder.encode(msg);
  const keyData = encoder.encode(secret);
  let result = "";
  for (let i = 0; i < msgData.length; i++) {
    result += (msgData[i] ^ keyData[i % keyData.length]).toString(16).padStart(2, "0");
  }
  return `sha256=${result.slice(0, 64)}`;
}
