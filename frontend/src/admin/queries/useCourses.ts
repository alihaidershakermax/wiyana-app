import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query'
import { computed, type Ref } from 'vue'
import { courseApi } from '@/admin/api/courseApi'
import type { CourseListFilters, CourseCreate, CourseUpdate, AssignTeacherPayload } from '@/admin/types/course'

export const courseKeys = {
  all: ['courses'] as const,
  lists: () => [...courseKeys.all, 'list'] as const,
  list: (filters: CourseListFilters) => [...courseKeys.lists(), filters] as const,
  details: () => [...courseKeys.all, 'detail'] as const,
  detail: (id: string) => [...courseKeys.details(), id] as const,
}

export function useCourses(filters: Ref<CourseListFilters> | CourseListFilters = {}) {
  const filtersVal = computed(() => (typeof filters === 'object' && 'value' in filters ? (filters as Ref<CourseListFilters>).value : filters as CourseListFilters))
  return useQuery({
    queryKey: computed(() => courseKeys.list(filtersVal.value)),
    queryFn: () => courseApi.getCourses(filtersVal.value),
    staleTime: 30_000,
    retry: 1,
  })
}

export function useCourse(id: Ref<string> | string) {
  const idVal = computed(() => (typeof id === 'string' ? id : id.value))
  return useQuery({
    queryKey: computed(() => courseKeys.detail(idVal.value)),
    queryFn: () => courseApi.getCourse(idVal.value),
    enabled: computed(() => !!idVal.value),
    staleTime: 30_000,
  })
}

export function useCreateCourse() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (payload: CourseCreate) => courseApi.createCourse(payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: courseKeys.lists() })
    },
  })
}

export function useUpdateCourse() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ id, payload }: { id: string; payload: CourseUpdate }) =>
      courseApi.updateCourse(id, payload),
    onSuccess: (_data, variables) => {
      queryClient.invalidateQueries({ queryKey: courseKeys.lists() })
      queryClient.invalidateQueries({ queryKey: courseKeys.detail(variables.id) })
    },
  })
}

export function useDeleteCourse() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (id: string) => courseApi.deleteCourse(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: courseKeys.lists() })
    },
  })
}

export function usePublishCourse() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (id: string) => courseApi.publishCourse(id),
    onSuccess: (_data, id) => {
      queryClient.invalidateQueries({ queryKey: courseKeys.lists() })
      queryClient.invalidateQueries({ queryKey: courseKeys.detail(id) })
    },
  })
}

export function useArchiveCourse() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (id: string) => courseApi.archiveCourse(id),
    onSuccess: (_data, id) => {
      queryClient.invalidateQueries({ queryKey: courseKeys.lists() })
      queryClient.invalidateQueries({ queryKey: courseKeys.detail(id) })
    },
  })
}

export function useAssignTeacher() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ id, payload }: { id: string; payload: AssignTeacherPayload }) =>
      courseApi.assignTeacher(id, payload),
    onSuccess: (_data, variables) => {
      queryClient.invalidateQueries({ queryKey: courseKeys.detail(variables.id) })
      queryClient.invalidateQueries({ queryKey: courseKeys.lists() })
    },
  })
}

export function useRemoveTeacher() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (id: string) => courseApi.removeTeacher(id),
    onSuccess: (_data, id) => {
      queryClient.invalidateQueries({ queryKey: courseKeys.detail(id) })
      queryClient.invalidateQueries({ queryKey: courseKeys.lists() })
    },
  })
}
