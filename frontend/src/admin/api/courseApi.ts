import { adminApiClient } from './adminApi'
import type { Course, CourseCreate, CourseUpdate, CourseListFilters, AssignTeacherPayload } from '@/admin/types/course'
import type { PaginatedResponse } from '@/admin/types/dashboard'

export const courseApi = {
  async getCourses(filters: CourseListFilters = {}): Promise<PaginatedResponse<Course>> {
    const params: Record<string, any> = {}
    if (filters.search) params.search = filters.search
    if (filters.subject_id) params.subject_id = filters.subject_id
    if (filters.teacher_id) params.teacher_id = filters.teacher_id
    if (filters.department) params.department = filters.department
    if (filters.stage) params.stage = filters.stage
    if (filters.status) params.status = filters.status
    params.page = filters.page ?? 1
    params.limit = filters.limit ?? 20

    const { data } = await adminApiClient.get<PaginatedResponse<Course>>('/courses/', { params })
    return data
  },

  async getCourse(id: string): Promise<Course> {
    const { data } = await adminApiClient.get<Course>(`/courses/${id}`)
    return data
  },

  async createCourse(payload: CourseCreate): Promise<Course> {
    const { data } = await adminApiClient.post<Course>('/courses/', payload)
    return data
  },

  async updateCourse(id: string, payload: CourseUpdate): Promise<Course> {
    const { data } = await adminApiClient.patch<Course>(`/courses/${id}`, payload)
    return data
  },

  async deleteCourse(id: string): Promise<{ success: boolean; message: string }> {
    const { data } = await adminApiClient.delete<{ success: boolean; message: string }>(
      `/courses/${id}`
    )
    return data
  },

  async publishCourse(id: string): Promise<Course> {
    const { data } = await adminApiClient.post<Course>(`/courses/${id}/publish`)
    return data
  },

  async archiveCourse(id: string): Promise<Course> {
    const { data } = await adminApiClient.post<Course>(`/courses/${id}/archive`)
    return data
  },

  async assignTeacher(id: string, payload: AssignTeacherPayload): Promise<Course> {
    const { data } = await adminApiClient.post<Course>(`/courses/${id}/teacher`, payload)
    return data
  },

  async removeTeacher(id: string): Promise<Course> {
    const { data } = await adminApiClient.delete<Course>(`/courses/${id}/teacher`)
    return data
  },
}
