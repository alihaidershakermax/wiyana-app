<script setup lang="ts">
import { computed } from "vue";
import { useRouter } from "vue-router";
import { useQuery, useQueryClient } from "@tanstack/vue-query";
import { Bookmark, Play, Trash2 } from "@lucide/vue";
import { fetchMyFavorites, removeCourseFromFavorites } from "@/api/favorites";
import { useUserStore } from "@/stores/useUserStore";
import SkeletonCard from "@/components/ui/SkeletonCard.vue";
import EmptyState from "@/components/ui/EmptyState.vue";
import ErrorState from "@/components/ui/ErrorState.vue";

const router = useRouter();
const queryClient = useQueryClient();
const userStore = useUserStore();

const { data: rawFavorites, isLoading, isError, refetch } = useQuery({
  queryKey: ["favorites"],
  queryFn: fetchMyFavorites,
  staleTime: 1000 * 60 * 2,
});

const favorites = computed(() => {
  return Array.isArray(rawFavorites.value) ? rawFavorites.value : [];
});

async function handleRemove(courseId: string) {
  try {
    await removeCourseFromFavorites(courseId);
    userStore.showToast("تمت إزالة المساق من المحفوظات");
    queryClient.invalidateQueries({ queryKey: ["favorites"] });
  } catch {
    userStore.showToast("تعذر حذف المساق", "error");
  }
}

function handleOpenCourse(courseId: string) {
  router.push(`/courses/${courseId}`);
}
</script>

<template>
  <div class="min-h-screen bg-white dark:bg-[#0F0F0F] text-[#111111] dark:text-[#FFFFFF] flex flex-col font-sans transition-colors pb-28 text-start">
    <!-- Header -->
    <header class="sticky top-0 z-30 w-full backdrop-blur-md bg-white/95 dark:bg-[#0F0F0F]/95 border-b border-[#E5E7EB] dark:border-[#27272A]">
      <div class="max-w-7xl mx-auto px-4 py-3.5 flex items-center justify-between">
        <div class="flex items-center gap-2.5">
          <div class="w-10 h-10 rounded-2xl bg-orange-500/10 text-orange-500 flex items-center justify-center">
            <Bookmark class="w-5 h-5" />
          </div>
          <div>
            <h1 class="text-xl sm:text-2xl font-black text-[#111111] dark:text-[#FFFFFF]">
              المحفوظات
            </h1>
            <p class="text-xs text-[#6B7280] dark:text-[#A1A1AA]">
              المساقات والمحاضرات المفضلة لديك
            </p>
          </div>
        </div>
      </div>
    </header>

    <!-- Loading State -->
    <div v-if="isLoading" class="max-w-7xl mx-auto px-4 py-6 w-full grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4">
      <SkeletonCard v-for="i in 6" :key="i" />
    </div>

    <!-- Error State -->
    <div v-else-if="isError" class="max-w-7xl mx-auto px-4 py-12 w-full">
      <ErrorState title="تعذر تحميل المحفوظات" message="حدث خطأ أثناء الاتصال بالخادم." @retry="refetch()" />
    </div>

    <!-- Empty State -->
    <div v-else-if="favorites.length === 0" class="max-w-7xl mx-auto px-4 py-16 w-full">
      <EmptyState
        title="لا توجد مساقات محفوظة بعد"
        description="انقر على أيقونة الحفظ في أي مساق للرجوع إليه بسرعة ومتابعة دراستك في أي وقت."
      />
    </div>

    <!-- Favorites Grid -->
    <main v-else class="flex-1 max-w-7xl mx-auto px-4 py-6 w-full grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4">
      <div
        v-for="item in favorites"
        :key="item.id"
        class="bg-white dark:bg-[#181818] border border-[#E5E7EB] dark:border-[#27272A] rounded-3xl p-4 flex flex-col justify-between gap-3 shadow-xs hover:border-orange-500/40 transition-all group"
      >
        <div class="flex items-start justify-between gap-3">
          <div class="flex flex-col gap-1 min-w-0">
            <span class="text-[11px] font-bold text-orange-500">
              {{ item.stage || item.department || 'مساق تعليمي' }}
            </span>
            <h3 class="text-sm font-black text-[#111111] dark:text-[#FFFFFF] group-hover:text-orange-500 transition-colors line-clamp-2">
              {{ item.title }}
            </h3>
            <span class="text-xs text-[#6B7280] dark:text-[#A1A1AA] mt-1">
              {{ item.teacher_name || item.department || 'أستاذ المساق' }}
            </span>
          </div>

          <button
            @click="handleRemove(item.id)"
            type="button"
            title="إزالة من المحفوظات"
            class="p-2 rounded-xl text-[#9CA3AF] hover:text-rose-500 hover:bg-rose-50 dark:hover:bg-rose-950/30 transition-colors shrink-0 cursor-pointer"
          >
            <Trash2 class="w-4 h-4" />
          </button>
        </div>

        <div class="flex items-center justify-between pt-2 border-t border-[#E5E7EB] dark:border-[#27272A]">
          <span class="text-xs font-bold text-[#6B7280] dark:text-[#A1A1AA]">
            {{ item.total_lessons ? `${item.total_lessons} محاضرة` : 'شامل الشروحات' }}
          </span>

          <button
            @click="handleOpenCourse(item.id)"
            type="button"
            class="px-4 py-2 rounded-2xl bg-orange-500 hover:bg-[#EA580C] active:scale-95 text-white font-bold text-xs shadow-sm flex items-center gap-1.5 transition-all cursor-pointer"
          >
            <Play class="w-3.5 h-3.5 fill-current" />
            <span>عرض المساق</span>
          </button>
        </div>
      </div>
    </main>
  </div>
</template>
