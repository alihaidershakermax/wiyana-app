import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query'
import { computed, type Ref } from 'vue'
import { subjectApi } from '@/admin/api/subjectApi'
import type { SubjectListFilters, SubjectCreate, SubjectUpdate } from '@/admin/types/subject'

export const subjectKeys = {
  all: ['subjects'] as const,
  lists: () => [...subjectKeys.all, 'list'] as const,
  list: (filters: SubjectListFilters) => [...subjectKeys.lists(), filters] as const,
  details: () => [...subjectKeys.all, 'detail'] as const,
  detail: (id: string) => [...subjectKeys.details(), id] as const,
}

export function useSubjects(filters: Ref<SubjectListFilters> | SubjectListFilters = {}) {
  const filtersVal = computed(() => (typeof filters === 'object' && 'value' in filters ? (filters as Ref<SubjectListFilters>).value : filters as SubjectListFilters))
  return useQuery({
    queryKey: computed(() => subjectKeys.list(filtersVal.value)),
    queryFn: () => subjectApi.getSubjects(filtersVal.value),
    staleTime: 30_000,
    retry: 1,
  })
}

export function useSubject(id: Ref<string> | string) {
  const idVal = computed(() => (typeof id === 'string' ? id : id.value))
  return useQuery({
    queryKey: computed(() => subjectKeys.detail(idVal.value)),
    queryFn: () => subjectApi.getSubject(idVal.value),
    enabled: computed(() => !!idVal.value),
    staleTime: 30_000,
  })
}

export function useCreateSubject() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (payload: SubjectCreate) => subjectApi.createSubject(payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: subjectKeys.lists() })
    },
  })
}

export function useUpdateSubject() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ id, payload }: { id: string; payload: SubjectUpdate }) =>
      subjectApi.updateSubject(id, payload),
    onSuccess: (_data, variables) => {
      queryClient.invalidateQueries({ queryKey: subjectKeys.lists() })
      queryClient.invalidateQueries({ queryKey: subjectKeys.detail(variables.id) })
    },
  })
}

export function useDeleteSubject() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (id: string) => subjectApi.deleteSubject(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: subjectKeys.lists() })
    },
  })
}
