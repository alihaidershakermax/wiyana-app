import { apiClient } from "./client";
import { adminApiClient, getAccessToken } from "@/admin/api/adminApi";

function getClient() {
  return getAccessToken() ? adminApiClient : apiClient;
}

export interface BannedUserInfo {
  user_id: string;
  points: number;
  violations: Record<string, number>;
  banned: boolean;
  ban?: {
    reason: string;
    duration_seconds: number;
    points: number;
    banned_at: number;
    until_ts: number;
    last_violation?: string;
  };
  last_events: Array<{ at: number; id: string }>;
}

export interface BannedUsersResponse {
  success: boolean;
  count: number;
  banned_users: BannedUserInfo[];
}

export async function fetchBannedUsers(): Promise<BannedUsersResponse> {
  const client = getClient();
  const endpoints = [
    "/banned-users",
    "/protection/banned-users",
    "/api/banned-users",
    "/api/protection/banned-users",
  ];
  for (const ep of endpoints) {
    try {
      const res = await client.get<BannedUsersResponse>(ep);
      if (res.data && Array.isArray(res.data.banned_users)) {
        return res.data;
      }
    } catch {
      /* continue */
    }
  }
  return { success: true, count: 0, banned_users: [] };
}

export async function unbanUserApi(userId: string): Promise<{ success: boolean; message: string }> {
  const client = getClient();
  const endpoints = [
    `/unban/${userId}`,
    `/protection/unban/${userId}`,
    `/api/unban/${userId}`,
    `/api/protection/unban/${userId}`,
  ];
  for (const ep of endpoints) {
    try {
      const res = await client.post<{ success: boolean; message: string }>(ep);
      return res.data;
    } catch {
      /* continue */
    }
  }
  return { success: false, message: "فشل رفع الحظر عن المستخدم" };
}

export async function banUserApi(userId: string, reason = "manual_admin"): Promise<{ success: boolean; message: string }> {
  const client = getClient();
  const endpoints = [
    `/ban/${userId}?reason=${encodeURIComponent(reason)}`,
    `/protection/ban/${userId}?reason=${encodeURIComponent(reason)}`,
    `/api/ban/${userId}?reason=${encodeURIComponent(reason)}`,
    `/api/protection/ban/${userId}?reason=${encodeURIComponent(reason)}`,
  ];
  for (const ep of endpoints) {
    try {
      const res = await client.post<{ success: boolean; message: string }>(ep);
      return res.data;
    } catch {
      /* continue */
    }
  }
  return { success: false, message: "فشل حظر المستخدم" };
}
