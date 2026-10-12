import { reportViolation } from "@/api/protection";
import { ScreenshotDetection } from "@/services/screenshotDetection";
import {
  SCREEN_RECORDING_EVENT,
  dispatchRecordingDetected,
  screenPrivacy,
} from "@/services/screenPrivacy";
import type { AccessStatus, ReportResult, ViolationType } from "@/types/protection";
import { detectPlatform, type Platform } from "@/security/platformCapabilities";

export interface ContentProtectionOptions {
  lessonId?: string;
  enabled?: boolean;
  closeOnCapture?: boolean;
  blockContextMenu?: boolean;
  onBanned?: (result: ReportResult) => void;
  onWarning?: (result: ReportResult) => void;
  onViolation?: (type: ViolationType, detectedBy: string) => void;
  onPrivacyPause?: (detectedBy: string) => void;
}

type ScreenshotPreventionClass = new (options?: {
  blurIntensity?: string;
  warningMessage?: string;
  preventCopy?: boolean;
  preventInspect?: boolean;
  recoveryDelay?: number;
  debug?: boolean;
  onAttempt?: (details: { method: string; count: number; timestamp: number; details?: string }) => void;
}) => { destroy(): void };

type AntiScreenshotClass = new (options?: {
  blockShortcuts?: boolean;
  detectRecording?: boolean;
  blockDevTools?: boolean;
  disableRightClick?: boolean;
  detectDOMManipulation?: boolean;
  targetElements?: string[];
  systemConfigs?: Record<string, Record<string, boolean>>;
}) => void;

const ANTI_SCREENSHOT_ALERT_MESSAGES: Record<string, ViolationType> = {
  "Screenshots are disabled!": "screenshot_attempt",
  "Screen recording detected!": "screen_recording",
  "DevTools detected!": "devtools_open",
};

export class ContentProtectionEngine {
  private readonly opts: Required<ContentProtectionOptions>;
  private started = false;
  private lastReport = 0;
  private playing = false;
  private extLib: ScreenshotDetection | null = null;
  private screenshotPrevention: { destroy(): void } | null = null;
  private antiScreenshot: unknown = null;
  private devtoolsTimer: ReturnType<typeof setInterval> | null = null;
  private captureTimer: ReturnType<typeof setTimeout> | null = null;
  private readonly platform: Platform;

  constructor(opts: ContentProtectionOptions = {}) {
    this.opts = {
      lessonId: opts.lessonId ?? "",
      enabled: opts.enabled ?? true,
      closeOnCapture: opts.closeOnCapture ?? true,
      blockContextMenu: opts.blockContextMenu ?? true,
      onBanned: opts.onBanned ?? (() => {}),
      onWarning: opts.onWarning ?? (() => {}),
      onViolation: opts.onViolation ?? (() => {}),
      onPrivacyPause: opts.onPrivacyPause ?? (() => {}),
    };
    this.platform = detectPlatform();
  }

  get enabled(): boolean {
    return this.opts.enabled;
  }

  get isPlaying(): boolean {
    return this.playing;
  }

  setPlaying(value: boolean): void {
    this.playing = value;
  }

  start(): void {
    if (this.started || !this.enabled) return;
    this.started = true;
    window.addEventListener("keydown", this.onKeyDown, { capture: true });
    window.addEventListener("contextmenu", this.onContextMenu, { capture: true });
    document.addEventListener("visibilitychange", this.onVisibilityChange);
    window.addEventListener("blur", this.onWindowBlur);
    window.addEventListener("focus", this.onWindowFocus);
    window.addEventListener("pagehide", this.onPageHide);
    document.addEventListener("fullscreenchange", this.onFullscreenChange);
    document.addEventListener("webkitfullscreenchange", this.onFullscreenChange);
    document.addEventListener("mozfullscreenchange", this.onFullscreenChange);
    document.addEventListener("MSFullscreenChange", this.onFullscreenChange);
    window.addEventListener(SCREEN_RECORDING_EVENT, this.onScreenRecordingDetected);
    screenPrivacy.init();
    screenPrivacy.enable();
    this.startDevtoolsHeuristic();
    this.loadDynamicLibrary();
    this.loadScreenshotPrevention();
    this.loadAntiScreenshot();
  }

