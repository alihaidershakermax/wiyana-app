import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query'
import { teacherManagementApi } from '@/admin/api/teacherManagementApi'
import type { TeacherCreate } from '@/admin/types/teacher'

export const teacherKeys = {
  all: ['teachers'] as const,
  lists: () => [...teacherKeys.all, 'list'] as const,
  details: () => [...teacherKeys.all, 'detail'] as const,
  detail: (id: string) => [...teacherKeys.details(), id] as const,
}

export function useTeachers() {
  return useQuery({
    queryKey: teacherKeys.lists(),
    queryFn: () => teacherManagementApi.getTeachers(),
    staleTime: 30_000,
    retry: 1,
  })
}

export function useTeacher(id: string) {
  return useQuery({
    queryKey: teacherKeys.detail(id),
    queryFn: () => teacherManagementApi.getTeacher(id),
    enabled: !!id,
    staleTime: 30_000,
    retry: 1,
  })
}

export function useCreateTeacher() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (payload: TeacherCreate) => teacherManagementApi.createTeacher(payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: teacherKeys.lists() })
    },
  })
}

export function useUpdateTeacherStatus() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ id, isActive }: { id: string; isActive: boolean }) =>
      teacherManagementApi.updateTeacherStatus(id, isActive),
    onSuccess: (_data, { id }) => {
      queryClient.invalidateQueries({ queryKey: teacherKeys.lists() })
      queryClient.invalidateQueries({ queryKey: teacherKeys.detail(id) })
    },
  })
}

export function useAssignTeacherToCourse() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ courseId, teacherId }: { courseId: string; teacherId: string }) =>
      teacherManagementApi.assignToCourse(courseId, teacherId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: teacherKeys.lists() })
    },
  })
}
