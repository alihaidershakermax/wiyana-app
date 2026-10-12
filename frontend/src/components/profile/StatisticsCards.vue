<script setup lang="ts">
import { UserStatistics } from "@/types/user";
import { BookOpen, CheckCircle2, Clock, Award } from "@lucide/vue";
import StatisticCard from "./StatisticCard.vue";

defineProps<{
  stats?: UserStatistics;
}>();

const emit = defineEmits<{
  (e: "clickCourses"): void;
  (e: "clickLectures"): void;
  (e: "clickWatchTime"): void;
  (e: "clickCompletion"): void;
}>();
</script>

<template>
  <div class="grid grid-cols-2 lg:grid-cols-4 gap-3 w-full">
      <StatisticCard
        label="المساقات المشتركة"
        :value="stats?.enrolledCoursesCount || 0"
        sub-value="مساقات"
        :icon="BookOpen"
        color-class="bg-orange-500/10 text-orange-500"
        @click="emit('clickCourses')"
      />

      <StatisticCard
        label="المحاضرات المكتملة"
        :value="stats?.completedLecturesCount || 0"
        sub-value="محاضرة"
        :icon="CheckCircle2"
        color-class="bg-orange-500/10 text-orange-500"
        @click="emit('clickLectures')"
      />

      <StatisticCard
        label="ساعات المشاهدة"
        :value="stats?.watchTimeHours || 0"
        sub-value="ساعة"
        :icon="Clock"
        color-class="bg-orange-500/10 text-orange-500"
        @click="emit('clickWatchTime')"
      />

      <StatisticCard
        label="نسبة الإنجاز العامة"
        :value="`%${stats?.overallCompletionPercentage || 0}`"
        :icon="Award"
        color-class="bg-orange-500/10 text-orange-500"
        @click="emit('clickCompletion')"
      />
  </div>
</template>
