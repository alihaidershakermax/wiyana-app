<script setup lang="ts">
import { X, FileText, CheckCircle2, ShieldCheck } from "@lucide/vue";

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
        class="w-full max-w-sm sm:max-w-md bg-white dark:bg-[#141414] border border-slate-200/80 dark:border-zinc-800/80 rounded-3xl p-6 shadow-2xl flex flex-col gap-4 my-8 text-start"
      >
        <!-- Header -->
        <div class="flex items-center justify-between border-b border-slate-100 dark:border-zinc-800/80 pb-3.5">
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 rounded-2xl bg-orange-500/10 text-orange-500 flex items-center justify-center">
              <FileText class="w-5 h-5" />
            </div>
            <div>
              <h3 class="text-base font-black text-slate-900 dark:text-zinc-100">
                شروط الاستخدام
              </h3>
              <p class="text-[11px] text-slate-500 dark:text-zinc-400">
                قواعد وأحكام المنصة
              </p>
            </div>
          </div>

          <button
            @click="emit('close')"
            type="button"
            class="p-2 rounded-xl text-slate-400 hover:text-slate-600 dark:hover:text-zinc-200 hover:bg-slate-100 dark:hover:bg-zinc-800 transition-colors cursor-pointer"
          >
            <X class="w-4 h-4" />
          </button>
        </div>

        <!-- Content -->
        <div class="flex flex-col gap-2.5 text-xs text-slate-600 dark:text-zinc-300 leading-relaxed max-h-80 overflow-y-auto pr-1">
          <div class="p-3.5 rounded-2xl bg-slate-50 dark:bg-zinc-800/40 border border-slate-100 dark:border-zinc-800/60">
            <h4 class="font-bold text-slate-900 dark:text-zinc-100 mb-0.5">1. حقوق الملكية الفكرية</h4>
            <p class="text-slate-500 dark:text-zinc-400">جميع المحاضرات والملفات ملكية حصرية لمنصة ويانا وأساتذتها، ويحظر إعادة نشرها أو تسجيلها.</p>
          </div>

          <div class="p-3.5 rounded-2xl bg-slate-50 dark:bg-zinc-800/40 border border-slate-100 dark:border-zinc-800/60">
            <h4 class="font-bold text-slate-900 dark:text-zinc-100 mb-0.5">2. استخدام كود الاشتراك</h4>
            <p class="text-slate-500 dark:text-zinc-400">كود الاشتراك مخصص للاستخدام الفردي فقط على جهاز الطالب المسجل لمنع حظر الحساب.</p>
          </div>

          <div class="p-3.5 rounded-2xl bg-slate-50 dark:bg-zinc-800/40 border border-slate-100 dark:border-zinc-800/60">
            <h4 class="font-bold text-slate-900 dark:text-zinc-100 mb-0.5">3. الالتزام بالاستخدام الأكاديمي</h4>
            <p class="text-slate-500 dark:text-zinc-400">يلتزم الطالب باستخدام المنصة للأغراض التعليمية ومتابعة تقدمه الدراسي بانتظام.</p>
          </div>
        </div>

        <!-- Footer -->
        <div class="pt-2 border-t border-slate-100 dark:border-zinc-800/80">
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