  stop(): void {
    if (!this.started) return;
    this.started = false;
    window.removeEventListener("keydown", this.onKeyDown, { capture: true });
    window.removeEventListener("contextmenu", this.onContextMenu, { capture: true });
    document.removeEventListener("visibilitychange", this.onVisibilityChange);
    window.removeEventListener("blur", this.onWindowBlur);
    window.removeEventListener("focus", this.onWindowFocus);
    window.removeEventListener("pagehide", this.onPageHide);
    document.removeEventListener("fullscreenchange", this.onFullscreenChange);
    document.removeEventListener("webkitfullscreenchange", this.onFullscreenChange);
    document.removeEventListener("mozfullscreenchange", this.onFullscreenChange);
    document.removeEventListener("MSFullscreenChange", this.onFullscreenChange);
    window.removeEventListener(SCREEN_RECORDING_EVENT, this.onScreenRecordingDetected);
    if (this.devtoolsTimer) {
      clearInterval(this.devtoolsTimer);
      this.devtoolsTimer = null;
    }
    if (this.captureTimer) {
      clearTimeout(this.captureTimer);
      this.captureTimer = null;
    }
    this.extLib?.destroy();
    this.extLib = null;
    this.screenshotPrevention?.destroy();
    this.screenshotPrevention = null;
    this.antiScreenshot = null;
  }

  async checkAccessNow(): Promise<AccessStatus> {
    const { checkAccess } = await import("@/api/protection");
    return checkAccess();
  }

  getTelegramUserId(): number | undefined {
    return (window as any).Telegram?.WebApp?.initDataUnsafe?.user?.id;
  }

  requestFullscreen(): boolean {
    const tg = (window as any).Telegram?.WebApp;
    if (tg?.requestFullscreen) {
      try {
        tg.requestFullscreen();
        return true;
      } catch {
        /* fall through to the browser API */
      }
    }
    const el = document.documentElement as any;
    try {
      if (el?.requestFullscreen) {
        void el.requestFullscreen();
        return true;
      }
      if (el?.webkitRequestFullscreen) {
        el.webkitRequestFullscreen();
        return true;
      }
    } catch {
      return false;
    }
    return false;
  }

  closeApp(): void {
    const tg = (window as any).Telegram?.WebApp;
    try {
      tg?.close?.();
    } catch {
      /* best effort */
    }
  }

  /**
   * Called when the forced fullscreen is lost while content is playing.
   * Spec: report a `fullscreen_exit` violation and close the app.
   */
  handleFullscreenExit(detectedBy = "fullscreenchange"): void {
    this.handleViolation("fullscreen_exit", detectedBy);
    this.closeApp();
  }

  handleViolation(type: ViolationType, detectedBy: string, meta?: Record<string, unknown>): void {
    if (!this.enabled) return;
    // Stop the video / show the notice on EVERY detection, even when a
    // server report was just sent (repeated signals within the throttle
    // window must still freeze playback).
    this.opts.onViolation(type, detectedBy);

    // Throttle only the server-side report to avoid flooding the ledger.
    const now = Date.now();
    if (now - this.lastReport < 2000) return;
    this.lastReport = now;
    void this.report(type, detectedBy, meta);
  }

  private async report(type: ViolationType, detectedBy: string, meta?: Record<string, unknown>): Promise<void> {
    try {
      const result = await reportViolation(type, {
        lessonId: this.opts.lessonId,
        meta: { detectedBy, ...meta },
      });
      if (result.banned) this.opts.onBanned(result);
      else if (result.status === "warning") this.opts.onWarning(result);
      if (this.opts.closeOnCapture) {
        this.closeApp();
      }
    } catch {
      // A transport failure must never become a local ban. The player is already
      // paused by onViolation; show a neutral warning and let the server decide
      // whether the event is a violation after connectivity returns.
      const fallbackResult: ReportResult = {
        success: false,
        user_id: String(this.getTelegramUserId() || "user"),
        violation_type: type,
        points: 0,
        status: "warning",
        banned: false,
        behavioral: false,
        warning_threshold: 0,
        ban_threshold: 0,
        ban_reason: null,
        ban_until: null,
        ban_days: null,
      };
      this.opts.onWarning(fallbackResult);
    }
  }

