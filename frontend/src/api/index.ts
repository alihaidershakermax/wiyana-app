import { adminApi } from "./domains/admin";
import { loginWithTelegram } from "./domains/account";
import { fetchCourseById, fetchCoursesList, fetchDepartments, fetchSubjects, fetchTeachers, fetchCourseLessons, fetchLessonById, fetchVideoById } from "./domains/content";
import { addCourseToFavorites, checkCourseFavorite, fetchMyFavorites, removeCourseFromFavorites, fetchMyNotifications, fetchUnreadNotificationsCount, markAllNotificationsRead, markNotificationRead, fetchCourseProgress, fetchLessonProgress, saveLessonProgress, fetchMyProfile, updateMyProfile } from "./domains/account";
import { fetchHomeDashboard } from "./domains/home";
import { checkAccess, createStreamGrant, reportViolation, createPlaybackSession, fetchAdminSecurityEvents, fetchAdminSessions, fetchUserSecurityHistory, reportSecurityTelemetry, revokeAllUserSessionsAdmin, revokeSessionAdmin, sendHeartbeat } from "./domains/security";
import { apiClient, adminApiClient, appApiClient, getApiError, resolveApiUrl } from "./client";

export { ApiError, type ApiErrorPayload, type ApiClientOptions, type TokenStore, createApiClient, getApiError } from "./client";
export { apiClient, appApiClient, adminApiClient, resolveApiUrl } from "./client";

/**
 * Stable domain facade. New screens should depend on this object instead of
 * importing Axios or constructing endpoint paths inside view components.
 */
export const api = {
  auth: { loginWithTelegram },
  home: { fetchHomeDashboard },
  courses: { fetchCoursesList, fetchCourseById, fetchDepartments, fetchSubjects, fetchTeachers },
  lessons: { fetchCourseLessons, fetchLessonById },
  videos: { fetchVideoById, createStreamGrant },
  favorites: { fetchMyFavorites, addCourseToFavorites, removeCourseFromFavorites, checkCourseFavorite },
  notifications: { fetchMyNotifications, fetchUnreadNotificationsCount, markNotificationRead, markAllNotificationsRead },
  progress: { fetchLessonProgress, fetchCourseProgress, saveLessonProgress },
  profile: { fetchMyProfile, updateMyProfile },
  security: { checkAccess, reportViolation, createPlaybackSession, sendHeartbeat, reportSecurityTelemetry },
  adminSecurity: { fetchAdminSecurityEvents, fetchAdminSessions, fetchUserSecurityHistory, revokeSessionAdmin, revokeAllUserSessionsAdmin },
  admin: adminApi,
};

export type UnifiedApi = typeof api;
