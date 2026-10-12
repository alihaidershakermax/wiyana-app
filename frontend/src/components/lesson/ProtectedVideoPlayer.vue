<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from "vue";
import Hls from "hls.js";
import {
  Play,
  Pause,
  RotateCcw,
  RotateCw,
  Gauge,
  Loader2,
  ShieldCheck,
  ShieldAlert,
  Lock,
  Maximize,
  Minimize,
} from "@lucide/vue";
import { usePlayerStore } from "@/stores/usePlayerStore";
import { saveLessonProgress } from "@/api/progress";
import { createStreamGrant, deviceFingerprint } from "@/api/protection";
import { resolveApiUrl } from "@/api/client";
import { useContentProtection } from "@/composables/useContentProtection";
import { formatClockTime } from "@/utils/format";
import WatermarkOverlay from "@/components/ui/WatermarkOverlay.vue";
import DynamicWatermark from "@/components/DynamicWatermark.vue";
import SecurityWarningModal from "@/components/SecurityWarningModal.vue";
import SecurityStatus from "@/components/SecurityStatus.vue";
import { useCaptureDetection } from "@/security/useCaptureDetection";
import { useSecurityGuard } from "@/security/useSecurityGuard";
import { PlaybackSessionClient } from "@/video/playbackSession";
import { detectPlatform } from "@/security/platformCapabilities";
import { startScreenShield, stopScreenShield } from "@/native/screenShield";
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

const containerRef = ref<HTMLDivElement | null>(null);
const videoRef = ref<HTMLVideoElement | null>(null);

const isPlaying = ref(false);
const isLoading = ref(true);
const errorMsg = ref("");
const currentTime = ref(0);
const duration = ref(0);
const playbackRate = ref(1);
const showControls = ref(true);
const isFullscreen = ref(false);

const warningInfo = ref<ReportResult | null>(null);
const violationNotice = ref<string | null>(null);
const bannedInfo = ref<ReportResult | null>(null);
const banRemaining = ref(0);

const sessionClient = new PlaybackSessionClient();
const watermarkConfig = ref({
  displayName: "",
  telegramId: "",
  privacyId: "",
  sessionId: "",
});

const {
  isBlocked: isHighRiskBlocked,
  warningModalOpen,
  warningLevel,
  warningTitle,
  warningMessage,
  warningDetail,
  requireRevalidation,
  showSuspiciousNotice,
  triggerHighRiskLockout,
  triggerCriticalLockout,
  showDetectionUnavailable,
  dismissWarning,
} = useSecurityGuard(detectPlatform());

let lastAdvisoryLevel = "NORMAL";

const {
  currentRiskScore,
  currentRiskLevel,
  capabilities,
  handleTamperDetected,
  resetRisk,
} = useCaptureDetection({
  videoId: props.video?.id || props.lessonId,
  sessionId: sessionClient.getSessionId() || "sess-init",
  platform: detectPlatform(),
  onRiskChange: (score, level) => {
    // SUSPICIOUS → non-accusatory advisory (never accuse the user, never pause).
    if (level === "SUSPICIOUS" && lastAdvisoryLevel !== "SUSPICIOUS" && !warningModalOpen.value) {
      showSuspiciousNotice();
    }
    lastAdvisoryLevel = level;
  },
  onHighRisk: (score) => {
    pauseVideo();
    triggerHighRiskLockout(
      `تم اكتشاف نشاط مشبوه أثناء تشغيل المحاضرة (درجة الخطر: ${score}). يرجى إيقاف أي برامج تسجيل أو مشاركة شاشة للمتابعة.`
    );
  },
  onCritical: (score) => {
    sessionClient.revoke();
    pauseVideo();
    triggerCriticalLockout(
      `تم إنهاء جلسة العرض بسبب مستوى خطر مرتفع (${score}). يرجى إعادة فتح المحاضرة للمتابعة.`
    );
  },
  onTampered: () => {
    pauseVideo();
    triggerHighRiskLockout(
      "تم رصد محاولة تلاعب أمني أو إزالة العلامة المائية. تم إيقاف المشغل لحماية المحتوى."
    );
  },
});

const progressPercent = computed(() =>
  !duration.value ? 0 : Math.min(100, Math.max(0, (currentTime.value / duration.value) * 100))
);
const currentTimeText = computed(() => formatClockTime(currentTime.value));
const durationText = computed(() => formatClockTime(duration.value));

