import { adminApiClient } from "./adminApi";
import type {
  Lesson,
  PlayerData,
  Section,
  PaginatedResult,
  Student,
  StudentCourse,
  LessonProgress,
  CourseProgressSummary,
  Attachment,
  Enrollment,
} from "@/admin/types/panel";

interface ApiResponse<T> {
  success: boolean;
  data: T;
}

// ---------------------------------------------------------------------------
// Query keys (used by TanStack Query)
// ---------------------------------------------------------------------------
export const lessonKeys = {
  all: ["lessons"] as const,
  course: (courseId: string) => ["lessons", "course", courseId] as const,
  section: (sectionId: string) => ["lessons", "section", sectionId] as const,
  detail: (lessonId: string) => ["lesson-detail", lessonId] as const,
  player: (lessonId: string) => ["lesson-player", lessonId] as const,
  attachments: (lessonId: string) => ["attachments", lessonId] as const,
};

export const studentKeys = {
  all: ["students"] as const,
  list: (params: Record<string, unknown>) => ["students", "list", params] as const,
  detail: (studentId: string) => ["student-detail", studentId] as const,
  courses: (studentId: string) => ["student-courses", studentId] as const,
  progress: (studentId: string) => ["student-progress", studentId] as const,
  courseProgress: (studentId: string, courseId: string) =>
    ["student-course-progress", studentId, courseId] as const,
  activity: (studentId: string) => ["student-activity", studentId] as const,
  courseStudents: (courseId: string) => ["course-students", courseId] as const,
};

export const sectionKeys = {
  course: (courseId: string) => ["sections", courseId] as const,
};

// ---------------------------------------------------------------------------
// Sections
// ---------------------------------------------------------------------------
export const sectionsApi = {
  async getByCourse(courseId: string): Promise<Section[]> {
    const { data } = await adminApiClient.get<ApiResponse<Section[]>>(
      `/courses/${courseId}/sections`
    );
    return data.data ?? data;
  },
};

// ---------------------------------------------------------------------------
// Lessons
// ---------------------------------------------------------------------------
export interface LessonInput {
  title: string;
  description?: string;
  section_id?: string;
  thumbnail?: string;
  youtube_video_id?: string;
  youtube_url?: string;
  duration?: number;
  order?: number;
  status?: string;
  is_free?: boolean;
}

