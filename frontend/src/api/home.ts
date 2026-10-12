import { apiClient } from "./client";

export async function fetchHomeDashboard(): Promise<any> {
  const response = await apiClient.get("/home/dashboard");
  return response.data;
}
