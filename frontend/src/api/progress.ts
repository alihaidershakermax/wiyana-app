import { apiClient } from "./client";

export interface ProgressUpdatePayload {
  progress: number;
  watched_seconds?: number;
  watch_time?: number;
  last_position?: number;
  duration?: number;
}

export async function saveLessonProgress(lessonId: string, payload: ProgressUpdatePayload): Promise<any> {
  const response = await apiClient.post(`/lessons/${lessonId}/progress`, {
    progress: payload.progress,
    watched_seconds: payload.watched_seconds ?? payload.watch_time ?? 0,
    last_position: payload.last_position ?? 0,
    duration: payload.duration,
  });
  return response.data;
}

export async function fetchLessonProgress(lessonId: string): Promise<any> {
  const response = await apiClient.get(`/lessons/${lessonId}/progress`);
  return response.data?.data ?? response.data;
}

export async function fetchCourseProgress(courseId: string): Promise<any> {
  const response = await apiClient.get(`/courses/${courseId}/progress`);
  return response.data;
}