export const lessonsApi = {
  async getByCourse(courseId: string): Promise<Lesson[]> {
    const { data } = await adminApiClient.get<ApiResponse<Lesson[]>>(
      `/courses/${courseId}/lessons`
    );
    return data.data ?? data;
  },

  async getBySection(sectionId: string): Promise<Lesson[]> {
    const { data } = await adminApiClient.get<ApiResponse<Lesson[]>>(
      `/sections/${sectionId}/lessons`
    );
    return data.data ?? data;
  },

  async getById(lessonId: string): Promise<Lesson> {
    const { data } = await adminApiClient.get<ApiResponse<Lesson>>(
      `/lessons/${lessonId}`
    );
    return data.data ?? data;
  },

  async getPlayer(lessonId: string): Promise<PlayerData> {
    const { data } = await adminApiClient.get<ApiResponse<PlayerData>>(
      `/lessons/${lessonId}/player`
    );
    return data.data ?? data;
  },

  async create(sectionId: string, input: LessonInput): Promise<Lesson> {
    const { data } = await adminApiClient.post<ApiResponse<Lesson>>(
      `/sections/${sectionId}/lessons`,
      input
    );
    return data.data ?? data;
  },

  async createForCourse(courseId: string, input: LessonInput): Promise<Lesson> {
    const { data } = await adminApiClient.post<ApiResponse<Lesson>>(
      `/courses/${courseId}/lessons`,
      input
    );
    return data.data ?? data;
  },

  async update(lessonId: string, input: Partial<LessonInput>): Promise<Lesson> {
    const { data } = await adminApiClient.patch<ApiResponse<Lesson>>(
      `/lessons/${lessonId}`,
      input
    );
    return data.data ?? data;
  },

  async delete(lessonId: string): Promise<boolean> {
    const { data } = await adminApiClient.delete<{ success: boolean }>(
      `/lessons/${lessonId}`
    );
    return data.success;
  },

  async publish(lessonId: string): Promise<Lesson> {
    const { data } = await adminApiClient.post<ApiResponse<Lesson>>(
      `/lessons/${lessonId}/publish`
    );
    return data.data ?? data;
  },

  async archive(lessonId: string): Promise<Lesson> {
    const { data } = await adminApiClient.post<ApiResponse<Lesson>>(
      `/lessons/${lessonId}/archive`
    );
    return data.data ?? data;
  },

  async reorder(sectionId: string, order: string[]): Promise<Lesson[]> {
    const { data } = await adminApiClient.post<ApiResponse<Lesson[]>>(
      `/sections/${sectionId}/lessons/reorder`,
      { order }
    );
    return data.data ?? data;
  },

  async uploadVideo(
    lessonId: string,
    file: File,
    opts: { title?: string; description?: string; privacyStatus?: string } = {},
    onProgress?: (percent: number) => void
  ): Promise<{ lesson: Lesson; youtube: Record<string, unknown> }> {
    const form = new FormData();
    form.append("file", file);
    if (opts.title) form.append("title", opts.title);
    if (opts.description) form.append("description", opts.description);
    form.append("privacy_status", opts.privacyStatus ?? "unlisted");

    const { data } = await adminApiClient.post<
      ApiResponse<Lesson> & { youtube: Record<string, unknown> }
    >(`/lessons/${lessonId}/upload-video`, form, {
      headers: { "Content-Type": "multipart/form-data" },
      timeout: 0,
      onUploadProgress: (e) => {
        if (e.total && onProgress) {
          onProgress(Math.round((e.loaded / e.total) * 100));
        }
      },
    });
    return { lesson: data.data ?? data, youtube: data.youtube ?? {} };
  },

  async getYouTubeStatus(): Promise<{
    configured: boolean;
    is_authenticated: boolean;
    token_file_exists: boolean;
    token_in_db: boolean;
    simulation_mode: boolean;
    channel_id?: string;
    channel_name?: string;
    channel_avatar?: string;
    subscribers?: string;
    videos_count?: string;
  }> {
    const { data } = await adminApiClient.get<
      ApiResponse<{
        configured: boolean;
        is_authenticated: boolean;
        token_file_exists: boolean;
        token_in_db: boolean;
        simulation_mode: boolean;
        channel_id?: string;
        channel_name?: string;
        channel_avatar?: string;
        subscribers?: string;
        videos_count?: string;
      }>
    >(`/youtube/status`);
    return data.data ?? data;
  },
};

export const youtubeApi = {
  async getStatus() {
    return lessonsApi.getYouTubeStatus();
  },
  async exchangeCode(code: string) {
    const { data } = await adminApiClient.post("/youtube/oauth/exchange-code", { code });
    return data;
  },
  async disconnect() {
    const { data } = await adminApiClient.delete("/youtube/disconnect");
    return data;
  },
};

// ---------------------------------------------------------------------------
// Attachments
// ---------------------------------------------------------------------------
export const attachmentsApi = {
  async getByLesson(lessonId: string): Promise<Attachment[]> {
    const { data } = await adminApiClient.get<ApiResponse<Attachment[]>>(
      `/lessons/${lessonId}/attachments`
    );
    return data.data ?? data;
  },

  async upload(
    lessonId: string,
    file: File,
    name?: string
  ): Promise<Attachment> {
    const form = new FormData();
    form.append("file", file);
    if (name) form.append("name", name);
    const { data } = await adminApiClient.post<ApiResponse<Attachment>>(
      `/lessons/${lessonId}/attachments`,
      form,
      { headers: { "Content-Type": "multipart/form-data" } }
    );
    return data.data ?? data;
  },

  async addLink(lessonId: string, url: string, name?: string): Promise<Attachment> {
    const form = new FormData();
    form.append("url", url);
    if (name) form.append("name", name);
    const { data } = await adminApiClient.post<ApiResponse<Attachment>>(
      `/lessons/${lessonId}/attachments`,
      form
    );
    return data.data ?? data;
  },

  async remove(attachmentId: string): Promise<boolean> {
    const { data } = await adminApiClient.delete<{ success: boolean }>(
      `/attachments/${attachmentId}`
    );
    return data.success;
  },
};

