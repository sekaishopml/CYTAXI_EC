export interface ApiClientConfig {
  baseUrl: string;
  getToken?: () => Promise<string | null>;
  onTokenRefresh?: () => Promise<string | null>;
  onError?: (error: ApiError) => void;
  retryConfig?: RetryConfig;
  cacheConfig?: CacheConfig;
  tracing?: boolean;
  timeout?: number;
}

export interface RetryConfig {
  maxRetries: number;
  baseDelayMs: number;
  maxDelayMs: number;
  retryOn?: number[];
}

export interface CacheConfig {
  enabled: boolean;
  ttlMs: number;
  maxEntries: number;
}

export interface ApiError {
  status: number;
  code: string;
  message: string;
  details?: unknown;
  requestId?: string;
  timestamp: string;
}

export interface ApiRequest<T = unknown> {
  method: "GET" | "POST" | "PUT" | "PATCH" | "DELETE";
  path: string;
  body?: T;
  params?: Record<string, string | number | boolean | undefined>;
  headers?: Record<string, string>;
  cache?: boolean;
  signal?: AbortSignal;
}

export interface ApiResponse<T> {
  data: T;
  status: number;
  headers: Record<string, string>;
  requestId: string;
  cached: boolean;
}

export interface AuthTokens {
  accessToken: string;
  refreshToken: string;
  expiresAt: number;
}

export interface UserSession {
  id: string;
  role: "passenger" | "driver" | "admin";
  email: string;
  name: string;
  phone: string;
}

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
  pageSize: number;
  totalPages: number;
}
