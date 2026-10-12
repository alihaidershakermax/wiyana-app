<script setup lang="ts">
import { ref, computed } from "vue";
import { useRouter } from "vue-router";
import { useQuery, useQueryClient } from "@tanstack/vue-query";
import { X, Bookmark, Play, Trash2, Clock, BookOpen } from "@lucide/vue";
import { fetchMyFavorites, removeCourseFromFavorites } from "@/api/favorites";
import { useUserStore } from "@/stores/useUserStore";

const props = defineProps<{
  isOpen: boolean;
}>();

const emit = defineEmits<{
  (e: "close"): void;
}>();

const router = useRouter();
const queryClient = useQueryClient();
const userStore = useUserStore();

const { data: rawFavorites, isLoading, refetch } = useQuery({
  queryKey: ["favorites"],
  queryFn: fetchMyFavorites,
  enabled: computed(() => props.isOpen),
  staleTime: 1000 * 60,
});

const favorites = computed(() => {
  return Array.isArray(rawFavorites.value) ? rawFavorites.value : [];
});

async function handleRemove(courseId: string) {
  try {
    await removeCourseFromFavorites(courseId);
    userStore.showToast("تمت إزالة المساق من المحفوظات");
    queryClient.invalidateQueries({ queryKey: ["favorites"] });
    refetch();
  } catch {
    userStore.showToast("تعذر حذف المساق", "error");
  }
}

function handleOpenCourse(courseId: string) {
  emit("close");
  router.push(`/courses/${courseId}`);
}
</script>

<template>
  <Transition
    enter-active-class="transition duration-200 ease-out"
    enter-from-class="opacity-0 scale-95"
    enter-to-class="opacity-100 scale-100"
    leave-active-class="transition duration-150 ease-in"
    leave-from-class="opacity-100 scale-100"
    leave-to-class="opacity-0 scale-95"
  >
    <div
      v-if="isOpen"
      class="fixed inset-0 z-50 bg-slate-950/80 backdrop-blur-md flex items-center justify-center p-4 overflow-y-auto"
      @click.self="emit('close')"
    >
      <div class="w-full max-w-lg bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl p-6 shadow-2xl flex flex-col gap-4 my-8">
        <!-- Header -->
        <div class="flex items-center justify-between border-b border-slate-100 dark:border-slate-800/80 pb-4">
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 rounded-2xl bg-orange-500/10 text-orange-500 flex items-center justify-center">
              <Bookmark class="w-5 h-5" />
            </div>
            <div>
              <h3 class="text-lg font-black text-slate-900 dark:text-slate-100">
                المساقات والمحاضرات المحفوظة
              </h3>
              <p class="text-xs text-slate-500 dark:text-slate-400">
                المساقات التي قمت بحفظها للرجوع إليها لاحقاً
              </p>
            </div>
          </div>

          <button
            @click="emit('close')"
            type="button"
            class="p-2 rounded-xl text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors cursor-pointer"
          >
            <X class="w-5 h-5" />
          </button>
        </div>

        <!-- Loading -->
        <div v-if="isLoading" class="py-12 flex items-center justify-center">
          <span class="text-xs font-bold text-slate-400">جاري تحميل المحفوظات...</span>
        </div>

        <!-- List -->
        <div v-else-if="favorites.length > 0" class="flex flex-col gap-3 max-h-96 overflow-y-auto pr-1">
          <div
            v-for="item in favorites"
            :key="item.id"
            class="p-4 rounded-2xl bg-slate-50 dark:bg-slate-950 border border-slate-100 dark:border-slate-800/60 flex flex-col gap-3 hover:border-orange-500/30 transition-all shadow-xs"
          >
            <div class="flex items-start justify-between gap-3">
              <div class="flex flex-col gap-1 min-w-0">
                <span class="text-[11px] font-bold text-orange-500">{{ item.department || item.stage || item.gradeName || 'المساق التعليمي' }}</span>
                <h4 class="text-xs sm:text-sm font-black text-slate-900 dark:text-slate-100 truncate">
                  {{ item.title }}
                </h4>
                <div class="flex items-center gap-2 text-[11px] text-slate-500 dark:text-slate-400 mt-1">
                  <span class="flex items-center gap-1">
                    <Clock class="w-3 h-3 text-orange-500" />
                    {{ item.teacher_name || 'أستاذ المساق' }}
                  </span>
                </div>
              </div>

              <button
                @click="handleRemove(item.id)"
                type="button"
                title="إزالة من المحفوظات"
                class="p-2 rounded-xl text-slate-400 hover:text-rose-500 hover:bg-rose-50 dark:hover:bg-rose-950/30 transition-colors shrink-0 cursor-pointer"
              >
                <Trash2 class="w-4 h-4" />
              </button>
            </div>

            <div class="flex items-center justify-end gap-2">
              <button
                @click="handleOpenCourse(item.id)"
                type="button"
                class="px-4 py-2 rounded-xl bg-orange-500 hover:bg-orange-600 text-white font-bold text-xs shadow-sm flex items-center gap-1.5 shrink-0 active:scale-95 transition-all cursor-pointer"
              >
                <Play class="w-3 h-3 fill-current" />
                عرض المساق
              </button>
            </div>
          </div>
        </div>

        <!-- Empty state -->
        <div v-else class="py-12 flex flex-col items-center justify-center gap-2 text-center">
          <BookOpen class="w-12 h-12 text-slate-300 dark:text-slate-700" />
          <h4 class="text-sm font-bold text-slate-700 dark:text-slate-300">لا توجد مساقات محفوظة</h4>
          <p class="text-xs text-slate-400">احفظ المساقات والمحاضرات أثناء التصفح للوصول السريع إليها من هنا</p>
        </div>

        <!-- Close -->
        <div class="pt-2 border-t border-slate-100 dark:border-slate-800/80">
          <button
            type="button"
            @click="emit('close')"
            class="w-full py-3 rounded-2xl border border-slate-200 dark:border-slate-800 text-slate-700 dark:text-slate-300 font-bold text-xs sm:text-sm hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors cursor-pointer"
          >
            إغلاق
          </button>
        </div>
      </div>
    </div>
  </Transition>
</template>
