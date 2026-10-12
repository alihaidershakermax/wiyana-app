import { apiClient } from "@/api/client";
import { Platform } from "@/security/platformCapabilities";

export interface PlaybackSessionResponse {
  success: boolean;
  session: {
    sessionId: string;
    userId: string;
    videoId: string;
    courseId: string;
    platform?: string;
    status: string;
    expiresAt: number;
  };
  token: string;
  watermark: {
    displayName: string;
    telegramId: string;
    privacyId: string;
    sessionId: string;
  };
}

export interface TelemetryResponse {
  success: boolean;
  eventId: string;
  riskScore: number;
  riskLevel: "NORMAL" | "SUSPICIOUS" | "HIGH_RISK" | "CRITICAL";
}

export interface HeartbeatResponse {
  success: boolean;
  status: string;
  riskScore?: number;
  riskLevel?: string;
}

export async function createPlaybackSession(
  videoId: string,
  courseId?: string,
  platform?: string
): Promise<PlaybackSessionResponse> {
  const telegramInitData = (window as any).Telegram?.WebApp?.initData || undefined;
  const res = await apiClient.post<PlaybackSessionResponse>("/playback/session", {
    videoId,
    courseId,
    platform,
    telegramInitData,
  });
  return res.data;
}

export async function sendHeartbeat(
  sessionId: string,
  videoId: string,
  currentTime: number,
  playing: boolean,
  riskScore: number = 0,
  platform?: string
): Promise<HeartbeatResponse> {
  const res = await apiClient.post<HeartbeatResponse>("/playback/heartbeat", {
    sessionId,
    videoId,
    currentTime,
    playing,
    riskScore,
    platform,
  });
  return res.data;
}

export async function reportSecurityTelemetry(
  eventType: string,
  videoId: string,
  sessionId: string,
  confidence: number = 1.0,
  metadata: Record<string, any> = {},
  courseId?: string
): Promise<TelemetryResponse> {
  const res = await apiClient.post<TelemetryResponse>("/security/telemetry", {
    eventType,
    videoId,
    sessionId,
    confidence,
    metadata,
    courseId,
  });
  return res.data;
}

export async function fetchAdminSecurityEvents(limit: number = 50): Promise<any[]> {
  const res = await apiClient.get<{ success: boolean; events: any[] }>("/security/admin/events", {
    params: { limit },
  });
  return res.data.events || [];
}

export async function fetchAdminSessions(): Promise<any[]> {
  const res = await apiClient.get<{ success: boolean; sessions: any[] }>("/security/admin/sessions");
  return res.data.sessions || [];
}

export async function revokeSessionAdmin(sessionId: string, reason: string = "admin_action"): Promise<boolean> {
  const res = await apiClient.post<{ success: boolean }>("/security/admin/revoke-session", {
    sessionId,
    reason,
  });
  return res.data.success;
}

export async function revokeAllUserSessionsAdmin(userId: string): Promise<boolean> {
  const res = await apiClient.post<{ success: boolean }>("/security/admin/revoke-all", { userId });
  return res.data.success;
}

export async function fetchUserSecurityHistory(userId: string, limit: number = 50): Promise<any[]> {
  const res = await apiClient.get<{ success: boolean; events: any[] }>(
    `/security/admin/users/${encodeURIComponent(userId)}/history`,
    { params: { limit } }
  );
  return res.data.events || [];
}