const banRemainingText = computed(() => {
  const s = Math.max(0, banRemaining.value);
  const h = Math.floor(s / 3600);
  const m = Math.floor((s % 3600) / 60);
  const sec = s % 60;
  const pad = (n: number) => String(n).padStart(2, "0");
  return h > 0 ? `${h}:${pad(m)}:${pad(sec)}` : `${m}:${pad(sec)}`;
});

let controlsTimer: ReturnType<typeof setTimeout> | null = null;
let saveTimer: ReturnType<typeof setInterval> | null = null;
let hls: Hls | null = null;
let flashTimer: ReturnType<typeof setTimeout> | null = null;
let tapTimer: ReturnType<typeof setTimeout> | null = null;
let warningTimer: ReturnType<typeof setTimeout> | null = null;
let banTimer: ReturnType<typeof setInterval> | null = null;
let lastTapTime = 0;
let pendingSingleTap = false;

const isBlackout = ref(false);
const privacyPauseNotice = ref<string | null>(null);

function triggerBlackout() {
  isBlackout.value = true;
  pauseVideo();
}

const { engine } = useContentProtection({
  lessonId: props.lessonId,
  closeOnCapture: true,
  onBanned: (result) => {
    bannedInfo.value = result;
    triggerBlackout();
    startBanCountdown(result);
    // Instant close app on ban / capture
    setTimeout(() => {
      engine.closeApp();
    }, 1500);
  },
  onWarning: (result) => {
    warningInfo.value = result;
    triggerBlackout();
    if (warningTimer) clearTimeout(warningTimer);
    // Close the app immediately upon violation/warning
    warningTimer = setTimeout(() => {
      engine.closeApp();
    }, 1000);
  },
  onViolation: (_type) => {
    triggerBlackout();
    violationNotice.value = "تم إيقاف التشغيل بعد رصد إشارة التقاط أو تلاعب أمني.";
    setTimeout(() => {
      engine.closeApp();
    }, 1200);
  },
  onPrivacyPause: () => {
    privacyPauseNotice.value = "تم إيقاف الفيديو مؤقتًا لحماية المحتوى بعد مغادرة نافذة العرض أو دخولها إلى الخلفية.";
    pauseVideo();
  },
});

function dismissViolationNotice() {
  violationNotice.value = null;
}

// ── playback ──────────────────────────────────────────────────────────────
function pauseVideo() {
  isPlaying.value = false;
  engine.setPlaying(false);
  videoRef.value?.pause();
}

function togglePlay() {
  const v = videoRef.value;
  if (!v || bannedInfo.value) return;
  if (v.paused) {
    void v.play().catch(() => {});
  } else {
    v.pause();
  }
}

function seekTo(sec: number) {
  if (videoRef.value) videoRef.value.currentTime = sec;
}

function skip(sec: number) {
  const t = Math.min(Math.max(0, currentTime.value + sec), duration.value || 0);
  seekTo(t);
  currentTime.value = t;
}

function cyclePlaybackRate() {
  const next = PLAYBACK_RATES[(PLAYBACK_RATES.indexOf(playbackRate.value) + 1) % PLAYBACK_RATES.length];
  playbackRate.value = next;
  if (videoRef.value) videoRef.value.playbackRate = next;
}

function onSeekInput(event: Event) {
  const val = Number((event.target as HTMLInputElement).value);
  currentTime.value = val;
  seekTo(val);
}

// ── fullscreen enforcement ────────────────────────────────────────────────
function isRealFullscreen() {
  const doc = document as any;
  return !!(
    doc.fullscreenElement ||
    doc.webkitFullscreenElement ||
    doc.mozFullScreenElement ||
    doc.msFullscreenElement
  );
}

function onFullscreenChange() {
  isFullscreen.value = isRealFullscreen();
}

function onVisibilityReturn() {
  if (document.visibilityState === "visible") {
    // On iOS and Telegram, returning from the background is not proof of capture.
    // Keep the player paused until the user explicitly resumes it.
    return;
  }
}

async function resumeAfterPrivacyPause() {
  privacyPauseNotice.value = null;
  togglePlay();
}

async function toggleFullscreen() {
  // Fullscreen is optional on iOS and Telegram; when supported, request it only
  // from an explicit user action rather than forcing it during mount.
  const el = containerRef.value as any;
  try {
    if (el?.requestFullscreen) await el.requestFullscreen();
    else if (el?.webkitRequestFullscreen) el.webkitRequestFullscreen();
  } catch {}
}

