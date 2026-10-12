<script setup lang="ts">
import { computed } from 'vue'
import { Loader2, CheckCircle, XCircle, UploadCloud, Clock } from '@lucide/vue'

const props = defineProps<{
  percent: number
  state: 'idle' | 'uploading' | 'processing' | 'completed' | 'failed'
  error?: string | null
  uploadedSize?: number
  totalSize?: number
}>()

const emit = defineEmits<{ cancel: [] }>()

const stateMeta = computed(() => {
  const map = {
    idle: { label: 'جاهز للرفع', color: 'text-slate-500', bar: 'bg-slate-300' },
    uploading: { label: 'جارٍ الرفع إلى YouTube', color: 'text-blue-500', bar: 'bg-blue-500' },
    processing: { label: 'جارٍ المعالجة', color: 'text-amber-500', bar: 'bg-amber-500' },
    completed: { label: 'اكتمل الرفع', color: 'text-emerald-500', bar: 'bg-emerald-500' },
    failed: { label: 'فشل الرفع', color: 'text-red-500', bar: 'bg-red-500' },
  }
  return map[props.state] ?? map.idle
})

const speed = computed(() => {
  // approximate: assume upload finished, show nothing extra
  return null
})

function formatBytes(b?: number): string {
  if (!b || b <= 0) return '0'
  const units = ['بايت', 'ك.ب', 'م.ب', 'ج.ب']
  let i = 0
  let v = b
  while (v >= 1024 && i < units.length - 1) {
    v /= 1024
    i++
  }
  return `${v.toFixed(v >= 10 || i === 0 ? 0 : 1)} ${units[i]}`
}
</script>

<template>
  <div v-if="state !== 'idle'" class="rounded-2xl border p-4 space-y-3" :class="state === 'failed' ? 'border-red-200 dark:border-red-900/60 bg-red-50/50 dark:bg-red-950/20' : 'border-slate-200 dark:border-slate-800 bg-slate-50/50 dark:bg-slate-900/50'">
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-2">
        <component
          :is="state === 'completed' ? CheckCircle : state === 'failed' ? XCircle : state === 'processing' ? Clock : Loader2"
          :class="['h-4 w-4', stateMeta.color, (state === 'uploading' || state === 'processing') ? 'animate-spin' : '']"
        />
        <span class="text-sm font-bold" :class="stateMeta.color">{{ stateMeta.label }}</span>
      </div>
      <div class="flex items-center gap-2">
        <span v-if="uploadedSize && totalSize" class="text-[11px] text-slate-400">
          {{ formatBytes(uploadedSize) }} / {{ formatBytes(totalSize) }}
        </span>
        <button
          v-if="state === 'uploading'"
          class="text-[11px] font-bold text-red-500 hover:text-red-600"
          @click="emit('cancel')"
        >
          إلغاء
        </button>
      </div>
    </div>

    <div class="h-2 rounded-full bg-slate-200 dark:bg-slate-700 overflow-hidden">
      <div
        class="h-full rounded-full transition-all duration-300"
        :class="stateMeta.bar"
        :style="{ width: `${Math.min(100, percent)}%` }"
      />
    </div>

    <div class="flex items-center justify-between">
      <span class="text-xs font-bold text-slate-500">{{ percent }}%</span>
      <span v-if="speed" class="text-[11px] text-slate-400">الوقت المتبقي: —</span>
    </div>

    <p v-if="error" class="text-xs text-red-500">{{ error }}</p>
  </div>

  <div v-else class="rounded-2xl border border-dashed border-slate-300 dark:border-slate-700 p-6 text-center">
    <UploadCloud class="h-8 w-8 text-slate-300 dark:text-slate-600 mx-auto mb-2" />
    <p class="text-sm text-slate-500 dark:text-slate-400">لم يتم رفع فيديو بعد</p>
    <p class="text-xs text-slate-400 mt-1">اختر ملف فيديو وسيتم رفعه إلى YouTube تلقائياً</p>
  </div>
</template>
