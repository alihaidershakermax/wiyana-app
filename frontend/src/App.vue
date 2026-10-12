<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount } from "vue";
import { useRoute } from "vue-router";
import { useQueryClient } from "@tanstack/vue-query";
import BottomNavigation from "@/components/layout/BottomNavigation.vue";
import TelegramOnlyGuard from "@/components/layout/TelegramOnlyGuard.vue";
import AppUpdateModal from "@/components/ui/AppUpdateModal.vue";
import { useUserStore } from "@/stores/useUserStore";
import { useAppFullscreen } from "@/composables/useAppFullscreen";

const route = useRoute();
const userStore = useUserStore();
const queryClient = useQueryClient();
const isAdminRoute = computed(() => route.path.startsWith("/admin"));
const isCourseDetail = computed(() => /^\/courses\/[^/]+$/.test(route.path));
const isLessonPlayer = computed(() => /^\/lessons\/[^/]+$/.test(route.path));
const showBottomNav = computed(() => !isAdminRoute.value && !isCourseDetail.value && !isLessonPlayer.value);
useAppFullscreen();

function pauseAllVideos() {
  try {
    document.querySelectorAll("video").forEach((v) => v.pause());
  } catch {}
}

function triggerPrivacyLockout() {
  if (isAdminRoute.value) return;
  // No black screen: just stop any playing video. The video player itself
  // shows the warning message overlay on detection.
  if (isLessonPlayer.value || document.querySelector("video")) {
    pauseAllVideos();
  }
}

function onGlobalContextMenu(event: MouseEvent) {
  if (isAdminRoute.value) return;
  const target = event.target as HTMLElement | null;
  if (target?.closest?.("input, textarea, [contenteditable='true']")) return;
  event.preventDefault();
}

function handleGlobalKeydown(event: KeyboardEvent) {
  if (isAdminRoute.value) return;
  const key = event.key ? event.key.toLowerCase() : "";
  const code = event.code ? event.code.toLowerCase() : "";
  const ctrl = event.ctrlKey;
  const meta = event.metaKey; // Windows Key or Mac Command Key
  const shift = event.shiftKey;
  const alt = event.altKey;

  const isScreenshotKey =
    key === "printscreen" ||
    code === "printscreen" ||
    key === "f12" ||
    code === "f12" ||
    (meta && shift && (key === "s" || key === "3" || key === "4" || key === "5" || code === "keys")) || // Win+Shift+S (Snipping Tool) or Mac Cmd+Shift+3/4/5
    (alt && (key === "printscreen" || code === "printscreen")) ||
    (ctrl && shift && (key === "i" || key === "j" || key === "c" || key === "s")) ||
    (ctrl && (key === "u" || key === "s" || key === "p"));

  if (isScreenshotKey) {
    event.preventDefault();
    event.stopPropagation();
    triggerPrivacyLockout();
  }
}

function onWindowBlur() {
  if (isAdminRoute.value) return;
  triggerPrivacyLockout();
}

function onVisibilityChange() {
  if (document.visibilityState === "hidden") {
    triggerPrivacyLockout();
  }
}

function onMobilePageHide() {
  triggerPrivacyLockout();
}

function handleTouchStart(event: TouchEvent) {
  // 3 or 4 finger touch gesture is the standard native mobile screenshot shortcut
  if (event.touches && event.touches.length >= 3) {
    try {
      event.preventDefault();
    } catch {}
    triggerPrivacyLockout();
  }
}

