import { adminApiClient } from './adminApi'
import type { Section, SectionCreate, SectionUpdate, SectionsReorderPayload } from '@/admin/types/section'

export const sectionApi = {
  async getSections(courseId: string): Promise<Section[]> {
    const { data } = await adminApiClient.get<Section[]>(`/courses/${courseId}/sections`)
    return data
  },

  async createSection(courseId: string, payload: SectionCreate): Promise<Section> {
    const { data } = await adminApiClient.post<Section>(`/courses/${courseId}/sections`, payload)
    return data
  },

  async updateSection(id: string, payload: SectionUpdate): Promise<Section> {
    const { data } = await adminApiClient.patch<Section>(`/sections/${id}`, payload)
    return data
  },

  async deleteSection(id: string): Promise<{ success: boolean }> {
    const { data } = await adminApiClient.delete<{ success: boolean }>(`/sections/${id}`)
    return data
  },

  async reorderSections(courseId: string, payload: SectionsReorderPayload): Promise<{ success: boolean; sections: Section[] }> {
    const { data } = await adminApiClient.post<{ success: boolean; sections: Section[] }>(
      `/courses/${courseId}/sections/reorder`,
      payload
    )
    return data
  },
}
