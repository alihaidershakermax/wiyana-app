<template>
  <div class="dash-player">
    <video ref="videoRef" controls playsinline class="w-full rounded-lg bg-black" />
    <div v-if="error" class="mt-2 text-sm text-red-500">{{ error }}</div>
    <div v-if="loading" class="mt-2 text-sm text-gray-500">جاري تحميل الفيديو...</div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from "vue";
import { useUserStore } from "@/stores/useUserStore";

const props = defineProps<{
  src: string; // URL to manifest.mpd on R2
}>();

const videoRef = ref<HTMLVideoElement | null>(null);
const error = ref("");
const loading = ref(true);
let player: any = null;

const userStore = useUserStore();

async function initPlayer() {
  if (!videoRef.value) return;
  try {
    const dashjs = await import("dashjs");
    player = dashjs.MediaPlayer().create();

    // ClearKey protection: MUST be set BEFORE initialize
    const token = userStore.token;
    player.setProtectionData({
      "org.w3.clearkey": {
        serverURL: `${import.meta.env.VITE_API_URL || "/api/v1"}/drm/license`,
        httpRequestHeaders: token ? { Authorization: `Bearer ${token}` } : {},
      },
    });

    player.initialize(videoRef.value, props.src, true);

    player.on(dashjs.MediaPlayer.events.ERROR, (e: any) => {
      console.error("[DASH]", e);
      if (e.error?.code === 25) {
        error.value = "تعذر فك تشفير الفيديو — تأكد من تسجيل الدخول";
      }
    });
    player.on(dashjs.MediaPlayer.events.STREAM_INITIALIZED, () => {
      loading.value = false;
    });
  } catch (e: any) {
    error.value = "تعذر تحميل المشغل";
    loading.value = false;
  }
}

onMounted(initPlayer);
onUnmounted(() => {
  if (player) {
    player.destroy();
    player = null;
  }
});
watch(() => props.src, () => {
  if (player) player.attachSource(props.src);
});
</script>

<style scoped>
.dash-player video {
  aspect-ratio: 16 / 9;
  max-height: 70vh;
}
</style>
