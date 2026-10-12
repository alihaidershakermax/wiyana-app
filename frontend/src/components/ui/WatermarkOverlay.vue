<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from "vue";
import { useUserStore } from "@/stores/useUserStore";

const userStore = useUserStore();

// Watermark display: Telegram Username or Telegram ID ONLY
const telegramWatermarkText = computed(() => {
  const u = userStore.currentUser;
  const tgUser = (window as any).Telegram?.WebApp?.initDataUnsafe?.user;
  const username = u?.username || tgUser?.username;
  const tgId = u?.telegram_id || tgUser?.id || u?.id;

  if (username) {
    return `@${username.replace(/^@/, '')}`;
  }
  if (tgId) {
    return `ID: ${tgId}`;
  }
  return u?.name || "طالب ويانا";
});

const positions = [
  { top: "10%", left: "8%" },
  { top: "18%", right: "10%" },
  { bottom: "22%", left: "12%" },
  { bottom: "18%", right: "12%" },
  { top: "40%", left: "30%" },
  { top: "50%", right: "30%" },
  { top: "12%", left: "45%" },
  { bottom: "10%", left: "40%" },
  { top: "70%", right: "25%" },
  { top: "30%", left: "60%" },
  { bottom: "35%", right: "50%" },
];

const currentPosIndex = ref(0);
const secondPosIndex = computed(() => (currentPosIndex.value + 5) % positions.length);
let posTimer: ReturnType<typeof setInterval> | null = null;

onMounted(() => {
  posTimer = setInterval(() => {
    currentPosIndex.value = (currentPosIndex.value + 1) % positions.length;
  }, 2500);
});

onUnmounted(() => {
  if (posTimer) clearInterval(posTimer);
});
</script>

<template>
  <div
    class="absolute inset-0 z-30 overflow-hidden select-none pointer-events-none"
    aria-hidden="true"
  >
    <!-- Primary Floating Watermark -->
    <div
      class="absolute transition-all duration-700 ease-in-out px-2.5 py-1 rounded text-xs sm:text-sm font-mono font-bold text-white/80 tracking-wider shadow-md bg-black/30 backdrop-blur-[1px]"
      :style="[
        positions[currentPosIndex],
        { textShadow: '0 1px 3px rgba(0,0,0,0.9), 0 0 2px rgba(0,0,0,0.8)' },
      ]"
    >
      <span>{{ telegramWatermarkText }}</span>
    </div>

    <!-- Secondary Floating Watermark (Opposite side) -->
    <div
      class="absolute transition-all duration-700 ease-in-out px-2 py-0.5 rounded text-[11px] sm:text-xs font-mono font-bold text-white/60 tracking-wider shadow-sm bg-black/20"
      :style="[
        positions[secondPosIndex],
        { textShadow: '0 1px 2px rgba(0,0,0,0.9)' },
      ]"
    >
      <span>{{ telegramWatermarkText }}</span>
    </div>

    <!-- Minimal diagonal repeated background watermark -->
    <div
      class="absolute inset-0 opacity-[0.14] flex flex-wrap justify-around items-center gap-16 p-6 rotate-[-18deg]"
      aria-hidden="true"
    >
      <span
        v-for="i in 24"
        :key="i"
        class="text-xs sm:text-sm font-bold text-white whitespace-nowrap tracking-widest uppercase"
      >
        {{ telegramWatermarkText }}
      </span>
    </div>
  </div>
</template>
