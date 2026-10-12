// Control Panel data types — mirror backend app/models/entities.py

export type LessonStatus = "draft" | "published" | "archived";
export type EnrollmentStatus = "active" | "expired" | "suspended" | "cancelled";

export interface Lesson {
  id: string;
  course_id: string;
  section_id?: string | null;
  title: string;
  description?: string | null;
  thumbnail?: string | null;
  youtube_video_id?: string | null;
  youtube_url?: string | null;
  duration?: number | null;
  order: number;
  status: LessonStatus;
  is_free: boolean;
  is_free_preview?: boolean;
  attachments?: string[];
  views: number;
  created_by?: string | null;
  published_at?: string | null;
  created_at: string;
  updated_at: string;
}

export interface Section {
  id: string;
  course_id: string;
  title: string;
  description?: string | null;
  order: number;
  lessons_count: number;
  created_at: string;
  updated_at?: string | null;
}

export type AttachmentType = "pdf" | "image" | "document" | "link" | "file";

export interface Attachment {
  id: string;
  lesson_id: string;
  name: string;
  type: AttachmentType;
  size?: number | null;
  url?: string | null;
  storage_key?: string | null;
  created_by?: string | null;
  created_at: string;
}

export interface Enrollment {
  id: string;
  student_id: string;
  course_id: string;
  subscription_id?: string | null;
  status: EnrollmentStatus;
  started_at: string;
  expires_at?: string | null;
  created_at: string;
  updated_at: string;
}

export interface LessonProgress {
  id: string;
  student_id: string;
  course_id: string;
  lesson_id: string;
  progress: number; // 0..100
  completed: boolean;
  watched_seconds: number;
  duration?: number | null;
  last_position: number;
  updated_at: string;
  created_at: string;
}

export interface CourseProgressSummary {
  student_id: string;
  course_id: string;
  total_lessons: number;
  completed_lessons: number;
  progress_percent: number;
  completed_lesson_ids: string[];
  last_watched_lesson_id?: string | null;
  updated_at: string;
}

export interface Student {
  id: string;
  telegram_id: number;
  username?: string | null;
  first_name?: string | null;
  last_name?: string | null;
  photo_url?: string | null;
  role: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
  // enriched by backend
  courses_count?: number;
  active_subscriptions?: number;
  last_active_at?: string | null;
  enrollment?: Enrollment;
  completed_lessons?: number;
  total_lessons?: number | null;
  enrollments?: Enrollment[];
}

export interface StudentCourse {
  course: {
    id: string;
    title: string;
    teacher_id?: string | null;
    status?: string;
    [k: string]: unknown;
  };
  teacher?: { id: string; first_name?: string | null; username?: string | null; photo_url?: string | null } | null;
  enrollment: Enrollment;
  progress: LessonProgress[];
}

export interface PlayerData {
  lesson: Lesson;
  course: { id: string; title: string; [k: string]: unknown };
  section?: Section | null;
  attachments: Attachment[];
  previous_lesson?: Lesson | null;
  next_lesson?: Lesson | null;
  progress?: LessonProgress | null;
}

export interface PaginatedMeta {
  total: number;
  page: number;
  limit: number;
  total_pages: number;
}

export interface PaginatedResult<T> {
  success: boolean;
  data: T[];
  meta: PaginatedMeta;
}
