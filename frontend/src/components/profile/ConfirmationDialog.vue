<script setup lang="ts">
import { LogOut, X } from "@lucide/vue";

defineProps<{
  isOpen: boolean;
  title?: string;
  message?: string;
  confirmText?: string;
  cancelText?: string;
}>();

const emit = defineEmits<{
  (e: "confirm"): void;
  (e: "cancel"): void;
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
      class="fixed inset-0 z-50 bg-slate-950/80 backdrop-blur-sm flex items-center justify-center p-4"
      @click.self="emit('cancel')"
    >
      <div class="w-full max-w-sm bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl p-6 shadow-2xl flex flex-col items-center text-center gap-4">
        <div class="w-14 h-14 rounded-2xl bg-rose-500/10 text-rose-500 flex items-center justify-center">
          <LogOut class="w-7 h-7" />
        </div>

        <div class="flex flex-col gap-1">
          <h3 class="text-lg font-bold text-slate-900 dark:text-slate-100">
            {{ title || "تأكيد تسجيل الخروج" }}
          </h3>
          <p class="text-xs sm:text-sm text-slate-500 dark:text-slate-400">
            {{ message || "هل أنت تأكد من رغبتك في تسجيل الخروج من المنصة؟" }}
          </p>
        </div>

        <div class="flex items-center gap-3 w-full mt-2">
          <button
            @click="emit('cancel')"
            type="button"
            class="flex-1 py-2.5 rounded-xl border border-slate-200 dark:border-slate-800 text-slate-700 dark:text-slate-300 font-bold text-xs sm:text-sm hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors"
          >
            {{ cancelText || "إلغاء" }}
          </button>

          <button
            @click="emit('confirm')"
            type="button"
            class="flex-1 py-2.5 rounded-xl bg-rose-600 hover:bg-rose-700 text-white font-bold text-xs sm:text-sm shadow-md shadow-rose-600/20 transition-colors"
          >
            {{ confirmText || "تأكيد الخروج" }}
          </button>
        </div>
      </div>
    </div>
  </Transition>
</template>
