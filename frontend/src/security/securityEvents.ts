export type SecurityEventType =
  // Lifecycle / player events
  | "VIDEO_STARTED"
  | "VIDEO_PAUSED"
  | "VIDEO_RESUMED"
  // Weak behavioral signals
  | "VISIBILITY_HIDDEN"
  | "WINDOW_BLUR"
  | "FULLSCREEN_CHANGED"
  | "PAGE_LIFECYCLE_CHANGE"
  // Capture / tampering signals
  | "CAPTURE_SIGNAL"
  | "PLAYER_TAMPER"
  | "WATERMARK_TAMPER"
  // Session & token anomalies
  | "MULTIPLE_SESSION"
  | "TOKEN_ANOMALY"
  | "TOKEN_REJECTED"
  | "SESSION_REVOKED"
  | "SESSION_EXPIRED"
  // Aggregated outcomes
  | "SUSPICIOUS_ACTIVITY"
  | "HIGH_RISK_ACTIVITY"
  // Legacy aliases retained for backward compatibility
  | "PLAYER_STARTED"
  | "PLAYER_PAUSED"
  | "PLAYER_RESUMED"
  | "WINDOW_FOCUS"
  | "VISIBILITY_VISIBLE"
  | "FULLSCREEN_ENTER"
  | "FULLSCREEN_EXIT"
  | "SUSPICIOUS_CAPTURE"
  | "HIGH_RISK_CAPTURE"
  | "MULTIPLE_SESSIONS"
  | "WATERMARK_RENDERED"
  | "WATERMARK_TAMPERED"
  | "PLAYBACK_RATE_CHANGED"
  | "DEVTOOLS_OPENED";

export interface SecurityEventPayload {
  eventId: string;
  userId: string;
  videoId: string;
  sessionId: string;
  eventType: SecurityEventType;
  riskScore: number;
  confidence: number;
  timestamp: number;
  platform?: string;
  metadata?: Record<string, any>;
}

type EventCallback = (event: SecurityEventPayload) => void;

class SecurityEventBus {
  private listeners: Map<string, Set<EventCallback>> = new Map();

  on(eventType: SecurityEventType | "*", callback: EventCallback): () => void {
    if (!this.listeners.has(eventType)) {
      this.listeners.set(eventType, new Set());
    }
    this.listeners.get(eventType)!.add(callback);

    return () => {
      this.listeners.get(eventType)?.delete(callback);
    };
  }

  emit(event: SecurityEventPayload): void {
    const specific = this.listeners.get(event.eventType);
    if (specific) {
      specific.forEach((cb) => cb(event));
    }
    const wildcard = this.listeners.get("*");
    if (wildcard) {
      wildcard.forEach((cb) => cb(event));
    }
  }
}

export const securityEventBus = new SecurityEventBus();

/** Map legacy event names to their spec-aligned equivalents. */
export function normalizeEventType(raw: string): SecurityEventType {
  const map: Record<string, SecurityEventType> = {
    PLAYER_STARTED: "VIDEO_STARTED",
    PLAYER_PAUSED: "VIDEO_PAUSED",
    PLAYER_RESUMED: "VIDEO_RESUMED",
    FULLSCREEN_EXIT: "FULLSCREEN_CHANGED",
    SUSPICIOUS_CAPTURE: "CAPTURE_SIGNAL",
    HIGH_RISK_CAPTURE: "HIGH_RISK_ACTIVITY",
    MULTIPLE_SESSIONS: "MULTIPLE_SESSION",
    WATERMARK_TAMPERED: "WATERMARK_TAMPER",
    SCREEN_RECORDING_DETECTED: "CAPTURE_SIGNAL",
  };
  return map[raw] ?? (raw as SecurityEventType);
}