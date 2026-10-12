/**
 * Capture-detection factory.
 *
 * Capability detection → platform-specific detector → common risk engine.
 * Picks the appropriate detector for the resolved platform.
 */
import {
  AbstractCaptureDetector,
  CaptureDetector,
  CaptureDetectorOptions,
  SecuritySignal,
  SIGNAL_TYPES,
} from "@/security/captureDetection/captureDetector";
import { detectCapabilities, resolvePlatformCapabilities } from "@/security/captureDetection/capabilities";
import { AndroidDetector } from "@/security/captureDetection/androidDetector";
import { IosDetector } from "@/security/captureDetection/iosDetector";
import { TelegramDetector } from "@/security/captureDetection/telegramDetector";
import { DesktopDetector } from "@/security/captureDetection/desktopDetector";
import { Platform, PlatformCapabilities } from "@/security/platformCapabilities";

export {
  AbstractCaptureDetector,
  SIGNAL_TYPES,
} from "@/security/captureDetection/captureDetector";
export type { SecuritySignal, CaptureDetector, CaptureDetectorOptions } from "@/security/captureDetection/captureDetector";
export { detectCapabilities, resolvePlatformCapabilities } from "@/security/captureDetection/capabilities";
export { detectPlatform } from "@/security/captureDetection/capabilities";
export type { Platform, PlatformCapabilities } from "@/security/platformCapabilities";

function detectorForPlatform(platform: Platform, options: CaptureDetectorOptions): CaptureDetector {
  switch (platform) {
    case "ANDROID_NATIVE":
    case "ANDROID_WEB":
    case "ANDROID_TELEGRAM_WEBVIEW":
      return new AndroidDetector(options);
    case "IOS_WEB":
    case "IOS_TELEGRAM_WEBVIEW":
      return new IosDetector(options);
    case "TELEGRAM_MINI_APP":
      return new TelegramDetector(options);
    case "DESKTOP_WEB":
      return new DesktopDetector(options);
    default:
      return new DesktopDetector({
        ...options,
        capabilities: detectCapabilities("UNKNOWN"),
      });
  }
}

/**
 * Create a capture detector for the current (or explicitly selected) platform.
 */
export function createCaptureDetector(
  platform?: Platform,
  options: CaptureDetectorOptions = {}
): CaptureDetector {
  const resolved: Platform = platform ?? resolvePlatformCapabilities().platform;
  return detectorForPlatform(resolved, { ...options });
}