export type ViolationType =
  | "screenshot_attempt"
  | "screen_recording"
  | "devtools_open"
  | "fullscreen_exit";

export type ViolationStatus = "ok" | "warning" | "banned";

export interface AccessStatus {
  allowed: boolean;
  banned: boolean;
  protection_enabled: boolean;
  points: number;
  warning: boolean;
  warning_threshold: number;
  ban_threshold: number;
  ban_reason: string | null;
  ban_until: number | null;
  ban_days: number | null;
}

export interface ReportResult {
  success: boolean;
  protection_enabled?: boolean;
  user_id?: string;
  violation_type?: string;
  points: number;
  status: ViolationStatus;
  banned: boolean;
  behavioral?: boolean;
  warning_threshold: number;
  ban_threshold: number;
  ban_reason: string | null;
  ban_until: number | null;
  ban_days: number | null;
}

export interface StreamGrant {
  success: boolean;
  mode: "hls" | "direct";
  video_id: string;
  playlist_url: string;
  expires_in: number;
}
