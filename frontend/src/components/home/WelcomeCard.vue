<script setup lang="ts">
import { Play, BookOpen, TrendingUp, GraduationCap } from "@lucide/vue";
import { UserProfile } from "@/types/user";

defineProps<{
  user?: UserProfile;
}>();

const emit = defineEmits<{
  (e: "clickContinue"): void;
}>();
</script>

<template>
  <div class="relative w-full rounded-2xl sm:rounded-3xl overflow-hidden bg-[#181818] border border-[#27272A] p-6 sm:p-7 text-white shadow-sm text-start">
    <div class="relative z-10 flex flex-col md:flex-row md:items-center justify-between gap-6">
      <!-- Greeting & Info -->
      <div class="flex flex-col gap-2 max-w-lg">
        <div class="inline-flex items-center gap-1.5 px-3 py-1 rounded-lg bg-orange-500/10 border border-orange-500/20 text-orange-400 text-xs font-bold w-max">
          <GraduationCap class="w-4 h-4" />
          <span>منصة ويانا التعليمية</span>
        </div>

        <h1 class="text-xl sm:text-2xl font-bold tracking-tight text-white mt-1">
          أهلاً بك مجدداً، {{ user?.name || "طالبنا العزيز" }}
        </h1>
        <p class="text-xs sm:text-sm text-[#A1A1AA] font-normal leading-relaxed">
          سعداء برؤيتك مجدداً. واصل خطتك الدراسية وحقق هدفك الأكاديمي اليوم.
        </p>

        <!-- Stats Chips -->
        <div class="flex items-center gap-3 mt-2">
          <div class="flex items-center gap-2.5 bg-[#222222] px-3.5 py-2 rounded-xl border border-[#27272A]">
            <div class="w-7 h-7 rounded-lg bg-orange-500/20 text-orange-400 flex items-center justify-center">
              <BookOpen class="w-4 h-4" />
            </div>
            <div class="flex flex-col">
              <span class="text-sm font-bold text-white">{{ user?.enrolledCoursesCount || 0 }}</span>
              <span class="text-[10px] text-[#A1A1AA] font-medium">مساقات مشتركة</span>
            </div>
          </div>

          <div class="flex items-center gap-2.5 bg-[#222222] px-3.5 py-2 rounded-xl border border-[#27272A]">
            <div class="w-7 h-7 rounded-lg bg-orange-500/20 text-orange-400 flex items-center justify-center">
              <TrendingUp class="w-4 h-4" />
            </div>
            <div class="flex flex-col">
              <span class="text-sm font-bold text-white">%{{ user?.overallCompletionPercentage || 0 }}</span>
              <span class="text-[10px] text-[#A1A1AA] font-medium">نسبة الإنجاز</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Action Button -->
      <div class="flex items-center">
        <button
          @click="emit('clickContinue')"
          type="button"
          class="w-full sm:w-auto inline-flex items-center justify-center gap-2 px-5 py-3 rounded-xl bg-orange-500 hover:bg-[#EA580C] active:scale-95 text-white font-bold text-sm shadow-sm shadow-orange-500/20 transition-all duration-150 focus:outline-none cursor-pointer"
        >
          <Play class="w-4 h-4 fill-white" />
          <span>متابعة التعلم</span>
        </button>
      </div>
    </div>
  </div>
</template>
