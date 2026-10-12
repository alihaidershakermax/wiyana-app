<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAdminStore } from '@/admin/stores/useAdminStore'
import {
  useDashboardOverview,
  useDashboardAnalytics,
  useDashboardActivity,
  useRecentCourses,
} from '@/admin/queries/useDashboard'
import { useToast } from '@/admin/composables/useToast'
import { youtubeApi } from '@/admin/api/panelApi'
import {
  Users, BookOpen, FolderOpen, FileText, TrendingUp, Eye,
  CreditCard, Key, Plus, ArrowLeft, RefreshCw, Activity,
  Clock, CheckCircle, Star, Zap, BookMarked, Video,
  ExternalLink, Unlink, Tv
} from '@lucide/vue'
import type { AnalyticsPeriod } from '@/admin/types/dashboard'

const adminStore = useAdminStore()
const router = useRouter()
const toast = useToast()

const selectedPeriod = ref<AnalyticsPeriod>('7d')
const periods = [
  { value: 'day' as const, label: 'اليوم' },
  { value: '7d' as const, label: '7 أيام' },
  { value: '30d' as const, label: '30 يوم' },
  { value: '3m' as const, label: '3 أشهر' },
  { value: '1y' as const, label: 'سنة' },
]

const { data: overview, isLoading: overviewLoading, isError: overviewError, refetch: refetchOverview } = useDashboardOverview()
const { data: analytics, isLoading: analyticsLoading } = useDashboardAnalytics(selectedPeriod)
const { data: activityData, isLoading: activityLoading } = useDashboardActivity()
const { data: recentCourses, isLoading: coursesLoading } = useRecentCourses(6)

const ytStatus = ref<{
  is_authenticated: boolean
  channel_name?: string
  channel_id?: string
  channel_avatar?: string
} | null>(null)
const ytLoading = ref(false)

async function fetchYtStatus() {
  ytLoading.value = true
  try {
    const res = await youtubeApi.getStatus()
    ytStatus.value = res
  } catch (err) {
    console.warn('Failed to fetch YouTube status', err)
  } finally {
    ytLoading.value = false
  }
}

function startGoogleOAuth() {
  const backendBase = import.meta.env.VITE_ADMIN_API_BASE_URL || '/api/v1'
  window.open(`${backendBase}/youtube/oauth/authorize`, '_blank')
}

fetchYtStatus()

const statCards = computed(() => {
  if (!overview.value) return []
  return [
    {
      label: 'إجمالي الطلاب',
      value: Number(overview.value.total_students ?? 0),
      sub: `${overview.value.active_students} نشط`,
      icon: Users,
      color: 'text-blue-400',
      bg: 'bg-blue-950/30',
      border: 'border-blue-900/40',
    },
    {
      label: 'المساقات',
      value: Number(overview.value.total_courses ?? 0),
      sub: 'إجمالي المساقات',
      icon: BookOpen,
      color: 'text-[#F97316]',
      bg: 'bg-[#1F0E03]',
      border: 'border-orange-900/40',
    },

    {
      label: 'المحاضرات',
      value: Number(overview.value.total_lessons ?? 0),
      sub: 'محاضرة منشورة',
      icon: FileText,
      color: 'text-emerald-400',
      bg: 'bg-emerald-950/30',
      border: 'border-emerald-900/40',
    },
    {
      label: 'الاشتراكات النشطة',
      value: Number(overview.value.active_subscriptions ?? 0),
      sub: 'اشتراك فعّال',
      icon: CreditCard,
      color: 'text-cyan-400',
      bg: 'bg-cyan-950/30',
      border: 'border-cyan-900/40',
    },
    {
      label: 'الأكواد المستخدمة',
      value: Number(overview.value.used_subscription_codes ?? 0),
      sub: 'كود تم استخدامه',
      icon: Key,
      color: 'text-rose-400',
      bg: 'bg-rose-950/30',
      border: 'border-rose-900/40',
    },
  ]
})

const maxAnalyticsValue = computed(() => {
  if (!analytics.value?.data?.length) return 1
  return Math.max(...analytics.value.data.map(d => d.students + d.subscriptions), 1)
})