// ---------------------------------------------------------------------------
// Students
// ---------------------------------------------------------------------------
export interface StudentFilters {
  [key: string]: unknown;
  search?: string;
  status?: string;
  course_id?: string;
  page?: number;
  limit?: number;
}

export const studentsApi = {
  async list(filters: StudentFilters = {}): Promise<PaginatedResult<Student>> {
    const { data } = await adminApiClient.get<PaginatedResult<Student>>(
      "/students/",
      { params: filters }
    );
    return data;
  },

  async getById(studentId: string): Promise<Student> {
    const { data } = await adminApiClient.get<ApiResponse<Student>>(
      `/students/${studentId}`
    );
    return data.data ?? data;
  },

  async getCourses(studentId: string): Promise<StudentCourse[]> {
    const { data } = await adminApiClient.get<ApiResponse<StudentCourse[]>>(
      `/students/${studentId}/courses`
    );
    return data.data ?? data;
  },

  async getProgress(studentId: string): Promise<LessonProgress[]> {
    const { data } = await adminApiClient.get<ApiResponse<LessonProgress[]>>(
      `/students/${studentId}/progress`
    );
    return data.data ?? data;
  },

  async getCourseProgress(studentId: string, courseId: string): Promise<CourseProgressSummary> {
    const { data } = await adminApiClient.get<ApiResponse<CourseProgressSummary>>(
      `/students/${studentId}/courses/${courseId}/progress`
    );
    return data.data ?? data;
  },

  async getActivity(studentId: string): Promise<LessonProgress[]> {
    const { data } = await adminApiClient.get<ApiResponse<LessonProgress[]>>(
      `/students/${studentId}/activity`
    );
    return data.data ?? data;
  },

  async getCourseStudents(courseId: string): Promise<Student[]> {
    const { data } = await adminApiClient.get<ApiResponse<Student[]>>(
      `/courses/${courseId}/students`
    );
    return data.data ?? data;
  },

  async suspend(studentId: string): Promise<Student> {
    const { data } = await adminApiClient.post<ApiResponse<Student>>(
      `/students/${studentId}/suspend`
    );
    return data.data ?? data;
  },

  async activate(studentId: string): Promise<Student> {
    const { data } = await adminApiClient.post<ApiResponse<Student>>(
      `/students/${studentId}/activate`
    );
    return data.data ?? data;
  },
};

// ---------------------------------------------------------------------------
// Progress
// ---------------------------------------------------------------------------
export interface ProgressInput {
  progress: number;
  watched_seconds: number;
  last_position: number;
  duration?: number;
}

export const progressApi = {
  async getLesson(studentId: string, lessonId: string): Promise<LessonProgress | null> {
    const { data } = await adminApiClient.get<ApiResponse<LessonProgress | null>>(
      `/lessons/${lessonId}/progress`
    );
    return data.data ?? null;
  },

  async save(lessonId: string, input: ProgressInput): Promise<LessonProgress> {
    const { data } = await adminApiClient.post<ApiResponse<LessonProgress>>(
      `/lessons/${lessonId}/progress`,
      input
    );
    return data.data ?? data;
  },
};

// ---------------------------------------------------------------------------
// Enrollments (read-only for panel usage)
// ---------------------------------------------------------------------------
export const enrollmentsApi = {
  async getByStudent(studentId: string): Promise<Enrollment[]> {
    const { data } = await adminApiClient.get<ApiResponse<Enrollment[]>>(
      `/students/${studentId}/enrollments`
    );
    return data.data ?? data;
  },
};