  private scheduleCaptureCheck(detectedBy: string): void {
    if (!this.started || !this.enabled) return;
    if (this.captureTimer) clearTimeout(this.captureTimer);
    // Only report when tab/app is explicitly hidden by user switching away
    this.captureTimer = setTimeout(() => {
      if (document.visibilityState === "hidden") {
        this.handleViolation("screen_recording", detectedBy);
      }
    }, 1000);
  }

  private onKeyDown = (event: KeyboardEvent): void => {
    const key = event.key;
    const ctrl = event.ctrlKey || event.metaKey;
    const shift = event.shiftKey;
    const blocked =
      key === "F12" ||
      key === "PrintScreen" ||
      (ctrl && shift && (key === "I" || key === "J" || key === "C")) ||
      (ctrl && key.toLowerCase() === "u") ||
      (ctrl && shift && key.toLowerCase() === "p");
    if (blocked) {
      event.preventDefault();
      event.stopPropagation();
      this.handleViolation("devtools_open", "keydown");
    }
  };

  private onContextMenu = (event: MouseEvent): void => {
    if (!this.opts.blockContextMenu) return;
    const target = event.target as HTMLElement | null;
    if (target?.closest?.("video") || this.playing) {
      event.preventDefault();
    }
  };

  private onVisibilityChange = (): void => {
    if (document.visibilityState !== "hidden") return;
    if (this.playing) {
      // iOS/Telegram WebViews do not expose trustworthy capture signals. Pause
      // conservatively without reporting or banning on an ambiguous lifecycle event.
      if (this.isConservativeCapturePlatform()) {
        this.pauseForPrivacy("visibility-hidden");
      } else {
        this.handleViolation("screen_recording", "visibility-hidden", { instant_ban: true });
      }
      return;
    }
    // Non-playing: only report if still hidden after 1s (avoids false positives).
    this.scheduleCaptureCheck("visibility-hidden");
  };

  private onWindowBlur = (): void => {
    // Native screen recorders (built-in recorder, OBS, floating bubble,
    // notification-shade controls) grab focus WITHOUT necessarily hiding the
    // document. While a video is actively playing, focus loss is the strongest
    // available signal that a recorder/control overlay is on screen.
    if (this.playing) {
      if (this.isConservativeCapturePlatform()) {
        this.pauseForPrivacy("window-blur-while-playing");
      } else {
        this.handleViolation("screen_recording", "window-blur-while-playing");
      }
      return;
    }
    // Fallback: when the tab is explicitly hidden by switching away.
    if (document.visibilityState === "hidden") {
      this.scheduleCaptureCheck("window-blur");
    }
  };

  private onWindowFocus = (): void => {
    /* window regained focus */
  };

  private onPageHide = (): void => {
    if (this.playing) {
      if (this.isConservativeCapturePlatform()) {
        this.pauseForPrivacy("pagehide");
      } else {
        this.handleViolation("screen_recording", "pagehide", { instant_ban: true });
      }
    }
  };

  private pauseForPrivacy(detectedBy: string): void {
    this.playing = false;
    this.opts.onPrivacyPause(detectedBy);
  }

  private isConservativeCapturePlatform(): boolean {
    return ["IOS_WEB", "IOS_TELEGRAM_WEBVIEW", "ANDROID_TELEGRAM_WEBVIEW", "TELEGRAM_MINI_APP"].includes(this.platform);
  }

  private onFullscreenChange = (): void => {
    const doc = document as any;
    const isFS = !!(doc.fullscreenElement || doc.webkitFullscreenElement || doc.mozFullScreenElement || doc.msFullscreenElement);
    if (!isFS && this.playing) {
      if (this.isConservativeCapturePlatform()) {
        this.pauseForPrivacy("fullscreenchange");
      } else {
        this.handleFullscreenExit("fullscreenchange");
      }
    }
  };

  private onScreenRecordingDetected = (event: Event): void => {
    const detail = (event as CustomEvent<{ source: string }>).detail;
    this.handleViolation("screen_recording", detail?.source || "screen-privacy");
  };

