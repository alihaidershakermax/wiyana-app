<script setup lang="ts">
import { computed } from "vue";
import { Clock, CalendarDays, GraduationCap, BookOpen } from "@lucide/vue";
import { formatDuration, formatDate } from "@/utils/format";
import type { Lesson, LessonInfoCourse } from "@/types/lesson";

const props = defineProps<{
  lesson: Lesson;
  course: LessonInfoCourse;
}>();

const durationText = computed(() => formatDuration(props.lesson.durationSeconds));
const addedDateText = computed(() => formatDate(props.lesson.createdAt));
</script>

<template>
  <div class="flex flex-col gap-1.5">
    <h2 class="text-base sm:text-lg font-black text-slate-900 dark:text-white leading-snug">
      {{ lesson.title }}
    </h2>

    <div class="flex flex-col gap-1 text-xs text-slate-500 dark:text-slate-400">
      <span class="inline-flex items-center gap-1.5 font-bold text-slate-700 dark:text-slate-200">
        <BookOpen class="w-3.5 h-3.5 text-orange-500 shrink-0" />
        <span class="truncate">{{ course.title }}</span>
      </span>

      <span v-if="course.subjectName || course.instructorName" class="inline-flex items-center gap-1.5">
        <GraduationCap class="w-3.5 h-3.5 text-orange-500 shrink-0" />
        <span class="truncate">
          <template v-if="course.subjectName">{{ course.subjectName }}</template>
          <template v-if="course.subjectName && course.instructorName"> · </template>
          <template v-if="course.instructorName">{{ course.instructorName }}</template>
        </span>
      </span>

      <span class="inline-flex items-center gap-3">
        <span v-if="durationText" class="inline-flex items-center gap-1">
          <Clock class="w-3.5 h-3.5 text-orange-500 shrink-0" />
          <span>{{ durationText }}</span>
        </span>
        <span v-if="addedDateText" class="inline-flex items-center gap-1">
          <CalendarDays class="w-3.5 h-3.5 text-orange-500 shrink-0" />
          <span>{{ addedDateText }}</span>
        </span>
      </span>
    </div>
  </div>
</template>
