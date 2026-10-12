<script setup lang="ts">
import { SlidersHorizontal, Search, BookOpen } from "@lucide/vue";

defineProps<{
  totalCount: number;
  activeFiltersCount: number;
  isSearchVisible?: boolean;
}>();

const emit = defineEmits<{
  (e: "toggleSearch"): void;
  (e: "openFilters"): void;
}>();
</script>

<template>
  <header class="w-full bg-white/90 dark:bg-[#0F0F0F]/90 backdrop-blur-xl border-b border-[#E5E7EB] dark:border-[#27272A] sticky top-0 z-30 transition-colors">
    <div class="max-w-7xl mx-auto px-4 py-3.5 flex items-center justify-between gap-4">
      <!-- Title & Courses Count Badge -->
      <div class="flex items-center gap-3">
        <h1 class="text-xl sm:text-2xl font-black text-[#111111] dark:text-[#FFFFFF] tracking-tight font-display">
          المساقات
        </h1>
        <span class="inline-flex items-center gap-1 px-3 py-1 rounded-full text-xs font-bold bg-orange-500/10 text-orange-500 border border-orange-500/20">
          <BookOpen class="w-3.5 h-3.5" />
          <span>{{ totalCount }} مساق</span>
        </span>
      </div>

      <!-- Action Buttons -->
      <div class="flex items-center gap-2">
        <!-- Search Toggle Button -->
        <button
          @click="emit('toggleSearch')"
          type="button"
          aria-label="البحث"
          :class="[
            'w-12 h-12 rounded-2xl flex items-center justify-center transition-all active:scale-90 focus:outline-none cursor-pointer',
            isSearchVisible
              ? 'bg-orange-500 text-white shadow-sm shadow-orange-500/20'
              : 'bg-[#F8F8F8] dark:bg-[#181818] hover:bg-[#E5E7EB] dark:hover:bg-[#222222] border border-[#E5E7EB] dark:border-[#27272A] text-[#6B7280] dark:text-[#A1A1AA]'
          ]"
        >
          <Search class="w-5 h-5" />
        </button>

        <!-- Filter Sheet Button -->
        <button
          @click="emit('openFilters')"
          type="button"
          aria-label="الفلاتر"
          class="relative w-12 h-12 rounded-2xl bg-[#F8F8F8] dark:bg-[#181818] hover:bg-[#E5E7EB] dark:hover:bg-[#222222] border border-[#E5E7EB] dark:border-[#27272A] text-[#6B7280] dark:text-[#A1A1AA] flex items-center justify-center transition-all active:scale-90 focus:outline-none cursor-pointer"
        >
          <SlidersHorizontal class="w-5 h-5" />
          <span
            v-if="activeFiltersCount > 0"
            class="absolute -top-1 -right-1 w-5.5 h-5.5 rounded-full bg-orange-500 text-white text-[10px] font-black flex items-center justify-center ring-2 ring-white dark:ring-[#181818]"
          >
            {{ activeFiltersCount }}
          </span>
        </button>
      </div>
    </div>
  </header>
</template>
