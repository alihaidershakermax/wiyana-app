import { defineStore } from "pinia";
import { ref, computed } from "vue";
import {
  adminApi,
  getAccessToken,
  getRefreshToken,
  setTokens,
  clearTokens,
} from "@/admin/api/adminApi";
import type { AdminUser } from "@/admin/types/auth";

export type AuthState =
  | "loading" // restoring session / fetching /auth/me
  | "authenticated"
  | "unauthenticated"
  | "unauthorized"
  | "session_expired"
  | "network_error"
  | "server_error";

export const useAdminStore = defineStore("admin", () => {
  // State
  const user = ref<AdminUser | null>(null);
  const authState = ref<AuthState>("loading");
  const isLoading = ref(false);
  const error = ref<string | null>(null);

  // Computed
  const isAuthenticated = computed(() => authState.value === "authenticated");
  const userRole = computed(() => user.value?.role ?? null);
  const userName = computed(() => user.value?.name ?? "");
  const userAvatar = computed(() => user.value?.avatar ?? null);
  const permissions = computed(() => user.value?.permissions ?? []);

  const isOwner = computed(() => userRole.value === "OWNER");
  const isTeacher = computed(() => userRole.value === "TEACHER");

  // Permission helpers
  const hasPermission = (permission: string): boolean => {
    if (isOwner.value) return true;
    return permissions.value.includes(permission) || permissions.value.includes("*");
  };

  const hasAnyPermission = (perms: string[]): boolean =>
    isOwner.value || perms.some((p) => hasPermission(p));

  const hasAllPermissions = (perms: string[]): boolean =>
    isOwner.value || perms.every((p) => hasPermission(p));

  const hasRole = (role: string | string[]): boolean => {
    if (isOwner.value) return true;
    const roles = Array.isArray(role) ? role : [role];
    return roles.includes(userRole.value ?? "");
  };

  // Actions
  function applySession(authResponse: { access_token: string; refresh_token: string; user: AdminUser }) {
    setTokens(authResponse.access_token, authResponse.refresh_token);
    user.value = authResponse.user;
    localStorage.setItem("wiyana_admin_user", JSON.stringify(authResponse.user));
    authState.value = "authenticated";
    error.value = null;
  }

  function setErrorState(state: AuthState, message: string) {
    authState.value = state;
    error.value = message;
  }

  async function login(initData: string): Promise<boolean> {
    isLoading.value = true;
    error.value = null;
    try {
      const response = await adminApi.loginTelegram(initData);
      applySession(response);
      return true;
    } catch (err: any) {
      const status = err?.response?.status;
      if (status === 403) {
        setErrorState("unauthorized", err?.response?.data?.detail || "ليس لديك صلاحية للوصول إلى لوحة التحكم");
      } else if (status === 401) {
        setErrorState("unauthorized", "فشل التحقق من هوية Telegram");
      } else if (err?.code === "ECONNABORTED" || !err?.response) {
        setErrorState("network_error", "تعذر الاتصال بالخادم. تأكد من تشغيل الـ Backend.");
      } else if (status >= 500) {
        setErrorState("server_error", "خطأ في الخادم. حاول مرة أخرى لاحقاً.");
      } else {
        setErrorState("server_error", err?.response?.data?.detail || "تعذر تسجيل الدخول");
      }
      return false;
    } finally {
      isLoading.value = false;
    }
  }

  async function loginIdentifier(identifier: string): Promise<boolean> {
    isLoading.value = true;
    error.value = null;
    try {
      const response = await adminApi.loginIdentifier(identifier);
      applySession(response);
      return true;
    } catch (err: any) {
      const status = err?.response?.status;
      if (status === 403) {
        setErrorState("unauthorized", err?.response?.data?.detail || "ليس لديك صلاحية للوصول إلى لوحة التحكم");
      } else if (status === 401) {
        setErrorState("unauthorized", err?.response?.data?.detail || "فشل التحقق من الهوية أو المستخدم غير موجود");
      } else {
        setErrorState("server_error", err?.response?.data?.detail || "تعذر تسجيل الدخول");
      }
      return false;
    } finally {
      isLoading.value = false;
    }
  }

  async function loginDev(role: string = "OWNER"): Promise<boolean> {
    isLoading.value = true;
    error.value = null;
    try {
      const response = await adminApi.loginDev(role);
      applySession(response);
      return true;
    } catch (err: any) {
      setErrorState("server_error", err?.response?.data?.detail || "تعذر الدخول التجريبي");
      return false;
    } finally {
      isLoading.value = false;
    }
  }

  async function fetchCurrentUser(): Promise<boolean> {
    if (!getAccessToken()) {
      authState.value = "unauthenticated";
      return false;
    }
    isLoading.value = true;
    try {
      const me = await adminApi.getMe();
      user.value = me;
      localStorage.setItem("wiyana_admin_user", JSON.stringify(me));
      authState.value = "authenticated";
      error.value = null;
      return true;
    } catch (err: any) {
      const status = err?.response?.status;
      if (status === 401) {
        authState.value = "session_expired";
      } else if (!err?.response) {
        authState.value = "network_error";
      } else if (status >= 500) {
        authState.value = "server_error";
      } else {
        authState.value = "unauthenticated";
      }
      return false;
    } finally {
      isLoading.value = false;
    }
  }

  async function logout(): Promise<void> {
    try {
      await adminApi.logout();
    } catch {
      // ignore network errors during logout — always clear local state
    } finally {
      user.value = null;
      clearTokens();
      authState.value = "unauthenticated";
      error.value = null;
    }
  }

  /**
   * Restore a persisted session (tokens in localStorage) and validate it
   * against the backend via GET /auth/me.
   */
  async function initializeAuth(): Promise<void> {
    if (authState.value !== "loading") return;
    const token = getAccessToken();
    if (!token) {
      authState.value = "unauthenticated";
      return;
    }
    // Optimistically restore cached user for instant layout render.
    try {
      const cached = localStorage.getItem("wiyana_admin_user");
      if (cached) user.value = JSON.parse(cached);
    } catch {
      /* ignore corrupt cache */
    }
    await fetchCurrentUser();
  }

  function clearError() {
    error.value = null;
  }

  function updateUser(updates: Partial<AdminUser>) {
    if (user.value) {
      user.value = { ...user.value, ...updates };
      localStorage.setItem("wiyana_admin_user", JSON.stringify(user.value));
    }
  }

  return {
    // State
    user,
    authState,
    isLoading,
    error,

    // Computed
    isAuthenticated,
    userRole,
    userName,
    userAvatar,
    permissions,
    isOwner,
    isTeacher,

    // Actions
    login,
    loginIdentifier,
    loginDev,
    fetchCurrentUser,
    logout,
    initializeAuth,
    clearError,
    updateUser,
    hasPermission,
    hasAnyPermission,
    hasAllPermissions,
    hasRole,
  };
});
