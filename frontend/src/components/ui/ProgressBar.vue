<script setup lang="ts">
import { computed } from "vue";

const props = withDefaults(
  defineProps<{
    percentage: number;
    showLabel?: boolean;
    heightClass?: string;
  }>(),
  {
    showLabel: false,
    heightClass: "h-2",
  }
);

const clampedPercentage = computed(() => {
  return Math.min(100, Math.max(0, props.percentage));
});
</script>

<template>
  <div class="w-full flex flex-col gap-1.5">
    <div v-if="showLabel" class="flex justify-between items-center text-xs text-[#6B7280] dark:text-[#A1A1AA] font-medium">
      <span>نسبة الإنجاز</span>
      <span class="text-orange-500 font-bold">%{{ clampedPercentage }}</span>
    </div>
    <div class="w-full bg-[#E5E7EB] dark:bg-[#27272A] rounded-full overflow-hidden" :class="heightClass">
      <div
        class="bg-orange-500 h-full rounded-full transition-all duration-500 ease-out"
        :style="{ width: `${clampedPercentage}%` }"
      ></div>
    </div>
  </div>
</template>
