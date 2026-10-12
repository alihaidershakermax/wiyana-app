import { adminApiClient } from './adminApi'
import type { Subject, SubjectCreate, SubjectUpdate, SubjectListFilters, DeleteSubjectResult } from '@/admin/types/subject'
import type { PaginatedResponse } from '@/admin/types/dashboard'

export const subjectApi = {
  async getSubjects(filters: SubjectListFilters = {}): Promise<PaginatedResponse<Subject>> {
    const params: Record<string, any> = {}
    if (filters.search) params.search = filters.search
    if (filters.department) params.department = filters.department
    if (filters.stage) params.stage = filters.stage
    if (filters.status) params.status = filters.status
    params.page = filters.page ?? 1
    params.limit = filters.limit ?? 20

    const { data } = await adminApiClient.get<PaginatedResponse<Subject>>('/subjects/', { params })
    return data
  },

  async getSubject(id: string): Promise<Subject> {
    const { data } = await adminApiClient.get<Subject>(`/subjects/${id}`)
    return data
  },

  async createSubject(payload: SubjectCreate): Promise<Subject> {
    const { data } = await adminApiClient.post<Subject>('/subjects/', payload)
    return data
  },

  async updateSubject(id: string, payload: SubjectUpdate): Promise<Subject> {
    const { data } = await adminApiClient.patch<Subject>(`/subjects/${id}`, payload)
    return data
  },

  async deleteSubject(id: string): Promise<{ success: boolean } & DeleteSubjectResult> {
    const { data } = await adminApiClient.delete<{ success: boolean } & DeleteSubjectResult>(
      `/subjects/${id}`
    )
    return data
  },
}
