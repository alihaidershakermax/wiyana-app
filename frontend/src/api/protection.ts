import { appApiClient } from "./client";
import type { AccessStatus, ReportResult, StreamGrant, ViolationType } from "@/types/protection";
import { playbackTokenManager } from "@/video/playbackToken";

export async function checkAccess(): Promise<AccessStatus> {
  const res = await appApiClient.get<AccessStatus>("/api/check-access");
  return res.data;
}

export interface ReportViolationOptions {
  lessonId?: string;
  meta?: Record<string, unknown>;
}

export async function reportViolation(
  violationType: ViolationType,
  opts: ReportViolationOptions = {}
): Promise<ReportResult> {
  const tgUser =
    typeof window !== "undefined" ? (window as any).Telegram?.WebApp?.initDataUnsafe?.user : undefined;
  const res = await appApiClient.post<ReportResult>("/api/report-violation", {
    violation_type: violationType,
    user_id: tgUser?.id ?? undefined,
    lesson_id: opts.lessonId,
    meta: opts.meta ?? {},
  });
  return res.data;
}

// Stable per-device fingerprint: stored once in sessionStorage so the server
// can bind the playback grant to this exact browser session. Captured links
// then only decrypt on the issuing device.
export function deviceFingerprint(): string {
  const key = "wiyana.device.id";
  let id = typeof window !== "undefined" ? sessionStorage.getItem(key) : null;
  if (!id) {
    const chars = "abcdefghijklmnopqrstuvwxyz0123456789";
    id = Array.from({ length: 32 }, () => chars[Math.floor(Math.random() * chars.length)]).join("");
    if (typeof window !== "undefined") sessionStorage.setItem(key, id);
  }
  return id;
}

export async function createStreamGrant(videoId: string): Promise<StreamGrant> {
  const playbackToken = playbackTokenManager.getToken();
  const res = await appApiClient.post<StreamGrant>(
    `/api/stream/hls/${encodeURIComponent(videoId)}`,
    undefined,
    { headers: { ...(playbackToken ? { "X-Playback-Token": playbackToken } : {}), "X-Device-ID": deviceFingerprint() } }
  );
  return res.data;
}
