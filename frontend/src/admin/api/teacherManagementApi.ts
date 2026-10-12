import { adminApiClient } from './adminApi'
import type { Teacher, TeacherCreate } from '@/admin/types/teacher'

export const teacherManagementApi = {
  async getTeachers(params?: { search?: string; limit?: number }): Promise<Teacher[]> {
    const { data } = await adminApiClient.get<Teacher[]>('/teachers/', {
      params: { search: params?.search, limit: params?.limit ?? 200 },
    })
    return Array.isArray(data) ? data : []
  },

  async getTeacher(id: string): Promise<Teacher> {
    const { data } = await adminApiClient.get<Teacher>(`/teachers/${id}`)
    return data
  },

  async createTeacher(payload: TeacherCreate): Promise<Teacher> {
    const { data } = await adminApiClient.post<Teacher>('/teachers/', {
      telegram_id: payload.telegram_id ?? undefined,
      username: payload.username?.trim() || undefined,
      name: payload.name?.trim() || undefined,
      bio: payload.bio?.trim() || undefined,
    })
    return data
  },

  async updateTeacherStatus(id: string, isActive: boolean): Promise<Teacher> {
    const { data } = await adminApiClient.patch<Teacher>(`/teachers/${id}/status`, {
      is_active: isActive,
    })
    return data
  },

  async assignToCourse(courseId: string, teacherId: string) {
    const { data } = await adminApiClient.post(`/courses/${courseId}/teacher`, {
      teacher_id: teacherId,
    })
    return data
  },
}
