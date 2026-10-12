/**
 * Weyana Platform Capability Layer.
 *
 * Detects which platform the protected player is running on and documents
 * what this platform can (and cannot) do. The rest of the security system
 * must NEVER assume native capabilities (FLAG_SECURE, iOS/Android APIs) are
 * reachable from a WebView / browser — this layer is the source of truth for
 * what signals are legitimately available.
 */

export const PLATFORM_IDS = [
  "TELEGRAM_MINI_APP",
  "ANDROID_NATIVE",
  "ANDROID_WEB",
  "ANDROID_TELEGRAM_WEBVIEW",
  "IOS_WEB",
  "IOS_TELEGRAM_WEBVIEW",
  "DESKTOP_WEB",
  "UNKNOWN",
] as const;

export type Platform = (typeof PLATFORM_IDS)[number];

export type CaptureDetectionConfidence = "NONE" | "LOW" | "MEDIUM" | "HIGH";

export interface PlatformCapabilities {
  /** Resolved platform profile. */
  platform: Platform;
  /** Whether reliable capture detection is exposed to this environment. */
  canDetectCapture: boolean;
  /** How trustworthy direct capture detection is on this platform. */
  captureDetectionConfidence: CaptureDetectionConfidence;
  /** visibilitychange / document.visibilityState are available. */
  supportsVisibilitySignals: boolean;
  /** fullscreen API (or Telegram fullscreen events) are available. */
  supportsFullscreenSignals: boolean;
  /** Page Lifecycle API (pagehide / pageshow / freeze events) is available. */
  supportsPageLifecycle: boolean;
  /** Capture-related browser APIs (e.g. getDisplayMedia) are legitimately exposed. */
  supportsScreenCaptureEvents: boolean;
  /** True only for native applications with FLAG_SECURE. Never true for web. */
  supportsNativeSecureWindow: boolean;
  /** Human-friendly message used when capture cannot be verified. */
  detectionUnavailableMessage: string;
}

function hasWindow(): boolean {
  return typeof window !== "undefined";
}

function getTg(): Record<string, any> | null {
  if (!hasWindow()) return null;
  return (window as any).Telegram?.WebApp ?? null;
}

/**
 * Detect the running platform using only signals legitimately available to
 * JavaScript. Telegram WebViews are recognized through the Telegram WebApp
 * SDK surface, never through native APIs.
 */
export function detectPlatform(): Platform {
  if (!hasWindow()) return "UNKNOWN";

  const tg = getTg();
  if (tg) {
    const tgPlatform = String(tg.platform || "").toLowerCase();
    if (tgPlatform.includes("android")) return "ANDROID_TELEGRAM_WEBVIEW";
    if (tgPlatform.includes("ios") || tgPlatform === "iphone" || tgPlatform.includes("ipad")) {
      return "IOS_TELEGRAM_WEBVIEW";
    }
    // In Telegram but the SDK did not report the OS — generic Mini App.
    return "TELEGRAM_MINI_APP";
  }

  const ua = (navigator?.userAgent || "").toLowerCase();
  const isAndroid = /android/.test(ua) || /\(linux; android/.test(ua);
  const isIos = /iphone|ipad|ipod/.test(ua);

  if (isAndroid) return "ANDROID_WEB";
  if (isIos) return "IOS_WEB";

  const isDesktop =
    typeof navigator !== "undefined" &&
    (typeof (navigator as any).userAgentData === "object" ||
      /windows|mac os|linux|x11|chrome os|crkey/.test(ua)) &&
    !isAndroid &&
    !isIos;

  if (isDesktop) return "DESKTOP_WEB";

  return "UNKNOWN";
}

export const PLATFORM_LABEL: Record<Platform, string> = {
  TELEGRAM_MINI_APP: "Telegram Mini App",
  ANDROID_NATIVE: "Android Native",
  ANDROID_WEB: "Android Web",
  ANDROID_TELEGRAM_WEBVIEW: "Android Telegram WebView",
  IOS_WEB: "iOS / iPadOS Web",
  IOS_TELEGRAM_WEBVIEW: "iOS Telegram WebView",
  DESKTOP_WEB: "Desktop Web",
  UNKNOWN: "غير معروف",
};