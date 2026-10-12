<script setup lang="ts">
import type { Component } from "vue";
import { ChevronLeft } from "@lucide/vue";

defineProps<{
  label: string;
  badgeText?: string;
  icon: Component;
  badgeColorClass?: string;
}>();

const emit = defineEmits<{
  (e: "click"): void;
}>();
</script>

<template>
  <button
    @click="emit('click')"
    type="button"
    class="w-full flex items-center justify-between p-3.5 hover:bg-slate-50 dark:hover:bg-slate-800/60 transition-colors text-right group focus:outline-none focus:bg-slate-50 dark:focus:bg-slate-800/60"
  >
    <div class="flex items-center gap-3 min-w-0">
      <div class="w-9 h-9 rounded-xl bg-slate-100 dark:bg-slate-800 text-slate-600 dark:text-slate-300 flex items-center justify-center shrink-0 group-hover:bg-orange-500/10 group-hover:text-orange-500 transition-colors">
        <component :is="icon" class="w-4 h-4" />
      </div>
      <span class="text-xs sm:text-sm font-semibold text-slate-900 dark:text-slate-100 truncate">
        {{ label }}
      </span>
    </div>

    <div class="flex items-center gap-2">
      <span
        v-if="badgeText"
        :class="[
          'px-2 py-0.5 rounded-full text-[10px] font-bold',
          badgeColorClass || 'bg-orange-500/10 text-orange-500'
        ]"
      >
        {{ badgeText }}
      </span>
      <ChevronLeft class="w-4 h-4 text-slate-400 group-hover:text-orange-500 transition-colors" />
    </div>
  </button>
</template>
