<script setup lang="ts">
defineProps<{
  hasVideo?: boolean;
}>();
</script>

<template>
  <div class="min-h-screen bg-slate-50 dark:bg-slate-950 text-slate-900 dark:text-slate-100 flex flex-col font-sans pb-20">
    <!-- Skeleton Header -->
    <header class="sticky top-0 z-30 w-full backdrop-blur-xl bg-slate-50/90 dark:bg-slate-950/90 border-b border-slate-200/80 dark:border-slate-800/80 px-4 py-3 flex items-center gap-3">
      <div class="w-11 h-11 rounded-2xl bg-slate-200 dark:bg-slate-800 animate-pulse shrink-0"></div>
      <div class="flex flex-col gap-2 min-w-0 flex-1">
        <div class="h-2.5 w-24 rounded-full bg-slate-200 dark:bg-slate-800 animate-pulse"></div>
        <div class="h-3.5 w-3/4 rounded-full bg-slate-200 dark:bg-slate-800 animate-pulse"></div>
      </div>
      <div class="w-11 h-11 rounded-2xl bg-slate-200 dark:bg-slate-800 animate-pulse shrink-0"></div>
    </header>

    <!-- Skeleton Video Area -->
    <div v-if="hasVideo" class="w-full aspect-video bg-slate-900 flex items-center justify-center">
      <div class="w-16 h-16 rounded-full bg-slate-800 animate-pulse"></div>
    </div>
    <div v-else class="w-full aspect-video bg-slate-200 dark:bg-slate-900 animate-pulse"></div>

    <!-- Skeleton Content -->
    <main class="flex-1 max-w-7xl mx-auto px-4 py-5 w-full flex flex-col gap-5">
      <div class="flex flex-col gap-2.5">
        <div class="h-4 w-2/3 rounded-full bg-slate-200 dark:bg-slate-800 animate-pulse"></div>
        <div class="h-3 w-1/2 rounded-full bg-slate-200 dark:bg-slate-800 animate-pulse"></div>
      </div>

      <div class="bg-white dark:bg-slate-900 border border-slate-200/80 dark:border-slate-800/80 rounded-2xl p-4">
        <div class="h-2.5 w-32 rounded-full bg-slate-200 dark:bg-slate-800 animate-pulse"></div>
        <div class="mt-3 h-2 w-full rounded-full bg-slate-200 dark:bg-slate-800 animate-pulse"></div>
      </div>

      <div class="grid grid-cols-2 gap-3">
        <div v-for="i in 2" :key="i" class="h-12 rounded-2xl bg-slate-200 dark:bg-slate-800 animate-pulse"></div>
      </div>

      <div class="flex gap-4 border-b border-slate-200 dark:border-slate-800">
        <div v-for="i in 4" :key="i" class="flex-1 h-11 rounded-xl bg-slate-200 dark:bg-slate-800 animate-pulse"></div>
      </div>

      <div class="flex flex-col gap-3">
        <div v-for="i in 4" :key="i" class="h-16 rounded-2xl bg-slate-200 dark:bg-slate-800 animate-pulse"></div>
      </div>
    </main>
  </div>
</template>
