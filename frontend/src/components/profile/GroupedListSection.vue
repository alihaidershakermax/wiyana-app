<script setup lang="ts">
defineProps<{
  title?: string;
}>();
</script>

<template>
  <div class="flex flex-col gap-2 w-full">
    <!-- iOS Section Header Label -->
    <h3
      v-if="title"
      class="text-xs font-bold text-slate-500 dark:text-zinc-400 px-3 tracking-wide select-none"
    >
      {{ title }}
    </h3>

    <!-- Grouped Card Container -->
    <div
      class="w-full bg-white dark:bg-[#141414] border border-slate-200/80 dark:border-zinc-800/80 rounded-2xl sm:rounded-3xl overflow-hidden shadow-[0_1px_3px_rgba(0,0,0,0.04)] divide-y divide-slate-100 dark:divide-zinc-800/60 transition-colors"
    >
      <slot />
    </div>
  </div>
</template>
