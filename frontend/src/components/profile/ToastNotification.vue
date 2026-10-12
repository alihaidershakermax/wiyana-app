<script setup lang="ts">
import { computed } from "vue";
import { CheckCircle2, AlertCircle, Info, AlertTriangle, X } from "@lucide/vue";

const props = defineProps<{
  show: boolean;
  message: string;
  type?: "success" | "info" | "warning" | "error";
}>();

const emit = defineEmits<{
  (e: "close"): void;
}>();

const iconComponent = computed(() => {
  switch (props.type) {
    case "error":
      return AlertCircle;
    case "warning":
      return AlertTriangle;
    case "info":
      return Info;
    default:
      return CheckCircle2;
  }
});

const colorClasses = computed(() => {
  switch (props.type) {
    case "error":
      return "border-rose-500/30 bg-white/95 dark:bg-slate-900/95 text-rose-600 dark:text-rose-400 shadow-rose-500/10";
      case "warning":
        return "border-orange-500/30 bg-white/95 dark:bg-slate-900/95 text-orange-600 dark:text-orange-400 shadow-orange-500/10";
      case "info":
        return "border-orange-500/30 bg-white/95 dark:bg-slate-900/95 text-orange-600 dark:text-orange-400 shadow-orange-500/10";
    default:
      return "border-orange-500/40 bg-white/95 dark:bg-slate-900/95 text-orange-600 dark:text-orange-400 shadow-orange-500/15";
  }
});
</script>

<template>
  <Transition
    enter-active-class="transform ease-out duration-300 transition"
    enter-from-class="translate-y-4 opacity-0 scale-95"
    enter-to-class="translate-y-0 opacity-100 scale-100"
    leave-active-class="transition ease-in duration-200"
    leave-from-class="opacity-100 scale-100"
    leave-to-class="opacity-0 scale-95"
  >
    <div
      v-if="show"
      class="fixed top-5 left-1/2 -translate-x-1/2 z-50 max-w-md w-[calc(100%-2rem)] pointer-events-auto"
    >
      <div
        :class="[
          'flex items-center justify-between gap-3 px-4 py-3.5 rounded-2xl border backdrop-blur-xl shadow-xl',
          colorClasses
        ]"
      >
        <div class="flex items-center gap-3 min-w-0">
          <div class="w-8 h-8 rounded-xl bg-orange-500/10 flex items-center justify-center shrink-0">
            <component :is="iconComponent" class="w-5 h-5 text-orange-500" />
          </div>
          <p class="text-xs sm:text-sm font-bold text-slate-800 dark:text-slate-100 leading-snug truncate">
            {{ message }}
          </p>
        </div>

        <button
          @click="emit('close')"
          type="button"
          class="p-1 rounded-lg text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 transition-colors shrink-0"
        >
          <X class="w-4 h-4" />
        </button>
      </div>
    </div>
  </Transition>
</template>
