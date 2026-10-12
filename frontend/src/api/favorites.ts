import { apiClient } from "./client";

export interface FavoriteItem {
  id: string;
  favorite_id?: string;
  title: string;
  subject_name?: string;
  teacher_name?: string;
  thumbnail_url?: string;
  total_lessons?: number;
  price?: number;
  department?: string;
  stage?: string;
  gradeName?: string;
}

export async function fetchMyFavorites(): Promise<FavoriteItem[]> {
  const response = await apiClient.get("/favorites/");
  return response.data?.data ?? response.data ?? [];
}

export async function addCourseToFavorites(courseId: string): Promise<any> {
  const response = await apiClient.post(`/favorites/${courseId}`);
  return response.data;
}

export async function removeCourseFromFavorites(courseId: string): Promise<any> {
  const response = await apiClient.delete(`/favorites/${courseId}`);
  return response.data;
}

export async function checkCourseFavorite(courseId: string): Promise<boolean> {
  const response = await apiClient.get<{ success: boolean; is_favorite: boolean }>(`/favorites/${courseId}/check`);
  return Boolean(response.data?.is_favorite);
}
