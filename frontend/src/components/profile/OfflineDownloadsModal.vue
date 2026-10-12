<script setup lang="ts">
import { ref } from "vue";
import { X, DownloadCloud, Trash2, HardDrive } from "@lucide/vue";

defineProps<{
  isOpen: boolean;
}>();

const emit = defineEmits<{
  (e: "close"): void;
  (e: "cleared"): void;
}>();

interface DownloadItem {
  id: string;
  title: string;
  courseTitle: string;
  sizeMb: number;
  quality: string;
  downloadedAt: string;
}

const downloads = ref<DownloadItem[]>([]);

function removeDownload(id: string) {
  downloads.value = downloads.value.filter(d => d.id !== id);
}

function clearAllDownloads() {
  downloads.value = [];
  emit("cleared");
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
              <DownloadCloud class="w-5 h-5" />
            </div>
            <div>
              <h3 class="text-lg font-black text-slate-900 dark:text-slate-100">
                التنزيلات والمشاهدة بدون إنترنت
              </h3>
              <p class="text-xs text-slate-500 dark:text-slate-400">
                المحاضرات المحفوظة محلياً على جهازك لتشغيلها أوفلاين
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

        <!-- Storage Card -->
        <div class="p-4 rounded-2xl bg-slate-50 dark:bg-slate-950 border border-slate-100 dark:border-slate-800/60 flex items-center justify-between">
          <div class="flex items-center gap-3">
            <HardDrive class="w-5 h-5 text-orange-500" />
            <div>
              <span class="text-xs font-bold text-slate-900 dark:text-slate-100">المساحة المستخدمة للدروس:</span>
              <p class="text-[11px] text-slate-500 dark:text-slate-400">
                {{ downloads.reduce((acc, curr) => acc + curr.sizeMb, 0) }} ميغابايت ({{ downloads.length }} محاضرات)
              </p>
            </div>
          </div>
          <span class="px-2.5 py-1 rounded-full bg-orange-500/10 text-orange-500 text-[10px] font-black">
            مساحة متوفرة
          </span>
        </div>

        <!-- Downloads list -->
        <div v-if="downloads.length > 0" class="flex flex-col gap-3 max-h-80 overflow-y-auto pr-1">
          <div
            v-for="dl in downloads"
            :key="dl.id"
            class="p-3.5 rounded-2xl bg-slate-50 dark:bg-slate-950 border border-slate-100 dark:border-slate-800/60 flex items-center justify-between gap-3"
          >
            <div class="flex flex-col gap-1 min-w-0">
              <span class="text-[10px] font-bold text-orange-500 truncate">{{ dl.courseTitle }}</span>
              <h4 class="text-xs font-black text-slate-900 dark:text-slate-100 truncate">
                {{ dl.title }}
              </h4>
              <div class="flex items-center gap-2 text-[10px] text-slate-500 dark:text-slate-400">
                <span>{{ dl.sizeMb }} MB</span>
                <span>•</span>
                <span>{{ dl.quality }}</span>
                <span>•</span>
                <span>تم التنزيل {{ dl.downloadedAt }}</span>
              </div>
            </div>

            <div class="flex items-center gap-1.5 shrink-0">
              <button
                @click="removeDownload(dl.id)"
                type="button"
                class="p-2 rounded-xl text-slate-400 hover:text-rose-500 hover:bg-rose-50 dark:hover:bg-rose-950/30 transition-colors cursor-pointer"
              >
                <Trash2 class="w-4 h-4" />
              </button>
            </div>
          </div>
        </div>

        <div v-else class="py-10 flex flex-col items-center justify-center gap-2 text-center">
          <DownloadCloud class="w-10 h-10 text-slate-300 dark:text-slate-700" />
          <h4 class="text-xs font-bold text-slate-700 dark:text-slate-300">لا توجد محاضرات منزلة</h4>
          <p class="text-[11px] text-slate-400">يمكنك الضغط على زر التنزيل داخل أي محاضرة لمشاهدتها بدون إنترنت</p>
        </div>

        <!-- Footer -->
        <div class="flex items-center gap-3 pt-2 border-t border-slate-100 dark:border-slate-800/80">
          <button
            v-if="downloads.length > 0"
            type="button"
            @click="clearAllDownloads"
            class="py-3 px-4 rounded-2xl border border-rose-200 dark:border-rose-900/50 text-rose-600 dark:text-rose-400 font-bold text-xs hover:bg-rose-50 dark:hover:bg-rose-950/30 transition-colors flex items-center justify-center gap-1.5 shrink-0 cursor-pointer"
          >
            <Trash2 class="w-4 h-4" />
            حذف كافة التنزيلات
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
