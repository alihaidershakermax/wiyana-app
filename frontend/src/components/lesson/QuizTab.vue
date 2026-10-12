<script setup lang="ts">
import { FileQuestion } from "@lucide/vue";
import QuizCard from "@/components/lesson/QuizCard.vue";
import EmptyState from "@/components/ui/EmptyState.vue";
import type { LessonQuiz } from "@/types/lesson";

defineProps<{
  quiz: LessonQuiz | null;
}>();
</script>

<template>
  <div v-motion :initial="{ opacity: 0, y: 12 }" :enter="{ opacity: 1, y: 0 }">
    <QuizCard v-if="quiz" :quiz="quiz" />
    <EmptyState
      v-else
      title="لا يوجد اختبار"
      description="لم يتم إضافة اختبار لهذه المحاضرة بعد."
    >
      <template #icon>
        <FileQuestion class="w-8 h-8" />
      </template>
    </EmptyState>
  </div>
</template>
