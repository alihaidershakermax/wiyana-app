import { adminApiClient } from './adminApi'
import type {
  DashboardOverview,
  DashboardAnalytics,
  DashboardActivity,
  AnalyticsPeriod,
  RecentCourse,
  RecentStudent,
} from '@/admin/types/dashboard'

export const dashboardApi = {
  async getOverview(): Promise<DashboardOverview> {
    const { data } = await adminApiClient.get<DashboardOverview>('/dashboard/overview')
    return data
  },

  async getAnalytics(period: AnalyticsPeriod = '7d'): Promise<DashboardAnalytics> {
    const { data } = await adminApiClient.get<DashboardAnalytics>('/dashboard/analytics', {
      params: { period },
    })
    return data
  },

  async getActivity(limit = 20): Promise<DashboardActivity> {
    const { data } = await adminApiClient.get<DashboardActivity>('/dashboard/activity', {
      params: { limit },
    })
    return data
  },

  async getRecentCourses(limit = 5): Promise<RecentCourse[]> {
    const { data } = await adminApiClient.get<{ courses: RecentCourse[] }>(
      '/dashboard/recent-courses',
      { params: { limit } }
    )
    return data.courses
  },

  async getRecentStudents(limit = 5): Promise<RecentStudent[]> {
    const { data } = await adminApiClient.get<{ students: RecentStudent[] }>(
      '/dashboard/recent-students',
      { params: { limit } }
    )
    return data.students
  },
}
