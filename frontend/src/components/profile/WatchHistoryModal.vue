<script setup lang="ts">
import { ref } from "vue";
import { useRouter } from "vue-router";
import { X, History, Play, Trash2, Clock, CheckCircle2 } from "@lucide/vue";

defineProps<{
  isOpen: boolean;
}>();

const emit = defineEmits<{
  (e: "close"): void;
}>();

const router = useRouter();

interface HistoryItem {
  id: string;
  courseId: string;
  lectureId: string;
  title: string;
  courseTitle: string;
  lastWatchedAt: string;
  stoppedAtTime: string;
  progressPercent: number;
}

const historyItems = ref<HistoryItem[]>([]);

function clearAllHistory() {
  historyItems.value = [];
}

function handleResume(courseId: string, lectureId: string) {
  emit("close");
  router.push(`/lessons/${lectureId}`);
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
              <History class="w-5 h-5" />
            </div>
            <div>
              <h3 class="text-lg font-black text-slate-900 dark:text-slate-100">
                سجل المشاهدة والاستماع
              </h3>
              <p class="text-xs text-slate-500 dark:text-slate-400">
                المحاضرات التي شاهدتها مؤخراً مع نقطة التوقف
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

        <!-- History list -->
        <div v-if="historyItems.length > 0" class="flex flex-col gap-3 max-h-96 overflow-y-auto pr-1">
          <div
            v-for="item in historyItems"
            :key="item.id"
            class="p-4 rounded-2xl bg-slate-50 dark:bg-slate-950 border border-slate-100 dark:border-slate-800/60 flex flex-col gap-3 hover:border-orange-500/30 transition-all"
          >
            <div class="flex items-start justify-between gap-3">
              <div class="flex flex-col gap-1">
                <span class="text-[11px] font-bold text-orange-500">{{ item.courseTitle }}</span>
                <h4 class="text-xs sm:text-sm font-black text-slate-900 dark:text-slate-100">
                  {{ item.title }}
                </h4>
                <div class="flex items-center gap-2 text-[11px] text-slate-500 dark:text-slate-400 mt-1">
                  <span class="flex items-center gap-1">
                    <Clock class="w-3 h-3 text-orange-500" />
                    توقف عند {{ item.stoppedAtTime }}
                  </span>
                  <span>•</span>
                  <span>{{ item.lastWatchedAt }}</span>
                </div>
              </div>

              <span
                v-if="item.progressPercent === 100"
                class="inline-flex items-center gap-1 px-2.5 py-1 rounded-full bg-emerald-500/10 text-emerald-500 text-[10px] font-black shrink-0"
              >
                <CheckCircle2 class="w-3 h-3" />
                مكتملة
              </span>
            </div>

            <div class="flex items-center gap-2">
              <div class="flex-1 bg-slate-200 dark:bg-slate-800 h-1.5 rounded-full overflow-hidden">
                <div
                  :class="['h-full rounded-full', item.progressPercent === 100 ? 'bg-emerald-500' : 'bg-orange-500']"
                  :style="{ width: `${item.progressPercent}%` }"
                ></div>
              </div>
              <button
                @click="handleResume(item.courseId, item.lectureId)"
                type="button"
                class="px-3.5 py-1.5 rounded-xl bg-orange-500 hover:bg-orange-600 text-white font-bold text-xs shadow-sm flex items-center gap-1.5 shrink-0 active:scale-95 cursor-pointer"
              >
                <Play class="w-3 h-3 fill-current" />
                {{ item.progressPercent === 100 ? 'إعادة المشاهدة' : 'استكمال' }}
              </button>
            </div>
          </div>
        </div>

        <!-- Empty -->
        <div v-else class="py-12 flex flex-col items-center justify-center gap-2 text-center">
          <History class="w-12 h-12 text-slate-300 dark:text-slate-700" />
          <h4 class="text-sm font-bold text-slate-700 dark:text-slate-300">سجل المشاهدة فارغ</h4>
          <p class="text-xs text-slate-400">ستظهر هنا أي محاضرة تبدأ بمشاهدتها</p>
        </div>

        <!-- Footer -->
        <div class="flex items-center gap-3 pt-2 border-t border-slate-100 dark:border-slate-800/80">
          <button
            v-if="historyItems.length > 0"
            type="button"
            @click="clearAllHistory"
            class="py-3 px-4 rounded-2xl border border-rose-200 dark:border-rose-900/50 text-rose-600 dark:text-rose-400 font-bold text-xs hover:bg-rose-50 dark:hover:bg-rose-950/30 transition-colors flex items-center justify-center gap-1.5 shrink-0 cursor-pointer"
          >
            <Trash2 class="w-4 h-4" />
            مسح السجل
          </button>
          <button
            type="button"
            @click="emit('close')"
            class="flex-1 py-3 rounded-2xl border border-slate-200 dark:border-slate-800 text-slate-700 dark:text-slate-300 font-bold text-xs sm:text-sm hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors cursor-pointer"
          >
            إغلاق
          </button>
        </div>
      </div>
    </div>
  </Transition>
</template>
