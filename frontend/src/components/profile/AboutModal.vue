<script setup lang="ts">
import { X, Sparkles, BookOpen, CheckCircle2, ShieldCheck } from "@lucide/vue";

defineProps<{
  isOpen: boolean;
}>();

const emit = defineEmits<{
  (e: "close"): void;
}>();
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
      class="fixed inset-0 z-50 bg-black/70 backdrop-blur-md flex items-center justify-center p-4 overflow-y-auto"
      @click.self="emit('close')"
    >
      <div
        class="w-full max-w-sm bg-white dark:bg-[#141414] border border-slate-200/80 dark:border-zinc-800/80 rounded-3xl p-6 shadow-2xl flex flex-col items-center text-center gap-4 my-8"
      >
        <!-- App Icon / Logo -->
        <div
          class="w-16 h-16 rounded-2xl bg-orange-500 text-white flex items-center justify-center shadow-lg shadow-orange-500/25"
        >
          <BookOpen class="w-8 h-8 stroke-[2.2]" />
        </div>

        <div>
          <h3 class="text-xl font-black text-slate-900 dark:text-zinc-100 tracking-tight">
            منصة WEYANA
          </h3>
          <p class="text-xs font-bold text-orange-500 mt-0.5">
            الإصدار 1.0.0
          </p>
        </div>

        <p class="text-xs text-slate-600 dark:text-zinc-300 leading-relaxed max-w-xs">
          منصة تعليمية تساعد الطلاب على الوصول إلى مساقاتهم ومحاضراتهم ومتابعة تقدمهم الدراسي.
        </p>

        <!-- Feature Badges -->
        <div class="grid grid-cols-2 gap-2 w-full text-xs pt-1">
          <div
            class="p-3 rounded-2xl bg-slate-50 dark:bg-zinc-800/40 border border-slate-100 dark:border-zinc-800/60 flex flex-col items-center"
          >
            <span class="text-[10px] text-slate-400 dark:text-zinc-500 font-bold">المشاهدة</span>
            <span class="font-bold text-slate-800 dark:text-zinc-200">بث فائق الدقة</span>
          </div>

          <div
            class="p-3 rounded-2xl bg-slate-50 dark:bg-zinc-800/40 border border-slate-100 dark:border-zinc-800/60 flex flex-col items-center"
          >
            <span class="text-[10px] text-slate-400 dark:text-zinc-500 font-bold">التكامل</span>
            <span class="font-bold text-orange-500 flex items-center gap-1">
              <CheckCircle2 class="w-3.5 h-3.5" />
              <span>تيليجرام رسمي</span>
            </span>
          </div>
        </div>

        <!-- Close Button -->
        <div class="w-full pt-2 border-t border-slate-100 dark:border-zinc-800/80 mt-1">
          <button
            type="button"
            @click="emit('close')"
            class="w-full py-3 rounded-2xl bg-orange-500 hover:bg-orange-600 text-white font-bold text-xs sm:text-sm shadow-sm shadow-orange-500/25 transition-all text-center active:scale-95 cursor-pointer"
          >
            إغلاق
          </button>
        </div>
      </div>
    </div>
  </Transition>
</template>
