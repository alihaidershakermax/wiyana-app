/**
 * ScreenShield bridge — native on-device screen-capture protection.
 *
 * Web (Telegram WebApp) cannot block OS-level recording; the web layers in
 * services/screenPrivacy.ts + the capture detectors handle what can be handled
 * there. This module activates the NATIVE layer when the app runs as a
 * Capacitor build:
 *
 *  Android: FLAG_SECURE (via this plugin AND @capacitor-community/privacy-screen)
 *           makes screenshots and the system screen recorder capture a BLACK
 *           screen — true OS-level prevention of internal recording.
 *           Android 14+ also delivers capture events for scoring.
 *  iOS:     no blocking API exists; UIScreen.isCaptured is polled during
 *           playback (recording / AirPlay mirroring) and screenshot
 *           notifications arrive as events. Response = instant pause + warn +
 *           backend violation report (10 pts → ban threshold).
 *
 * All native calls are best-effort: on web or when the plugin is missing,
 * every method rejects and is swallowed — the web protection layers keep
 * working unchanged.
 */
import { Capacitor, registerPlugin } from "@capacitor/core";
import { reportViolation } from "@/api/protection";
import { dispatchRecordingDetected } from "@/services/screenPrivacy";

export interface ScreenShieldPlugin {
  enableSecure(): Promise<void>;
  disableSecure(): Promise<void>;
  isCaptured(): Promise<{ captured: boolean }>;
  addListener(
    eventName: "screenCaptureChanged" | "screenshotTaken",
    listener: (info: { captured?: boolean; platform?: string }) => void
  ): Promise<{ remove: () => void }>;
}

const ScreenShield = registerPlugin<ScreenShieldPlugin>("ScreenShield");

export interface ScreenShieldOptions {
  lessonId?: string;
  /** Called the moment a capture is confirmed (pause video, lock UI). */
  onCaptured?: () => void;
  /** Called when a screenshot is detected. */
  onScreenshot?: () => void;
}

let started = false;
let pollTimer: number | null = null;
let listeners: Array<{ remove: () => void }> = [];

export function isNativeApp(): boolean {
  try {
    return Capacitor.isNativePlatform();
  } catch {
    return false;
  }
}

/** True only when running natively AND the ScreenShield plugin resolved. */
export async function isShieldAvailable(): Promise<boolean> {
  if (!isNativeApp()) return false;
  try {
    await ScreenShield.isCaptured();
    return true;
  } catch {
    return false;
  }
}

async function handleCapture(opts: ScreenShieldOptions) {
  dispatchRecordingDetected("native_screen_capture");
  try {
    opts.onCaptured?.();
  } catch {}
  try {
    await reportViolation("screen_recording", { lessonId: opts.lessonId });
  } catch {}
}

async function handleScreenshot(opts: ScreenShieldOptions) {
  dispatchRecordingDetected("native_screenshot");
  try {
    opts.onScreenshot?.();
  } catch {}
  try {
    await reportViolation("screenshot_attempt", { lessonId: opts.lessonId });
  } catch {}
}

export async function startScreenShield(opts: ScreenShieldOptions = {}): Promise<boolean> {
  if (started) return true;
  started = true;
  if (!isNativeApp()) return false;

  // OS-level block on Android (no-op on iOS, kept for parity).
  try {
    await ScreenShield.enableSecure();
  } catch {}

  try {
    const l1 = await ScreenShield.addListener("screenCaptureChanged", (info) => {
      if (info?.captured !== false) void handleCapture(opts);
    });
    const l2 = await ScreenShield.addListener("screenshotTaken", () => {
      void handleScreenshot(opts);
    });
    listeners.push(l1, l2);
  } catch {}

  // iOS: poll isCaptured while playback is active (recording / mirroring).
  // Also acts as a second signal path on Android 14+.
  try {
    pollTimer = window.setInterval(async () => {
      try {
        const res = await ScreenShield.isCaptured();
        if (res?.captured) void handleCapture(opts);
      } catch {}
    }, 1500);
  } catch {}

  return true;
}

export async function stopScreenShield(): Promise<void> {
  started = false;
  if (pollTimer !== null) {
    clearInterval(pollTimer);
    pollTimer = null;
  }
  for (const l of listeners) {
    try {
      l.remove();
    } catch {}
  }
  listeners = [];
}
