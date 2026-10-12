<script setup lang="ts">
import { AlertCircle, RotateCw } from "@lucide/vue";

withDefaults(
  defineProps<{
    title?: string;
    message?: string;
    retryText?: string;
  }>(),
  {
    title: "حدث خطأ أثناء تحميل البيانات",
    message: "تعذر الاتصال بالخادم مؤقتاً. يرجى التحقق من اتصالك بالإنترنت والمحاولة مجدداً.",
    retryText: "إعادة المحاولة",
  }
);

const emit = defineEmits<{
  (e: "retry"): void;
}>();
</script>

<template>
  <div class="flex flex-col items-center justify-center text-center p-8 bg-rose-50 dark:bg-rose-950/20 border border-rose-200 dark:border-rose-900/40 rounded-3xl my-6">
    <div class="w-14 h-14 rounded-2xl bg-rose-500/10 text-rose-500 flex items-center justify-center mb-3">
      <AlertCircle class="w-7 h-7" />
    </div>
    <h3 class="text-base sm:text-lg font-bold text-rose-950 dark:text-rose-200">
      {{ title }}
    </h3>
    <p class="text-xs sm:text-sm text-rose-700 dark:text-rose-400/90 max-w-md mt-1 mb-5">
      {{ message }}
    </p>
    <button
      @click="emit('retry')"
      type="button"
      class="inline-flex items-center gap-2 px-5 py-2.5 rounded-xl bg-rose-600 hover:bg-rose-700 text-white font-bold text-xs sm:text-sm shadow-md shadow-rose-600/20 transition-all duration-200"
    >
      <RotateCw class="w-4 h-4" />
      <span>{{ retryText }}</span>
    </button>
  </div>
</template>
