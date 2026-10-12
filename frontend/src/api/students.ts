import { apiClient } from "./client";

export async function fetchStudentProfile(): Promise<any> {
  const response = await apiClient.get("/me");
  return response.data;
}

export async function fetchStudentEnrolledCourses(studentId: string): Promise<any[]> {
  const response = await apiClient.get(`/students/${studentId}/courses`);
  return response.data?.data ?? response.data ?? [];
}

export async function fetchStudentProgress(studentId: string): Promise<any> {
  const response = await apiClient.get(`/students/${studentId}/progress`);
  return response.data?.data ?? response.data;
}
