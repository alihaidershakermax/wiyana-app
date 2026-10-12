<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch } from "vue";
import { saveLessonProgress } from "@/api/progress";
import {
  Play,
  Pause,
  RotateCcw,
  RotateCw,
  Minimize,
  Maximize,
  PictureInPicture2,
  Gauge,
  Loader2,
  Settings2,
  Subtitles,
  MonitorSmartphone,
  Check,
  X,
  ShieldCheck,
  ShieldAlert,
  Lock,
} from "@lucide/vue";
import { usePlayerStore } from "@/stores/usePlayerStore";
import { useContentProtection } from "@/composables/useContentProtection";
import { formatClockTime } from "@/utils/format";
import WatermarkOverlay from "@/components/ui/WatermarkOverlay.vue";
import type { LessonVideo } from "@/types/lesson";
import type { ReportResult } from "@/types/protection";

const props = withDefaults(
  defineProps<{
    video: LessonVideo | null;
    lessonId: string;
    poster?: string;
  }>(),
  { poster: "" }
);

const playerStore = usePlayerStore();
const PLAYBACK_RATES = [0.5, 0.75, 1, 1.25, 1.5, 2];

const warningInfo = ref<ReportResult | null>(null);
const violationNotice = ref<string | null>(null);
const bannedInfo = ref<ReportResult | null>(null);
const banRemaining = ref(0);

const banRemainingText = computed(() => {
  const s = Math.max(0, banRemaining.value);
  const h = Math.floor(s / 3600);
  const m = Math.floor((s % 3600) / 60);
  const sec = s % 60;
  const pad = (n: number) => String(n).padStart(2, "0");
  return h > 0 ? `${h}:${pad(m)}:${pad(sec)}` : `${m}:${pad(sec)}`;
});

const { engine } = useContentProtection({
  lessonId: props.lessonId,
  closeOnCapture: false,
  onBanned: (result) => {
    bannedInfo.value = result;
    pausePlayback();
    startBanCountdown(result);
  },
  onWarning: (result) => {
    warningInfo.value = result;
    pausePlayback();
    if (warningTimer) clearTimeout(warningTimer);
    warningTimer = setTimeout(() => {
      engine.closeApp();
    }, 3000);
  },
  onViolation: (_type) => {
    // Any detected screenshot / recording attempt: stop the video and show
    // a clear warning message (no black screen).
    violationNotice.value = "تم رصد محاولة تصوير أو تسجيل للشاشة، لذلك تم إيقاف الفيديو.";
    pausePlayback();
  },
});

function dismissViolationNotice() {
  violationNotice.value = null;
}

const YT_QUALITY_LABELS: Record<string, string> = {
  default: "تلقائي",
  tiny: "144p",
  small: "240p",
  medium: "360p",
  large: "480p",
  hd720: "720p",
  hd1080: "1080p",
  hd1440: "1440p",
  hd2160: "2160p",
  highres: "أعلى جودة",
};

const containerRef = ref<HTMLDivElement | null>(null);
const stageRef = ref<HTMLDivElement | null>(null);
const videoRef = ref<HTMLVideoElement | null>(null);
const youtubeMountRef = ref<HTMLDivElement | null>(null);

const isPlaying = ref(false);
const isLoading = ref(true);
const currentTime = ref(0);
const duration = ref(0);
const playbackRate = ref(1);
const isFullscreen = ref(false);
const isLandscape = ref(false);
const viewportTick = ref(0);
const isPiP = ref(false);
const showControls = ref(true);
const settingsOpen = ref(false);
const ytQuality = ref("default");
const availableQualities = ref<{ value: string; label: string }[]>([]);
const isCaptionsOn = ref(false);
const captionsAvailable = ref(false);
const seekFlash = ref<number>(0);

let youtubePlayer: any = null;
let youtubeReady = false;
let controlsTimer: ReturnType<typeof setTimeout> | null = null;
let saveTimer: ReturnType<typeof setInterval> | null = null;
let ytInterval: ReturnType<typeof setInterval> | null = null;
let flashTimer: ReturnType<typeof setTimeout> | null = null;
let tapTimer: ReturnType<typeof setTimeout> | null = null;
let warningTimer: ReturnType<typeof setTimeout> | null = null;
let banTimer: ReturnType<typeof setInterval> | null = null;
let lastTapTime = 0;
let pendingSingleTap = false;

const isTouchDevice =
  typeof window !== "undefined" && ("ontouchstart" in window || (navigator.maxTouchPoints ?? 0) > 0);

const isYouTube = computed(() => props.video?.kind === "youtube");
const supportsPiP = computed(
  () =>
    !isYouTube.value &&
    typeof document !== "undefined" &&
    "pictureInPictureEnabled" in document &&
    document.pictureInPictureEnabled &&
    videoRef.value !== null
);
const progressPercent = computed(() =>
  !duration.value ? 0 : Math.min(100, Math.max(0, (currentTime.value / duration.value) * 100))
);
const currentTimeText = computed(() => formatClockTime(currentTime.value));
const durationText = computed(() => formatClockTime(duration.value));

// ── Fullscreen & Landscape helpers ──────────────────────────────────────────
// "Immersive mode": either real Fullscreen API or the CSS landscape overlay.
const playerMode = computed(() => isFullscreen.value || isLandscape.value);

function viewportDims() {
  if (typeof window === "undefined") return { w: 0, h: 0 };
  const vv = (window as any).visualViewport;
  return {
    w: Math.round(vv?.width || window.innerWidth || 0),
    h: Math.round(vv?.height || window.innerHeight || 0),
  };
}

const viewport = computed(() => {
  void viewportTick.value;
  return viewportDims();
});

const isPortraitView = computed(() => viewport.value.w < viewport.value.h);

// In landscape mode the player container itself becomes a fixed layer that
// covers the exact visible viewport (visualViewport-aware), so no app chrome
// or empty space is visible around the video.
const landscapeContainerStyle = computed<Record<string, string> | undefined>(() => {
  if (!isLandscape.value) return undefined;
  const { w, h } = viewport.value;
  return {
    position: "fixed",
    top: "0",
    left: "0",
    width: `${w}px`,
    height: `${h}px`,
    zIndex: "9999",
    background: "#000",
  };
});

