import { apiClient } from "./client";

export async function fetchMyNotifications(): Promise<any[]> {
  const response = await apiClient.get("/notifications/my-notifications");
  return response.data?.data ?? response.data ?? [];
}

export async function markNotificationRead(notificationId: string): Promise<any> {
  const response = await apiClient.post(`/notifications/${notificationId}/read`);
  return response.data;
}

export async function markAllNotificationsRead(): Promise<any> {
  const response = await apiClient.post("/notifications/read-all");
  return response.data;
}

export async function fetchUnreadNotificationsCount(): Promise<number> {
  const response = await apiClient.get<{ success: boolean; unread_count: number }>("/notifications/unread");
  return response.data?.unread_count ?? 0;
}
