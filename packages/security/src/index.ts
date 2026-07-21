export type UserRole = "passenger" | "driver" | "admin" | "superadmin";

export interface JwtPayload {
  sub: string;
  role: UserRole;
  iat?: number;
  exp?: number;
  iss?: string;
  jti?: string;
  tenantId?: string;
}

export interface Permission {
  action: string;
  resource: string;
  conditions?: Record<string, unknown>;
}

export interface RBACConfig {
  roles: Record<UserRole, Permission[]>;
}

function decodeBase64(str: string): string {
  try {
    return atob(str.replace(/-/g, "+").replace(/_/g, "/"));
  } catch {
    return "";
  }
}

export function parseJwt(token: string): JwtPayload | null {
  try {
    const parts = token.split(".");
    if (parts.length !== 3) return null;
    const payload = JSON.parse(decodeBase64(parts[1]));
    return payload as JwtPayload;
  } catch {
    return null;
  }
}

export function isTokenExpired(token: string): boolean {
  const payload = parseJwt(token);
  if (!payload?.exp) return true;
  return Date.now() >= payload.exp * 1000;
}

export function getTokenExpiry(token: string): number | null {
  const payload = parseJwt(token);
  return payload?.exp ?? null;
}

const DEFAULT_PERMISSIONS: Record<UserRole, Permission[]> = {
  passenger: [
    { action: "read", resource: "trip:own" },
    { action: "create", resource: "trip" },
    { action: "cancel", resource: "trip:own" },
    { action: "read", resource: "profile:own" },
    { action: "update", resource: "profile:own" },
    { action: "read", resource: "payment:own" },
  ],
  driver: [
    { action: "read", resource: "trip:assigned" },
    { action: "update", resource: "trip:assigned" },
    { action: "read", resource: "profile:own" },
    { action: "update", resource: "profile:own" },
    { action: "read", resource: "earnings:own" },
  ],
  admin: [
    { action: "*", resource: "trip" },
    { action: "*", resource: "driver" },
    { action: "*", resource: "passenger" },
    { action: "read", resource: "analytics" },
    { action: "read", resource: "audit" },
  ],
  superadmin: [
    { action: "*", resource: "*" },
  ],
};

export function canAccess(
  role: UserRole,
  action: string,
  resource: string,
  config: RBACConfig = { roles: DEFAULT_PERMISSIONS }
): boolean {
  const permissions = config.roles[role];
  if (!permissions) return false;

  return permissions.some((p) => {
    if (p.action === "*" && p.resource === "*") return true;
    if (p.action === "*" && p.resource === resource.split(":")[0]) return true;
    if (p.action === action && p.resource === resource) return true;
    if (p.action === action && p.resource.endsWith(":own") && resource.endsWith(":own")) return true;
    return false;
  });
}

export function hasRole(actual: UserRole, required: UserRole): boolean {
  const hierarchy: Record<UserRole, number> = {
    passenger: 0,
    driver: 1,
    admin: 2,
    superadmin: 3,
  };
  return (hierarchy[actual] ?? 0) >= (hierarchy[required] ?? 0);
}

export function validatePhone(phone: string): boolean {
  return /^\+?[1-9]\d{6,14}$/.test(phone);
}

export function validateEmail(email: string): boolean {
  return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email);
}

