<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { TrendingUp, Users, BookOpen, Clock, Eye, CheckCircle2, RefreshCw, Award } from '@lucide/vue'
import { adminApiClient } from '@/admin/api/adminApi'
import { useAdminStore } from '@/admin/stores/useAdminStore'
import StatePanel from '@/admin/components/StatePanel.vue'

interface AnalyticsData {
  time_range: string
  total_students?: number
  active_students?: number
  new_students?: number
  total_courses: number
  total_lessons: number
  active_subscriptions?: number
  views: number
  watch_time_hours: number
  completion_rate: number
  top_courses?: Array<{ id: string; title: string; students_count: number; status: string }>
}

const adminStore = useAdminStore()
const isOwner = computed(() => adminStore.isOwner)

const selectedRange = ref('30d')
const data = ref<AnalyticsData | null>(null)
const loading = ref(false)
const errorMsg = ref<string | null>(null)

const rangeOptions = [
  { label: 'اليوم', value: 'today' },
  { label: '7 أيام', value: '7d' },
  { label: '30 يوماً', value: '30d' },
  { label: '3 أشهر', value: '3m' },
  { label: '6 أشهر', value: '6m' },
  { label: 'سنة', value: '1y' },
]

async function fetchAnalytics() {
  loading.value = true
  errorMsg.value = null
  try {
    const endpoint = isOwner.value ? '/analytics/owner' : '/analytics/teacher'
    const res = await adminApiClient.get(endpoint, { params: { range: selectedRange.value } })
    data.value = res.data.data ?? res.data
  } catch (err: unknown) {
    errorMsg.value = (err as Error)?.message || 'فشل تحميل التحليلات'
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchAnalytics()
})
</script>

<template>
  <div class="space-y-6">
    <div class="flex flex-wrap items-center justify-between gap-3">
      <div>
        <h1 class="text-xl font-extrabold flex items-center gap-2">
          <TrendingUp class="h-5 w-5 text-orange-500" />
          {{ isOwner ? 'تحليلات المنصة الشاملة' : 'تحليلات المساقات والطلاب' }}
        </h1>
        <p class="text-xs text-slate-500 dark:text-slate-400 mt-1">
          مؤشرات الأداء، المشاهدات، وقت المشاهدة، ونسبة الإكمال
        </p>
      </div>

      <div class="flex items-center gap-2">
        <select v-model="selectedRange" class="input text-xs w-36" @change="fetchAnalytics">
          <option v-for="opt in rangeOptions" :key="opt.value" :value="opt.value">
            {{ opt.label }}
          </option>
        </select>
        <button class="btn-secondary text-xs py-2" @click="fetchAnalytics">
          <RefreshCw class="h-4 w-4" />
        </button>
      </div>
    </div>

    <StatePanel
      :loading="loading"
      :error="errorMsg"
      :empty="!data"
      empty-title="لا توجد بيانات تحليلات"
      :skeleton-count="4"
      @retry="fetchAnalytics"
    >
      <template v-if="data">
        <!-- Stat Grid -->
        <div class="grid grid-cols-2 sm:grid-cols-4 gap-4">
          <div class="card p-4 space-y-2">
            <div class="flex items-center justify-between text-slate-400">
              <span class="text-xs font-bold">إجمالي الطلاب</span>
              <Users class="h-4 w-4 text-orange-500" />
            </div>
            <p class="text-2xl font-black">{{ data.total_students ?? 0 }}</p>
            <p v-if="data.new_students !== undefined" class="text-[11px] text-emerald-600 font-bold">
              +{{ data.new_students }} طالب جديد
            </p>
          </div>

          <div class="card p-4 space-y-2">
            <div class="flex items-center justify-between text-slate-400">
              <span class="text-xs font-bold">عدد المشاهدات</span>
              <Eye class="h-4 w-4 text-blue-500" />
            </div>
            <p class="text-2xl font-black">{{ data.views }}</p>
            <p class="text-[11px] text-slate-400">إجمالي المشاهدات</p>
          </div>

          <div class="card p-4 space-y-2">
            <div class="flex items-center justify-between text-slate-400">
              <span class="text-xs font-bold">ساعات المشاهدة</span>
              <Clock class="h-4 w-4 text-purple-500" />
            </div>
            <p class="text-2xl font-black">{{ data.watch_time_hours }} س</p>
            <p class="text-[11px] text-slate-400">ساعات المشاهدة الفعلية</p>
          </div>

          <div class="card p-4 space-y-2">
            <div class="flex items-center justify-between text-slate-400">
              <span class="text-xs font-bold">نسبة الإكمال</span>
              <CheckCircle2 class="h-4 w-4 text-emerald-500" />
            </div>
            <p class="text-2xl font-black text-emerald-600 dark:text-emerald-400">{{ data.completion_rate }}%</p>
            <p class="text-[11px] text-slate-400">متوسط إكمال المحاضرات</p>
          </div>
        </div>

        <!-- Top Courses Section -->
        <div v-if="data.top_courses && data.top_courses.length" class="card p-5 space-y-4">
          <h3 class="text-sm font-extrabold flex items-center gap-2">
            <Award class="h-4 w-4 text-orange-500" />
            المساقات الأكثر إقبالاً
          </h3>
          <div class="divide-y divide-slate-100 dark:divide-slate-800">
            <div v-for="c in data.top_courses" :key="c.id" class="py-3 flex items-center justify-between">
              <div class="flex items-center gap-3">
                <BookOpen class="h-4 w-4 text-slate-400" />
                <span class="text-xs font-bold">{{ c.title }}</span>
              </div>
              <span class="text-xs font-mono font-bold text-orange-600 bg-orange-50 dark:bg-orange-950/40 px-2 py-0.5 rounded-full">
                {{ c.students_count }} طالب
              </span>
            </div>
          </div>
        </div>
      </template>
    </StatePanel>
  </div>
</template>
