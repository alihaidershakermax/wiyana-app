<script setup lang="ts">
import type { Component } from "vue";

defineProps<{
  label: string;
  value: string | number;
  subValue?: string;
  icon: Component;
  colorClass?: string;
}>();

const emit = defineEmits<{
  (e: "click"): void;
}>();
</script>

<template>
  <div
    @click="emit('click')"
    class="bg-white dark:bg-slate-900 border border-slate-200/80 dark:border-slate-800/80 rounded-2xl p-4 flex items-center gap-3.5 shadow-xs hover:border-orange-500/40 hover:shadow-md transition-all cursor-pointer group active:scale-98"
  >
    <div
      :class="[
        'w-11 h-11 rounded-2xl flex items-center justify-center shrink-0 transition-transform group-hover:scale-110 duration-200',
        colorClass || 'bg-orange-500/10 text-orange-500'
      ]"
    >
      <component :is="icon" class="w-5 h-5" />
    </div>

    <div class="flex flex-col">
      <span class="text-xs text-slate-500 dark:text-slate-400 font-medium group-hover:text-slate-700 dark:group-hover:text-slate-300 transition-colors">
        {{ label }}
      </span>
      <div class="flex items-baseline gap-1 mt-0.5">
        <span class="text-lg font-black text-slate-900 dark:text-slate-100 font-display">
          {{ value }}
        </span>
        <span v-if="subValue" class="text-[10px] font-bold text-slate-400">
          {{ subValue }}
        </span>
      </div>
    </div>
  </div>
</template>
