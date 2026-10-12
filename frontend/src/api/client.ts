import axios, {
  AxiosError,
  AxiosInstance,
  InternalAxiosRequestConfig,
} from "axios";

export interface TokenStore {
  getAccessToken(): string | null;
  getRefreshToken(): string | null;
  setTokens(accessToken: string, refreshToken?: string | null): void;
  clear(): void;
}

export interface ApiErrorPayload {
  success?: boolean;
  error?: string | { code?: string; message?: string; details?: unknown };
  detail?: string;
  message?: string;
  request_id?: string;
  requestId?: string;
  [key: string]: unknown;
}

export class ApiError extends Error {
  readonly status?: number;
  readonly code?: string;
  readonly details?: unknown;
  readonly requestId?: string;
  readonly cause: unknown;

  constructor(message: string, options: {
    status?: number;
    code?: string;
    details?: unknown;
    requestId?: string;
    cause?: unknown;
  } = {}) {
    super(message);
    this.name = "ApiError";
    this.status = options.status;
    this.code = options.code;
    this.details = options.details;
    this.requestId = options.requestId;
    this.cause = options.cause;
  }
}

const rawBaseUrl = import.meta.env.VITE_API_URL || import.meta.env.VITE_API_BASE_URL || "/api/v1";
export const API_BASE_URL = normalizeBaseUrl(rawBaseUrl);

export const TOKEN_KEY = "wiyana_access_token";
export const ADMIN_ACCESS_TOKEN_KEY = "wiyana_admin_access_token";
export const ADMIN_REFRESH_TOKEN_KEY = "wiyana_admin_refresh_token";

function normalizeBaseUrl(value: string): string {
  return value.endsWith("/") ? value.slice(0, -1) : value;
}

function getStorage(): Storage | null {
  if (typeof window === "undefined") return null;
  try {
    return window.localStorage;
  } catch {
    return null;
  }
}

function localTokenStore(accessKey: string, refreshKey?: string, cleanupKeys: string[] = []): TokenStore {
  return {
    getAccessToken() {
      return getStorage()?.getItem(accessKey) ?? null;
    },
    getRefreshToken() {
      return refreshKey ? getStorage()?.getItem(refreshKey) ?? null : null;
    },
    setTokens(accessToken, refreshToken) {
      const storage = getStorage();
      if (!storage) return;
      storage.setItem(accessKey, accessToken);
      if (refreshKey && refreshToken) storage.setItem(refreshKey, refreshToken);
    },
    clear() {
      const storage = getStorage();
      if (!storage) return;
      storage.removeItem(accessKey);
      if (refreshKey) storage.removeItem(refreshKey);
      cleanupKeys.forEach((key) => storage.removeItem(key));
    },
  };
}

export const userTokenStore = localTokenStore(TOKEN_KEY);
export const adminTokenStore = localTokenStore(
  ADMIN_ACCESS_TOKEN_KEY,
  ADMIN_REFRESH_TOKEN_KEY,
  ["wiyana_admin_user"],
);

export function getAuthToken(): string | null {
  return userTokenStore.getAccessToken();
}

export function setAuthToken(token: string): void {
  userTokenStore.setTokens(token);
}

export function removeAuthToken(): void {
  userTokenStore.clear();
}

export function getApiError(error: unknown): ApiError {
  if (error instanceof ApiError) return error;
  if (!axios.isAxiosError(error)) {
    return new ApiError(error instanceof Error ? error.message : "تعذر الاتصال بالخادم", { cause: error });
  }

  const payload = error.response?.data as ApiErrorPayload | undefined;
  const nested = payload?.error && typeof payload.error === "object" ? payload.error : undefined;
  const message = nested?.message || payload?.detail || payload?.message || (typeof payload?.error === "string" ? payload.error : undefined) || error.message || "تعذر الاتصال بالخادم";
  return new ApiError(message, {
    status: error.response?.status,
    code: nested?.code,
    details: nested?.details ?? payload,
    requestId: payload?.request_id || payload?.requestId || error.response?.headers?.["x-request-id"],
    cause: error,
  });
}

type RetryableRequestConfig = InternalAxiosRequestConfig & {
  _wiyanaRetry?: boolean;
};

export interface ApiClientOptions {
  baseURL: string;
  tokenStore?: TokenStore;
  refreshEnabled?: boolean;
  authEventPrefix?: "wiyana" | "admin";
  attachTelegramInitData?: boolean;
  timeout?: number;
}

function createRequestId(): string {
  if (typeof crypto !== "undefined" && typeof crypto.randomUUID === "function") {
    return crypto.randomUUID();
  }
  return `${Date.now()}-${Math.random().toString(36).slice(2, 10)}`;
}

