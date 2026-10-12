<script setup lang="ts">
import { computed } from 'vue'
import { CheckCircle2, Circle, TrendingUp } from '@lucide/vue'
import type { CourseProgressSummary } from '@/admin/types/panel'

const props = defineProps<{
  summary: CourseProgressSummary | null
  loading?: boolean
}>()

const percent = computed(() => Math.min(100, Math.max(0, props.summary?.progress_percent ?? 0)))
</script>

<template>
  <div class="rounded-2xl border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 p-4 space-y-4">
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-2">
        <TrendingUp class="h-4 w-4 text-orange-500" />
        <h4 class="text-sm font-bold">التقدم في المساق</h4>
      </div>
      <span class="text-xs font-bold text-orange-500">{{ percent }}%</span>
    </div>

    <div class="h-2.5 rounded-full bg-slate-100 dark:bg-slate-800 overflow-hidden">
      <div class="h-full rounded-full bg-gradient-to-l from-orange-400 to-orange-600 transition-all duration-500" :style="{ width: `${percent}%` }" />
    </div>

    <div class="flex items-center justify-around text-center">
      <div>
        <p class="text-xl font-extrabold text-slate-800 dark:text-slate-100">{{ summary?.completed_lessons ?? 0 }}</p>
        <p class="text-[11px] text-slate-400">مكتملة</p>
      </div>
      <div class="h-8 w-px bg-slate-200 dark:bg-slate-800" />
      <div>
        <p class="text-xl font-extrabold text-slate-800 dark:text-slate-100">{{ summary?.total_lessons ?? 0 }}</p>
        <p class="text-[11px] text-slate-400">إجمالي المحاضرات</p>
      </div>
      <div class="h-8 w-px bg-slate-200 dark:bg-slate-800" />
      <div>
        <p class="text-xl font-extrabold text-slate-800 dark:text-slate-100">{{ summary?.completed_lesson_ids?.length ?? 0 }}</p>
        <p class="text-[11px] text-slate-400">مشاهدة كاملة</p>
      </div>
    </div>

    <div v-if="summary?.last_watched_lesson_id" class="flex items-center gap-2 text-[11px] text-slate-400">
      <CheckCircle2 class="h-3.5 w-3.5 text-emerald-500" />
      آخر محاضرة مشاهدة: {{ summary.last_watched_lesson_id }}
    </div>
  </div>
</template>
