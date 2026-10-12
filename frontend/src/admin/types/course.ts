export type CourseStatus = 'draft' | 'published' | 'archived'

export interface Course {
  id: string
  title: string
  description?: string | null
  thumbnail_url?: string | null
  subject_id?: string | null
  subject_name?: string | null
  department?: string | null
  department_id?: string | null
  stage?: string | null
  teacher_id?: string | null
  teacher_name?: string | null
  teacher_avatar?: string | null
  teacher_bio?: string | null
  price: number
  status: CourseStatus
  is_published: boolean
  students_count: number
  lessons_count: number
  semester_id?: string | null
  published_at?: string | null
  created_by?: string | null
  created_at: string
  updated_at: string
}

export interface CourseCreate {
  title: string
  description?: string
  thumbnail_url?: string
  subject_id?: string
  department?: string
  stage?: string
  teacher_id?: string
  price?: number
  status?: CourseStatus
}

export interface CourseUpdate {
  title?: string
  description?: string
  thumbnail_url?: string
  subject_id?: string
  department?: string
  stage?: string
  price?: number
  status?: CourseStatus
}

export interface CourseListFilters {
  search?: string
  subject_id?: string
  teacher_id?: string
  department?: string
  stage?: string
  status?: CourseStatus
  page?: number
  limit?: number
}

export interface AssignTeacherPayload {
  teacher_id: string
}

export const COURSE_STATUS_LABELS: Record<CourseStatus, string> = {
  draft: 'مسودة',
  published: 'منشور',
  archived: 'مؤرشف',
}

export const COURSE_STATUS_COLORS: Record<CourseStatus, string> = {
  draft: 'bg-amber-100 text-amber-700 dark:bg-amber-950/40 dark:text-amber-400',
  published: 'bg-emerald-100 text-emerald-700 dark:bg-emerald-950/40 dark:text-emerald-400',
  archived: 'bg-slate-100 text-slate-600 dark:bg-slate-800 dark:text-slate-400',
}