// The stage (video + controls) is rotated 90° when the device is portrait so
// the video fills the whole screen in landscape without cropping or empty
// space. On an already-landscape viewport no rotation is needed.
const landscapeStageStyle = computed<Record<string, string> | undefined>(() => {
  if (!isLandscape.value) return undefined;
  const { w, h } = viewport.value;
  if (!isPortraitView.value) {
    return { top: "0", left: "0", width: "100%", height: "100%", transform: "none", transformOrigin: "center" };
  }
  return {
    top: "50%",
    left: "50%",
    width: `${h}px`,
    height: `${w}px`,
    transform: "translate(-50%, -50%) rotate(90deg)",
    transformOrigin: "center",
  };
});

// ── YouTube helpers ────────────────────────────────────────────────────────
function extractYouTubeId(url: string): string | null {
  const patterns = [
    /(?:youtube\.com\/watch\?(?:.*&)?v=)([\w-]{11})/,
    /(?:youtube\.com\/embed\/)([\w-]{11})/,
    /(?:youtu\.be\/)([\w-]{11})/,
    /(?:youtube\.com\/shorts\/)([\w-]{11})/,
  ];
  for (const p of patterns) {
    const m = url.match(p);
    if (m) return m[1];
  }
  return null;
}

function loadYouTubeApi(): Promise<void> {
  return new Promise((resolve) => {
    if (window.YT?.Player) {
      resolve();
      return;
    }
    const existing = document.getElementById("youtube-iframe-api");
    if (existing) {
      const check = () => {
        if (window.YT?.Player) resolve();
        else setTimeout(check, 150);
      };
      check();
      return;
    }
    const script = document.createElement("script");
    script.id = "youtube-iframe-api";
    script.src = "https://www.youtube.com/iframe_api";
    window.onYouTubeIframeAPIReady = () => resolve();
    document.head.appendChild(script);
  });
}

async function initYouTube() {
  const videoId = props.video?.url ? extractYouTubeId(props.video.url) : null;
  if (!videoId || !youtubeMountRef.value) return;
  await loadYouTubeApi();
  const YT = window.YT;
  if (!YT?.Player) return;

  youtubePlayer = new YT.Player(youtubeMountRef.value, {
    videoId,
    playerVars: {
      playsinline: 1,
      rel: 0,
      controls: 0,
      disablekb: 1,
      modestbranding: 1,
      showinfo: 0,
      iv_load_policy: 3,
      fs: 0,
      origin: window.location.origin,
      widget_referrer: window.location.origin,
    },
    events: {
      onReady: () => {
        youtubeReady = true;
        isLoading.value = false;
        const d = youtubePlayer.getDuration?.() || 0;
        if (d > 0) duration.value = d;
        const resume = playerStore.getResume(props.lessonId);
        if (resume > 0 && resume < d - 5) youtubePlayer.seekTo(resume, true);
        syncYtQualities();
        try {
          youtubePlayer.loadModule?.("captions");
        } catch {}
      },
      onStateChange: (event: any) => {
        const code = event?.data;
        if (code === 1) {
          isPlaying.value = true;
          engine.setPlaying(true);
          startYtPolling();
          wakeControls();
        } else if (code === 2 || code === 0) {
          isPlaying.value = false;
          engine.setPlaying(false);
          stopYtPolling();
          wakeControls();
          if (code === 0) playerStore.markCompleted(props.lessonId);
        }
        if (youtubeReady && (code === 1 || code === 2)) {
          const d = youtubePlayer.getDuration?.() || duration.value;
          if (d > 0) duration.value = d;
        }
      },
      onPlaybackQualityChange: (event: any) => {
        const q = event?.data;
        if (q) {
          ytQuality.value = q;
          syncYtQualities();
        }
      },
      onApiChange: () => syncYtCaptions(),
    },
  });
}

function syncYtQualities() {
  if (!youtubeReady || !youtubePlayer?.getAvailableQualityLevels) return;
  try {
    const levels: string[] = youtubePlayer.getAvailableQualityLevels() || [];
    availableQualities.value = levels.map((v) => ({
      value: v,
      label: YT_QUALITY_LABELS[v] || v,
    }));
  } catch {
    // Keep the last known list.
  }
}

function syncYtCaptions() {
  if (!youtubePlayer) return;
  try {
    const tracks = youtubePlayer.getOption?.("captions", "tracklist");
    captionsAvailable.value = Array.isArray(tracks) && tracks.length > 0;
    const cur = youtubePlayer.getOption?.("captions", "track");
    isCaptionsOn.value = Boolean(cur && cur.languageCode);
  } catch {
    captionsAvailable.value = false;
  }
}

function startYtPolling() {
  stopYtPolling();
  ytInterval = setInterval(() => {
    if (!youtubeReady || !youtubePlayer) return;
    const t = youtubePlayer.getCurrentTime?.() || 0;
    currentTime.value = t;
    const d = youtubePlayer.getDuration?.() || duration.value;
    if (d > 0) duration.value = d;
    const resume = playerStore.getResume(props.lessonId);
    if (t - resume >= 5) {
      playerStore.setResume(props.lessonId, t);
      const pct = d > 0 ? Math.min(100, Math.round((t / d) * 100)) : 0;
      saveLessonProgress(props.lessonId, {
        progress: pct,
        watched_seconds: Math.round(t),
        last_position: t,
        duration: Math.round(d),
      }).catch(() => {});
    }
    if (d > 0 && t >= d - 10) playerStore.markCompleted(props.lessonId);
  }, 1000);
}

function stopYtPolling() {
  if (ytInterval) {
    clearInterval(ytInterval);
    ytInterval = null;
  }
}

function syncResumeTimer() {
  if (saveTimer) return;
  saveTimer = setInterval(async () => {
    if (isYouTube.value) return;
    if (videoRef.value && !videoRef.value.paused && duration.value > 0) {
      const pos = currentTime.value;
      const pct = Math.min(100, Math.round((pos / duration.value) * 100));
      playerStore.setResume(props.lessonId, pos);
      saveLessonProgress(props.lessonId, {
        progress: pct,
        watched_seconds: Math.round(pos),
        last_position: pos,
        duration: Math.round(duration.value),
      }).catch(() => {});
    }
  }, 5000);
}

