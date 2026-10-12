/**
 * Platform-specific response policies for the Weyana protection system.
 *
 * The same risk level leads to different actions depending on whether the
 * platform can back that level with reliable capture evidence. This file picks
 * weights per platform and produces platform-appropriate Arabic messages.
 */
import { Platform } from "@/security/platformCapabilities";
import { RiskLevel } from "@/security/riskEngine";

export type ResponseAction = "continue" | "warn" | "pause" | "revoke";

export interface RiskBand {
  min: number;
  max: number;
  action: ResponseAction;
}

export interface PlatformWarningMessages {
  normal?: string;
  suspicious: string;
  highRisk: string;
  critical: string;
  detectionUnavailable: string;
}

export interface PlatformPolicy {
  platform: Platform;
  /** Signal weights for this platform. */
  weights: Record<string, number>;
  bands: RiskBand[];
  messages: PlatformWarningMessages;
}

const SUSPICIOUS_MSG = "تنبيه: تم اكتشاف نشاط غير اعتيادي أثناء تشغيل المحاضرة.";
const PAUSED_MSG = "تم إيقاف المحاضرة مؤقتًا";
const PAUSED_DETAIL = "يمنع تسجيل أو تصوير أو مشاركة محاضرات Weyana.";
const CRITICAL_MSG = "تم إنهاء جلسة العرض لحماية المحتوى";
const CRITICAL_DETAIL =
  "يُرجى إعادة فتح المحاضرة للتحقق من جلسة العرض والمتابعة.";
const UNAVAILABLE_MSG = "تعذر التحقق من حالة تسجيل الشاشة على هذا الجهاز.";

const DEFAULT_BANDS: RiskBand[] = [
  { min: 0, max: 29, action: "continue" },
  { min: 30, max: 59, action: "warn" },
  { min: 60, max: 79, action: "pause" },
  { min: 80, max: 100, action: "revoke" },
];

/**
 * Spec weights:
 *  reliable native capture +70, reliable browser capture +60, watermark
 *  tamper +50, repeated visibility +15, single blur +5, single
 *  visibilitychange +5. Platforms without reliable capture lower the capture
 *  weight so a mis-fired browser event can't escalate to HIGH/CRITICAL.
 */
const BASE_WEIGHTS: Record<string, number> = {
  CAPTURE_SIGNAL: 60,
  VISIBILITY_HIDDEN: 5,
  REPEATED_VISIBILITY_CHANGE: 15,
  WINDOW_BLUR: 5,
  REPEATED_WINDOW_BLUR: 15,
  FULLSCREEN_CHANGED: 15,
  PAGE_LIFECYCLE_CHANGE: 10,
  PLAYER_TAMPER: 50,
  WATERMARK_TAMPER: 50,
  MULTIPLE_SESSION: 40,
  TOKEN_ANOMALY: 60,
};

function policy(
  platform: Platform,
  overrides: Partial<Pick<PlatformPolicy, "weights" | "messages">> = {}
): PlatformPolicy {
  return {
    platform,
    weights: { ...BASE_WEIGHTS, ...(overrides.weights || {}) },
    bands: DEFAULT_BANDS,
    messages: {
      suspicious: SUSPICIOUS_MSG,
      highRisk: PAUSED_MSG,
      critical: CRITICAL_MSG,
      detectionUnavailable: UNAVAILABLE_MSG,
      ...(overrides.messages || {}),
    },
  };
}

const POLICIES: Record<Platform, PlatformPolicy> = {
  ANDROID_NATIVE: policy("ANDROID_NATIVE", {
    weights: { ...BASE_WEIGHTS, CAPTURE_SIGNAL: 70 },
  }),
  ANDROID_WEB: policy("ANDROID_WEB", {
    weights: { ...BASE_WEIGHTS, CAPTURE_SIGNAL: 45 },
  }),
  ANDROID_TELEGRAM_WEBVIEW: policy("ANDROID_TELEGRAM_WEBVIEW", {
    weights: { ...BASE_WEIGHTS, CAPTURE_SIGNAL: 45 },
  }),
  IOS_WEB: policy("IOS_WEB", {
    weights: { ...BASE_WEIGHTS, CAPTURE_SIGNAL: 40 },
  }),
  IOS_TELEGRAM_WEBVIEW: policy("IOS_TELEGRAM_WEBVIEW", {
    weights: { ...BASE_WEIGHTS, CAPTURE_SIGNAL: 40 },
  }),
  TELEGRAM_MINI_APP: policy("TELEGRAM_MINI_APP", {
    weights: { ...BASE_WEIGHTS, CAPTURE_SIGNAL: 40 },
  }),
  DESKTOP_WEB: policy("DESKTOP_WEB"),
  UNKNOWN: policy("UNKNOWN", {
    weights: { ...BASE_WEIGHTS, CAPTURE_SIGNAL: 40 },
  }),
};

export function getPlatformPolicy(platform: Platform): PlatformPolicy {
  return POLICIES[platform] ?? POLICIES.UNKNOWN;
}

export function actionForLevel(policy: PlatformPolicy, level: RiskLevel): ResponseAction {
  if (level === "CRITICAL") return "revoke";
  if (level === "HIGH_RISK") return "pause";
  if (level === "SUSPICIOUS") return "warn";
  return "continue";
}

export function actionForScore(policy: PlatformPolicy, score: number): ResponseAction {
  const band = policy.bands.find(
    (b) => score >= b.min && score <= b.max
  );
  return band?.action ?? "continue";
}

export function messageForLevel(
  policy: PlatformPolicy,
  level: RiskLevel,
  detail = false
): string {
  switch (level) {
    case "SUSPICIOUS":
      return policy.messages.suspicious;
    case "HIGH_RISK":
      return detail ? PAUSED_DETAIL : policy.messages.highRisk;
    case "CRITICAL":
      return detail ? CRITICAL_DETAIL : policy.messages.critical;
    default:
      return "";
  }
}

/** Conservative Arabic message used when the platform cannot verify capture. */
export function unavailableMessage(policy: PlatformPolicy): string {
  return policy.messages.detectionUnavailable;
}