onMounted(async () => {
  document.addEventListener("contextmenu", onGlobalContextMenu);
  window.addEventListener("keydown", handleGlobalKeydown, { capture: true });
  window.addEventListener("keyup", handleGlobalKeydown, { capture: true });
  window.addEventListener("blur", onWindowBlur);
  window.addEventListener("pagehide", onMobilePageHide);
  window.addEventListener("freeze", onMobilePageHide);
  window.addEventListener("wiyana:protection-lockout", triggerPrivacyLockout);
  document.addEventListener("visibilitychange", onVisibilityChange);
  document.addEventListener("touchstart", handleTouchStart, { passive: false, capture: true });

  // Telegram Mini App Security Binding for screen protection
  const tg = (window as any).Telegram?.WebApp;
  if (tg) {
    try {
      tg.ready();
      tg.expand();
      tg.requestFullscreen?.();
      tg.disableVerticalSwipes?.();
      tg.enableClosingConfirmation?.();
    } catch {}

    // ── Telegram event bindings for screen-protection ─────────────────
    // These fire ONLY on strong "user left / minimized / lost fullscreen"
    // signals while a video is on screen. Guarded to avoid the false
    // positives that previously forced removal of this binding.
    const hasActiveVideo = () => isLessonPlayer.value || !!document.querySelector("video");
    // Ignore events fired during the initial mount (expand/fullscreen settle).
    let tgSettled = false;
    setTimeout(() => {
      tgSettled = true;
    }, 3000);

    try {
      // viewportChanged: lock only when the Mini App is collapsed/minimized
      // (isExpanded === false), never on ordinary resize/expand events.
      tg.onEvent("viewportChanged", (evt: any) => {
        if (!tgSettled) return;
        if (evt?.isStateChanged && !tg.isExpanded && hasActiveVideo()) {
          triggerPrivacyLockout();
        }
      });

      // fullscreenChanged: lock when the forced fullscreen is lost while a
      // video is present (screen-recorder overlay / app-switch signal).
      tg.onEvent("fullscreenChanged", (evt: any) => {
        if (!tgSettled) return;
        if (evt?.isFullscreen === false && hasActiveVideo()) {
          triggerPrivacyLockout();
        }
      });
    } catch {}
  }

  // Sync dark mode from saved settings
  if (userStore.settings.darkMode) {
    document.documentElement.classList.add("dark");
  } else {
    document.documentElement.classList.remove("dark");
  }
  await userStore.initAuth();
  queryClient.invalidateQueries();
});

onBeforeUnmount(() => {
  document.removeEventListener("contextmenu", onGlobalContextMenu);
  window.removeEventListener("keydown", handleGlobalKeydown, { capture: true });
  window.removeEventListener("keyup", handleGlobalKeydown, { capture: true });
  window.removeEventListener("blur", onWindowBlur);
  window.removeEventListener("pagehide", onMobilePageHide);
  window.removeEventListener("freeze", onMobilePageHide);
  window.removeEventListener("wiyana:protection-lockout", triggerPrivacyLockout);
  document.removeEventListener("visibilitychange", onVisibilityChange);
  document.removeEventListener("touchstart", handleTouchStart, { capture: true });
});
</script>

<template>
  <div
    dir="rtl"
    :class="[
      'app-shell min-h-screen bg-white dark:bg-[#0F0F0F] text-[#111111] dark:text-[#FFFFFF] font-sans selection:bg-orange-500 selection:text-white antialiased overscroll-y-none transition-colors duration-200',
      { 'pb-20': showBottomNav },
      { 'app-shell--admin': isAdminRoute },
    ]"
  >
    <!-- Telegram Only Access Guard Screen -->
    <TelegramOnlyGuard />

    <!-- App Update Features Announcement Modal -->
    <AppUpdateModal />

    <!-- Router View for Pages: one consistent, interruptible page transition -->
    <RouterView v-slot="{ Component, route: viewRoute }">
      <Transition name="page" mode="out-in">
        <component :is="Component" :key="viewRoute.fullPath" />
      </Transition>
    </RouterView>

    <!-- Fixed Native iOS Bottom Navigation Bar -->
    <BottomNavigation v-if="showBottomNav" />
  </div>
</template>

<style>
.page-enter-active,
.page-leave-active {
  transition: opacity 180ms cubic-bezier(0.23, 1, 0.32, 1), transform 180ms cubic-bezier(0.23, 1, 0.32, 1);
}

.page-enter-from {
  opacity: 0;
  transform: translateY(6px);
}

.page-leave-to {
  opacity: 0;
  transform: translateY(-4px);
}

/* iOS Safe Area Insets for Mobile Devices & Telegram Mini App */
.pb-safe {
  padding-bottom: env(safe-area-inset-bottom, 12px);
}

/* Content protection: prevent dragging images/saving them via long-press */
img,
video {
  -webkit-user-drag: none;
  user-select: none;
  -webkit-touch-callout: none;
}
</style>
