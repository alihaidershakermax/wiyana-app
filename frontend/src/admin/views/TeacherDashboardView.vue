<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAdminStore } from '@/admin/stores/useAdminStore'
import {
  useTeacherDashboard,
  useTeacherAnalytics,
  useTeacherCourses,
} from '@/admin/queries/useTeacherDashboard'
import {
  BookOpen, Users, FileText, TrendingUp, Eye, Award,
  RefreshCw, ArrowLeft, Clock, Upload
} from '@lucide/vue'
import type { AnalyticsPeriod } from '@/admin/types/dashboard'
import type { Course } from '@/admin/types/course'
import { COURSE_STATUS_LABELS, COURSE_STATUS_COLORS } from '@/admin/types/course'

const adminStore = useAdminStore()
const router = useRouter()
const selectedPeriod = ref<AnalyticsPeriod>('7d')

const periods = [
  { value: 'day' as const, label: 'اليوم' },
  { value: '7d' as const, label: '7 أيام' },
  { value: '30d' as const, label: '30 يوم' },
]

const { data: dashboard, isLoading, isError, refetch } = useTeacherDashboard()
const { data: analytics, isLoading: analyticsLoading } = useTeacherAnalytics(selectedPeriod)
const { data: courses, isLoading: coursesLoading } = useTeacherCourses()

const statCards = [
  { key: 'my_courses', label: 'مساقاتي', icon: BookOpen, color: 'text-orange-500', bg: 'bg-orange-50 dark:bg-orange-950/40', border: 'border-orange-200 dark:border-orange-900/60' },
  { key: 'my_students', label: 'طلابي', icon: Users, color: 'text-blue-500', bg: 'bg-blue-50 dark:bg-blue-950/40', border: 'border-blue-200 dark:border-blue-900/60' },
  { key: 'my_lessons', label: 'محاضراتي', icon: FileText, color: 'text-purple-500', bg: 'bg-purple-50 dark:bg-purple-950/40', border: 'border-purple-200 dark:border-purple-900/60' },
  { key: 'total_views', label: 'إجمالي المشاهدات', icon: Eye, color: 'text-emerald-500', bg: 'bg-emerald-50 dark:bg-emerald-950/40', border: 'border-emerald-200 dark:border-emerald-900/60' },
]

function getStatValue(key: string): number {
  if (!dashboard.value) return 0
  return (dashboard.value as any)[key] ?? 0
}

function navigateTo(name: string, params?: Record<string, string>) {
  router.push({ name, params })
}

const maxVal = (data: any[]) => Math.max(...data.map(d => d.subscriptions + d.students), 1)
</script>

