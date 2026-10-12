<script setup lang="ts">
import { Course } from "@/types/course";
import SectionHeader from "@/components/ui/SectionHeader.vue";
import CourseCard from "@/components/ui/CourseCard.vue";
import EmptyState from "@/components/ui/EmptyState.vue";

defineProps<{
  courses: Course[];
}>();

const emit = defineEmits<{
  (e: "selectCourse", courseId: string): void;
  (e: "viewAll"): void;
}>();
</script>

<template>
  <div class="flex flex-col gap-4">
    <SectionHeader
      title="المساقات المقترحة لك"
      subtitle="تم اختيارها بناءً على تخصصك واهتماماتك الدراسية"
      action-text="عرض الكل"
      @click-action="emit('viewAll')"
    />

    <div v-if="courses && courses.length > 0" class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-3">
      <CourseCard
        v-for="course in courses"
        :key="course.id"
        :course="course"
        @click="emit('selectCourse', course.id)"
        @action="(id) => emit('selectCourse', id)"
        class="cursor-pointer"
      />
    </div>

    <EmptyState v-else title="لا توجد مساقات مقترحة حالياً" />
  </div>
</template>
