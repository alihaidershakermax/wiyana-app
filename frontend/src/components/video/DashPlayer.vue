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

const props = defineProps<{
  src: string; // URL to manifest.mpd on R2 (no DRM)
}>();

const videoRef = ref<HTMLVideoElement | null>(null);
const error = ref("");
let player: any = null;

async function initPlayer() {
  if (!videoRef.value || !props.src) return;
  error.value = "";
  try {
    // dash.js reference player (like reference.dashif.org) — no DRM
    const dashjs = await import("dashjs");
    player = dashjs.MediaPlayer().create();
    player.initialize(videoRef.value, props.src, false);

    player.on(dashjs.MediaPlayer.events.ERROR, (e: any) => {
      console.error("[DASH]", e);
      error.value = "تعذر تشغيل الفيديو";
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
