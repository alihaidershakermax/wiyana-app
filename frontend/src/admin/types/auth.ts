// Control Panel (external) auth types — mirror the backend /auth/* responses.

export type PanelRole = "OWNER" | "TEACHER";

export type PanelStatus = "active" | "suspended";

export interface AdminUser {
  id: string;
  telegramId: number;
  name: string;
  username?: string | null;
  avatar?: string | null;
  role: PanelRole;
  permissions: string[];
  status: PanelStatus;
}

export interface AdminAuthResponse {
  access_token: string;
  refresh_token: string;
  token_type: string;
  user: AdminUser;
}

export interface LoginRequest {
  initData: string;
}

export interface RefreshTokenResponse {
  access_token: string;
  refresh_token: string;
  token_type: string;
}

export interface MeResponse {
  user: AdminUser;
}

export interface LogoutResponse {
  success: boolean;
  message: string;
}

// Permission keys — must match backend app/models/enums.py
export const PERMISSIONS = {
  COURSES_READ: "courses.read",
  COURSES_CREATE: "courses.create",
  COURSES_UPDATE: "courses.update",
  COURSES_DELETE: "courses.delete",
  LESSONS_READ: "lessons.read",
  LESSONS_CREATE: "lessons.create",
  LESSONS_UPDATE: "lessons.update",
  STUDENTS_READ: "students.read",
  STUDENTS_MANAGE: "students.manage",
  TEACHERS_READ: "teachers.read",
  TEACHERS_MANAGE: "teachers.manage",
  SUBSCRIPTIONS_READ: "subscriptions.read",
  SUBSCRIPTIONS_MANAGE: "subscriptions.manage",
  ANALYTICS_READ: "analytics.read",
  SETTINGS_MANAGE: "settings.manage",
} as const;