function stopResumeTimer() {
  if (saveTimer) {
    clearInterval(saveTimer);
    saveTimer = null;
  }
}

function scheduleControlHide(delay = 2500) {
  if (controlsTimer) clearTimeout(controlsTimer);
  if (!isFsActive() && !isPlaying.value) return;
  controlsTimer = setTimeout(() => {
    showControls.value = false;
  }, delay);
}

function wakeControls() {
  showControls.value = true;
  scheduleControlHide();
}

function flashSeek(sec: number) {
  seekFlash.value = sec;
  if (flashTimer) clearTimeout(flashTimer);
  flashTimer = setTimeout(() => {
    seekFlash.value = 0;
  }, 600);
}

function handleSurfaceTap(event: PointerEvent) {
  const revealOnly = !showControls.value;
  wakeControls();
  if (event.target && (event.target as HTMLElement).closest("button, input, a")) return;

  if (settingsOpen.value) {
    settingsOpen.value = false;
    return;
  }

  const now = Date.now();
  const isDoubleTap = now - lastTapTime < 300;
  lastTapTime = now;

  if (isDoubleTap) {
    pendingSingleTap = false;
    if (tapTimer) {
      clearTimeout(tapTimer);
      tapTimer = null;
    }
    const rect = containerRef.value?.getBoundingClientRect();
    const half = (rect?.width || 0) / 2;
    const delta = event.clientX < half ? -10 : 10;
    skip(delta);
    flashSeek(delta);
    return;
  }

  if (revealOnly) {
    pendingSingleTap = false;
    return;
  }

  pendingSingleTap = true;
  if (tapTimer) clearTimeout(tapTimer);
  tapTimer = setTimeout(() => {
    tapTimer = null;
    if (pendingSingleTap) {
      pendingSingleTap = false;
      togglePlay();
    }
  }, 300);
}

// ── settings ───────────────────────────────────────────────────────────────
function toggleSettings() {
  settingsOpen.value = !settingsOpen.value;
}

function setYtQuality(val: string) {
  ytQuality.value = val;
  settingsOpen.value = false;
  if (isYouTube.value && youtubePlayer?.setPlaybackQuality) {
    try {
      youtubePlayer.setPlaybackQuality(val);
    } catch {}
  }
}

function setPlaybackRate(rate: number) {
  playbackRate.value = rate;
  settingsOpen.value = false;
  if (isYouTube.value) {
    if (youtubePlayer?.setPlaybackRate) {
      try {
        youtubePlayer.setPlaybackRate(rate);
      } catch {}
    }
  } else if (videoRef.value) {
    videoRef.value.playbackRate = rate;
  }
}

function toggleCaptions() {
  if (isYouTube.value) {
    if (!youtubePlayer) return;
    try {
      const cur = youtubePlayer.getOption?.("captions", "track");
      if (cur && cur.languageCode) {
        youtubePlayer.setOption("captions", "track", {});
        isCaptionsOn.value = false;
      } else {
        youtubePlayer.setOption("captions", "track", { languageCode: "ar", reload: true });
        isCaptionsOn.value = true;
      }
    } catch {}
    return;
  }
  const v = videoRef.value;
  if (!v?.textTracks?.length) return;
  const anyOn = Array.from(v.textTracks).some((t) => t.mode === "showing");
  for (const t of Array.from(v.textTracks)) t.mode = anyOn ? "hidden" : "showing";
  isCaptionsOn.value = !anyOn;
}

// ── controls ───────────────────────────────────────────────────────────────
function togglePlay() {
  if (isYouTube.value) {
    if (!youtubePlayer) return;
    if (isPlaying.value) youtubePlayer.pauseVideo();
    else youtubePlayer.playVideo();
    return;
  }
  const v = videoRef.value;
  if (!v) return;
  if (v.paused) v.play();
  else v.pause();
}

// ── content protection ────────────────────────────────────────────────────
function pausePlayback() {
  isPlaying.value = false;
  engine.setPlaying(false);
  if (isYouTube.value) {
    youtubePlayer?.pauseVideo();
  } else {
    videoRef.value?.pause();
  }
}

function startBanCountdown(info: ReportResult) {
  const until = Number(info.ban_until || 0);
  if (!until) return;
  const update = () => {
    banRemaining.value = Math.max(0, Math.ceil(until - Date.now() / 1000));
  };
  update();
  if (banTimer) clearInterval(banTimer);
  banTimer = setInterval(update, 1000);
}

async function recheckAccess() {
  try {
    const status = await engine.checkAccessNow();
    if (!status.banned && status.allowed) {
      bannedInfo.value = null;
      if (banTimer) {
        clearInterval(banTimer);
        banTimer = null;
      }
    }
  } catch {
    /* keep the ban overlay */
  }
}

function onProtectionFullscreenChange() {
  if (isRealFullscreen()) {
    isFullscreen.value = true;
  } else {
    isFullscreen.value = false;
    if (isLandscape.value) setLandscapeActive(false);
    // Leaving forced fullscreen mid-playback is a violation per spec.
    if (isPlaying.value && !bannedInfo.value) {
      engine.handleFullscreenExit();
    }
  }
}

function onProtectionVisibilityReturn() {
  if (document.visibilityState === "visible") engine.requestFullscreen();
}

function seekTo(sec: number) {
  if (isYouTube.value) {
    if (youtubePlayer) youtubePlayer.seekTo(sec, true);
    return;
  }
  if (videoRef.value) videoRef.value.currentTime = sec;
}
function skip(sec: number) {
  const t = Math.min(Math.max(0, currentTime.value + sec), duration.value || 0);
  seekTo(t);
  currentTime.value = t;
}

function onSeekInput(event: Event) {
  const val = Number((event.target as HTMLInputElement).value);
  currentTime.value = val;
  if (isYouTube.value) {
    if (youtubePlayer) youtubePlayer.seekTo(val, true);
  } else if (videoRef.value) videoRef.value.currentTime = val;
}

function cyclePlaybackRate() {
  const next = PLAYBACK_RATES[(PLAYBACK_RATES.indexOf(playbackRate.value) + 1) % PLAYBACK_RATES.length];
  setPlaybackRate(next);
}