  private startDevtoolsHeuristic(): void {
    // Skip DevTools window-size heuristic on mobile/touch/Telegram devices
    const isTouch = typeof window !== "undefined" && ("ontouchstart" in window || (navigator.maxTouchPoints ?? 0) > 0);
    const isTelegram = typeof window !== "undefined" && Boolean((window as any).Telegram?.WebApp);
    if (isTouch || isTelegram) return;

    this.devtoolsTimer = setInterval(() => {
      const threshold = 220;
      const widthGap = window.outerWidth - window.innerWidth;
      const heightGap = window.outerHeight - window.innerHeight;
      if (widthGap > threshold || heightGap > threshold) {
        this.handleViolation("devtools_open", "window-size-heuristic");
      }
    }, 5000);
  }

  private loadDynamicLibrary(): void {
    this.extLib = new ScreenshotDetection({
      showOverlay: true,
      onAttempt: (details) => {
        // Strong capture signals also flow through the screen-privacy channel.
        dispatchRecordingDetected(details.details || details.method);
        this.handleViolation(
          details.method === "screenCapture" ? "screen_recording" : "screenshot_attempt",
          `screenshot-detection:${details.method}`
        );
      },
    });
  }

  private async loadScreenshotPrevention(): Promise<void> {
    try {
      const mod: any = await import("screenshot-prevention");
      const EnhancedScreenshotPrevention: ScreenshotPreventionClass =
        mod.default ?? mod.EnhancedScreenshotPrevention;
      this.screenshotPrevention = new EnhancedScreenshotPrevention({
        blurIntensity: "30px",
        warningMessage: "تسجيل الشاشة غير مسموح به",
        preventCopy: true,
        preventInspect: true,
        recoveryDelay: 2000,
        debug: false,
        onAttempt: (details) => {
          this.handleViolation("screenshot_attempt", `screenshot-prevention:${details.method}`);
        },
      });
    } catch (e) {
      console.warn("screenshot-prevention failed to load:", e);
    }
  }

  private async loadAntiScreenshot(): Promise<void> {
    try {
      // The published d.ts describes only static members, but the runtime class
      // is constructed with options. Load it loosely.
      const mod: any = await import("anti-screenshot-lib");
      const AntiScreenshot: AntiScreenshotClass = mod.default ?? mod;
      // Its `alert()` calls are blocked inside Telegram; route them into the
      // violation pipeline instead of relying on native dialogs.
      this.installAlertBridge();
      this.antiScreenshot = new AntiScreenshot({
        blockShortcuts: true,
        // `detectRecording` re-hooks getDisplayMedia and calls
        // `location.reload()` on any screen-share attempt → that destroys the
        // lesson session instead of just stopping the video. Our own layers
        // (screenPrivacy + ScreenshotDetection + screenshot-prevention) already
        // intercept getDisplayMedia and pause the player, so we disable it here.
        detectRecording: false,
        // `blockDevTools` polls outerWidth-innerWidth every 1s and produces
        // false positives on mobile/Telegram; our own heuristic is guarded for
        // touch devices, so this one is disabled.
        blockDevTools: false,
        // `detectDOMManipulation` observes body mutations → false-positives
        // under Vue's virtual DOM. `disableRightClick` we handle ourselves.
        detectDOMManipulation: false,
        disableRightClick: false,
        targetElements: ["#video-container"],
      });
    } catch (e) {
      console.warn("anti-screenshot-lib failed to load:", e);
    }
  }

  private installAlertBridge(): void {
    if ((window as any).__wiyanaAlertBridge) return;
    (window as any).__wiyanaAlertBridge = true;
    const original = window.alert.bind(window);
    window.alert = (message?: unknown) => {
      const text = typeof message === "string" ? message : "";
      const type = ANTI_SCREENSHOT_ALERT_MESSAGES[text];
      if (type) {
        this.handleViolation(type, "anti-screenshot-lib:alert");
      } else {
        original(text);
      }
    };
  }
}

export function createContentProtection(opts: ContentProtectionOptions = {}): ContentProtectionEngine {
  return new ContentProtectionEngine(opts);
}
