import { apiClient } from "./client";

export async function fetchMyProfile(): Promise<any> {
  const response = await apiClient.get("/users/me");
  return response.data;
}

export async function updateMyProfile(updateData: Record<string, any>): Promise<any> {
  const response = await apiClient.patch("/users/me", updateData);
  return response.data;
}
