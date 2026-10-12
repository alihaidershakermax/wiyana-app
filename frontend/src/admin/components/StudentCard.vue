<script setup lang="ts">
import { computed } from 'vue'
import { ChevronLeft, BookOpen, Clock, CheckCircle, Eye, Activity } from '@lucide/vue'
import type { StudentCourse, Student } from '@/admin/types/panel'
import { formatDate } from '@/utils/format'

const props = defineProps<{
  student: Student
  courses: StudentCourse[]
}>()

const emit = defineEmits<{ viewCourse: [courseId: string] }>()

const statusMeta: Record<string, { label: string; cls: string }> = {
  active: { label: 'نشط', cls: 'bg-emerald-100 text-emerald-700 dark:bg-emerald-950/40 dark:text-emerald-400' },
  expired: { label: 'منتهي', cls: 'bg-slate-100 text-slate-600 dark:bg-slate-800 dark:text-slate-400' },
  suspended: { label: 'موقوف', cls: 'bg-red-100 text-red-700 dark:bg-red-950/40 dark:text-red-400' },
  cancelled: { label: 'ملغي', cls: 'bg-amber-100 text-amber-700 dark:bg-amber-950/40 dark:text-amber-400' },
}

function courseProgress(c: StudentCourse): number {
  if (!c.progress?.length) return 0
  const completed = c.progress.filter((p) => p.completed).length
  return Math.round((completed / c.progress.length) * 100)
}
</script>

<template>
  <div class="space-y-3">
    <div
      v-for="c in courses"
      :key="c.course.id"
      class="rounded-2xl border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 p-4 hover:border-orange-300 dark:hover:border-orange-800 transition-colors"
    >
      <div class="flex items-start justify-between gap-3">
        <div class="min-w-0 flex-1">
          <p class="text-sm font-bold text-slate-800 dark:text-slate-100 truncate">{{ c.course.title }}</p>
          <p class="text-[11px] text-slate-400 mt-0.5">
            {{ c.teacher?.first_name || 'بدون أستاذ' }}
          </p>
        </div>
        <span :class="['text-[10px] font-bold px-2 py-0.5 rounded-full shrink-0', statusMeta[c.enrollment.status]?.cls ?? statusMeta.active.cls]">
          {{ statusMeta[c.enrollment.status]?.label ?? c.enrollment.status }}
        </span>
      </div>

      <div class="mt-3 space-y-2">
        <div class="flex items-center justify-between text-[11px] text-slate-400">
          <span>التقدم: {{ courseProgress(c) }}%</span>
          <span class="flex items-center gap-1">
            <CheckCircle class="h-3 w-3" />
            {{ c.progress.filter((p) => p.completed).length }} / {{ c.progress.length }} محاضرة
          </span>
        </div>
        <div class="h-1.5 rounded-full bg-slate-100 dark:bg-slate-800 overflow-hidden">
          <div class="h-full rounded-full bg-orange-500 transition-all" :style="{ width: `${courseProgress(c)}%` }" />
        </div>
      </div>

      <div class="mt-3 flex items-center justify-between text-[11px] text-slate-400">
        <span class="flex items-center gap-1">
          <Clock class="h-3 w-3" />
          يبدأ: {{ formatDate(c.enrollment.started_at) ?? '—' }}
        </span>
        <span v-if="c.enrollment.expires_at" class="flex items-center gap-1">
          ينتهي: {{ formatDate(c.enrollment.expires_at) }}
        </span>
        <button
          class="flex items-center gap-1 text-orange-500 hover:text-orange-600 font-bold"
          @click="emit('viewCourse', c.course.id)"
        >
          <Eye class="h-3.5 w-3.5" />
          التفاصيل
        </button>
      </div>
    </div>

    <div v-if="!courses.length" class="rounded-2xl border border-dashed border-slate-200 dark:border-slate-800 py-8 text-center">
      <BookOpen class="h-6 w-6 text-slate-300 dark:text-slate-600 mx-auto mb-1" />
      <p class="text-xs text-slate-400">لا توجد مساقات مسجل بها الطالب</p>
    </div>
  </div>
</template>
