import { apiClient } from "./client";

export interface VideoApiResponse {
  id: string;
  title?: string;
  youtube_video_id?: string | null;
  youtube_url?: string | null;
  file_path?: string | null;
  file_size?: number | null;
  status?: string;
  duration_seconds?: number | null;
  created_at?: string;
}

export async function fetchVideoById(videoId: string): Promise<VideoApiResponse> {
  const response = await apiClient.get<VideoApiResponse>(`/videos/${videoId}`);
  return response.data;
}