// ── HLS / stream ──────────────────────────────────────────────────────────
function attachHls(url: string) {
  const video = videoRef.value;
  if (!video) return;
  if (Hls.isSupported()) {
    // Bind every segment and key request to the issuing device so a captured
    // playlist cannot be replayed from another browser.
    const dev = deviceFingerprint();
    hls = new Hls({
      enableWorker: true,
      maxBufferLength: 30,
      maxMaxBufferLength: 60,
      xhrSetup(xhr, url) {
        xhr.setRequestHeader("X-Device-ID", dev);
        void url;
      },
      fetchSetup(_context, initParams) {
        const headers = new Headers(initParams.headers);
        headers.set("X-Device-ID", dev);
        initParams.headers = headers;
        return initParams;
      },
    });
    hls.loadSource(url);
    hls.attachMedia(video);
    hls.on(Hls.Events.MANIFEST_PARSED, () => {
      isLoading.value = false;
      void video.play().catch(() => {});
    });
    hls.on(Hls.Events.ERROR, (_event, data) => {
      if (!data.fatal) return;
      if (data.type === Hls.ErrorTypes.NETWORK_ERROR) {
        hls?.startLoad();
      } else if (data.type === Hls.ErrorTypes.MEDIA_ERROR) {
        hls?.recoverMediaError();
      } else {
        errorMsg.value = "حدث خطأ أثناء تشغيل الفيديو";
        isLoading.value = false;
      }
    });
    return;
  }
  // iOS Safari native HLS
  if (video.canPlayType("application/vnd.apple.mpegurl")) {
    video.src = url;
    isLoading.value = false;
    void video.play().catch(() => {});
    return;
  }
  errorMsg.value = "متصفحك لا يدعم البث المحمي";
  isLoading.value = false;
}

async function initStream() {
  const v = props.video;
  if (!v) {
    isLoading.value = false;
    return;
  }
  isLoading.value = true;
  errorMsg.value = "";
  try {
    // 1. Start unique Playback Session & Heartbeat
    try {
      const sessRes = await sessionClient.startSession(v.id, undefined, detectPlatform());
      if (sessRes?.watermark) {
        watermarkConfig.value = sessRes.watermark;
      }
      sessionClient.startHeartbeat(
        () => currentTime.value,
        () => isPlaying.value,
        () => currentRiskScore.value,
        () => {
          pauseVideo();
          triggerHighRiskLockout("انتهت صلاحية جلسة البث أو تم إغلاقها من الخادم.");
        }
      );
    } catch {
      /* best effort fallback */
    }

    const grant = await createStreamGrant(v.id);
    if (!grant.success) throw new Error("grant failed");
    const url = resolveApiUrl(grant.playlist_url);
    if (grant.mode === "direct") {
      const video = videoRef.value;
      if (video) {
        video.src = url;
        isLoading.value = false;
        void video.play().catch(() => {});
      }
      return;
    }
    attachHls(url);
  } catch (e: any) {
    errorMsg.value = e?.response?.data?.detail || "تعذر تحميل الفيديو";
    isLoading.value = false;
  }
}

// ── video events ──────────────────────────────────────────────────────────
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

function syncResumeTimer() {
  if (saveTimer) return;
  saveTimer = setInterval(() => {
    const v = videoRef.value;
    if (!v || v.paused || duration.value <= 0) return;
    const pos = currentTime.value;
    const pct = Math.min(100, Math.round((pos / duration.value) * 100));
    playerStore.setResume(props.lessonId, pos);
    saveLessonProgress(props.lessonId, {
      progress: pct,
      watched_seconds: Math.round(pos),
      last_position: pos,
      duration: Math.round(duration.value),
    }).catch(() => {});
  }, 5000);
}

function saveFinalPosition() {
  const v = videoRef.value;
  if (v) playerStore.setResume(props.lessonId, v.currentTime);
}

// ── controls UX ───────────────────────────────────────────────────────────
function scheduleControlHide(delay = 2500) {
  if (controlsTimer) clearTimeout(controlsTimer);
  if (!isFullscreen.value && !isPlaying.value) return;
  controlsTimer = setTimeout(() => {
    showControls.value = false;
  }, delay);
}

function wakeControls() {
  showControls.value = true;
  scheduleControlHide();
}

function flashSeek(sec: number) {
  if (flashTimer) clearTimeout(flashTimer);
  // Flash overlay intentionally omitted to keep the protected surface minimal.
  void sec;
}

