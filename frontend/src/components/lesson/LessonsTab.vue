<script setup lang="ts">
import { ListVideo } from "@lucide/vue";
import LessonItem from "@/components/lesson/LessonItem.vue";
import EmptyState from "@/components/ui/EmptyState.vue";
import type { Lesson } from "@/types/lesson";

defineProps<{
  lessons: Lesson[];
  currentLessonId: string;
  isEnrolled: boolean;
}>();

const emit = defineEmits<{
  (e: "select", lessonId: string): void;
}>();
</script>

<template>
  <div v-if="lessons.length === 0" v-motion :initial="{ opacity: 0, y: 12 }" :enter="{ opacity: 1, y: 0 }">
    <EmptyState
      title="لا توجد محاضرات"
      description="لم تتم إضافة محاضرات لهذا الكورس بعد."
    />
  </div>

  <div v-else class="flex flex-col gap-2.5">
    <div
      v-for="(lesson, index) in lessons"
      :key="lesson.id"
      v-motion
      :initial="{ opacity: 0, y: 10 }"
      :enter="{ opacity: 1, y: 0, transition: { delay: Math.min(index * 40, 300) } }"
    >
      <LessonItem
        :lesson="lesson"
        :is-current="lesson.id === currentLessonId"
        :is-enrolled="isEnrolled"
        @select="(id) => emit('select', id)"
      />
    </div>
  </div>
</template>
