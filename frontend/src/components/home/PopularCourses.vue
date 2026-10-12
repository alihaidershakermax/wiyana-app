<script setup lang="ts">
import { Course } from "@/types/course";
import SectionHeader from "@/components/ui/SectionHeader.vue";
import CourseCard from "@/components/ui/CourseCard.vue";

defineProps<{
  courses: Course[];
}>();

const emit = defineEmits<{
  (e: "selectCourse", courseId: string): void;
  (e: "viewAll"): void;
}>();
</script>

<template>
  <div v-if="courses && courses.length > 0" class="flex flex-col gap-4">
    <SectionHeader
      title="الأكثر مشاهدة وإقبالاً"
      subtitle="المساقات الأعلى تقييماً وإقبالاً من الطلاب"
      action-text="عرض الكل"
      @click-action="emit('viewAll')"
    />

    <div class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-3">
      <CourseCard
        v-for="course in courses"
        :key="course.id"
        :course="course"
        @click="emit('selectCourse', course.id)"
        @action="(id) => emit('selectCourse', id)"
        class="cursor-pointer"
      />
    </div>
  </div>
</template>