<template>
  <div dir="rtl" class="space-y-8 animate-fade-in">
    <!-- Header -->
    <div class="flex items-center justify-between flex-wrap gap-4">
      <div>
        <h1 class="text-2xl font-extrabold text-slate-900 dark:text-white">لوحتي</h1>
        <p class="mt-1 text-sm text-slate-500 dark:text-slate-400">
          مرحباً، {{ adminStore.userName }} 👋
        </p>
      </div>

      <button
        type="button"
        @click="navigateTo('admin.teacher.lessons')"
        class="px-5 py-2.5 rounded-2xl bg-orange-500 hover:bg-[#EA580C] text-white font-black text-xs sm:text-sm flex items-center gap-2 shadow-sm shadow-orange-500/20 cursor-pointer"
      >
        <Upload class="w-4 h-4" />
        <span>رفع وإدارة المحاضرات</span>
      </button>
    </div>

    <!-- Error -->
    <div v-if="isError" class="card p-6 text-center border-red-200 bg-red-50 dark:bg-red-950/20">
      <p class="text-red-600 dark:text-red-400 font-medium mb-3">حدث خطأ أثناء تحميل البيانات</p>
      <button @click="refetch()" class="btn-secondary text-red-500 border-red-200">إعادة المحاولة</button>
    </div>

    <!-- Stats -->
    <div class="grid grid-cols-2 lg:grid-cols-4 gap-4">
      <div v-for="stat in statCards" :key="stat.key">
        <div v-if="isLoading" class="card p-5 animate-pulse space-y-3">
          <div class="h-3 bg-slate-200 dark:bg-slate-700 rounded w-2/3" />
          <div class="h-7 bg-slate-200 dark:bg-slate-700 rounded w-1/2" />
        </div>
        <div v-else class="card p-5 flex items-start gap-3 hover:shadow-md transition-shadow">
          <div :class="['h-11 w-11 rounded-xl border flex items-center justify-center shrink-0', stat.bg, stat.border]">
            <component :is="stat.icon" :class="['h-5 w-5', stat.color]" />
          </div>
          <div>
            <p class="text-xs text-slate-500 dark:text-slate-400">{{ stat.label }}</p>
            <p class="text-xl font-extrabold tabular-nums text-slate-900 dark:text-white mt-0.5">
              {{ getStatValue(stat.key).toLocaleString('ar') }}
            </p>
          </div>
        </div>
      </div>
    </div>

    <!-- Analytics -->
    <div class="card p-6 space-y-4">
      <div class="flex items-center justify-between flex-wrap gap-3">
        <h2 class="text-sm font-bold text-slate-700 dark:text-slate-200 flex items-center gap-2">
          <TrendingUp class="h-4 w-4 text-orange-500" />
          تحليلات الاشتراكات
        </h2>
        <div class="flex items-center gap-1 bg-slate-100 dark:bg-slate-800 rounded-xl p-1">
          <button
            v-for="p in periods"
            :key="p.value"
            @click="selectedPeriod = p.value"
            :class="[
              'text-xs px-3 py-1.5 rounded-lg font-medium transition-all',
              selectedPeriod === p.value
                ? 'bg-white dark:bg-slate-700 shadow text-orange-600 dark:text-orange-400 font-bold'
                : 'text-slate-500 hover:text-slate-700 dark:hover:text-slate-300',
            ]"
          >
            {{ p.label }}
          </button>
        </div>
      </div>

      <div v-if="analyticsLoading" class="h-36 animate-pulse bg-slate-100 dark:bg-slate-800 rounded-xl" />
      <div v-else-if="analytics?.data?.length" class="h-36 flex items-end gap-2">
        <div
          v-for="point in analytics.data"
          :key="point.label"
          class="flex-1 min-w-6 flex flex-col items-center gap-1"
        >
          <div
            :style="{ height: `${Math.max(4, (point.subscriptions / maxVal(analytics.data)) * 120)}px` }"
            class="w-full bg-orange-400 dark:bg-orange-500 rounded-t-md transition-all duration-500"
          />
          <span class="text-[9px] text-slate-400">{{ point.label }}</span>
        </div>
      </div>
    </div>

    <!-- My Courses -->
    <div class="card p-6 space-y-4">
      <div class="flex items-center justify-between">
        <h2 class="text-sm font-bold text-slate-700 dark:text-slate-200 flex items-center gap-2">
          <BookOpen class="h-4 w-4 text-orange-500" />
          مساقاتي
        </h2>
        <button
          @click="navigateTo('admin.teacher.courses')"
          class="text-xs text-orange-500 hover:text-orange-600 flex items-center gap-1 font-medium"
        >
          عرض الكل <ArrowLeft class="h-3 w-3" />
        </button>
      </div>

      <div v-if="coursesLoading" class="space-y-3">
        <div v-for="i in 3" :key="i" class="animate-pulse flex gap-4 py-3">
          <div class="h-10 w-10 rounded-xl bg-slate-200 dark:bg-slate-700 shrink-0" />
          <div class="flex-1 space-y-2 pt-1">
            <div class="h-3.5 bg-slate-200 dark:bg-slate-700 rounded w-2/3" />
            <div class="h-3 bg-slate-200 dark:bg-slate-700 rounded w-1/3" />
          </div>
        </div>
      </div>

      <div v-else-if="courses?.length" class="divide-y divide-slate-100 dark:divide-slate-800">
        <div
          v-for="course in courses"
          :key="course.id"
          @click="router.push({ name: 'admin.course-detail', params: { id: course.id } })"
          class="flex items-center gap-4 py-3 cursor-pointer hover:bg-slate-50 dark:hover:bg-slate-800/50 -mx-2 px-2 rounded-xl transition-colors"
        >
          <div class="h-10 w-10 rounded-xl bg-orange-50 dark:bg-orange-950/40 border border-orange-200 dark:border-orange-900/60 flex items-center justify-center shrink-0">
            <BookOpen class="h-5 w-5 text-orange-500" />
          </div>
          <div class="flex-1 min-w-0">
            <p class="text-sm font-semibold text-slate-900 dark:text-white truncate">{{ course.title }}</p>
            <div class="flex items-center gap-3 mt-0.5 text-xs text-slate-400">
              <span class="flex items-center gap-1"><Users class="h-3 w-3" /> {{ course.students_count }}</span>
              <span class="flex items-center gap-1"><FileText class="h-3 w-3" /> {{ course.lessons_count }} محاضرة</span>
            </div>
          </div>
          <span :class="['text-[10px] font-bold px-2 py-0.5 rounded-full', COURSE_STATUS_COLORS[course.status as keyof typeof COURSE_STATUS_COLORS] ?? COURSE_STATUS_COLORS.draft]">
            {{ COURSE_STATUS_LABELS[course.status as keyof typeof COURSE_STATUS_LABELS] ?? course.status }}
          </span>
        </div>
      </div>

      <div v-else class="text-center py-10">
        <Award class="h-10 w-10 text-slate-200 dark:text-slate-700 mx-auto mb-3" />
        <p class="text-sm text-slate-400">لم يتم تعيينك في أي مساق بعد</p>
      </div>
    </div>
  </div>
</template>