const activityTypeIcons: Record<string, string> = {
  NEW_STUDENT: '👤',
  NEW_COURSE: '📚',
  NEW_SUBJECT: '📁',
  NEW_LESSON: '🎬',
  NEW_SUBSCRIPTION: '🎫',
  CODE_USED: '🔑',
  TEACHER_ASSIGNED: '👨‍🏫',
  COURSE_PUBLISHED: '✅',
  COURSE_ARCHIVED: '📦',
}

const courseStatusLabels: Record<string, string> = {
  draft: 'مسودة',
  published: 'منشور',
  archived: 'مؤرشف',
}

const courseStatusClasses: Record<string, string> = {
  draft: 'bg-amber-100 text-amber-700 dark:bg-amber-950/40 dark:text-amber-400',
  published: 'bg-emerald-100 text-emerald-700 dark:bg-emerald-950/40 dark:text-emerald-400',
  archived: 'bg-slate-100 text-slate-600 dark:bg-slate-800 dark:text-slate-400',
}

function formatDate(dateStr: string) {
  return new Date(dateStr).toLocaleDateString('ar', { month: 'short', day: 'numeric' })
}

function formatTime(dateStr: string) {
  const d = new Date(dateStr)
  const diff = Date.now() - d.getTime()
  const mins = Math.floor(diff / 60000)
  if (mins < 1) return 'الآن'
  if (mins < 60) return `منذ ${mins} دقيقة`
  const hours = Math.floor(mins / 60)
  if (hours < 24) return `منذ ${hours} ساعة`
  return `منذ ${Math.floor(hours / 24)} يوم`
}

function navigateTo(name: string, params?: Record<string, string>) {
  router.push({ name, params })
}
</script>

