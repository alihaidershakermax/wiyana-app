<script setup lang="ts">
defineProps<{
  title?: string;
}>();
</script>

<template>
  <header
    class="sticky top-0 z-30 w-full backdrop-blur-xl bg-white/85 dark:bg-[#0F0F0F]/85 border-b border-slate-200/80 dark:border-zinc-800/80 transition-colors"
  >
    <div class="max-w-xl mx-auto px-4 py-3.5 flex items-center justify-between">
      <h1 class="text-lg sm:text-xl font-black text-slate-900 dark:text-zinc-100 tracking-tight">
        {{ title || "حسابي" }}
      </h1>
      
      <!-- Brand badge on trailing corner -->
      <span class="text-[11px] font-black text-orange-500 tracking-wider uppercase bg-orange-500/10 px-2.5 py-1 rounded-full">
        WEYANA
      </span>
    </div>
  </header>
</template>
