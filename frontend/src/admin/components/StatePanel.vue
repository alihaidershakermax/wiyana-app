<script setup lang="ts">
import { Loader2, AlertCircle, Inbox, RefreshCw } from '@lucide/vue'

defineProps<{
  loading?: boolean
  error?: string | null
  empty?: boolean
  emptyTitle?: string
  emptyMessage?: string
  skeletonCount?: number
  height?: string
}>()

const emit = defineEmits<{ retry: [] }>()

const defaultSkeleton = [
  'h-14',
  'h-14',
  'h-14',
  'h-14',
  'h-14',
]
</script>

<template>
  <div v-if="loading" class="space-y-3" :style="{ minHeight: height }" role="status" aria-live="polite" aria-label="جاري تحميل البيانات">
    <div
      v-for="(h, i) in (skeletonCount ? Array.from({ length: skeletonCount }) : defaultSkeleton) as string[]"
      :key="i"
      class="skeleton-shimmer rounded-xl"
      :class="skeletonCount ? 'h-16' : h"
    />
  </div>

  <div
    v-else-if="error"
    class="state-panel text-center space-y-3 py-10" role="alert"
  >
    <div class="h-14 w-14 rounded-2xl bg-red-50 dark:bg-red-950/40 border border-red-200 dark:border-red-900/60 flex items-center justify-center">
      <AlertCircle class="h-7 w-7 text-red-500" />
    </div>
    <div>
      <p class="text-sm font-bold text-red-600 dark:text-red-400">حدث خطأ أثناء تحميل البيانات</p>
      <p class="text-xs text-slate-500 dark:text-slate-400 mt-1">{{ error }}</p>
    </div>
    <button class="btn-secondary text-sm" @click="emit('retry')">
      <RefreshCw class="h-4 w-4" />
      إعادة المحاولة
    </button>
  </div>

  <div
    v-else-if="empty"
    class="state-panel text-center space-y-3 py-12" aria-live="polite"
  >
    <div class="h-14 w-14 rounded-2xl bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-800 flex items-center justify-center">
      <Inbox class="h-7 w-7 text-slate-400 dark:text-slate-600" />
    </div>
    <div>
      <p class="text-sm font-bold text-slate-600 dark:text-slate-300">{{ emptyTitle ?? 'لا توجد بيانات' }}</p>
      <p v-if="emptyMessage" class="text-xs text-slate-400 mt-1">{{ emptyMessage }}</p>
    </div>
    <slot name="empty-action" />
  </div>

  <slot v-else />
</template>

<style scoped>
.skeleton {
  background: linear-gradient(
    90deg,
    rgb(241 245 249) 25%,
    rgb(226 232 240) 50%,
    rgb(241 245 249) 75%
  );
  background-size: 200% 100%;
  animation: shimmer 1.4s infinite;
}
.dark .skeleton {
  background: linear-gradient(
    90deg,
    rgb(30 41 59) 25%,
    rgb(51 65 85) 50%,
    rgb(30 41 59) 75%
  );
  background-size: 200% 100%;
  animation: shimmer 1.4s infinite;
}
@keyframes shimmer {
  0% { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}
</style>
