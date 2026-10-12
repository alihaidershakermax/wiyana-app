export type RiskLevel = "NORMAL" | "SUSPICIOUS" | "HIGH_RISK" | "CRITICAL";

export interface RiskConfig {
  normalThreshold: number; // 0 - 29
  suspiciousThreshold: number; // 30 - 59
  highThreshold: number; // 60 - 79
  weights: Record<string, number>;
  /** Signal types whose weight is bumped when repeated inside the window. */
  repeatTypes: Record<string, number>;
  /** How many occurrences inside the window count as "repeated". */
  repeatCount: number;
}

export const DEFAULT_RISK_CONFIG: RiskConfig = {
  normalThreshold: 29,
  suspiciousThreshold: 59,
  highThreshold: 79,
  weights: {
    CAPTURE_SIGNAL: 60,
    VISIBILITY_HIDDEN: 5,
    WINDOW_BLUR: 5,
    FULLSCREEN_CHANGED: 15,
    PAGE_LIFECYCLE_CHANGE: 10,
    PLAYER_TAMPER: 50,
    WATERMARK_TAMPER: 50,
    MULTIPLE_SESSION: 40,
    TOKEN_ANOMALY: 60,
    // Legacy aliases kept for backward compatibility.
    PLAYER_STARTED: 0,
    PLAYER_PAUSED: 0,
    WINDOW_FOCUS: 0,
    VISIBILITY_VISIBLE: 0,
    FULLSCREEN_ENTER: 0,
    PLAYBACK_RATE_CHANGED: 10,
    DEVTOOLS_OPENED: 35,
    REPEATED_FOCUS_SWITCH: 15,
    INVALID_PLAYBACK_TOKEN: 50,
    SCREEN_RECORDING_DETECTED: 60,
    MULTIPLE_SESSIONS: 40,
    WATERMARK_TAMPERED: 50,
    SESSION_EXPIRED: 30,
    TOKEN_REJECTED: 50,
  },
  repeatTypes: {
    VISIBILITY_HIDDEN: 15,
    WINDOW_BLUR: 15,
  },
  repeatCount: 3,
};

const WINDOW_MS = 300000; // 5 minutes — signals decay out of the score.

export class ClientRiskEngine {
  private config: RiskConfig;
  private platformWeights: Record<string, number> | null;
  private currentScore = 0;
  private history: Array<{ eventType: string; timestamp: number; weight: number }> = [];

  constructor(
    platformWeights?: Record<string, number> | null,
    config?: Partial<RiskConfig>
  ) {
    this.platformWeights = platformWeights ?? null;
    this.config = {
      ...DEFAULT_RISK_CONFIG,
      ...config,
      weights: { ...DEFAULT_RISK_CONFIG.weights, ...(config?.weights || {}) },
      repeatTypes: { ...DEFAULT_RISK_CONFIG.repeatTypes, ...(config?.repeatTypes || {}) },
    };
  }

  private weightFor(eventType: string, now: number): number {
    let weight = this.platformWeights?.[eventType] ?? this.config.weights[eventType] ?? 5;

    const repeatWeight = this.config.repeatTypes[eventType];
    if (repeatWeight !== undefined) {
      const recentCount =
        this.history.filter(
          (item) => item.eventType === eventType && now - item.timestamp < WINDOW_MS
        ).length + 1; // +1 includes the current signal
      if (recentCount >= this.config.repeatCount) {
        weight = repeatWeight;
      }
    }
    return weight;
  }

  addSignal(eventType: string, confidence = 1.0): number {
    const now = Date.now();
    const weight = Math.round(
      this.weightFor(eventType, now) * Math.max(0.1, Math.min(1.0, confidence))
    );

    this.history.push({ eventType, timestamp: now, weight });
    this.history = this.history.filter((item) => now - item.timestamp < WINDOW_MS);

    this.currentScore = Math.min(100, this.history.reduce((acc, curr) => acc + curr.weight, 0));
    return this.currentScore;
  }

  getScore(): number {
    return this.currentScore;
  }

  getLevel(): RiskLevel {
    if (this.currentScore <= this.config.normalThreshold) return "NORMAL";
    if (this.currentScore <= this.config.suspiciousThreshold) return "SUSPICIOUS";
    if (this.currentScore <= this.config.highThreshold) return "HIGH_RISK";
    return "CRITICAL";
  }

  reset(): void {
    this.currentScore = 0;
    this.history = [];
  }
}