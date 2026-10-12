import { apiClient } from "./client";

export async function fetchMySubscriptions(): Promise<any[]> {
  try {
    const response = await apiClient.get("/subscriptions/my-subscriptions");
    const data = response.data;
    if (Array.isArray(data)) {
      return data;
    }
    if (data && Array.isArray(data.data)) {
      return data.data;
    }
    if (data && Array.isArray(data.subscriptions)) {
      return data.subscriptions;
    }
  } catch (err) {
    console.warn("fetch /subscriptions/my-subscriptions error, checking /users/me fallback:", err);
  }

  // Fallback to /users/me if endpoint differed
  try {
    const meRes = await apiClient.get("/users/me");
    const u = meRes.data;
    if (Array.isArray(u?.subscriptions) && u.subscriptions.length > 0) {
      return u.subscriptions;
    }
  } catch (err) {
    console.warn("fetch /users/me fallback error:", err);
  }

  return [];
}

export async function redeemSubscriptionCode(code: string): Promise<any> {
  const response = await apiClient.post("/subscriptions/redeem", { code });
  return response.data;
}
