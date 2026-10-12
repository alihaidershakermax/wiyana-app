/**
 * iOS / iPadOS detector.
 *
 * Safari WebViews expose no direct capture events. This detector collects the
 * same weak signals as other web environments and stays conservative:
 * visibility/focus/fullscreen/lifecycle only. It never claims screen-recording
 * was detected — the warning layer uses non-accusatory Arabic messages.
 */
import {
  AbstractCaptureDetector,
  CaptureDetectorOptions,
} from "@/security/captureDetection/captureDetector";
import { detectCapabilities } from "@/security/captureDetection/capabilities";

export class IosDetector extends AbstractCaptureDetector {
  constructor(options: CaptureDetectorOptions = {}) {
    super({
      ...options,
      capabilities: options.capabilities ?? detectCapabilities("IOS_WEB"),
    });
  }

  protected attachPlatformListeners(): void {
    // Prefer page-lifecycle events which Safari fires more predictably than
    // focus/blur in fullscreen WebViews.
    // (Wiring is inherited from the common listeners; nothing Android-specific.)
  }
}