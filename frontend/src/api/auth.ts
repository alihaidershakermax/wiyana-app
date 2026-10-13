import { apiClient, setAuthToken } from "./client";

export interface TelegramLoginPayload {
  init_data: string;
}

export interface AuthResponse {
  access_token: string;
  token_type: string;
  user: any;
}

export async function loginWithTelegram(initData: string): Promise<AuthResponse> {
  const response = await apiClient.post<AuthResponse>("/auth/telegram-login", {
    init_data: initData,
  });

  if (response.data?.access_token) {
    setAuthToken(response.data.access_token);
  }

  return response.data;
}

export async function loginWithToken(token: string): Promise<AuthResponse> {
  const response = await apiClient.post<AuthResponse>("/auth/token-login", {
    token,
  });

  if (response.data?.access_token) {
    setAuthToken(response.data.access_token);
  }

  return response.data;
}