function handleSurfaceTap(event: PointerEvent) {
  const revealOnly = !showControls.value;
  wakeControls();
  if (event.target && (event.target as HTMLElement).closest("button, input, a")) return;

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
    skip(event.clientX < half ? -10 : 10);
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
      await initStream();
    }
  } catch {
    /* keep the ban overlay */
  }
}

async function revalidateSession() {
  dismissWarning();
  sessionClient.revoke();
  resetRisk();
  await initStream();
}

// ── lifecycle ─────────────────────────────────────────────────────────────
onMounted(() => {
  const platform = detectPlatform();
  if (platform === "DESKTOP_WEB" || platform === "UNKNOWN") engine.requestFullscreen();
  document.addEventListener("fullscreenchange", onFullscreenChange);
  document.addEventListener("webkitfullscreenchange", onFullscreenChange);
  document.addEventListener("visibilitychange", onVisibilityReturn);
  syncResumeTimer();
  void initStream();

  // Native on-device screen-capture shield (Capacitor builds): Android gets
  // OS-level FLAG_SECURE blocking; iOS gets capture/screenshot detection with
  // instant pause + lockout + violation report. No-op on web.
  void startScreenShield({
    lessonId: props.video?.id || props.lessonId,
    onCaptured: () => {
      sessionClient.revoke();
      pauseVideo();
      triggerCriticalLockout(
        "تم رصد تسجيل الشاشة على جهازك. تم إيقاف المشغل وإنهاء الجلسة لحماية المحتوى."
      );
    },
  });

  // Platform cannot verify capture → non-accusatory advisory. Deterrence
  // (watermark + session tracking) keeps working regardless.
  if (!capabilities.canDetectCapture) {
    showDetectionUnavailable();
  }
});

onBeforeUnmount(() => {
  saveFinalPosition();
  void stopScreenShield();
  if (saveTimer) {
    clearInterval(saveTimer);
    saveTimer = null;
  }
  if (controlsTimer) clearTimeout(controlsTimer);
  if (flashTimer) clearTimeout(flashTimer);
  if (tapTimer) clearTimeout(tapTimer);
  if (warningTimer) clearTimeout(warningTimer);
  if (banTimer) clearInterval(banTimer);
  if (hls) {
    hls.destroy();
    hls = null;
  }
  document.removeEventListener("fullscreenchange", onFullscreenChange);
  document.removeEventListener("webkitfullscreenchange", onFullscreenChange);
  document.removeEventListener("visibilitychange", onVisibilityReturn);
});
</script>