// ── Fullscreen and Horizontal Rotation (YouTube-style) ─────────────────────
function isFsActive() {
  return playerMode.value;
}

function syncBodyClass() {
  if (typeof document !== "undefined") {
    document.body.classList.toggle("wiyana-player-fullscreen", isLandscape.value);
  }
}

function setLandscapeActive(active: boolean) {
  isLandscape.value = active;
  syncBodyClass();
}

function refreshViewport() {
  viewportTick.value++;
}

async function enterFullscreen() {
  const el = containerRef.value as any;
  if (!el) return false;
  try {
    if (el.requestFullscreen) {
      await el.requestFullscreen();
      return true;
    }
    if (el.webkitRequestFullscreen) {
      await el.webkitRequestFullscreen();
      return true;
    }
    if (el.mozRequestFullScreen) {
      await el.mozRequestFullScreen();
      return true;
    }
    if (el.msRequestFullscreen) {
      await el.msRequestFullscreen();
      return true;
    }
    if (videoRef.value && (videoRef.value as any).webkitEnterFullscreen) {
      (videoRef.value as any).webkitEnterFullscreen();
      return true;
    }
  } catch {
    return false;
  }
  return false;
}

async function exitFullscreen() {
  const doc = document as any;
  try {
    if (doc.exitFullscreen) await doc.exitFullscreen();
    else if (doc.webkitExitFullscreen) await doc.webkitExitFullscreen();
    else if (doc.mozCancelFullScreen) await doc.mozCancelFullScreen();
    else if (doc.msExitFullscreen) await doc.msExitFullscreen();
  } catch {}
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

async function enterLandscape() {
  settingsOpen.value = false;
  refreshViewport();
  setLandscapeActive(true);
  isFullscreen.value = await enterFullscreen();
  showControls.value = true;
  scheduleControlHide(3000);
}

async function exitLandscape() {
  settingsOpen.value = false;
  setLandscapeActive(false);
  isFullscreen.value = false;
  showControls.value = true;
  scheduleControlHide();
  if (isRealFullscreen()) await exitFullscreen();
}

function togglePlayerFullscreen() {
  if (isFsActive()) exitLandscape();
  else enterLandscape();
}

function toggleFullscreen() {
  togglePlayerFullscreen();
}

function toggleLandscape() {
  togglePlayerFullscreen();
}

// Physical rotation (mobile): entering landscape auto-enters the immersive
// mode, coming back to portrait restores the normal page (YouTube behavior).
function isLandscapeOrientation() {
  const so = (screen as any)?.orientation;
  if (so?.type) return so.type.startsWith("landscape");
  if (typeof window !== "undefined" && typeof window.orientation === "number") {
    return Math.abs(window.orientation) === 90;
  }
  return !isPortraitView.value;
}

function onDeviceOrientationChange() {
  refreshViewport();
  if (!isTouchDevice) return;
  const landscapeNow = isLandscapeOrientation();
  if (landscapeNow && !isLandscape.value && props.video) {
    enterLandscape();
  } else if (!landscapeNow && isLandscape.value) {
    exitLandscape();
  }
}

async function togglePiP() {
  const video = videoRef.value;
  if (!video || !document.pictureInPictureEnabled) return;
  try {
    if (document.pictureInPictureElement) await document.exitPictureInPicture();
    else await video.requestPictureInPicture();
  } catch {
    /* not available */
  }
}

// ── video events ───────────────────────────────────────────────────────────
function onVideoLoadedMetadata() {
  const v = videoRef.value;
  if (!v) return;
  if (Number.isFinite(v.duration)) duration.value = v.duration;
  const resume = playerStore.getResume(props.lessonId);
  if (resume > 0 && resume < duration.value - 5) v.currentTime = resume;
  isLoading.value = false;
}
function onVideoTimeUpdate() {
  const v = videoRef.value;
  if (!v) return;
  currentTime.value = v.currentTime;
  if (duration.value > 0 && v.currentTime >= duration.value - 10) playerStore.markCompleted(props.lessonId);
}
function onVideoEnded() {
  isPlaying.value = false;
  engine.setPlaying(false);
  playerStore.markCompleted(props.lessonId);
}

function saveFinalPosition() {
  if (isYouTube.value) {
    if (youtubeReady && youtubePlayer) playerStore.setResume(props.lessonId, youtubePlayer.getCurrentTime?.() || 0);
    return;
  }
  if (videoRef.value) playerStore.setResume(props.lessonId, videoRef.value.currentTime);
}

function onPiPChange() {
  isPiP.value = !!document.pictureInPictureElement;
}

function onDocumentKeydown(event: KeyboardEvent) {
  if (event.code === "Space") {
    const t = event.target as HTMLElement;
    if (t && (t.tagName === "INPUT" || t.tagName === "TEXTAREA" || t.isContentEditable)) return;
    event.preventDefault();
    togglePlay();
  }
}

onMounted(() => {
  if (isYouTube.value) initYouTube();
  else syncResumeTimer();
  engine.requestFullscreen();
  document.addEventListener("fullscreenchange", onProtectionFullscreenChange);
  document.addEventListener("webkitfullscreenchange", onProtectionFullscreenChange);
  document.addEventListener("visibilitychange", onProtectionVisibilityReturn);
  document.addEventListener("enterpictureinpicture", onPiPChange);
  document.addEventListener("leavepictureinpicture", onPiPChange);
  document.addEventListener("keydown", onDocumentKeydown);
  window.addEventListener("orientationchange", onDeviceOrientationChange);
  window.addEventListener("resize", refreshViewport);
  window.matchMedia?.("(orientation: portrait)")?.addEventListener?.("change", onDeviceOrientationChange);
  window.visualViewport?.addEventListener?.("resize", refreshViewport);
});

onBeforeUnmount(() => {
  setLandscapeActive(false);
  saveFinalPosition();
  stopResumeTimer();
  stopYtPolling();
  if (flashTimer) clearTimeout(flashTimer);
  if (tapTimer) clearTimeout(tapTimer);
  if (warningTimer) clearTimeout(warningTimer);
  if (banTimer) clearInterval(banTimer);
  if (youtubePlayer && typeof youtubePlayer.destroy === "function") youtubePlayer.destroy();
  document.removeEventListener("fullscreenchange", onProtectionFullscreenChange);
  document.removeEventListener("webkitfullscreenchange", onProtectionFullscreenChange);
  document.removeEventListener("visibilitychange", onProtectionVisibilityReturn);
  document.removeEventListener("enterpictureinpicture", onPiPChange);
  document.removeEventListener("leavepictureinpicture", onPiPChange);
  document.removeEventListener("keydown", onDocumentKeydown);
  window.removeEventListener("orientationchange", onDeviceOrientationChange);
  window.removeEventListener("resize", refreshViewport);
  window.matchMedia?.("(orientation: portrait)")?.removeEventListener?.("change", onDeviceOrientationChange);
  window.visualViewport?.removeEventListener?.("resize", refreshViewport);
});

watch(
  () => props.lessonId,
  () => {
    setLandscapeActive(false);
    isFullscreen.value = false;
    currentTime.value = 0;
    duration.value = 0;
    isPlaying.value = false;
    isLoading.value = true;
    settingsOpen.value = false;
    showControls.value = true;
    seekFlash.value = 0;
    availableQualities.value = [];
    ytQuality.value = "default";
    isCaptionsOn.value = false;
    captionsAvailable.value = false;
    if (isYouTube.value) {
      if (youtubePlayer) {
        stopYtPolling();
        youtubePlayer.destroy();
        youtubePlayer = null;
      }
      youtubeReady = false;
      initYouTube();
    } else {
      syncResumeTimer();
    }
  }
);
</script>

<template>
  <div
    ref="containerRef"
    dir="ltr"
    :class="[
      'bg-black overflow-hidden group select-none',
      isLandscape ? '' : 'relative w-full aspect-video transition-all',
      isFsActive() ? 'is-fullscreen' : '',
    ]"
    :style="landscapeContainerStyle"
    @pointermove="wakeControls"
    @pointerleave="() => scheduleControlHide()"
    @contextmenu.prevent
  >
    <!-- Rotated stage: the whole video surface (media + overlays + controls)
         rotates 90° into landscape inside fullscreen (YouTube style). -->
    <div
      ref="stageRef"
      :class="isLandscape ? 'absolute' : 'relative w-full h-full'"
      :style="landscapeStageStyle"
    >
    <!-- Native HTML5 Video for direct uploaded files -->
    <video
      v-if="!isYouTube && video"
      ref="videoRef"
      class="w-full h-full object-contain select-none"
      :poster="poster || undefined"
      :src="video.url"
      preload="metadata"
      playsinline
      draggable="false"
      :controls="false"
      :controlslist="'nodownload noremoteplayback'"
      disablepictureinpicture
      @dragstart.prevent
      @contextmenu.prevent
      @pointerdown="handleSurfaceTap"
      @loadedmetadata="onVideoLoadedMetadata"
      @timeupdate="onVideoTimeUpdate"
      @ended="onVideoEnded"
      @play="isPlaying = true; engine.setPlaying(true); wakeControls()"
      @pause="isPlaying = false; engine.setPlaying(false); wakeControls(); saveFinalPosition()"
      @waiting="isLoading = true"
      @playing="isLoading = false; syncResumeTimer()"
      @error="isLoading = false"
    ></video>

    <!-- Clean YouTube Embed: Full Video Frame Visible with Masked Channel Logo & Links -->
    <div v-else-if="isYouTube" class="w-full h-full relative overflow-hidden bg-black select-none" @contextmenu.prevent>
      <div class="w-full h-full relative">
        <div ref="youtubeMountRef" class="w-full h-full pointer-events-none"></div>
      </div>

      <!-- Transparent click interceptor: Toggles play/pause & autohides controls without opening YouTube -->
      <div
        @pointerdown="handleSurfaceTap"
        class="absolute inset-0 z-10 cursor-pointer"
        aria-label="تشغيل أو إيقاف"
      ></div>

      <!-- Total protective masks: hide YouTube channel title, share button, and channel logo -->
      <div class="absolute top-0 inset-x-0 h-14 bg-gradient-to-b from-black via-black/70 to-transparent pointer-events-none z-15"></div>
      <div class="absolute top-0 left-0 w-64 h-16 bg-transparent z-15 pointer-events-auto" @pointerdown.stop="handleSurfaceTap"></div>
      <div class="absolute top-0 right-0 w-64 h-16 bg-transparent z-15 pointer-events-auto" @pointerdown.stop="handleSurfaceTap"></div>
      <div class="absolute bottom-0 right-0 w-48 h-14 bg-transparent z-15 pointer-events-auto" @pointerdown.stop="handleSurfaceTap"></div>
    </div>

    <!-- No Video State -->
    <div v-else class="w-full h-full flex flex-col items-center justify-center gap-3 bg-slate-900 text-slate-400">
      <Play class="w-12 h-12 text-slate-500" />
      <span class="text-xs font-bold" dir="rtl">لا يوجد فيديو متاح لهذه المحاضرة</span>
    </div>

    <!-- Dynamic Watermark Protection & Anti-Recording Layer (Must be AFTER video and iframe elements!) -->
    <WatermarkOverlay />

    <!-- Permanent protected badge (always visible, independent of auto-hiding controls) -->
    <div
      v-if="video && !bannedInfo"
      class="absolute top-2 left-2 z-40 inline-flex items-center gap-1 px-2 py-1 rounded-lg bg-black/60 backdrop-blur-sm border border-emerald-500/40 pointer-events-none select-none"
      dir="rtl"
      title="المحتوى محمي ضد التصوير والتسجيل"
    >
      <ShieldCheck class="w-3.5 h-3.5 text-emerald-400" />
      <span class="text-[10px] font-black text-emerald-300">محمي</span>
    </div>

    <!-- Loading spinner -->
    <div v-if="isLoading && video" class="absolute inset-0 flex items-center justify-center bg-black/50 pointer-events-none z-20">
      <Loader2 class="w-12 h-12 text-orange-500 animate-spin" />
    </div>

    <!-- Center play button overlay when paused -->
    <button
      v-if="video && !isPlaying && !isLoading && !bannedInfo && (!isFsActive() || showControls)"
      @click="togglePlay"
      type="button"
      aria-label="تشغيل المحاضرة"
      class="absolute inset-0 flex items-center justify-center bg-black/30 z-20 transition-opacity cursor-pointer"
    >
      <span class="w-16 h-16 sm:w-20 sm:h-20 rounded-full bg-orange-500 text-white flex items-center justify-center shadow-2xl shadow-orange-500/50 scale-in hover:scale-110 active:scale-95 transition-all">
        <Play class="w-8 h-8 sm:w-10 sm:h-10 fill-white ml-1" />
      </span>
    </button>

    <!-- Double-tap seek flash -->
    <Transition name="seek-flash">
      <div
        v-if="seekFlash !== 0 && video"
        class="absolute inset-0 flex items-center justify-center pointer-events-none z-20"
      >
        <div
          class="px-4 py-2.5 rounded-2xl bg-black/70 backdrop-blur-sm border border-white/15 text-white text-sm font-black tabular-nums inline-flex items-center gap-1.5"
          dir="ltr"
        >
          <template v-if="seekFlash < 0">
            <RotateCcw class="w-4 h-4 text-orange-400" />
            <span>10s</span>
          </template>
          <template v-else>
            <span>10s</span>
            <RotateCw class="w-4 h-4 text-orange-400" />
          </template>
        </div>
      </div>
    </Transition>

    <!-- Custom Controls Bar -->
    <div
      v-if="video && !bannedInfo"
      :class="[
        'absolute inset-x-0 bottom-0 z-30 bg-gradient-to-t from-black/95 via-black/60 to-transparent transition-all duration-300',
        showControls ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-2 pointer-events-none',
        isFsActive() ? 'px-5 pb-4 pt-16' : 'px-4 pb-3 pt-12'
      ]"
    >
      <!-- Timeline scrubber range (LTR) -->
      <div class="w-full flex items-center gap-3" dir="ltr">
        <span class="text-xs font-black text-white tabular-nums shrink-0">{{ currentTimeText }}</span>
        <input
          type="range"
          min="0"
          :max="duration || 0"
          step="1"
          :value="currentTime"
          @input="onSeekInput"
          aria-label="شريط التقدم"
          class="player-range flex-1"
        />
        <span class="text-xs font-bold text-white/70 tabular-nums shrink-0">{{ durationText }}</span>
      </div>

      <!-- Controls row (RTL order) -->
      <div class="flex items-center gap-1.5 mt-2" dir="rtl">
        <!-- Right side: Play/Pause, Rewind, Fast-forward, Speed -->
        <button @click="togglePlay" type="button" aria-label="تشغيل أو إيقاف" class="player-btn cursor-pointer">
          <Pause v-if="isPlaying" class="w-5 h-5 fill-white" />
          <Play v-else class="w-5 h-5 fill-white mr-0.5" />
        </button>

        <button @click="skip(-10)" type="button" aria-label="تأخير 10 ثوانٍ" class="player-btn cursor-pointer" title="تأخير 10 ثوانٍ">
          <RotateCcw class="w-4.5 h-4.5" />
        </button>

        <button @click="skip(10)" type="button" aria-label="تقديم 10 ثوانٍ" class="player-btn cursor-pointer" title="تقديم 10 ثوانٍ">
          <RotateCw class="w-4.5 h-4.5" />
        </button>

        <button @click="cyclePlaybackRate" type="button" aria-label="سرعة التشغيل" class="player-btn gap-1 px-2.5 cursor-pointer" style="width:auto;" title="سرعة التشغيل">
          <Gauge class="w-4 h-4 shrink-0 text-orange-400" />
          <span class="text-[11px] font-black tabular-nums">{{ playbackRate }}x</span>
        </button>

        <!-- Spacer -->
        <div class="flex-1"></div>

        <!-- Protected badge -->
        <span class="inline-flex items-center gap-1 px-2 py-1 rounded-lg bg-white/10 text-[10px] font-black text-white/70" dir="rtl">
          <ShieldCheck class="w-3.5 h-3.5 text-emerald-400" />
          محمي
        </span>

        <!-- Settings -->
        <button
          @click="toggleSettings"
          type="button"
          aria-label="الإعدادات"
          :class="['player-btn cursor-pointer', settingsOpen ? 'text-orange-400 bg-white/10' : '']"
          title="الإعدادات"
        >
          <Settings2 class="w-4.5 h-4.5" />
        </button>

        <button
          v-if="supportsPiP"
          @click="togglePiP"
          type="button"
          :aria-label="isPiP ? 'خروج من PiP' : 'صورة داخل صورة'"
          class="player-btn cursor-pointer hidden min-[480px]:inline-flex"
          title="صورة داخل صورة"
        >
          <PictureInPicture2 :class="['w-4.5 h-4.5', isPiP ? 'text-orange-400' : '']" />
        </button>

        <!-- Fullscreen -->
        <button
          @click="toggleFullscreen"
          type="button"
          :aria-label="isFsActive() ? 'إنهاء ملء الشاشة' : 'ملء الشاشة'"
          class="player-btn cursor-pointer"
          title="ملء الشاشة"
        >
          <Minimize v-if="isFsActive()" class="w-4.5 h-4.5 text-orange-400" />
          <Maximize v-else class="w-4.5 h-4.5" />
        </button>

        <!-- Landscape (fullscreen + rotate to landscape) -->
        <button
          @click="toggleLandscape"
          type="button"
          aria-label="شاشة أفقية"
          class="player-btn cursor-pointer bg-white/10 hover:bg-orange-500 hover:text-white transition-colors"
          title="قلب الفيديو أفقياً"
        >
          <MonitorSmartphone class="w-4.5 h-4.5" />
        </button>
      </div>
    </div>
    </div>
    <!-- Mobile settings backdrop -->
    <Transition name="fade-fast">
      <div
        v-if="settingsOpen"
        class="sm:hidden fixed inset-0 bg-black/60 backdrop-blur-[2px] z-40"
        @pointerdown="settingsOpen = false"
      ></div>
    </Transition>

    <!-- Settings sheet: bottom sheet on mobile, anchored panel on desktop -->
    <Transition name="sheet">
      <div
        v-if="settingsOpen && video"
        dir="rtl"
        class="fixed inset-x-0 bottom-0 z-50 max-h-[75vh] overflow-y-auto bg-[#1a1a1a]/95 backdrop-blur-xl border-t border-white/10 rounded-t-2xl shadow-2xl sm:absolute sm:inset-x-auto sm:bottom-16 sm:right-3 sm:left-auto sm:rounded-2xl sm:border sm:max-h-[70vh] sm:w-72 sm:overflow-hidden"
        @pointerdown.stop
      >
        <!-- Sheet header -->
        <div class="flex items-center justify-between px-4 py-3 border-b border-white/10 sticky top-0 bg-[#1a1a1a]/95 backdrop-blur-xl z-10">
          <span class="text-sm font-black text-white inline-flex items-center gap-2">
            <Settings2 class="w-4 h-4 text-orange-400" />
            إعدادات التشغيل
          </span>
          <button
            @click="settingsOpen = false"
            type="button"
            aria-label="إغلاق الإعدادات"
            class="w-8 h-8 rounded-xl bg-white/10 text-white flex items-center justify-center hover:bg-white/20 transition-colors cursor-pointer"
          >
            <X class="w-4 h-4" />
          </button>
        </div>

        <div class="p-2.5 flex flex-col gap-4">
          <!-- Quality -->
          <div class="flex flex-col gap-1">
            <span class="px-2 text-[10px] font-bold text-white/40 tracking-wide">الجودة</span>
            <template v-if="isYouTube">
              <button
                @click="setYtQuality('default')"
                type="button"
                :class="[
                  'w-full flex items-center justify-between px-3 py-2 text-xs font-bold rounded-xl transition-colors cursor-pointer',
                  ytQuality === 'default'
                    ? 'text-orange-400 bg-orange-500/20'
                    : 'text-white/80 hover:bg-white/10'
                ]"
              >
                <span>تلقائي</span>
                <Check v-if="ytQuality === 'default'" class="w-4 h-4 text-orange-400" />
              </button>
              <button
                v-for="q in availableQualities"
                :key="q.value"
                @click="setYtQuality(q.value)"
                type="button"
                :class="[
                  'w-full flex items-center justify-between px-3 py-2 text-xs font-bold rounded-xl transition-colors cursor-pointer',
                  ytQuality === q.value && q.value !== 'default'
                    ? 'text-orange-400 bg-orange-500/20'
                    : 'text-white/80 hover:bg-white/10'
                ]"
              >
                <span>{{ q.label }}</span>
                <Check v-if="ytQuality === q.value" class="w-4 h-4 text-orange-400" />
              </button>
              <div v-if="availableQualities.length === 0" class="px-3 py-2 text-xs font-bold text-white/40">
                {{ isLoading ? 'جارٍ تحميل الجودة…' : 'جودة تلقائية' }}
              </div>
            </template>
            <div v-else class="px-3 py-2 text-xs font-bold text-white/50">جودة تلقائية</div>
          </div>

          <!-- Speed -->
          <div class="flex flex-col gap-1">
            <span class="px-2 text-[10px] font-bold text-white/40 tracking-wide">سرعة التشغيل</span>
            <button
              v-for="rate in PLAYBACK_RATES"
              :key="rate"
              @click="setPlaybackRate(rate)"
              type="button"
              :class="[
                'w-full flex items-center justify-between px-3 py-2 text-xs font-bold rounded-xl transition-colors cursor-pointer',
                playbackRate === rate
                  ? 'text-orange-400 bg-orange-500/20'
                  : 'text-white/80 hover:bg-white/10'
              ]"
            >
              <span>{{ rate }}x</span>
              <Check v-if="playbackRate === rate" class="w-4 h-4 text-orange-400" />
            </button>
          </div>

          <!-- Captions -->
          <div
            v-if="!isYouTube || captionsAvailable"
            class="flex flex-col gap-1"
          >
            <span class="px-2 text-[10px] font-bold text-white/40 tracking-wide">الترجمات</span>
            <button
              @click="toggleCaptions"
              type="button"
              :class="[
                'w-full flex items-center justify-between px-3 py-2 text-xs font-bold rounded-xl transition-colors cursor-pointer',
                isCaptionsOn
                  ? 'text-orange-400 bg-orange-500/20'
                  : 'text-white/80 hover:bg-white/10'
              ]"
            >
              <span class="inline-flex items-center gap-2">
                <Subtitles class="w-4 h-4" />
                {{ isCaptionsOn ? 'الترجمات مفعّلة' : 'الترجمات معطّلة' }}
              </span>
              <Check v-if="isCaptionsOn" class="w-4 h-4 text-orange-400" />
            </button>
          </div>

          <!-- Screen orientation -->
          <div class="flex flex-col gap-1">
            <span class="px-2 text-[10px] font-bold text-white/40 tracking-wide">الشاشة</span>
            <button
              @click="toggleLandscape"
              type="button"
              class="w-full flex items-center justify-between px-3 py-2 text-xs font-bold rounded-xl text-white/80 hover:bg-white/10 transition-colors cursor-pointer"
            >
              <span class="inline-flex items-center gap-2">
                <MonitorSmartphone class="w-4 h-4" />
                شاشة أفقية
              </span>
              <span v-if="isFsActive()" class="text-[10px] font-bold text-orange-400">مفعلة</span>
            </button>
          </div>
        </div>
      </div>
    </Transition>

    <!-- Instant violation warning notice (stops the video, no black screen) -->
    <div
      v-if="violationNotice && !bannedInfo"
      dir="rtl"
      class="absolute inset-0 z-50 flex items-center justify-center bg-black/85 backdrop-blur-sm px-6 text-center"
    >
      <div class="flex flex-col items-center gap-3 max-w-sm">
        <div class="w-14 h-14 rounded-full bg-amber-500/15 flex items-center justify-center">
          <ShieldAlert class="w-7 h-7 text-amber-400" />
        </div>
        <h3 class="text-base font-black text-white">تنبيه: محاولة التقاط محتوى</h3>
        <p class="text-xs font-medium text-white/70 leading-relaxed">
          {{ violationNotice }}
        </p>
        <p class="text-xs font-bold text-white/50">
          تكرار المحاولة يؤدي إلى تسجيل نقاط مخالفة وحظر الحساب تلقائياً.
        </p>
        <button
          type="button"
          @click="dismissViolationNotice"
          class="mt-1 px-5 py-2.5 rounded-xl bg-orange-500 hover:bg-orange-600 active:scale-95 transition-all text-white text-sm font-black cursor-pointer"
        >
          متابعة الدرس
        </button>
      </div>
    </div>

    <!-- Violation warning overlay -->
    <div
      v-if="warningInfo"
      dir="rtl"
      class="absolute inset-0 z-50 flex items-center justify-center bg-black/85 backdrop-blur-sm px-6 text-center"
    >
      <div class="flex flex-col items-center gap-3 max-w-sm">
        <div class="w-14 h-14 rounded-full bg-amber-500/15 flex items-center justify-center">
          <ShieldAlert class="w-7 h-7 text-amber-400" />
        </div>
        <h3 class="text-base font-black text-white">تنبيه: محاولة التقاط محتوى</h3>
        <p class="text-xs font-medium text-white/70 leading-relaxed">
          تم رصد محاولة التقاط شاشة أو تصوير للمحتوى. الحساب يرصد النقاط، وسيُغلق التطبيق الآن، ويُحظر الحساب تلقائياً عند تجاوز الحد المسموح.
        </p>
        <p class="text-xs font-black text-amber-400">
          نقاط المخالفات: {{ warningInfo.points }} / {{ warningInfo.ban_threshold }}
        </p>
        <p class="text-xs font-bold text-white/50">سيُغلق التطبيق خلال ثوانٍ…</p>
      </div>
    </div>

    <!-- Ban overlay -->
    <div
      v-if="bannedInfo"
      dir="rtl"
      class="absolute inset-0 z-50 flex items-center justify-center bg-black/95 px-6 text-center"
    >
      <div class="flex flex-col items-center gap-4 max-w-sm">
        <div class="w-16 h-16 rounded-full bg-red-500/15 flex items-center justify-center">
          <Lock class="w-8 h-8 text-red-500" />
        </div>
        <h3 class="text-lg font-black text-white">تم حظر الحساب مؤقتاً</h3>
        <p class="text-xs font-medium text-white/70 leading-relaxed">
          تم تجاوز الحد المسموح من محاولات التقاط المحتوى، لذلك أُوقف وصولك إلى المحاضرات بشكل مؤقت.
        </p>
        <p v-if="banRemaining > 0" class="text-sm font-black text-red-400 tabular-nums" dir="ltr">
          {{ banRemainingText }}
        </p>
        <div class="flex flex-col gap-2 w-full">
          <button
            @click="recheckAccess"
            type="button"
            class="px-5 py-2.5 rounded-xl bg-orange-500 hover:bg-orange-600 active:scale-95 transition-all text-white text-sm font-black cursor-pointer"
          >
            إعادة التحقق من الوصول
          </button>
          <button
            @click="engine.closeApp()"
            type="button"
            class="px-5 py-2.5 rounded-xl bg-white/10 hover:bg-white/20 active:scale-95 transition-all text-white text-sm font-black cursor-pointer"
          >
            إغلاق
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.player-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 2.25rem;
  height: 2.25rem;
  padding: 0 0.4rem;
  border-radius: 0.75rem;
  color: #ffffff;
  transition: transform 0.12s ease, background-color 0.12s ease;
  flex-shrink: 0;
}
.player-btn:hover {
  background-color: rgb(255 255 255 / 0.18);
}
.player-btn:active {
  transform: scale(0.9);
}

