<script setup lang="ts">
import { computed } from "vue";
import { CheckCircle2, Circle, TrendingUp } from "@lucide/vue";
import ProgressBar from "@/components/ui/ProgressBar.vue";
import { usePlayerStore } from "@/stores/usePlayerStore";
import type { Lesson } from "@/types/lesson";

const props = defineProps<{
  lessons: Lesson[];
}>();

const playerStore = usePlayerStore();

const totalCount = computed(() => props.lessons.length);

const completedCount = computed(
  () =>
    props.lessons.filter(
      (l) => l.isCompleted === true || playerStore.isCompleted(l.id)
    ).length
);

const remainingCount = computed(() => Math.max(0, totalCount.value - completedCount.value));

const percentage = computed(() => {
  if (!totalCount.value) return 0;
  return Math.round((completedCount.value / totalCount.value) * 100);
});
</script>

<template>
  <div class="bg-white dark:bg-[#181818] border border-[#E5E7EB] dark:border-[#27272A] rounded-2xl p-4 text-start">
    <div class="flex items-center justify-between mb-3">
      <span class="inline-flex items-center gap-1.5 text-xs font-bold text-[#111111] dark:text-[#FFFFFF]">
        <TrendingUp class="w-4 h-4 text-orange-500" />
        تقدمك في الكورس
      </span>
      <span class="text-xs font-black text-orange-500">%{{ percentage }}</span>
    </div>

    <ProgressBar :percentage="percentage" />

    <div class="grid grid-cols-2 gap-2 mt-3.5">
      <div class="flex items-center gap-2 rounded-xl bg-[#F8F8F8] dark:bg-[#222222] border border-[#E5E7EB] dark:border-[#27272A] px-3 py-2">
        <CheckCircle2 class="w-4 h-4 text-orange-500 shrink-0" />
        <span class="text-[10px] font-bold text-[#6B7280] dark:text-[#A1A1AA]">
          {{ completedCount }} مكتملة من {{ totalCount }}
        </span>
      </div>
      <div class="flex items-center gap-2 rounded-xl bg-[#F8F8F8] dark:bg-[#222222] border border-[#E5E7EB] dark:border-[#27272A] px-3 py-2">
        <Circle class="w-4 h-4 text-[#9CA3AF] dark:text-[#71717A] shrink-0" />
        <span class="text-[10px] font-bold text-[#6B7280] dark:text-[#A1A1AA]">
          {{ remainingCount }} متبقية
        </span>
      </div>
    </div>
  </div>
</template>
