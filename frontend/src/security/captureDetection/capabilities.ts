/**
 * Capability detection for each Weyana platform.
 *
 * Every platform gets a capability profile describing which security signals
 * are legitimately exposed. Detection never assumes capabilities — it derives
 * them from what the environment actually provides.
 */
import {
  detectPlatform,
  Platform,
  PlatformCapabilities,
} from "@/security/platformCapabilities";

export { detectPlatform } from "@/security/platformCapabilities";
export type { Platform, PlatformCapabilities } from "@/security/platformCapabilities";

function hasVisibilitySignals(): boolean {
  return typeof document !== "undefined" && "visibilityState" in document;
}

function hasFullscreenSignals(): boolean {
  const doc = document as any;
  return Boolean(
    doc &&
      (typeof doc.fullscreenEnabled === "boolean" || "webkitFullscreenEnabled" in doc)
  ) || typeof window !== "undefined" && Boolean((window as any).Telegram?.WebApp);
}

function hasPageLifecycleSignals(): boolean {
  return (
    typeof document !== "undefined" &&
    (typeof document.addEventListener === "function") &&
    ("visibilitychange" in document || true)
  ) && typeof window !== "undefined" && (
    "pagehide" in window || "pageshow" in window || typeof document.hidden === "boolean"
  );
}

function hasScreenCaptureEvents(): boolean {
  return (
    typeof navigator !== "undefined" &&
    Boolean((navigator as any).mediaDevices?.getDisplayMedia)
  );
}

const DETECTION_UNAVAILABLE = "تعذر التحقق من حالة تسجيل الشاشة على هذا الجهاز.";

function base(platform: Platform): PlatformCapabilities {
  return {
    platform,
    canDetectCapture: false,
    captureDetectionConfidence: "NONE",
    supportsVisibilitySignals: hasVisibilitySignals(),
    supportsFullscreenSignals: hasFullscreenSignals(),
    supportsPageLifecycle: hasPageLifecycleSignals(),
    supportsScreenCaptureEvents: hasScreenCaptureEvents(),
    supportsNativeSecureWindow: false,
    detectionUnavailableMessage: DETECTION_UNAVAILABLE,
  };
}

/**
 * Derive the full capability profile for a platform using only real browser
 * signals. Native secure windows are never assumed for web environments.
 */
export function detectCapabilities(platform: Platform = detectPlatform()): PlatformCapabilities {
  const caps = base(platform);

  const tg = (window as any)?.Telegram?.WebApp;
  const inTelegram = Boolean(tg);

  switch (platform) {
    case "ANDROID_NATIVE":
      // Native Android with an explicit secure-window bridge (FLAG_SECURE).
      // Only reachable through a native host — never through WebView JS.
      caps.canDetectCapture = true;
      caps.captureDetectionConfidence = "HIGH";
      caps.supportsNativeSecureWindow = Boolean((tg as any)?._wiyanaNativeSecureWindow);
      return caps;

    case "ANDROID_WEB":
    case "ANDROID_TELEGRAM_WEBVIEW":
      caps.captureDetectionConfidence = "LOW";
      caps.canDetectCapture = false;
      caps.supportsScreenCaptureEvents = false; // Android WebViews do not expose this.
      caps.supportsNativeSecureWindow = false; // FLAG_SECURE is NOT available here.
      return caps;

    case "IOS_WEB":
    case "IOS_TELEGRAM_WEBVIEW":
      caps.captureDetectionConfidence = "LOW";
      caps.canDetectCapture = false;
      caps.detectionUnavailableMessage = "تعذر التحقق من حالة تسجيل الشاشة على هذا الجهاز.";
      return caps;

    case "DESKTOP_WEB":
      caps.captureDetectionConfidence = caps.supportsScreenCaptureEvents ? "MEDIUM" : "LOW";
      caps.canDetectCapture = caps.supportsScreenCaptureEvents;
      return caps;

    case "TELEGRAM_MINI_APP":
      caps.captureDetectionConfidence = inTelegram ? "LOW" : "NONE";
      caps.canDetectCapture = false;
      return caps;

    default:
      caps.captureDetectionConfidence = "NONE";
      caps.canDetectCapture = false;
      return caps;
  }
}

/**
 * Convenience: resolve platform + capabilities in one call and tag the
 * capabilities with the resolved platform.
 */
export function resolvePlatformCapabilities(): PlatformCapabilities {
  const platform = detectPlatform();
  return detectCapabilities(platform);
}