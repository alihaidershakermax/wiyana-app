import { apiClient } from "./client";

export interface LessonApiResponse {
  id: string;
  course_id: string;
  title: string;
  description?: string | null;
  order: number;
  video_id?: string | null;
  attachments: string[];
  is_free_preview: boolean;
  created_at: string;
}

export async function fetchCourseLessons(courseId: string): Promise<LessonApiResponse[]> {
  const response = await apiClient.get<LessonApiResponse[]>(`/lessons/course/${courseId}`);
  return response.data;
}

export async function fetchLessonById(lessonId: string): Promise<LessonApiResponse> {
  const response = await apiClient.get<LessonApiResponse>(`/lessons/${lessonId}`);
  return response.data;
}
