<script setup lang="ts">
import { Course } from "@/types/course";
import CourseCard from "@/components/ui/CourseCard.vue";

defineProps<{
  courses: Course[];
}>();

const emit = defineEmits<{
  (e: "selectCourse", courseId: string): void;
  (e: "courseAction", courseId: string): void;
}>();
</script>

<template>
  <div class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-3 w-full">
    <CourseCard
      v-for="course in courses"
      :key="course.id"
      :course="course"
      @click="emit('selectCourse', course.id)"
      @action="(id) => emit('courseAction', id)"
      class="cursor-pointer"
    />
  </div>
</template>
