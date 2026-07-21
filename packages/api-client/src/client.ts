import type { ApiClientConfig, ApiRequest, ApiResponse, ApiError, RetryConfig } from "./types";

let _requestIdCounter = 0;
const uid = () => `req_${++_requestIdCounter}_${Date.now().toString(36)}`;

const DEFAULT_RETRY: RetryConfig = {
  maxRetries: 2,
  baseDelayMs: 500,
  maxDelayMs: 5000,
  retryOn: [408, 429, 500, 502, 503, 504],
};

export class ApiClient {
  private config: ApiClientConfig & { retryConfig: RetryConfig };
  private cache = new Map<string, { data: unknown; expires: number }>();

  constructor(config: ApiClientConfig) {
    this.config = { ...config, retryConfig: { ...DEFAULT_RETRY, ...config.retryConfig } };
  }

  setConfig(partial: Partial<ApiClientConfig>): void {
    this.config = { ...this.config, ...partial, retryConfig: { ...this.config.retryConfig, ...partial.retryConfig } };
  }

  async request<T>(req: ApiRequest): Promise<ApiResponse<T>> {
    const requestId = uid();
    const url = this.buildUrl(req);
    const cacheKey = req.method === "GET" && req.cache !== false ? url : null;

    if (cacheKey) {
      const cached = this.cache.get(cacheKey);
      if (cached && cached.expires > Date.now()) {
        return { data: cached.data as T, status: 200, headers: {}, requestId, cached: true };
      }
    }

    let lastError: Error | null = null;
    let attempts = 0;

    while (attempts <= this.config.retryConfig.maxRetries) {
      try {
        const response = await this.execute<T>(req, url, requestId, attempts);
        if (cacheKey && response.status < 400) {
          this.cache.set(cacheKey, { data: response.data, expires: Date.now() + (this.config.cacheConfig?.ttlMs ?? 30000) });
          if (this.cache.size > (this.config.cacheConfig?.maxEntries ?? 100)) {
            const first = this.cache.keys().next().value;
            if (first) this.cache.delete(first);
          }
        }
        return response;
      } catch (err) {
        lastError = err instanceof Error ? err : new Error(String(err));
        const apiErr = err as ApiError;
        if (this.config.retryConfig.retryOn?.includes(apiErr.status) && attempts < this.config.retryConfig.maxRetries) {
          const delay = Math.min(this.config.retryConfig.baseDelayMs * Math.pow(2, attempts), this.config.retryConfig.maxDelayMs);
          await new Promise((r) => setTimeout(r, delay));
          attempts++;
          continue;
        }
        throw err;
      }
    }

    throw lastError ?? new Error("Request failed");
  }

  private async execute<T>(req: ApiRequest, url: string, requestId: string, attempt: number): Promise<ApiResponse<T>> {
    const headers: Record<string, string> = {
      "Content-Type": "application/json",
      "X-Request-Id": requestId,
      ...req.headers,
    };

    if (this.config.tracing) {
      headers["X-Trace-Id"] = requestId;
      headers["X-Attempt"] = String(attempt);
    }

    const token = await this.config.getToken?.();
    if (token) headers["Authorization"] = `Bearer ${token}`;

    const controller = new AbortController();
    const timeoutMs = this.config.timeout ?? 15000;
    const timer = setTimeout(() => controller.abort(), timeoutMs);

    const signal = req.signal ? combineSignals(req.signal, controller.signal) : controller.signal;

    try {
      const res = await fetch(url, {
        method: req.method,
        headers,
        body: req.body ? JSON.stringify(req.body) : undefined,
        signal,
      });

      clearTimeout(timer);

      if (res.status === 401 && this.config.onTokenRefresh) {
        const newToken = await this.config.onTokenRefresh();
        if (newToken) {
          headers["Authorization"] = `Bearer ${newToken}`;
          const retryRes = await fetch(url, { method: req.method, headers, body: req.body ? JSON.stringify(req.body) : undefined });
          return this.toResponse<T>(retryRes, requestId);
        }
      }

      return this.toResponse<T>(res, requestId);
    } catch (err) {
      clearTimeout(timer);
      throw this.normalizeError(err, requestId);
    }
  }

  private async toResponse<T>(res: Response, requestId: string): Promise<ApiResponse<T>> {
    const data = res.status === 204 ? null : await res.json();
    if (!res.ok) {
      const apiError: ApiError = {
        status: res.status,
        code: data?.code || `HTTP_${res.status}`,
        message: data?.message || data?.error || res.statusText,
        details: data?.details,
        requestId,
        timestamp: new Date().toISOString(),
      };
      this.config.onError?.(apiError);
      throw apiError;
    }
    const headers: Record<string, string> = {};
    res.headers.forEach((v, k) => { headers[k] = v; });
    return { data: data as T, status: res.status, headers, requestId, cached: false };
  }

  private buildUrl(req: ApiRequest): string {
    let url = `${this.config.baseUrl}${req.path}`;
    if (req.params) {
      const search = new URLSearchParams();
      Object.entries(req.params).forEach(([k, v]) => {
        if (v !== undefined) search.set(k, String(v));
      });
      const qs = search.toString();
      if (qs) url += `?${qs}`;
    }
    return url;
  }

  private normalizeError(err: unknown, requestId: string): ApiError {
    if (err && typeof err === "object" && "status" in err) return err as ApiError;
    const message = err instanceof Error ? err.message : String(err);
    return {
      status: 0,
      code: "NETWORK_ERROR",
      message,
      requestId,
      timestamp: new Date().toISOString(),
    };
  }

  clearCache(): void { this.cache.clear(); }
}

function combineSignals(...signals: AbortSignal[]): AbortSignal {
  const controller = new AbortController();
  signals.forEach((s) => {
    if (s.aborted) controller.abort(s.reason);
    else s.addEventListener("abort", () => controller.abort(s.reason), { once: true });
  });
  return controller.signal;
}
