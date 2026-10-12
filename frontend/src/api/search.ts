import { apiClient } from "./client";

export async function searchPlatform(query: string): Promise<any> {
  const response = await apiClient.get("/search/", { params: { q: query } });
  return response.data;
}
