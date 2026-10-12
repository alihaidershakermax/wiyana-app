import { AxiosProgressEvent, CancelTokenSource } from "axios";
import { apiClient } from "@/api/client";

export interface UploadItem {
  id: string;
  file: File;
  title: string;
  description?: string;
  courseId: string;
  courseTitle?: string;
  subjectId?: string;
  lectureNumber: number;
  isFreePreview: boolean;
  type: "video" | "pdf" | "attachment";
  progress: number;
  speedMbps: number;
  uploadedBytes: number;
  totalBytes: number;
  status: "queued" | "uploading" | "processing" | "completed" | "failed" | "paused";
  error?: string;
  youtubeId?: string;
  streamUrl?: string;
  cancelTokenSource?: CancelTokenSource;
  createdAt: string;
}

export interface MediaAsset {
  id: string;
  title: string;
  filename: string;
  fileSize: number;
  fileType: string;
  durationSeconds?: number;
  status: "pending" | "uploading" | "processing" | "completed" | "failed";
  courseTitle: string;
  uploaderName: string;
  youtubeVideoId?: string;
  streamUrl?: string;
  pdfAttachmentUrl?: string;
  isFreePreview: boolean;
  createdAt: string;
}

export const uploadService = {
  // 1. Fetch all media assets & processing videos
  async getMediaAssets(): Promise<MediaAsset[]> {
    const response = await apiClient.get<MediaAsset[]>("/videos/");
    return response.data;
  },

  // 2. Upload video with live progress reporting
  async uploadVideo(
    file: File,
    title: string,
    onProgress?: (progressEvent: AxiosProgressEvent) => void,
    cancelSource?: CancelTokenSource
  ): Promise<{ video_id: string; status: string }> {
    const formData = new FormData();
    formData.append("file", file);
    formData.append("title", title);

    const response = await apiClient.post("/videos/upload", formData, {
      headers: {
        "Content-Type": "multipart/form-data",
      },
      onUploadProgress: onProgress,
      cancelToken: cancelSource?.token,
    });
    return response.data;
  },

  // 3. Upload PDF / attachment
  async uploadAttachment(
    file: File,
    onProgress?: (progressEvent: AxiosProgressEvent) => void
  ): Promise<{ url: string; filename: string }> {
    const formData = new FormData();
    formData.append("file", file);

    const response = await apiClient.post("/uploads/file", formData, {
      headers: {
        "Content-Type": "multipart/form-data",
      },
      onUploadProgress: onProgress,
    });
    return response.data;
  },

  // 4. Retry video processing
  async retryProcessing(videoId: string): Promise<boolean> {
    await apiClient.post(`/videos/${videoId}/retry`);
    return true;
  },

  // 5. Delete video/asset
  async deleteAsset(videoId: string): Promise<boolean> {
    await apiClient.delete(`/videos/${videoId}`);
    return true;
  },
};
