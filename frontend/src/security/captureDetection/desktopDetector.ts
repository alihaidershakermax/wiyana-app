/**
 * Desktop browser detector.
 *
 * Desktop browsers are the only web environments where capture-related APIs
 * (getDisplayMedia) can exist. This detector hooks getDisplayMedia — when a
 * display-capture stream becomes active and ends, a CAPTURE_SIGNAL is emitted
 * with MEDIUM confidence. The remaining wiring relies on visibility/focus/
 * fullscreen/lifecycle. It never claims universal recorder detection.
 */
import {
  AbstractCaptureDetector,
  CaptureDetectorOptions,
  SIGNAL_TYPES,
} from "@/security/captureDetection/captureDetector";
import { detectCapabilities } from "@/security/captureDetection/capabilities";

export class DesktopDetector extends AbstractCaptureDetector {
  private originalGetDisplayMedia: ((opts?: any) => Promise<MediaStream>) | null = null;
  private captureActive = false;

  constructor(options: CaptureDetectorOptions = {}) {
    super({
      ...options,
      capabilities: options.capabilities ?? detectCapabilities("DESKTOP_WEB"),
    });
  }

  protected attachPlatformListeners(): void {
    // Guard: never patch when the API is absent.
    const md = (navigator as any)?.mediaDevices;
    const target = md?.getDisplayMedia;
    if (typeof target !== "function") return;

    this.originalGetDisplayMedia = target;

    md.getDisplayMedia = ((...args: any[]) => {
      const result = target.apply(md, args);

      if (result && typeof result.then === "function") {
        result
          .then((stream: MediaStream) => {
            this.captureActive = true;
            stream.getTracks().forEach((track) => {
              track.addEventListener?.("ended", () => {
                this.captureActive = false;
                this.emitSignal(SIGNAL_TYPES.CAPTURE_SIGNAL, 0.6, {
                  source: "getDisplayMedia",
                  state: "ended",
                });
              });
            });
          })
          .catch(() => {
            this.captureActive = false;
          });
      }
      return result;
    }) as any;
  }

  destroy(): void {
    super.destroy();
    if (this.originalGetDisplayMedia) {
      try {
        (navigator as any).mediaDevices.getDisplayMedia = this.originalGetDisplayMedia;
      } catch {
        /* best effort restore */
      }
      this.originalGetDisplayMedia = null;
    }
  }
}