function isAuthPath(url = ""): boolean {
  return /\/auth\/(telegram|telegram-login|admin-login|dev-login|refresh|logout)/.test(url);
}

export function createApiClient(options: ApiClientOptions): AxiosInstance {
  const tokenStore = options.tokenStore;
  const eventPrefix = options.authEventPrefix || "wiyana";
  let refreshPromise: Promise<string | null> | null = null;

  const instance = axios.create({
    baseURL: options.baseURL,
    timeout: options.timeout ?? 120000,
    headers: {
      Accept: "application/json",
    },
  });

  instance.interceptors.request.use(
    (config: InternalAxiosRequestConfig) => {
      const token = tokenStore?.getAccessToken();
      if (token) config.headers.Authorization = `Bearer ${token}`;

      if (options.attachTelegramInitData !== false && typeof window !== "undefined") {
        const initData = (window as Window & { Telegram?: { WebApp?: { initData?: string } } }).Telegram?.WebApp?.initData;
        if (initData) config.headers["X-Telegram-Init-Data"] = initData;
      }

      config.headers["X-Request-ID"] = config.headers["X-Request-ID"] || createRequestId();
      if (typeof FormData !== "undefined" && config.data instanceof FormData) {
        delete config.headers["Content-Type"];
      } else if (!config.headers["Content-Type"]) {
        config.headers["Content-Type"] = "application/json";
      }
      return config;
    },
    (error: AxiosError) => Promise.reject(error),
  );

  async function refreshAccessToken(): Promise<string | null> {
    if (!tokenStore) return null;
    const refreshToken = tokenStore.getRefreshToken();
    if (!refreshToken) return null;
    if (refreshPromise) return refreshPromise;

    refreshPromise = axios
      .post<{ access_token: string; refresh_token?: string }>(
        `${options.baseURL}/auth/refresh`,
        { refresh_token: refreshToken },
        { headers: { Accept: "application/json", "Content-Type": "application/json" }, timeout: options.timeout ?? 120000 },
      )
      .then(({ data }) => {
        if (!data.access_token) return null;
        tokenStore.setTokens(data.access_token, data.refresh_token ?? refreshToken);
        return data.access_token;
      })
      .catch(() => null)
      .finally(() => {
        refreshPromise = null;
      });

    return refreshPromise;
  }

  instance.interceptors.response.use(
    (response) => response,
    async (error: AxiosError) => {
      const original = error.config as RetryableRequestConfig | undefined;
      const status = error.response?.status;

      if (status === 401 && options.refreshEnabled !== false && original && !original._wiyanaRetry && !isAuthPath(original.url)) {
        original._wiyanaRetry = true;
        const nextToken = await refreshAccessToken();
        if (nextToken) {
          original.headers.Authorization = `Bearer ${nextToken}`;
          return instance(original);
        }
        tokenStore?.clear();
        if (typeof window !== "undefined") {
          window.dispatchEvent(new CustomEvent(`${eventPrefix}:session-expired`));
          if (eventPrefix === "wiyana") window.dispatchEvent(new CustomEvent("wiyana:auth-required"));
        }
      }

      if (status === 403 && typeof window !== "undefined") {
        window.dispatchEvent(new CustomEvent(`${eventPrefix}:access-denied`, { detail: error.response?.data }));
      }

      return Promise.reject(error);
    },
  );

  return instance;
}

export const apiClient = createApiClient({
  baseURL: API_BASE_URL,
  tokenStore: userTokenStore,
  refreshEnabled: true,
  authEventPrefix: "wiyana",
  attachTelegramInitData: true,
});

const API_ORIGIN = /^https?:\/\//.test(API_BASE_URL) ? new URL(API_BASE_URL).origin + "/" : "/";

export const appApiClient = createApiClient({
  baseURL: API_ORIGIN,
  tokenStore: userTokenStore,
  refreshEnabled: false,
  authEventPrefix: "wiyana",
  attachTelegramInitData: true,
  timeout: 20000,
});

export const adminApiClient = createApiClient({
  baseURL: normalizeBaseUrl(import.meta.env.VITE_ADMIN_API_BASE_URL || API_BASE_URL),
  tokenStore: adminTokenStore,
  refreshEnabled: true,
  authEventPrefix: "admin",
  attachTelegramInitData: false,
});

export function resolveApiUrl(path: string): string {
  if (/^https?:\/\//.test(path)) return path;
  const origin = /^https?:\/\//.test(API_BASE_URL)
    ? new URL(API_BASE_URL).origin
    : typeof window !== "undefined"
      ? window.location.origin
      : "";
  return origin ? new URL(path, origin).href : path;
}
