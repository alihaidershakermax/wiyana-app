import { useQuery } from '@tanstack/vue-query'
import { computed, type Ref } from 'vue'
import { dashboardApi } from '@/admin/api/dashboardApi'
import type { AnalyticsPeriod } from '@/admin/types/dashboard'

export const dashboardKeys = {
  all: ['dashboard'] as const,
  overview: () => [...dashboardKeys.all, 'overview'] as const,
  analytics: (period: string) => [...dashboardKeys.all, 'analytics', period] as const,
  activity: () => [...dashboardKeys.all, 'activity'] as const,
  recentCourses: () => [...dashboardKeys.all, 'recent-courses'] as const,
  recentStudents: () => [...dashboardKeys.all, 'recent-students'] as const,
}

export function useDashboardOverview() {
  return useQuery({
    queryKey: dashboardKeys.overview(),
    queryFn: () => dashboardApi.getOverview(),
    staleTime: 30_000,
    refetchInterval: 30_000,
    retry: 2,
  })
}

export function useDashboardAnalytics(period: Ref<AnalyticsPeriod> | AnalyticsPeriod = '7d') {
  const periodVal = computed(() => (typeof period === 'string' ? period : period.value))
  return useQuery({
    queryKey: computed(() => dashboardKeys.analytics(periodVal.value)),
    queryFn: () => dashboardApi.getAnalytics(periodVal.value),
    staleTime: 60_000,
    refetchInterval: 60_000,
    retry: 1,
  })
}

export function useDashboardActivity(limit = 20) {
  return useQuery({
    queryKey: dashboardKeys.activity(),
    queryFn: () => dashboardApi.getActivity(limit),
    staleTime: 30_000,
    refetchInterval: 30_000,
  })
}

export function useRecentCourses(limit = 5) {
  return useQuery({
    queryKey: dashboardKeys.recentCourses(),
    queryFn: () => dashboardApi.getRecentCourses(limit),
    staleTime: 60_000,
    refetchInterval: 60_000,
  })
}

export function useRecentStudents(limit = 5) {
  return useQuery({
    queryKey: dashboardKeys.recentStudents(),
    queryFn: () => dashboardApi.getRecentStudents(limit),
    staleTime: 60_000,
    refetchInterval: 60_000,
  })
}