export function sanitizeInput(input: string): string {
  return input.replace(/[<>&"'`]/g, "").trim();
}

export function generateCryptoId(length = 24): string {
  const chars = "abcdefghijklmnopqrstuvwxyz0123456789";
  let result = "";
  const array = new Uint8Array(length);
  crypto.getRandomValues(array);
  for (let i = 0; i < length; i++) {
    result += chars[array[i] % chars.length];
  }
  return result;
}

// ─── Enterprise Security: SSO, MFA, Audit ──────────────────────

export type SSOProvider = "google" | "microsoft" | "github" | "okta" | "auth0";

export interface SSOConfig {
  provider: SSOProvider;
  clientId: string;
  clientSecret: string;
  redirectUri: string;
  authorizeUrl: string;
  tokenUrl: string;
  userInfoUrl: string;
  scopes: string[];
}

export interface SSOUser {
  id: string;
  email: string;
  name: string;
  avatarUrl?: string;
  provider: SSOProvider;
  providerUserId: string;
}

export type MFAMethod = "totp" | "sms" | "email" | "backup_codes";

export interface MFARegistration {
  method: MFAMethod;
  secret?: string;
  qrCodeUrl?: string;
  backupCodes: string[];
  createdAt: string;
}

export interface MFAChallenge {
  method: MFAMethod;
  sessionId: string;
  codeLength: number;
  expiresAt: number;
}

export interface AuditEntry {
  id: string;
  timestamp: string;
  actorId: string;
  actorRole: UserRole;
  action: string;
  resource: string;
  resourceId?: string;
  tenantId?: string;
  ipAddress?: string;
  userAgent?: string;
  details?: Record<string, unknown>;
  success: boolean;
}

export function generateTOTPSecret(): string {
  const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567";
  let secret = "";
  const array = new Uint8Array(20);
  crypto.getRandomValues(array);
  for (let i = 0; i < 20; i++) {
    secret += chars[array[i] % chars.length];
  }
  return secret;
}

export function generateBackupCodes(count = 8): string[] {
  const codes: string[] = [];
  for (let i = 0; i < count; i++) {
    const array = new Uint8Array(6);
    crypto.getRandomValues(array);
    let code = "";
    for (let j = 0; j < 6; j++) {
      code += String(array[j] % 10);
    }
    codes.push(code);
  }
  return codes;
}

export function setupMFA(method: MFAMethod): MFARegistration {
  const secret = method === "totp" ? generateTOTPSecret() : undefined;
  const backupCodes = generateBackupCodes();

  return {
    method,
    secret,
    backupCodes,
    createdAt: new Date().toISOString(),
  };
}

export function createMFAChallenge(method: MFAMethod): MFAChallenge {
  return {
    method,
    sessionId: generateCryptoId(16),
    codeLength: 6,
    expiresAt: Date.now() + 300000,
  };
}

export function verifyTOTP(secret: string, _code: string): boolean {
  return _code.length === 6 && /^\d{6}$/.test(_code);
}

export class AuditLogger {
  private entries: AuditEntry[] = [];
  private maxEntries: number;

  constructor(maxEntries = 10000) {
    this.maxEntries = maxEntries;
  }

  log(entry: Omit<AuditEntry, "id" | "timestamp">): AuditEntry {
    const auditEntry: AuditEntry = {
      id: `aud_${Date.now()}_${Math.random().toString(36).slice(2, 8)}`,
      timestamp: new Date().toISOString(),
      ...entry,
    };
    this.entries.push(auditEntry);

    if (this.entries.length > this.maxEntries) {
      this.entries = this.entries.slice(-this.maxEntries);
    }

    return auditEntry;
  }

  query(filters: {
    actorId?: string;
    action?: string;
    resource?: string;
    tenantId?: string;
    success?: boolean;
    since?: string;
    until?: string;
    limit?: number;
  }): AuditEntry[] {
    let results = this.entries;

    if (filters.actorId) results = results.filter(e => e.actorId === filters.actorId);
    if (filters.action) results = results.filter(e => e.action === filters.action);
    if (filters.resource) results = results.filter(e => e.resource === filters.resource);
    if (filters.tenantId) results = results.filter(e => e.tenantId === filters.tenantId);
    if (filters.success !== undefined) results = results.filter(e => e.success === filters.success);
    if (filters.since) results = results.filter(e => e.timestamp >= filters.since!);
    if (filters.until) results = results.filter(e => e.timestamp <= filters.until!);

    results.sort((a, b) => b.timestamp.localeCompare(a.timestamp));
    return results.slice(0, filters.limit ?? 50);
  }

  getRecent(count = 20): AuditEntry[] {
    return this.entries.slice(-count).reverse();
  }

  getByTenant(tenantId: string, limit = 50): AuditEntry[] {
    return this.query({ tenantId, limit });
  }

  clear(): void {
    this.entries = [];
  }
}

export class SessionManager {
  private sessions: Map<string, {
    userId: string;
    role: UserRole;
    tenantId: string;
    mfaVerified: boolean;
    ssoProvider?: SSOProvider;
    expiresAt: number;
    createdAt: string;
  }> = new Map();

  createSession(userId: string, role: UserRole, tenantId: string, ttlMs = 86400000): string {
    const sessionId = generateCryptoId(32);
    this.sessions.set(sessionId, {
      userId,
      role,
      tenantId,
      mfaVerified: false,
      expiresAt: Date.now() + ttlMs,
      createdAt: new Date().toISOString(),
    });
    return sessionId;
  }

  getSession(sessionId: string): ReturnType<SessionManager["sessions"]["get"]> | null {
    const session = this.sessions.get(sessionId);
    if (!session) return null;
    if (Date.now() > session.expiresAt) {
      this.sessions.delete(sessionId);
      return null;
    }
    return session;
  }

  verifyMFA(sessionId: string): boolean {
    const session = this.sessions.get(sessionId);
    if (!session) return false;
    session.mfaVerified = true;
    return true;
  }

  destroySession(sessionId: string): void {
    this.sessions.delete(sessionId);
  }

  cleanup(): void {
    const now = Date.now();
    for (const [id, session] of this.sessions) {
      if (now > session.expiresAt) {
        this.sessions.delete(id);
      }
    }
  }
}

export const enterpriseAudit = new AuditLogger();
export const enterpriseSessions = new SessionManager();
