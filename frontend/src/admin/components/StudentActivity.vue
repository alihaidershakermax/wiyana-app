<script setup lang="ts">
import { computed } from 'vue'
import { Activity as ActivityIcon, PlayCircle, CheckCircle, UserPlus, FileText, Clock } from '@lucide/vue'
import { formatDate, formatTimeOfDay } from '@/utils/format'

interface ActivityItem {
  type?: string
  lesson_id?: string | null
  course_id?: string | null
  progress?: number
  completed?: boolean
  watched_seconds?: number
  last_position?: number
  updated_at?: string | null
}

const props = defineProps<{
  activities: ActivityItem[]
  loading?: boolean
}>()

function typeMeta(type?: string) {
  const map: Record<string, { label: string; icon: unknown; cls: string }> = {
    NEW_STUDENT: { label: 'تسجيل جديد', icon: UserPlus, cls: 'text-emerald-500' },
    LESSON_WATCH: { label: 'مشاهدة محاضرة', icon: PlayCircle, cls: 'text-blue-500' },
    LESSON_COMPLETED: { label: 'إكمال محاضرة', icon: CheckCircle, cls: 'text-emerald-500' },
    ATTACHMENT: { label: 'تحميل مرفق', icon: FileText, cls: 'text-amber-500' },
  }
  return map[type ?? ''] ?? { label: type ?? 'نشاط', icon: Clock, cls: 'text-slate-400' }
}

const items = computed(() =>
  [...props.activities].sort(
    (a, b) => new Date(b.updated_at ?? 0).getTime() - new Date(a.updated_at ?? 0).getTime()
  )
)
</script>

<template>
  <div class="space-y-3">
    <div
      v-for="(act, i) in items"
      :key="i"
      class="flex items-start gap-3 rounded-xl border border-slate-100 dark:border-slate-800 p-3"
    >
      <div class="h-8 w-8 rounded-lg bg-slate-50 dark:bg-slate-800 flex items-center justify-center shrink-0">
        <component :is="typeMeta(act.type).icon" :class="['h-4 w-4', typeMeta(act.type).cls]" />
      </div>
      <div class="flex-1 min-w-0">
        <p class="text-xs font-bold text-slate-700 dark:text-slate-200">{{ typeMeta(act.type).label }}</p>
        <p class="text-[11px] text-slate-400 mt-0.5">
          <span v-if="act.progress !== undefined">التقدم: {{ act.progress }}% · </span>
          <span v-if="act.completed">مكتمل · </span>
          <span v-if="act.watched_seconds">مدة المشاهدة: {{ Math.round(act.watched_seconds / 60) }} د · </span>
          <span v-if="act.updated_at">{{ formatDate(act.updated_at) }} {{ formatTimeOfDay(act.updated_at) }}</span>
        </p>
      </div>
    </div>

    <div v-if="!loading && !items.length" class="rounded-xl border border-dashed border-slate-200 dark:border-slate-800 py-8 text-center">
      <ActivityIcon class="h-6 w-6 text-slate-300 dark:text-slate-600 mx-auto mb-1" />
      <p class="text-xs text-slate-400">لا يوجد نشاط حديث للطالب</p>
    </div>
  </div>
</template>
