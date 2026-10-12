/**
 * Android detector.
 *
 * Android Web / WebView capabilities: visibility, focus, fullscreen and page
 * lifecycle signals only. FLAG_SECURE is NOT available to JavaScript — this
 * detector never references it and prefers LOW/MEDIUM confidence. In Telegram
 * Android WebViews the platform is handled by the Telegram detector profile.
 */
import {
  AbstractCaptureDetector,
  CaptureDetectorOptions,
} from "@/security/captureDetection/captureDetector";
import { detectCapabilities } from "@/security/captureDetection/capabilities";

export class AndroidDetector extends AbstractCaptureDetector {
  constructor(options: CaptureDetectorOptions = {}) {
    super({
      ...options,
      capabilities: options.capabilities ?? detectCapabilities("ANDROID_WEB"),
    });
  }

  protected attachPlatformListeners(): void {
    // Android WebViews expose no capture events; the base visibility/focus/
    // lifecycle wiring is the complete signal surface. Everything else is left
    // to watermark + session protection (deterrence, not detection).
  }
}