import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query'
import { computed, type Ref } from 'vue'
import { sectionApi } from '@/admin/api/sectionApi'
import type { SectionCreate, SectionUpdate, SectionsReorderPayload } from '@/admin/types/section'

export const sectionKeys = {
  all: ['sections'] as const,
  byCourse: (courseId: string) => [...sectionKeys.all, 'course', courseId] as const,
  detail: (id: string) => [...sectionKeys.all, 'detail', id] as const,
}

export function useSections(courseId: Ref<string> | string) {
  const idVal = computed(() => {
    const raw = typeof courseId === 'string' ? courseId : courseId.value
    return raw || ''
  })
  return useQuery({
    queryKey: computed(() => sectionKeys.byCourse(idVal.value)),
    queryFn: () => sectionApi.getSections(idVal.value),
    enabled: computed(() => !!idVal.value),
    staleTime: 30_000,
  })
}

export function useCreateSection() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ courseId, payload }: { courseId: string; payload: SectionCreate }) =>
      sectionApi.createSection(courseId, payload),
    onSuccess: (_data, variables) => {
      queryClient.invalidateQueries({ queryKey: sectionKeys.byCourse(variables.courseId) })
    },
  })
}

export function useUpdateSection() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ id, courseId, payload }: { id: string; courseId: string; payload: SectionUpdate }) =>
      sectionApi.updateSection(id, payload),
    onSuccess: (_data, variables) => {
      queryClient.invalidateQueries({ queryKey: sectionKeys.byCourse(variables.courseId) })
    },
  })
}

export function useDeleteSection() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ id }: { id: string; courseId: string }) => sectionApi.deleteSection(id),
    onSuccess: (_data, variables) => {
      queryClient.invalidateQueries({ queryKey: sectionKeys.byCourse(variables.courseId) })
    },
  })
}

export function useReorderSections() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ courseId, payload }: { courseId: string; payload: SectionsReorderPayload }) =>
      sectionApi.reorderSections(courseId, payload),
    onSuccess: (_data, variables) => {
      queryClient.invalidateQueries({ queryKey: sectionKeys.byCourse(variables.courseId) })
    },
  })
}
