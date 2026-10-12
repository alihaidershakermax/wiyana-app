/**
 * Telegram Mini App detector.
 *
 * Runs inside a Telegram WebView. The Mini App can only use Telegram WebApp
 * surface + standard web signals — never native Android/iOS APIs, never
 * FLAG_SECURE. Detection confidence stays LOW/MEDIUM and the deterrence layer
 * (watermark, session tracking, risk scoring) carries the real weight.
 */
import {
  AbstractCaptureDetector,
  CaptureDetectorOptions,
  SIGNAL_TYPES,
} from "@/security/captureDetection/captureDetector";
import { detectCapabilities } from "@/security/captureDetection/capabilities";
import { Platform } from "@/security/platformCapabilities";

export class TelegramDetector extends AbstractCaptureDetector {
  private tg: Record<string, any> | null = null;

  constructor(options: CaptureDetectorOptions = {}) {
    const platform: Platform = options.capabilities?.platform ?? "TELEGRAM_MINI_APP";
    super({
      ...options,
      capabilities: options.capabilities ?? detectCapabilities(platform),
    });
  }

  protected attachPlatformListeners(): void {
    this.tg = (window as any)?.Telegram?.WebApp ?? null;

    // A WebView banned by Telegram reports a "cant_open" event for the Mini App.
    if (this.tg?.onEvent) {
      this.tg.onEvent("viewportChanged", () => {
        if (this.tg?.isExpanded === false) {
          this.emitSignal(SIGNAL_TYPES.PAGE_LIFECYCLE_CHANGE, 0.5, { event: "telegram_collapsed" });
        }
      });
    }

    // Nothing else: Telegram provides no capture signals to JavaScript.
  }
}