<template>
  <div
    ref="containerRef"
    id="video-container"
    dir="ltr"
    class="relative w-full aspect-video bg-black overflow-hidden group select-none"
    @pointermove="wakeControls"
    @pointerleave="() => scheduleControlHide()"
    @contextmenu.prevent
  >
    <!-- Native video element (no native controls reachable) -->
    <video
      v-if="video"
      ref="videoRef"
      :class="[
        'w-full h-full object-cover sm:object-contain select-none transition-opacity duration-150',
        isBlackout ? 'opacity-0 pointer-events-none filter blur-3xl' : 'opacity-100'
      ]"
      :poster="poster || undefined"
      preload="metadata"
      playsinline
      draggable="false"
      :controls="false"
      :controlslist="'nodownload noremoteplayback noplaybackrate'"
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

    <!-- Total Blackout Overlay when capture or violation is detected -->
    <div
      v-if="isBlackout"
      class="absolute inset-0 z-[100] bg-black flex flex-col items-center justify-center p-6 text-center select-none pointer-events-auto"
      dir="rtl"
    >
      <ShieldAlert class="w-16 h-16 text-red-500 mb-3 animate-pulse" />
      <h3 class="text-lg font-black text-white mb-2">🚫 تم حظر تسجيل الشاشة</h3>
      <p class="text-xs text-slate-400 max-w-xs leading-relaxed">
        تم تسويد الشاشة وتفعيل الحظر الأمني الآلي بسبب محاولة تصوير أو تسجيل المحتوى.
      </p>
    </div>

    <!-- Conservative privacy pause for iOS/Telegram lifecycle signals -->
    <div
      v-if="privacyPauseNotice && !isBlackout && !bannedInfo"
      class="absolute inset-0 z-[80] flex items-center justify-center bg-black/90 backdrop-blur-sm px-6 text-center"
      dir="rtl"
      role="alert"
      aria-live="polite"
    >
      <div class="max-w-sm space-y-4">
        <ShieldCheck class="mx-auto h-12 w-12 text-emerald-400" />
        <p class="text-sm font-bold leading-relaxed text-white">{{ privacyPauseNotice }}</p>
        <p class="text-xs leading-relaxed text-white/60">لا يمكن لتطبيق ويب منع التصوير الخارجي بشكل مطلق؛ العلامة المائية والرابط الموقّع يحدّان من إساءة الاستخدام.</p>
        <button
          type="button"
          class="rounded-xl bg-orange-500 px-5 py-2.5 text-sm font-black text-white transition hover:bg-orange-600"
          @click="resumeAfterPrivacyPause"
        >
          متابعة التشغيل
        </button>
      </div>
    </div>

    <!-- Dynamic Watermark Protection & Anti-Recording Layer -->
    <DynamicWatermark :config="watermarkConfig" @tampered="handleTamperDetected" />
    <WatermarkOverlay />

    <!-- Security Warning Modal -->
    <SecurityWarningModal
      :open="warningModalOpen"
      :level="warningLevel"
      :title="warningTitle"
      :message="warningMessage"
      :detail="warningDetail"
      :require-revalidation="requireRevalidation"
      @dismiss="dismissWarning"
      @revalidate="revalidateSession"
    />

    <!-- Transparent click interceptor: blocks direct access to the media
         surface / native context menu and toggles playback -->
    <div
      class="absolute inset-0 z-10 cursor-pointer"
      aria-label="تشغيل أو إيقاف"
      @pointerdown="handleSurfaceTap"
      @contextmenu.prevent
    ></div>

    <!-- No video state -->
    <div v-if="!video" class="absolute inset-0 flex flex-col items-center justify-center gap-3 bg-slate-900 text-slate-400 z-10">
      <Play class="w-12 h-12 text-slate-500" />
      <span class="text-xs font-bold" dir="rtl">لا يوجد فيديو متاح لهذه المحاضرة</span>
    </div>

    <!-- Loading spinner -->
    <div v-if="isLoading && video && !errorMsg" class="absolute inset-0 flex items-center justify-center bg-black/50 pointer-events-none z-20">
      <Loader2 class="w-12 h-12 text-orange-500 animate-spin" />
    </div>

    <!-- Error state -->
    <div v-if="errorMsg && !isLoading" class="absolute inset-0 flex flex-col items-center justify-center gap-3 bg-black/80 px-6 text-center z-20">
      <ShieldAlert class="w-10 h-10 text-orange-400" />
      <span class="text-sm font-bold text-white" dir="rtl">{{ errorMsg }}</span>
    </div>

    <!-- Center play button overlay when paused -->
    <button
      v-if="video && !isPlaying && !isLoading && !errorMsg && !bannedInfo"
      @click="togglePlay"
      type="button"
      aria-label="تشغيل المحاضرة"
      class="absolute inset-0 flex items-center justify-center bg-black/30 z-20 transition-opacity cursor-pointer"
    >
      <span class="w-16 h-16 sm:w-20 sm:h-20 rounded-full bg-orange-500 text-white flex items-center justify-center shadow-2xl shadow-orange-500/50 scale-in hover:scale-110 active:scale-95 transition-all">
        <Play class="w-8 h-8 sm:w-10 sm:h-10 fill-white ml-1" />
      </span>
    </button>

    <!-- Custom Controls Bar -->
    <div
      v-if="video && !bannedInfo"
      :class="[
        'absolute inset-x-0 bottom-0 z-30 bg-gradient-to-t from-black/95 via-black/60 to-transparent transition-all duration-300',
        showControls ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-2 pointer-events-none',
        isFullscreen ? 'px-5 pb-4 pt-16' : 'px-4 pb-3 pt-12'
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

        <div class="flex-1"></div>

        <!-- Security Risk Status -->
        <SecurityStatus :score="currentRiskScore" :level="currentRiskLevel" />

        <!-- Protected badge -->
        <span class="inline-flex items-center gap-1 px-2 py-1 rounded-lg bg-white/10 text-[10px] font-black text-white/70" dir="rtl">
          <ShieldCheck class="w-3.5 h-3.5 text-emerald-400" />
          محمي
        </span>

        <!-- Fullscreen -->
        <button
          @click="toggleFullscreen"
          type="button"
          :aria-label="isFullscreen ? 'إنهاء ملء الشاشة' : 'ملء الشاشة'"
          class="player-btn cursor-pointer"
          title="ملء الشاشة"
        >
          <Minimize v-if="isFullscreen" class="w-4.5 h-4.5 text-orange-400" />
          <Maximize v-else class="w-4.5 h-4.5" />
        </button>
      </div>
    </div>

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
</style>
