export interface ScreenshotDetectionAttempt {
  count: number;
  method: string;
  timestamp: number;
  details?: string;
}

export interface ScreenshotDetectionOptions {
  onAttempt?: (details: ScreenshotDetectionAttempt) => void;
  showOverlay?: boolean;
}

const HOOKED_DISPLAY_MEDIA = "__wiyanaDisplayMediaHooked";

export class ScreenshotDetection {
  private count = 0;
  private destroyed = false;
  private overlay: HTMLDivElement | null = null;
  private cleanupFns: (() => void)[] = [];
  private readonly options: Required<ScreenshotDetectionOptions>;

  constructor(options: ScreenshotDetectionOptions = {}) {
    this.options = {
      onAttempt: options.onAttempt ?? (() => {}),
      showOverlay: options.showOverlay ?? true,
    };
    this.setup();
  }

  private setup(): void {
    if (typeof window === "undefined") return;

    // Hook getDisplayMedia — the strongest signal the user is screen-sharing
    // or recording the tab. Hooked once per page.
    const mediaDevices = (navigator as any).mediaDevices as any;
    const original = mediaDevices?.getDisplayMedia;
    if (original && !(window as any)[HOOKED_DISPLAY_MEDIA]) {
      (window as any)[HOOKED_DISPLAY_MEDIA] = true;
      mediaDevices.getDisplayMedia = (...args: unknown[]) => {
        this.attempt("screenCapture", "getDisplayMedia");
        return original.apply(mediaDevices, args);
      };
    }

    // PrintScreen key.
    const onKeyUp = (event: KeyboardEvent) => {
      if (event.key === "PrintScreen") this.attempt("keyboard", "PrintScreen");
    };
    window.addEventListener("keyup", onKeyUp);
    this.cleanupFns.push(() => window.removeEventListener("keyup", onKeyUp));

    // When the app goes hidden mid-playback, darken the page so a later
    // screenshot of the tab cannot reveal the video frame. The overlay is
    // removed as soon as the user returns to the tab.
    const onVisibility = () => {
      if (document.visibilityState === "hidden") this.showOverlay();
      else this.hideOverlay();
    };
    document.addEventListener("visibilitychange", onVisibility);
    this.cleanupFns.push(() => document.removeEventListener("visibilitychange", onVisibility));
  }

  private attempt(method: string, details?: string): void {
    if (this.destroyed) return;
    this.count += 1;
    this.options.onAttempt({ count: this.count, method, timestamp: Date.now(), details });
  }

  private showOverlay(): void {
    if (!this.options.showOverlay || this.overlay || this.destroyed) return;
    const el = document.createElement("div");
    el.setAttribute("aria-hidden", "true");
    Object.assign(el.style, {
      position: "fixed",
      inset: "0",
      zIndex: "2147483646",
      background: "rgba(0,0,0,0.98)",
      pointerEvents: "none",
      userSelect: "none",
      display: "flex",
      alignItems: "center",
      justifyContent: "center",
      color: "#ffffff",
      fontSize: "14px",
      fontFamily: "system-ui, sans-serif",
    });
    el.textContent = "المحتوى محمي — لا يمكن التقاط الشاشة";
    document.body.appendChild(el);
    this.overlay = el;
  }

  private hideOverlay(): void {
    if (!this.overlay) return;
    this.overlay.remove();
    this.overlay = null;
  }

  destroy(): void {
    if (this.destroyed) return;
    this.destroyed = true;
    this.cleanupFns.forEach((fn) => fn());
    this.cleanupFns = [];
    if (this.overlay) {
      this.overlay.remove();
      this.overlay = null;
    }
  }
}
