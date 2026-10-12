import { apiClient } from "./client";

export interface FetchCoursesQueryParams {
  department_id?: string;
  subject_id?: string;
  teacher_id?: string;
  is_published?: boolean;
  skip?: number;
  limit?: number;
}

export async function fetchCoursesList(params?: FetchCoursesQueryParams): Promise<any[]> {
  const response = await apiClient.get("/courses/", { params });
  return response.data;
}

export async function fetchCourseById(courseId: string): Promise<any> {
  const response = await apiClient.get(`/courses/${courseId}`);
  return response.data;
}

export async function fetchDepartments(): Promise<any[]> {
  const response = await apiClient.get("/departments/");
  return response.data;
}

export async function fetchSubjects(): Promise<any[]> {
  const response = await apiClient.get("/subjects/");
  return response.data;
}

export async function fetchTeachers(): Promise<any[]> {
  const response = await apiClient.get("/teachers/");
  return response.data;
}
