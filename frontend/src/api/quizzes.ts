import { apiClient } from "./client";

export interface QuizApiResponse {
  id: string;
  lesson_id: string;
  title: string;
  questions_count: number;
  duration_minutes?: number;
}

export async function fetchLessonQuiz(lessonId: string): Promise<QuizApiResponse | null> {
  try {
    const response = await apiClient.get<QuizApiResponse>(`/quizzes/lesson/${lessonId}`);
    return response.data;
  } catch {
    return null;
  }
}
