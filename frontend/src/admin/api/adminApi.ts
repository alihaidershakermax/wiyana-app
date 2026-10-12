import type {
  AdminAuthResponse,
  AdminUser,
  LogoutResponse,
  MeResponse,
} from "@/admin/types/auth";
import {
  adminApiClient,
  adminTokenStore,
} from "@/api/client";

export { adminApiClient };

export function getAccessToken(): string | null {
  return adminTokenStore.getAccessToken();
}

export function getRefreshToken(): string | null {
  return adminTokenStore.getRefreshToken();
}

export function setTokens(access: string, refresh: string): void {
  adminTokenStore.setTokens(access, refresh);
}

export function clearTokens(): void {
  adminTokenStore.clear();
}

export const adminApi = {
  async loginTelegram(initData: string): Promise<AdminAuthResponse> {
    const { data } = await adminApiClient.post<AdminAuthResponse>(
      "/auth/telegram",
      { init_data: initData },
    );
    setTokens(data.access_token, data.refresh_token);
    return data;
  },

  async loginIdentifier(identifier: string): Promise<AdminAuthResponse> {
    const { data } = await adminApiClient.post<AdminAuthResponse>(
      "/auth/admin-login",
      { identifier },
    );
    setTokens(data.access_token, data.refresh_token);
    return data;
  },

  async loginDev(role = "OWNER"): Promise<AdminAuthResponse> {
    const { data } = await adminApiClient.post<AdminAuthResponse>(
      "/auth/dev-login",
      { role },
    );
    setTokens(data.access_token, data.refresh_token);
    return data;
  },

  async getMe(): Promise<AdminUser> {
    const { data } = await adminApiClient.get<MeResponse>("/auth/me");
    return data.user;
  },

  async logout(): Promise<void> {
    try {
      await adminApiClient.post<LogoutResponse>("/auth/logout");
    } finally {
      clearTokens();
    }
  },
};
