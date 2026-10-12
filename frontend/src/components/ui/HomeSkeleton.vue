<script setup lang="ts">
import SkeletonCard from "./SkeletonCard.vue";
</script>

<template>
  <div class="flex flex-col gap-6 animate-pulse p-4 max-w-7xl mx-auto w-full">
    <!-- Header Skeleton -->
    <div class="flex justify-between items-center py-3 border-b border-slate-200 dark:border-slate-800">
      <div class="flex items-center gap-3">
        <div class="w-10 h-10 rounded-xl bg-slate-200 dark:bg-slate-800"></div>
        <div class="h-5 w-24 bg-slate-200 dark:bg-slate-800 rounded"></div>
      </div>
      <div class="flex items-center gap-2">
        <div class="w-9 h-9 rounded-xl bg-slate-200 dark:bg-slate-800"></div>
        <div class="w-9 h-9 rounded-xl bg-slate-200 dark:bg-slate-800"></div>
        <div class="w-9 h-9 rounded-full bg-slate-200 dark:bg-slate-800"></div>
      </div>
    </div>

    <!-- Welcome Card Skeleton -->
    <div class="w-full h-44 rounded-3xl bg-slate-200 dark:bg-slate-900"></div>

    <!-- Continue Learning Skeleton -->
    <div class="flex flex-col gap-3">
      <div class="h-6 w-36 bg-slate-200 dark:bg-slate-800 rounded"></div>
      <div class="flex gap-4 overflow-hidden">
        <div v-for="i in 2" :key="i" class="w-72 sm:w-80 h-36 rounded-2xl bg-slate-200 dark:bg-slate-900 shrink-0"></div>
      </div>
    </div>

    <!-- Category Chips Skeleton -->
    <div class="flex gap-2">
      <div v-for="i in 5" :key="i" class="h-9 w-24 rounded-xl bg-slate-200 dark:bg-slate-800"></div>
    </div>

    <!-- Grid Skeleton -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
      <SkeletonCard v-for="i in 6" :key="i" />
    </div>
  </div>
</template>
