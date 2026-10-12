<script setup lang="ts">
import { BookOpen } from "@lucide/vue";

withDefaults(
  defineProps<{
    title?: string;
    description?: string;
    actionText?: string;
  }>(),
  {
    title: "لا توجد مساقات متاحة حالياً",
    description: "لم نجد أي نتائج في هذا التصنيف. يمكنك تجربة اختيار تصنيف آخر.",
  }
);

const emit = defineEmits<{
  (e: "clickAction"): void;
}>();
</script>

<template>
  <div class="flex flex-col items-center justify-center text-center p-8 bg-white dark:bg-[#181818] border border-[#E5E7EB] dark:border-[#27272A] rounded-3xl my-6">
    <div class="w-16 h-16 rounded-2xl bg-orange-500/10 dark:bg-orange-500/15 flex items-center justify-center text-orange-500 mb-4">
      <slot name="icon">
        <BookOpen class="w-8 h-8" />
      </slot>
    </div>
    <h3 class="text-base sm:text-lg font-bold text-[#111111] dark:text-[#FFFFFF]">
      {{ title }}
    </h3>
    <p class="text-xs sm:text-sm text-[#6B7280] dark:text-[#A1A1AA] max-w-sm mt-1 mb-5">
      {{ description }}
    </p>
    <button
      v-if="actionText"
      @click="emit('clickAction')"
      type="button"
      class="px-5 py-2.5 rounded-xl bg-orange-500 hover:bg-[#EA580C] text-white font-bold text-xs sm:text-sm shadow-sm shadow-orange-500/20 active:scale-95 transition-all cursor-pointer"
    >
      {{ actionText }}
    </button>
  </div>
</template>
