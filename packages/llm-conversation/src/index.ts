export type IntentKind =
  | "greeting"
  | "trip_request"
  | "trip_status"
  | "support"
  | "cancel"
  | "change"
  | "unknown";

export interface Intent {
  kind: IntentKind;
  description: string;
}

export interface Place {
  name: string;
  address?: string;
  lat?: number;
  lng?: number;
  isCurrent?: boolean;
}

export interface Schedule {
  isNow: boolean;
  date?: string;
  time?: string;
  flexible?: boolean;
}

export interface Entities {
  origin?: Place;
  destination?: Place;
  passengers?: number;
  luggage?: string;
  schedule?: Schedule;
  vehicleType?: string;
  preferences?: string[];
  tripId?: string;
  phone?: string;
  name?: string;
}

export interface LLMResponse {
  intent: Intent;
  entities: Entities;
  confidence: number;
  pendingQuestions: string[];
  needsClarification: boolean;
  clarificationQuestion: string;
  rawInput: string;
}

export interface ConversationMessage {
  role: "user" | "assistant" | "system";
  content: string;
  timestamp: string;
  response?: LLMResponse;
}

export interface ConversationSession {
  phone: string;
  sessionId: string;
  messages: ConversationMessage[];
  context: Record<string, string>;
  createdAt: string;
  expiresAt: string;
}

export function hasHighConfidence(resp: LLMResponse): boolean {
  return resp.confidence >= 0.7 && !resp.needsClarification;
}

export function hasCompleteTripRequest(resp: LLMResponse): boolean {
  return (
    resp.intent.kind === "trip_request" &&
    resp.entities.origin !== undefined &&
    resp.entities.destination !== undefined
  );
}

export function parseLLMResponse(json: string): LLMResponse | null {
  try {
    const parsed = JSON.parse(json);
    if (!parsed.intent || !parsed.entities) return null;
    return parsed as LLMResponse;
  } catch {
    return null;
  }
}
