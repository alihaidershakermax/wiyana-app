/**
 * Common capture detector interface + base implementation.
 *
 * A CaptureDetector is responsible for a single platform's worth of signal
 * collection. Every detector MUST clean up its listeners on destroy() and
 * expose the raw SecuritySignals it collected since initialize().
 */
import { PlatformCapabilities } from "@/security/platformCapabilities";
import { detectCapabilities } from "@/security/captureDetection/capabilities";

export interface SecuritySignal {
  /** One of SIGNAL_TYPES.values (custom platform signals allowed). */
  type: string;
  /** 0.0 - 1.0 signal confidence. */
  confidence: number;
  /** ISO-8601 timestamp. */
  timestamp: string;
  metadata: Record<string, any>;
}

export const SIGNAL_TYPES = {
  CAPTURE_SIGNAL: "CAPTURE_SIGNAL",
  VISIBILITY_HIDDEN: "VISIBILITY_HIDDEN",
  WINDOW_BLUR: "WINDOW_BLUR",
  FULLSCREEN_CHANGED: "FULLSCREEN_CHANGED",
  PAGE_LIFECYCLE_CHANGE: "PAGE_LIFECYCLE_CHANGE",
  PLAYER_TAMPER: "PLAYER_TAMPER",
  WATERMARK_TAMPER: "WATERMARK_TAMPER",
  MULTIPLE_SESSION: "MULTIPLE_SESSION",
  TOKEN_ANOMALY: "TOKEN_ANOMALY",
} as const;

export type SignalType = (typeof SIGNAL_TYPES)[keyof typeof SIGNAL_TYPES];

export interface CaptureDetector {
  initialize(): void;
  destroy(): void;
  getCapabilities(): PlatformCapabilities;
  getSignals(): SecuritySignal[];
}

export interface CaptureDetectorOptions {
  onSignal?: (signal: SecuritySignal) => void;
  capabilities?: PlatformCapabilities;
}

const MAX_SIGNALS = 200;

export abstract class AbstractCaptureDetector implements CaptureDetector {
  protected capabilities: PlatformCapabilities;
  protected signals: SecuritySignal[] = [];
  protected onSignal?: (signal: SecuritySignal) => void;
  protected initialized = false;

  protected blurTimes: number[] = [];
  protected hiddenChanges: number[] = [];
  protected fullscreenExpected = false;

  private boundListeners: Array<[EventTarget, string, EventListener]> = [];

  constructor(options: CaptureDetectorOptions = {}) {
    this.capabilities = options.capabilities ?? detectCapabilities();
    this.onSignal = options.onSignal;
  }

  initialize(): void {
    if (this.initialized) return;
    this.initialized = true;
    this.attachCommonListeners();
    this.attachPlatformListeners();
  }

  destroy(): void {
    if (!this.initialized) return;
    this.initialized = false;
    for (const [target, type, handler] of this.boundListeners) {
      target.removeEventListener(type, handler);
    }
    this.boundListeners = [];
    this.blurTimes = [];
    this.hiddenChanges = [];
  }

  getCapabilities(): PlatformCapabilities {
    return this.capabilities;
  }

  getSignals(): SecuritySignal[] {
    return [...this.signals];
  }

  /** Platform-specific listeners — subclasses override. */
  protected abstract attachPlatformListeners(): void;

  protected emitSignal(type: string, confidence: number, metadata: Record<string, any> = {}): void {
    const signal: SecuritySignal = {
      type,
      confidence: Math.min(1, Math.max(0, confidence)),
      timestamp: new Date().toISOString(),
      metadata,
    };
    this.signals.push(signal);
    if (this.signals.length > MAX_SIGNALS) {
      this.signals.splice(0, this.signals.length - MAX_SIGNALS);
    }
    this.onSignal?.(signal);
  }

  protected addListener(target: EventTarget | null, type: string, handler: EventListener): void {
    if (!target || typeof target.addEventListener !== "function") return;
    target.addEventListener(type, handler, { passive: true });
    this.boundListeners.push([target, type, handler]);
  }

  protected pushTimestamp(now: number, bucket: number[]): void {
    bucket.push(now);
    // Keep only entries within the last 15 seconds.
    const cutoff = now - 15000;
    let i = 0;
    while (i < bucket.length && bucket[i] < cutoff) i++;
    bucket.splice(0, i);
  }

  // ── common web lifecycle wiring ────────────────────────────────────────
  private attachCommonListeners(): void {
    this.addListener(document, "visibilitychange", this.handleVisibility);
    this.addListener(window, "blur", this.handleBlur);
    this.addListener(window, "focus", this.handleFocus);
    this.addListener(document, "fullscreenchange", this.handleFullscreen);
    this.addListener(document, "webkitfullscreenchange", this.handleFullscreen);
    this.addListener(document, "mozfullscreenchange", this.handleFullscreen);
    this.addListener(document, "MSFullscreenChange", this.handleFullscreen);

    if (this.capabilities.supportsPageLifecycle) {
      this.addListener(window, "pagehide", () =>
        this.emitSignal(SIGNAL_TYPES.PAGE_LIFECYCLE_CHANGE, 0.5, { event: "pagehide" })
      );
      this.addListener(window, "pageshow", () =>
        this.emitSignal(SIGNAL_TYPES.PAGE_LIFECYCLE_CHANGE, 0.2, { event: "pageshow" })
      );
      this.addListener(document, "freeze", () =>
        this.emitSignal(SIGNAL_TYPES.PAGE_LIFECYCLE_CHANGE, 0.5, { event: "freeze" })
      );
      this.addListener(document, "resume", () =>
        this.emitSignal(SIGNAL_TYPES.PAGE_LIFECYCLE_CHANGE, 0.2, { event: "resume" })
      );
    }
  }

  private readonly handleVisibility: EventListener = () => {
    if (typeof document === "undefined") return;
    if (document.visibilityState === "hidden") {
      const now = Date.now();
      this.pushTimestamp(now, this.hiddenChanges);
      const count = this.hiddenChanges.length;
      // Confidence rises with suspicious (repeated) switches.
      const confidence = count >= 3 ? 1.0 : 0.5;
      this.emitSignal(SIGNAL_TYPES.VISIBILITY_HIDDEN, confidence, {
        state: "hidden",
        changesInWindow: count,
      });
    }
  };

  private readonly handleBlur: EventListener = () => {
    const now = Date.now();
    this.pushTimestamp(now, this.blurTimes);
    const count = this.blurTimes.length;
    const confidence = count >= 3 ? 1.0 : 0.4;
    this.emitSignal(SIGNAL_TYPES.WINDOW_BLUR, confidence, {
      blursInWindow: count,
    });
  };

  private readonly handleFocus: EventListener = () => {
    // Return to window is not evidence of capture — no signal emitted.
  };

  private readonly handleFullscreen: EventListener = () => {
    const doc = document as any;
    const entered = Boolean(
      doc.fullscreenElement ||
        doc.webkitFullscreenElement ||
        doc.mozFullScreenElement ||
        doc.msFullscreenElement
    );
    if (entered) {
      this.fullscreenExpected = true;
      return;
    }
    const wasExpected = this.fullscreenExpected;
    this.fullscreenExpected = false;
    if (wasExpected) {
      this.emitSignal(SIGNAL_TYPES.FULLSCREEN_CHANGED, 0.6, { entered: false });
    }
  };
}