<script setup lang="ts">
import { X, ShieldCheck, Lock, EyeOff, Server } from "@lucide/vue";

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
              <ShieldCheck class="w-5 h-5" />
            </div>
            <div>
              <h3 class="text-base font-black text-slate-900 dark:text-zinc-100">
                سياسة الخصوصية
              </h3>
              <p class="text-[11px] text-slate-500 dark:text-zinc-400">
                حماية البيانات وسرية الحساب
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
          <div class="p-3.5 rounded-2xl bg-slate-50 dark:bg-zinc-800/40 border border-slate-100 dark:border-zinc-800/60 flex items-start gap-3">
            <Lock class="w-4 h-4 text-orange-500 shrink-0 mt-0.5" />
            <div>
              <h4 class="font-bold text-slate-900 dark:text-zinc-100 mb-0.5">1. تشفير البيانات والجلسات</h4>
              <p class="text-slate-500 dark:text-zinc-400">تشفير كامل لبيانات تسجيل الدخول عبر بروتوكولات أمان متقدمة لضمان خصوصية حسابك.</p>
            </div>
          </div>

          <div class="p-3.5 rounded-2xl bg-slate-50 dark:bg-zinc-800/40 border border-slate-100 dark:border-zinc-800/60 flex items-start gap-3">
            <EyeOff class="w-4 h-4 text-orange-500 shrink-0 mt-0.5" />
            <div>
              <h4 class="font-bold text-slate-900 dark:text-zinc-100 mb-0.5">2. عدم مشاركة البيانات</h4>
              <p class="text-slate-500 dark:text-zinc-400">نلتزم بعدم بيع أو مشاركة بياناتك الأكاديمية أو رقم هاتفك أو اسمك مع أي طرف ثالث.</p>
            </div>
          </div>

          <div class="p-3.5 rounded-2xl bg-slate-50 dark:bg-zinc-800/40 border border-slate-100 dark:border-zinc-800/60 flex items-start gap-3">
            <Server class="w-4 h-4 text-orange-500 shrink-0 mt-0.5" />
            <div>
              <h4 class="font-bold text-slate-900 dark:text-zinc-100 mb-0.5">3. أمان المحتوى والمشاهدة</h4>
              <p class="text-slate-500 dark:text-zinc-400">جميع المحاضرات والدروس محمية ضد التسريب لضمان استمرارية جودة التجربة التعليمية.</p>
            </div>
          </div>
        </div>

        <!-- Footer -->
        <div class="pt-2 border-t border-slate-100 dark:border-zinc-800/80">
          <button
            type="button"
            @click="emit('close')"
            class="w-full py-3 rounded-2xl bg-orange-500 hover:bg-orange-600 text-white font-bold text-xs sm:text-sm shadow-sm shadow-orange-500/25 transition-all text-center active:scale-95 cursor-pointer"
          >
            فهمت وموافق
          </button>
        </div>
      </div>
    </div>
  </Transition>
</template>
