import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query'
import { computed, type Ref } from 'vue'
import { lessonsApi, lessonKeys, type LessonInput } from '@/admin/api/panelApi'

export function useCourseLessons(courseId: Ref<string> | string) {
  const idVal = computed(() => (typeof courseId === 'string' ? courseId : courseId.value))
  return useQuery({
    queryKey: computed(() => lessonKeys.course(idVal.value)),
    queryFn: () => lessonsApi.getByCourse(idVal.value),
    enabled: computed(() => !!idVal.value),
    staleTime: 30_000,
    retry: 1,
  })
}

export function useSectionLessons(sectionId: Ref<string> | string) {
  const idVal = computed(() => (typeof sectionId === 'string' ? sectionId : sectionId.value))
  return useQuery({
    queryKey: computed(() => lessonKeys.section(idVal.value)),
    queryFn: () => lessonsApi.getBySection(idVal.value),
    enabled: computed(() => !!idVal.value),
    staleTime: 30_000,
    retry: 1,
  })
}

export function useLesson(id: Ref<string> | string) {
  const idVal = computed(() => (typeof id === 'string' ? id : id.value))
  return useQuery({
    queryKey: computed(() => lessonKeys.detail(idVal.value)),
    queryFn: () => lessonsApi.getById(idVal.value),
    enabled: computed(() => !!idVal.value),
    staleTime: 30_000,
    retry: 1,
  })
}

export function useCreateLesson() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ sectionId, courseId, input }: { sectionId?: string; courseId?: string; input: LessonInput }) => {
      if (sectionId) return lessonsApi.create(sectionId, input)
      if (courseId) return lessonsApi.createForCourse(courseId, input)
      throw new Error('يجب تحديد معرف المساق أو القسم')
    },
    onSuccess: (_data, variables) => {
      queryClient.invalidateQueries({ queryKey: lessonKeys.all })
    },
  })
}

export function useUpdateLesson() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ lessonId, input }: { lessonId: string; input: Partial<LessonInput> }) =>
      lessonsApi.update(lessonId, input),
    onSuccess: (_data, variables) => {
      queryClient.invalidateQueries({ queryKey: lessonKeys.detail(variables.lessonId) })
      queryClient.invalidateQueries({ queryKey: lessonKeys.all })
    },
  })
}

export function useDeleteLesson() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (lessonId: string) => lessonsApi.delete(lessonId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: lessonKeys.all })
    },
  })
}

export function usePublishLesson() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (lessonId: string) => lessonsApi.publish(lessonId),
    onSuccess: (_data, variables) => {
      queryClient.invalidateQueries({ queryKey: lessonKeys.detail(variables) })
      queryClient.invalidateQueries({ queryKey: lessonKeys.all })
    },
  })
}

export function useArchiveLesson() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (lessonId: string) => lessonsApi.archive(lessonId),
    onSuccess: (_data, variables) => {
      queryClient.invalidateQueries({ queryKey: lessonKeys.detail(variables) })
      queryClient.invalidateQueries({ queryKey: lessonKeys.all })
    },
  })
}

export function useReorderLessons() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ sectionId, order }: { sectionId: string; order: string[] }) =>
      lessonsApi.reorder(sectionId, order),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: lessonKeys.all })
    },
  })
}
