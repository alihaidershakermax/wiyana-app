import { ref, computed, onMounted, onBeforeUnmount } from "vue";

export function useAppFullscreen() {
  const isFullscreen = ref(false);

  const isDesktop = computed(
    () =>
      typeof window !== "undefined" &&
      typeof window.matchMedia === "function" &&
      window.matchMedia("(pointer: fine)").matches
  );

  const isTelegramEnv = computed(
    () => typeof window !== "undefined" && !!(window as any).Telegram?.WebApp
  );

  function getTg() {
    return typeof window !== "undefined" ? (window as any).Telegram?.WebApp : null;
  }

  // Telegram Mini App SDK fullscreen (v8+): works inside the Telegram clients
  // (desktop & mobile) where the browser Fullscreen API is blocked.
  function requestTgFullscreen(): Promise<boolean> {
    const tg = getTg();
    if (!tg || typeof tg.requestFullscreen !== "function") {
      return Promise.resolve(false);
    }
    return new Promise((resolve) => {
      let settled = false;
      const finish = (ok: boolean) => {
        if (settled) return;
        settled = true;
        tg.offEvent?.("fullscreen_changed", onChange);
        tg.offEvent?.("fullscreen_failed", onFail);
        resolve(ok);
      };
      const onChange = (e: any) => finish(!!e?.is_fullscreen);
      const onFail = () => finish(false);
      tg.onEvent?.("fullscreen_changed", onChange);
      tg.onEvent?.("fullscreen_failed", onFail);
      try {
        const result = tg.requestFullscreen();
        if (result && typeof result.then === "function") {
          result.then(() => finish(true)).catch(() => finish(false));
        }
      } catch {
        finish(false);
      }
      setTimeout(() => finish(false), 2500);
    });
  }

  function isRealFullscreen() {
    const doc = document as any;
    return !!(
      doc.fullscreenElement ||
      doc.webkitFullscreenElement ||
      doc.mozFullScreenElement ||
      doc.msFullscreenElement
    );
  }

  async function requestBrowserFullscreen() {
    const el = document.documentElement as any;
    try {
      if (el.requestFullscreen) {
        await el.requestFullscreen();
        return isRealFullscreen();
      }
      if (el.webkitRequestFullscreen) {
        el.webkitRequestFullscreen();
        return isRealFullscreen();
      }
      if (el.mozRequestFullScreen) {
        el.mozRequestFullScreen();
        return isRealFullscreen();
      }
      if (el.msRequestFullscreen) {
        el.msRequestFullscreen();
        return isRealFullscreen();
      }
    } catch {
      return false;
    }
    return false;
  }

  async function enterFullscreen() {
    if (await requestTgFullscreen()) return true;
    return requestBrowserFullscreen();
  }

  function onFullscreenChange() {
    isFullscreen.value = isRealFullscreen();
  }

  function onTgFullscreenChanged(e: any) {
    isFullscreen.value = !!e?.is_fullscreen;
  }

  let autoAttempted = false;
  let gestureCleanup: (() => void) | null = null;

  function cleanupGestureListeners() {
    if (gestureCleanup) {
      gestureCleanup();
      gestureCleanup = null;
    }
  }

  function requestAutoFullscreen() {
    if (autoAttempted || !isDesktop.value || !isTelegramEnv.value) return;
    autoAttempted = true;
    if (isFullscreen.value || isRealFullscreen()) return;

    // Try to enter fullscreen right away on load.
    void enterFullscreen().then((entered) => {
      if (entered || isFullscreen.value || isRealFullscreen()) return;

      // Fallback: the platform requires a user gesture, so enter fullscreen on
      // the very first interaction (once per session only).
      const attempt = () => {
        cleanupGestureListeners();
        if (isFullscreen.value || isRealFullscreen()) return;
        void enterFullscreen();
      };
      window.addEventListener("pointerdown", attempt, { once: true, passive: true });
      window.addEventListener("keydown", attempt, { once: true, passive: true });
      gestureCleanup = () => {
        window.removeEventListener("pointerdown", attempt);
        window.removeEventListener("keydown", attempt);
      };
    });
  }

  onMounted(() => {
    const tg = getTg();
    isFullscreen.value = Boolean(tg?.isFullscreen) || isRealFullscreen();
    document.addEventListener("fullscreenchange", onFullscreenChange);
    document.addEventListener("webkitfullscreenchange", onFullscreenChange);
    tg?.onEvent?.("fullscreen_changed", onTgFullscreenChanged);
    requestAutoFullscreen();
  });

  onBeforeUnmount(() => {
    document.removeEventListener("fullscreenchange", onFullscreenChange);
    document.removeEventListener("webkitfullscreenchange", onFullscreenChange);
    getTg()?.offEvent?.("fullscreen_changed", onTgFullscreenChanged);
    cleanupGestureListeners();
  });

  return { isFullscreen, isDesktop, isTelegramEnv };
}