/* Bigger, touch-friendly controls while in fullscreen (YouTube mobile style) */
.is-fullscreen .player-btn {
  min-width: 2.75rem;
  height: 2.75rem;
  border-radius: 0.875rem;
}
.is-fullscreen .player-btn svg {
  width: 1.3rem !important;
  height: 1.3rem !important;
}
.is-fullscreen .player-range {
  height: 6px;
}
.is-fullscreen .player-range::-webkit-slider-thumb {
  width: 18px;
  height: 18px;
}
.is-fullscreen .player-range::-moz-range-thumb {
  width: 18px;
  height: 18px;
}

/* Progress range slider */
.player-range {
  -webkit-appearance: none;
  appearance: none;
  height: 5px;
  border-radius: 9999px;
  background: linear-gradient(
    to right,
    #f97316 0%,
    #f97316 v-bind(progressPercent + "%"),
    rgb(255 255 255 / 0.25) v-bind(progressPercent + "%"),
    rgb(255 255 255 / 0.25) 100%
  );
  cursor: pointer;
  outline: none;
}
.player-range::-webkit-slider-thumb {
  -webkit-appearance: none;
  appearance: none;
  width: 14px;
  height: 14px;
  border-radius: 9999px;
  background: #f97316;
  border: 2px solid #fff;
  box-shadow: 0 2px 6px rgb(0 0 0 / 0.5);
}
.player-range::-moz-range-thumb {
  width: 14px;
  height: 14px;
  border-radius: 9999px;
  background: #f97316;
  border: 2px solid #fff;
}

