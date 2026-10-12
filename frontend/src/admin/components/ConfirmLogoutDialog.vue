<script setup lang="ts">
import { LogOut, X, Loader2 } from "@lucide/vue";

defineProps<{
  open: boolean;
  loading?: boolean;
}>();

const emit = defineEmits<{
  (e: "confirm"): void;
  (e: "cancel"): void;
}>();
</script>

<template>
  <Teleport to="body">
    <Transition name="fade">
      <div v-if="open" class="fixed inset-0 z-[100] flex items-center justify-center p-4">
        <!-- Backdrop -->
        <div
          class="absolute inset-0 bg-slate-950/60 backdrop-blur-sm"
          @click="!loading && emit('cancel')"
        ></div>

        <!-- Dialog -->
        <div
          dir="rtl"
          class="relative w-full max-w-sm rounded-2xl bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 shadow-2xl p-6 animate-fade-in"
        >
          <button
            type="button"
            class="absolute top-4 left-4 p-1.5 rounded-lg text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 transition-colors"
            @click="!loading && emit('cancel')"
            aria-label="إغلاق"
          >
            <X class="h-4 w-4" />
          </button>

          <div
            class="mx-auto h-14 w-14 rounded-2xl bg-red-50 dark:bg-red-950/40 border border-red-200 dark:border-red-900/60 flex items-center justify-center mb-4"
          >
            <LogOut class="h-7 w-7 text-red-500" />
          </div>

          <h2 class="text-center text-lg font-extrabold">تسجيل الخروج</h2>
          <p class="mt-2 text-center text-sm text-slate-500 dark:text-slate-400 leading-relaxed">
            هل أنت متأكد من رغبتك في تسجيل الخروج من لوحة التحكم؟
          </p>

          <div class="mt-6 grid grid-cols-2 gap-3">
            <button
              type="button"
              class="btn-secondary py-2.5 text-sm"
              :disabled="loading"
              @click="emit('cancel')"
            >
              إلغاء
            </button>
            <button
              type="button"
              class="inline-flex items-center justify-center gap-2 rounded-xl bg-red-500 text-white text-sm font-bold py-2.5 hover:bg-red-600 disabled:opacity-50 transition-colors"
              :disabled="loading"
              @click="emit('confirm')"
            >
              <Loader2 v-if="loading" class="h-4 w-4 animate-spin" />
              <LogOut v-else class="h-4 w-4" />
              خروج
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.15s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
