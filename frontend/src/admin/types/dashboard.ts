export interface DashboardOverview {
  total_students: number
  active_students: number
  total_courses: number
  total_subjects: number
  total_lessons: number
  total_views: number
  active_subscriptions: number
  used_subscription_codes: number
}

export type AnalyticsPeriod = 'day' | '7d' | '30d' | '3m' | '1y'

export interface AnalyticsPoint {
  label: string
  students: number
  courses: number
  views: number
  subscriptions: number
}

export interface DashboardAnalytics {
  period: AnalyticsPeriod
  data: AnalyticsPoint[]
}

export interface ActivityItem {
  id: string
  type: string
  message: string
  meta?: Record<string, any>
  created_at: string
}

export interface DashboardActivity {
  activities: ActivityItem[]
}

export interface RecentCourse {
  id: string
  title: string
  teacher_name?: string | null
  subject_name?: string | null
  status: string
  students_count: number
  created_at: string
}

export interface RecentStudent {
  id: string
  name: string
  username?: string | null
  avatar?: string | null
  created_at: string
}

export interface TeacherDashboard {
  my_courses: number
  my_students: number
  my_lessons: number
  total_views: number
  completion_rate: number
}

export interface PaginatedMeta {
  total: number
  page: number
  limit: number
  total_pages: number
}

export interface PaginatedResponse<T> {
  success: boolean
  data: T[]
  meta: PaginatedMeta
}
