<script setup lang="ts">
import { computed } from "vue";
import { SkipBack, SkipForward } from "@lucide/vue";
import type { Lesson } from "@/types/lesson";

const props = defineProps<{
  lessons: Lesson[];
  currentLessonId: string;
}>();

const emit = defineEmits<{
  (e: "navigate", lessonId: string): void;
}>();

const currentIndex = computed(() =>
  props.lessons.findIndex((l) => l.id === props.currentLessonId)
);

const prevLesson = computed(() =>
  currentIndex.value > 0 ? props.lessons[currentIndex.value - 1] : null
);

const nextLesson = computed(() =>
  currentIndex.value >= 0 && currentIndex.value < props.lessons.length - 1
    ? props.lessons[currentIndex.value + 1]
    : null
);
</script>

<template>
  <div class="grid grid-cols-2 gap-3">
    <button
      @click="prevLesson && emit('navigate', prevLesson.id)"
      type="button"
      :disabled="!prevLesson"
      class="inline-flex items-center justify-center gap-2 h-12 rounded-2xl border border-[#E5E7EB] dark:border-[#27272A] bg-white dark:bg-[#181818] text-[#111111] dark:text-[#FFFFFF] font-bold text-xs active:scale-95 transition-all disabled:opacity-40 disabled:cursor-not-allowed cursor-pointer"
    >
      <SkipBack class="w-4 h-4" />
      <span>المحاضرة السابقة</span>
    </button>

    <button
      @click="nextLesson && emit('navigate', nextLesson.id)"
      type="button"
      :disabled="!nextLesson"
      class="inline-flex items-center justify-center gap-2 h-12 rounded-2xl bg-orange-500 hover:bg-[#EA580C] text-white font-bold text-xs shadow-sm shadow-orange-500/20 active:scale-95 transition-all disabled:opacity-40 disabled:cursor-not-allowed cursor-pointer"
    >
      <span>المحاضرة التالية</span>
      <SkipForward class="w-4 h-4" />
    </button>
  </div>
</template>
