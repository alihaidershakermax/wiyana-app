import { adminApiClient } from './adminApi'
import type { TeacherDashboard, DashboardAnalytics, AnalyticsPeriod } from '@/admin/types/dashboard'
import type { Course } from '@/admin/types/course'

export const teacherApi = {
  async getDashboard(): Promise<TeacherDashboard> {
    const { data } = await adminApiClient.get<TeacherDashboard>('/teacher/dashboard')
    return data
  },

  async getAnalytics(period: AnalyticsPeriod = '7d'): Promise<DashboardAnalytics> {
    const { data } = await adminApiClient.get<DashboardAnalytics>('/teacher/analytics', {
      params: { period },
    })
    return data
  },

  async getCourses(): Promise<Course[]> {
    const { data } = await adminApiClient.get<{ courses: Course[] }>('/teacher/courses')
    return data.courses
  },
}
