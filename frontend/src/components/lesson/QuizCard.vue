<script setup lang="ts">
import { computed } from "vue";
import { FileQuestion, ListChecks, Play, Clock } from "@lucide/vue";
import type { LessonQuiz } from "@/types/lesson";

const props = defineProps<{
  quiz: LessonQuiz;
}>();

const emit = defineEmits<{
  (e: "start"): void;
}>();

const durationText = computed(() => {
  if (!props.quiz.durationMinutes) return null;
  return `${props.quiz.durationMinutes} دقيقة`;
});
</script>

<template>
  <div class="bg-white dark:bg-[#181818] border border-[#E5E7EB] dark:border-[#27272A] rounded-2xl p-5 text-start">
    <div class="flex items-center gap-3.5 mb-4">
      <div class="w-12 h-12 rounded-2xl bg-orange-500/10 text-orange-500 flex items-center justify-center shrink-0">
        <FileQuestion class="w-6 h-6" />
      </div>
      <div class="flex flex-col min-w-0">
        <span class="text-xs font-bold text-orange-500">اختبار المحاضرة</span>
        <h3 class="text-sm font-black text-[#111111] dark:text-[#FFFFFF] truncate">{{ quiz.title }}</h3>
      </div>
    </div>

    <div class="flex items-center gap-4 text-[11px] font-bold text-[#6B7280] dark:text-[#A1A1AA] mb-5">
      <span class="inline-flex items-center gap-1.5">
        <ListChecks class="w-4 h-4 text-orange-500" />
        {{ quiz.questionsCount }} سؤال
      </span>
      <span v-if="durationText" class="inline-flex items-center gap-1.5">
        <Clock class="w-4 h-4 text-orange-500" />
        {{ durationText }}
      </span>
    </div>

    <button
      @click="emit('start')"
      type="button"
      class="w-full h-12 inline-flex items-center justify-center gap-2 rounded-2xl bg-orange-500 hover:bg-[#EA580C] text-white font-bold text-sm shadow-sm shadow-orange-500/20 active:scale-95 transition-all cursor-pointer"
    >
      <Play class="w-4 h-4 fill-white" />
      بدء الاختبار
    </button>
  </div>
</template>
