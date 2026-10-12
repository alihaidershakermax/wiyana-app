// Local replacement for the `@macrotom/screen-privacy` package, which does
// NOT exist on npm (verified: `npm view @macrotom/screen-privacy` → 404).
//
// It exposes the same surface the original library provides (init/enable/
// disable/isRecording + the `screenRecordingDetected` event) so call sites
// stay identical. Real screen-recording detection on the web is limited:
//  - Android: a native secure layer (FLAG_SECURE) is not reachable from JS,
//    so `isRecording()` always returns `false` there.
//  - iOS: only Safari's active screen-capture can be observed, which this
//    shim reports via the shared `screenRecordingDetected` CustomEvent that
//    the protection engine also feeds its own heuristics into.
export const SCREEN_RECORDING_EVENT = "screenRecordingDetected";

export interface ScreenRecordingEventDetail {
  source: string;
  timestamp: number;
}

export const screenPrivacy = {
  init(): void {
    void this.enable();
    this.interceptScreenCaptureAPIs();
  },

  /** Enable native FLAG_SECURE on Android & PrivacyScreen on iOS + Telegram Fullscreen */
  enable(): boolean {
    if (typeof window === "undefined") return false;
    try {
      const tg = (window as any).Telegram?.WebApp;
      if (tg) {
        try {
          tg.requestFullscreen?.();
          tg.disableVerticalSwipes?.();
          tg.enableClosingConfirmation?.();
        } catch {}
      }

      const cap = (window as any).Capacitor;
      if (cap?.Plugins?.PrivacyScreen?.enable) {
        void cap.Plugins.PrivacyScreen.enable();
      }
      return true;
    } catch {
      return false;
    }
  },

  disable(): void {
    try {
      const cap = (window as any).Capacitor;
      if (cap?.Plugins?.PrivacyScreen?.disable) {
        void cap.Plugins.PrivacyScreen.disable();
      }
    } catch {}
  },

  isRecording(): boolean {
    return false;
  },

  interceptScreenCaptureAPIs(): void {
    if (typeof window === "undefined" || !navigator.mediaDevices) return;

    // Intercept getDisplayMedia to block web/browser screen recorders instantly
    try {
      if (navigator.mediaDevices.getDisplayMedia) {
        navigator.mediaDevices.getDisplayMedia = async function (_options?: DisplayMediaStreamOptions) {
          dispatchRecordingDetected("getDisplayMedia_attempt");
          throw new DOMException("تسجيل الشاشة حظر أمني لمنع السرقة على منصة ويانا.", "NotAllowedError");
        };
      }
    } catch {}

    // Intercept HTMLMediaElement captureStream
    try {
      const proto = HTMLMediaElement.prototype as any;
      if (proto.captureStream) {
        proto.captureStream = function () {
          dispatchRecordingDetected("captureStream_attempt");
          return new MediaStream();
        };
      }
    } catch {}

    // Intercept MediaRecorder: blocks in-page recording of any media stream
    // (e.g. a page script trying to record the <video> element). The
    // constructor throws and any start() call on a surviving instance is
    // neutered, both paths dispatching the recording-detected event.
    try {
      const NativeMediaRecorder = window.MediaRecorder;
      if (NativeMediaRecorder) {
        (window as any).MediaRecorder = function (...args: any[]) {
          dispatchRecordingDetected("mediaRecorder_attempt");
          throw new DOMException(
            "تسجيل الشاشة حظر أمني لمنع السرقة على منصة ويانا.",
            "NotAllowedError"
          );
        };
        (window as any).MediaRecorder.prototype = NativeMediaRecorder.prototype;
        try {
          (window as any).MediaRecorder.isTypeSupported =
            NativeMediaRecorder.isTypeSupported.bind(NativeMediaRecorder);
        } catch {}
      }
    } catch {}
  },
};

export function dispatchRecordingDetected(source: string): void {
  if (typeof window === "undefined") return;
  window.dispatchEvent(
    new CustomEvent<ScreenRecordingEventDetail>(SCREEN_RECORDING_EVENT, {
      detail: { source, timestamp: Date.now() },
    })
  );
  window.dispatchEvent(new CustomEvent("wiyana:protection-lockout"));
}
