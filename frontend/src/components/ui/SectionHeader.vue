<script setup lang="ts">
import { ChevronLeft } from "@lucide/vue";

defineProps<{
  title: string;
  subtitle?: string;
  actionText?: string;
}>();

const emit = defineEmits<{
  (e: "clickAction"): void;
}>();
</script>

<template>
  <div class="flex items-center justify-between gap-4 mb-4">
    <div>
      <h2 class="text-lg sm:text-xl font-bold text-slate-900 dark:text-slate-100 tracking-tight">
        {{ title }}
      </h2>
      <p v-if="subtitle" class="text-xs sm:text-sm text-slate-500 dark:text-slate-400 mt-0.5 font-normal">
        {{ subtitle }}
      </p>
    </div>

    <button
      v-if="actionText"
      @click="emit('clickAction')"
      type="button"
      class="inline-flex items-center gap-1 text-xs sm:text-sm font-semibold text-orange-500 hover:text-orange-600 transition-colors focus:outline-none focus-visible:ring-2 focus-visible:ring-orange-500/50 rounded-lg px-2 py-1"
    >
      <span>{{ actionText }}</span>
      <ChevronLeft class="w-4 h-4 rtl:rotate-0 ltr:rotate-180" />
    </button>
  </div>
</template>