/* Center play button animation */
@keyframes scaleIn {
  from {
    opacity: 0;
    transform: scale(0.75);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}
.scale-in {
  animation: scaleIn 0.22s cubic-bezier(0.16, 1, 0.3, 1) forwards;
}

/* Settings sheet (bottom sheet on mobile, panel on desktop) */
.sheet-enter-active,
.sheet-leave-active {
  transition: transform 0.25s cubic-bezier(0.16, 1, 0.3, 1), opacity 0.2s ease;
}
.sheet-enter-from,
.sheet-leave-to {
  opacity: 0;
  transform: translateY(16px);
}
@media (min-width: 640px) {
  .sheet-enter-from,
  .sheet-leave-to {
    transform: translateY(8px) scale(0.96);
  }
}

/* Seek flash */
.seek-flash-enter-active {
  transition: opacity 0.15s ease, transform 0.15s ease;
}
.seek-flash-leave-active {
  transition: opacity 0.3s ease, transform 0.3s ease;
}
.seek-flash-enter-from,
.seek-flash-leave-to {
  opacity: 0;
  transform: scale(0.9);
}

/* Fade (backdrop) */
.fade-fast-enter-active,
.fade-fast-leave-active {
  transition: opacity 0.2s ease;
}
.fade-fast-enter-from,
.fade-fast-leave-to {
  opacity: 0;
}
</style>