<template>
  <div dir="rtl" class="space-y-8 animate-fade-in">
    <!-- Header -->
    <div class="flex items-center justify-between flex-wrap gap-4">
      <div>
        <h1 class="text-2xl font-extrabold tracking-tight text-slate-900 dark:text-white">
          لوحة التحكم
        </h1>
        <p class="mt-1 text-sm text-slate-500 dark:text-slate-400">
          مرحباً، {{ adminStore.userName }} 👋
        </p>
      </div>
      <div class="flex items-center gap-3">
        <button
          @click="navigateTo('admin.courses')"
          class="btn-primary gap-2"
        >
          <Plus class="h-4 w-4" />
          مساق جديد
        </button>
      </div>
    </div>

    <!-- Error State -->
    <div
      v-if="overviewError"
      class="card p-6 border-red-200 dark:border-red-900/60 bg-red-50 dark:bg-red-950/20 text-center"
    >
      <p class="text-red-600 dark:text-red-400 font-medium mb-3">
        حدث خطأ أثناء تحميل البيانات
      </p>
      <button @click="refetchOverview()" class="btn-secondary text-red-500 border-red-200">
        إعادة المحاولة
      </button>
    </div>

    <!-- Stats Grid Skeleton -->
    <div v-if="overviewLoading" class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
      <div
        v-for="i in 6"
        :key="i"
        class="card p-6 animate-pulse space-y-3"
      >
        <div class="h-4 bg-slate-200 dark:bg-slate-700 rounded w-1/2" />
        <div class="h-8 bg-slate-200 dark:bg-slate-700 rounded w-1/3" />
        <div class="h-3 bg-slate-200 dark:bg-slate-700 rounded w-2/3" />
      </div>
    </div>

    <!-- Stats Grid -->
    <div v-else-if="overview" class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
      <div
        v-for="stat in statCards"
        :key="stat.label"
        class="card p-6 flex items-start gap-4 hover:shadow-md transition-shadow"
      >
        <div :class="['h-12 w-12 rounded-2xl border flex items-center justify-center shrink-0', stat.bg, stat.border]">
          <component :is="stat.icon" :class="['h-6 w-6', stat.color]" />
        </div>
        <div class="min-w-0">
          <p class="text-xs text-slate-500 dark:text-slate-400 mb-1">{{ stat.label }}</p>
          <p class="text-2xl font-extrabold tabular-nums text-slate-900 dark:text-white">
            {{ Number(stat.value ?? 0).toLocaleString('ar') }}
          </p>
          <p class="text-xs text-slate-400 dark:text-slate-500 mt-0.5">{{ stat.sub }}</p>
        </div>
      </div>
    </div>

    <!-- Quick Actions -->
    <div class="card p-6">
      <h2 class="text-sm font-bold text-slate-500 dark:text-slate-400 mb-4 flex items-center gap-2">
        <Zap class="h-4 w-4 text-orange-500" />
        إجراءات سريعة
      </h2>
      <div class="grid grid-cols-2 sm:grid-cols-4 gap-3">
        <button
          @click="navigateTo('admin.course-create')"
          class="flex flex-col items-center gap-2 p-4 rounded-2xl border border-slate-200 dark:border-slate-700 hover:border-orange-300 hover:bg-orange-50 dark:hover:bg-orange-950/20 transition-all group"
        >
          <div class="h-10 w-10 rounded-xl bg-orange-100 dark:bg-orange-950/40 flex items-center justify-center group-hover:scale-110 transition-transform">
            <BookOpen class="h-5 w-5 text-orange-500" />
          </div>
          <span class="text-xs font-semibold text-center">إنشاء مساق</span>
        </button>

        <button
          @click="navigateTo('admin.subscription-codes')"
          class="flex flex-col items-center gap-2 p-4 rounded-2xl border border-slate-200 dark:border-slate-700 hover:border-cyan-300 hover:bg-cyan-50 dark:hover:bg-cyan-950/20 transition-all group"
        >
          <div class="h-10 w-10 rounded-xl bg-cyan-100 dark:bg-cyan-950/40 flex items-center justify-center group-hover:scale-110 transition-transform">
            <Key class="h-5 w-5 text-cyan-500" />
          </div>
          <span class="text-xs font-semibold text-center">كود اشتراك</span>
        </button>

        <button
          @click="navigateTo('admin.teachers')"
          class="flex flex-col items-center gap-2 p-4 rounded-2xl border border-slate-200 dark:border-slate-700 hover:border-emerald-300 hover:bg-emerald-50 dark:hover:bg-emerald-950/20 transition-all group"
        >
          <div class="h-10 w-10 rounded-xl bg-emerald-100 dark:bg-emerald-950/40 flex items-center justify-center group-hover:scale-110 transition-transform">
            <Users class="h-5 w-5 text-emerald-500" />
          </div>
          <span class="text-xs font-semibold text-center">إضافة أستاذ</span>
        </button>

        <button
          @click="navigateTo('admin.settings')"
          class="flex flex-col items-center gap-2 p-4 rounded-2xl border border-slate-200 dark:border-slate-700 hover:border-red-300 hover:bg-red-50 dark:hover:bg-red-950/20 transition-all group"
        >
          <div class="h-10 w-10 rounded-xl bg-red-100 dark:bg-red-950/40 flex items-center justify-center group-hover:scale-110 transition-transform">
            <Video class="h-5 w-5 text-red-500" />
          </div>
          <span class="text-xs font-semibold text-center">YouTube</span>
        </button>
      </div>
    </div>

    <!-- YouTube Connection Card -->
    <div class="card p-6">
      <div class="flex items-center justify-between flex-wrap gap-3">
        <h2 class="text-sm font-bold text-slate-700 dark:text-slate-200 flex items-center gap-2">
          <Video class="h-4 w-4 text-red-500" />
          ربط YouTube
        </h2>
        <button
          @click="fetchYtStatus"
          :disabled="ytLoading"
          class="btn-ghost border border-slate-200 dark:border-slate-700 gap-2 text-xs"
        >
          <RefreshCw :class="['h-3.5 w-3.5', ytLoading ? 'animate-spin' : '']" />
          تحديث
        </button>
      </div>

      <div
        v-if="ytStatus?.is_authenticated"
        class="mt-4 flex items-center justify-between p-4 rounded-2xl bg-emerald-950/20 border border-emerald-500/30"
      >
        <div class="flex items-center gap-3">
          <img
            v-if="ytStatus.channel_avatar"
            :src="ytStatus.channel_avatar"
            class="h-10 w-10 rounded-full object-cover ring-2 ring-emerald-500/30"
          />
          <div v-else class="h-10 w-10 rounded-full bg-emerald-500/20 flex items-center justify-center">
            <Tv class="h-5 w-5 text-emerald-400" />
          </div>
          <div>
            <p class="text-sm font-bold text-slate-900 dark:text-white">{{ ytStatus.channel_name ?? 'قناة YouTube' }}</p>
            <p class="text-xs text-emerald-500 font-medium">متصل بنجاح</p>
          </div>
        </div>
        <button
          @click="startGoogleOAuth"
          class="btn-ghost border border-emerald-200 dark:border-emerald-900/50 text-emerald-500 gap-2 text-xs"
        >
          <ExternalLink class="h-3.5 w-3.5" />
          إعادة ربط
        </button>
      </div>

      <div v-else class="mt-4 flex items-center justify-between p-4 rounded-2xl bg-slate-50 dark:bg-slate-800/50 border border-slate-200 dark:border-slate-700">
        <div class="flex items-center gap-3">
          <div class="h-10 w-10 rounded-full bg-red-50 dark:bg-red-950/30 flex items-center justify-center">
            <Video class="h-5 w-5 text-red-400" />
          </div>
          <div>
            <p class="text-sm font-bold text-slate-900 dark:text-white">لم يتم ربط YouTube</p>
            <p class="text-xs text-slate-400">اربط قناتك لرفع الفيديوهات تلقائياً</p>
          </div>
        </div>
        <button
          @click="startGoogleOAuth"
          class="btn-primary gap-2 text-xs bg-gradient-to-r from-red-600 to-rose-600 hover:from-red-500 hover:to-rose-500 border-none"
        >
          <ExternalLink class="h-3.5 w-3.5" />
          ربط الآن
        </button>
      </div>
    </div>

    <!-- Analytics + Activity Row -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Analytics Chart -->
      <div class="lg:col-span-2 card p-6 space-y-4">
        <div class="flex items-center justify-between flex-wrap gap-3">
          <h2 class="text-sm font-bold text-slate-700 dark:text-slate-200 flex items-center gap-2">
            <TrendingUp class="h-4 w-4 text-orange-500" />
            التحليلات
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

        <!-- Chart Skeleton -->
        <div v-if="analyticsLoading" class="h-44 animate-pulse bg-slate-100 dark:bg-slate-800 rounded-xl" />

        <!-- Chart -->
        <div v-else-if="analytics?.data?.length" class="h-44 flex items-end gap-2 overflow-x-auto pb-2">
          <div
            v-for="point in analytics.data"
            :key="point.label"
            class="flex-1 min-w-8 flex flex-col items-center gap-1"
          >
            <div class="w-full flex flex-col items-center gap-0.5">
              <div
                :style="{ height: `${Math.max(4, (point.students / maxAnalyticsValue) * 140)}px` }"
                class="w-full bg-orange-400 dark:bg-orange-500 rounded-t-md transition-all duration-500"
                :title="`${point.students} طالب`"
              />
              <div
                :style="{ height: `${Math.max(2, (point.subscriptions / maxAnalyticsValue) * 60)}px` }"
                class="w-full bg-blue-400/60 dark:bg-blue-500/60 rounded-t-sm transition-all duration-500 -mt-0.5"
                :title="`${point.subscriptions} اشتراك`"
              />
            </div>
            <span class="text-[10px] text-slate-400 whitespace-nowrap">{{ point.label }}</span>
          </div>
        </div>
        <div class="flex items-center gap-4 text-xs text-slate-500 dark:text-slate-400">
          <div class="flex items-center gap-1.5">
            <div class="h-3 w-3 rounded-sm bg-orange-400" />
            <span>طلاب</span>
          </div>
          <div class="flex items-center gap-1.5">
            <div class="h-3 w-3 rounded-sm bg-blue-400/60" />
            <span>اشتراكات</span>
          </div>
        </div>
      </div>

      <!-- Activity Feed -->
      <div class="card p-6 space-y-4">
        <h2 class="text-sm font-bold text-slate-700 dark:text-slate-200 flex items-center gap-2">
          <Activity class="h-4 w-4 text-orange-500" />
          آخر النشاطات
        </h2>

        <div v-if="activityLoading" class="space-y-3">
          <div v-for="i in 5" :key="i" class="animate-pulse flex gap-3">
            <div class="h-8 w-8 rounded-full bg-slate-200 dark:bg-slate-700 shrink-0" />
            <div class="flex-1 space-y-1.5 pt-1">
              <div class="h-3 bg-slate-200 dark:bg-slate-700 rounded w-3/4" />
              <div class="h-2.5 bg-slate-200 dark:bg-slate-700 rounded w-1/2" />
            </div>
          </div>
        </div>

        <div v-else-if="activityData?.activities?.length" class="space-y-3">
          <div
            v-for="activity in activityData.activities.slice(0, 8)"
            :key="activity.id"
            class="flex items-start gap-3 group"
          >
            <div class="h-8 w-8 rounded-full bg-slate-100 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 flex items-center justify-center shrink-0 text-sm">
              {{ activityTypeIcons[activity.type] ?? '📌' }}
            </div>
            <div class="flex-1 min-w-0">
              <p class="text-xs text-slate-700 dark:text-slate-300 leading-snug">{{ activity.message }}</p>
              <p class="text-[10px] text-slate-400 mt-0.5">{{ formatTime(activity.created_at) }}</p>
            </div>
          </div>
        </div>

        <div v-else class="text-center py-6">
          <Clock class="h-8 w-8 text-slate-300 dark:text-slate-600 mx-auto mb-2" />
          <p class="text-xs text-slate-400">لا توجد نشاطات حديثة</p>
        </div>
      </div>
    </div>

    <!-- Recent Courses -->
    <div class="card p-6 space-y-4">
      <div class="flex items-center justify-between">
        <h2 class="text-sm font-bold text-slate-700 dark:text-slate-200 flex items-center gap-2">
          <BookMarked class="h-4 w-4 text-orange-500" />
          أحدث المساقات
        </h2>
        <button
          @click="navigateTo('admin.courses')"
          class="text-xs text-orange-500 hover:text-orange-600 flex items-center gap-1 font-medium"
        >
          عرض الكل
          <ArrowLeft class="h-3 w-3" />
        </button>
      </div>

      <div v-if="coursesLoading" class="space-y-3">
        <div v-for="i in 4" :key="i" class="animate-pulse flex gap-4 p-3 rounded-xl">
          <div class="h-10 w-10 rounded-xl bg-slate-200 dark:bg-slate-700 shrink-0" />
          <div class="flex-1 space-y-2 pt-1">
            <div class="h-3.5 bg-slate-200 dark:bg-slate-700 rounded w-2/3" />
            <div class="h-3 bg-slate-200 dark:bg-slate-700 rounded w-1/2" />
          </div>
        </div>
      </div>

      <div v-else-if="recentCourses?.length" class="divide-y divide-slate-100 dark:divide-slate-800">
        <div
          v-for="course in recentCourses"
          :key="course.id"
          @click="navigateTo('admin.course-detail', { id: course.id })"
          class="flex items-center gap-4 py-3 cursor-pointer hover:bg-slate-50 dark:hover:bg-slate-800/50 -mx-2 px-2 rounded-xl transition-colors"
        >
          <div class="h-10 w-10 rounded-xl bg-orange-50 dark:bg-orange-950/40 border border-orange-200 dark:border-orange-900/60 flex items-center justify-center shrink-0">
            <BookOpen class="h-5 w-5 text-orange-500" />
          </div>
          <div class="flex-1 min-w-0">
            <p class="text-sm font-semibold text-slate-900 dark:text-white truncate">{{ course.title }}</p>
            <p class="text-xs text-slate-400 dark:text-slate-500 truncate">
              {{ course.teacher_name ?? 'بدون أستاذ' }} • {{ course.subject_name ?? 'بدون مادة' }}
            </p>
          </div>
          <div class="flex items-center gap-3 shrink-0">
            <span class="flex items-center gap-1 text-xs text-slate-500">
              <Users class="h-3.5 w-3.5" />
              {{ course.students_count }}
            </span>
            <span
              :class="['text-[10px] font-bold px-2 py-0.5 rounded-full', courseStatusClasses[course.status] ?? courseStatusClasses.draft]"
            >
              {{ courseStatusLabels[course.status] ?? course.status }}
            </span>
          </div>
        </div>
      </div>

      <div v-else class="text-center py-10">
        <BookOpen class="h-10 w-10 text-slate-200 dark:text-slate-700 mx-auto mb-3" />
        <p class="text-sm text-slate-400 mb-4">لا توجد مساقات بعد</p>
        <button @click="navigateTo('admin.course-create')" class="btn-primary text-sm">
          إنشاء أول مساق
        </button>
      </div>
    </div>
  </div>
</template>
