<script setup lang="ts">
import { computed } from "vue";
import { Course } from "@/types/course";
import { BookOpen, Clock, Users, User, Play, CheckCircle2 } from "@lucide/vue";
import ProgressBar from "./ProgressBar.vue";
import PriceBadge from "./PriceBadge.vue";

const props = defineProps<{
  course: Course;
}>();

const emit = defineEmits<{
  (e: "click"): void;
  (e: "action", courseId: string): void;
}>();

const formattedDuration = computed(() => {
  const hours = Math.floor(props.course.durationMinutes / 60);
  const minutes = props.course.durationMinutes % 60;
  if (hours > 0 && minutes > 0) return `${hours}س ${minutes}د`;
  if (hours > 0) return `${hours} ساعة`;
  return `${minutes} دقيقة`;
});
</script>

<template>
  <div
    @click="emit('click'); emit('action', course.id)"
    data-glass="card"
    class="glass-card group flex flex-col rounded-xl overflow-hidden hover:border-orange-500/50 transition-all duration-300 cursor-pointer"
  >
    <!-- Cover Image Container -->
    <div class="relative w-full aspect-[16/9] bg-[#F8F8F8] dark:bg-[#222222] overflow-hidden">
      <img
        :src="course.coverImageUrl"
        :alt="course.title"
        loading="lazy"
        class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500 ease-out"
      />
      <div class="absolute inset-0 bg-gradient-to-t from-black/60 via-transparent to-transparent opacity-60"></div>
      
      <!-- Price / Enrollment Badge -->
      <div class="absolute top-2 right-2">
        <PriceBadge
          :price="course.price"
          :currency="course.currency"
          :is-free="course.isFree"
          :enrollment-status="course.enrollmentStatus"
        />
      </div>
    </div>

    <!-- Details Section -->
    <div class="p-2.5 flex flex-col flex-1 justify-between gap-2 text-start">
      <div>
        <div class="flex items-center gap-1 mb-0.5">
          <span class="text-[9px] sm:text-[10px] font-bold uppercase text-orange-500 tracking-wider">
            {{ course.gradeName || course.departmentName || 'مساق تعليمي' }}
          </span>
          <span v-if="course.gradeName && course.departmentName" class="text-[9px] text-[#6B7280] dark:text-[#A1A1AA] font-medium">
            • {{ course.departmentName }}
          </span>
        </div>
        <h3 class="text-[11px] sm:text-xs font-bold text-[#111111] dark:text-[#FFFFFF] line-clamp-2 leading-snug group-hover:text-orange-500 transition-colors">
          {{ course.title }}
        </h3>
      </div>

      <!-- Instructor Info -->
      <div class="flex items-center gap-1.5 pt-1 border-t border-[#E5E7EB] dark:border-[#27272A]">
        <div class="w-4 h-4 rounded-full bg-[#F8F8F8] dark:bg-[#222222] overflow-hidden flex items-center justify-center shrink-0">
          <img
            v-if="course.instructor.avatarUrl"
            :src="course.instructor.avatarUrl"
            :alt="course.instructor.name"
            loading="lazy"
            class="w-full h-full object-cover"
          />
          <User v-else class="w-2.5 h-2.5 text-[#6B7280] dark:text-[#A1A1AA]" />
        </div>
        <span class="text-[10px] text-[#6B7280] dark:text-[#A1A1AA] font-medium truncate">
          {{ course.instructor.name }}
        </span>
      </div>

      <!-- Stats Line -->
      <div class="flex items-center justify-between text-[9px] sm:text-[10px] text-[#6B7280] dark:text-[#A1A1AA]">
        <div class="flex items-center gap-2">
          <span class="flex items-center gap-0.5">
            <BookOpen class="w-2.5 h-2.5 text-[#9CA3AF] dark:text-[#71717A]" />
            <span>{{ course.lecturesCount }}</span>
          </span>
          <span class="flex items-center gap-0.5">
            <Clock class="w-2.5 h-2.5 text-[#9CA3AF] dark:text-[#71717A]" />
            <span>{{ formattedDuration }}</span>
          </span>
        </div>

        <span class="flex items-center gap-0.5 text-[#9CA3AF] dark:text-[#71717A]">
          <Users class="w-2.5 h-2.5" />
          <span>{{ course.studentsCount }}</span>
        </span>
      </div>

      <!-- State Action / Progress Footer -->
      <div v-if="course.enrollmentStatus !== 'not_enrolled'" class="pt-1.5 border-t border-[#E5E7EB] dark:border-[#27272A] flex items-center justify-between gap-2">
        <div class="flex-1">
          <ProgressBar :percentage="course.completionPercentage || (course.enrollmentStatus === 'completed' ? 100 : 0)" show-label />
        </div>

        <button
          v-if="course.enrollmentStatus === 'completed'"
          @click.stop="emit('action', course.id)"
          type="button"
          class="inline-flex items-center gap-1 px-2.5 py-1 rounded-lg bg-orange-500/10 text-orange-500 hover:bg-orange-500/20 font-bold text-[10px] transition-colors cursor-pointer shrink-0"
        >
          <CheckCircle2 class="w-3 h-3" />
          <span>إعادة</span>
        </button>

        <button
          v-else
          @click.stop="emit('action', course.id)"
          type="button"
          class="inline-flex items-center gap-1 px-2.5 py-1 rounded-lg bg-orange-500 hover:bg-[#EA580C] text-white font-bold text-[10px] shadow-sm shadow-orange-500/20 transition-all shrink-0 active:scale-95 cursor-pointer"
        >
          <Play class="w-2 h-2 fill-white" />
          <span>متابعة</span>
        </button>
      </div>
    </div>
  </div>
</template>
