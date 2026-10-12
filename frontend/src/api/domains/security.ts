export { checkAccess, reportViolation, createStreamGrant, type ReportViolationOptions } from "../protection";
export type { AccessStatus } from "@/types/protection";
export {
  createPlaybackSession,
  sendHeartbeat,
  reportSecurityTelemetry,
  fetchAdminSecurityEvents,
  fetchAdminSessions,
  revokeSessionAdmin,
  revokeAllUserSessionsAdmin,
  fetchUserSecurityHistory,
  type PlaybackSessionResponse,
  type TelemetryResponse,
  type HeartbeatResponse,
} from "../securityApi";
export { fetchBannedUsers, banUserApi, unbanUserApi, type BannedUserInfo, type BannedUsersResponse } from "../protectionAdmin";
