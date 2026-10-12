import { useQuery } from '@tanstack/vue-query'
import { computed, type Ref } from 'vue'
import { teacherApi } from '@/admin/api/teacherApi'
import type { AnalyticsPeriod } from '@/admin/types/dashboard'

export const teacherKeys = {
  all: ['teacher'] as const,
  dashboard: () => [...teacherKeys.all, 'dashboard'] as const,
  analytics: (period: string) => [...teacherKeys.all, 'analytics', period] as const,
  courses: () => [...teacherKeys.all, 'courses'] as const,
}

export function useTeacherDashboard() {
  return useQuery({
    queryKey: teacherKeys.dashboard(),
    queryFn: () => teacherApi.getDashboard(),
    staleTime: 30_000,
    refetchInterval: 30_000,
  })
}

export function useTeacherAnalytics(period: Ref<AnalyticsPeriod> | AnalyticsPeriod = '7d') {
  const periodVal = computed(() => (typeof period === 'string' ? period : period.value))
  return useQuery({
    queryKey: computed(() => teacherKeys.analytics(periodVal.value)),
    queryFn: () => teacherApi.getAnalytics(periodVal.value),
    staleTime: 60_000,
    refetchInterval: 60_000,
  })
}

export function useTeacherCourses() {
  return useQuery({
    queryKey: teacherKeys.courses(),
    queryFn: () => teacherApi.getCourses(),
    staleTime: 30_000,
    refetchInterval: 30_000,
  })
}
