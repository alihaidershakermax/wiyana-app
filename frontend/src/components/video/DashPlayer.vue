<template>
  <div class="dash-player w-full">
    <video
      ref="videoRef"
      class="w-full h-full bg-black"
      controls
      playsinline
      preload="metadata"
    />
    <div v-if="error" class="p-3 text-sm text-red-500 text-center">{{ error }}</div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from "vue";
import { useUserStore } from "@/stores/useUserStore";

const props = defineProps<{
  src: string; // URL to manifest.mpd on R2 (ClearKey encrypted)
}>();

const videoRef = ref<HTMLVideoElement | null>(null);
const error = ref("");
let player: any = null;

const userStore = useUserStore();

async function initPlayer() {
  if (!videoRef.value || !props.src) return;
  error.value = "";
  try {
    const dashjs = await import("dashjs");
    player = dashjs.MediaPlayer().create();

    // ClearKey DRM: protection data MUST be set BEFORE initialize
    const token = userStore.token;
    player.setProtectionData({
      "org.w3.clearkey": {
        serverURL: `${import.meta.env.VITE_API_URL || "/api/v1"}/drm/license`,
        httpRequestHeaders: token ? { Authorization: `Bearer ${token}` } : {},
      },
    });

    player.initialize(videoRef.value, props.src, false);

    player.on(dashjs.MediaPlayer.events.ERROR, (e: any) => {
      console.error("[DASH]", e);
      if (e.error?.code === 25) {
        error.value = "تعذر فك تشفير الفيديو — تأكد من تسجيل الدخول";
      } else {
        error.value = "تعذر تشغيل الفيديو";
      }
    });
  } catch (e: any) {
    console.error("[DASH] init failed", e);
    error.value = "تعذر تحميل المشغل";
  }
}

function destroyPlayer() {
  if (player) {
    try { player.destroy(); } catch {}
    player = null;
  }
}

onMounted(initPlayer);
onUnmounted(destroyPlayer);
watch(() => props.src, () => {
  destroyPlayer();
  initPlayer();
});
</script>

<style scoped>
.dash-player {
  position: relative;
  width: 100%;
  aspect-ratio: 16 / 9;
  background: #000;
}
.dash-player video {
  width: 100%;
  height: 100%;
  object-fit: contain;
}
</style>
