import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query'
import { computed, type Ref } from 'vue'
import {
  studentsApi,
  studentKeys,
  lessonKeys,
  enrollmentsApi,
  progressApi,
  type StudentFilters,
} from '@/admin/api/panelApi'

export function useStudents(filters: Ref<StudentFilters> | StudentFilters = {}) {
  const filtersVal = computed(() =>
    typeof filters === 'object' && 'value' in filters
      ? (filters as Ref<StudentFilters>).value
      : (filters as StudentFilters)
  )
  return useQuery({
    queryKey: computed(() => studentKeys.list(filtersVal.value)),
    queryFn: () => studentsApi.list(filtersVal.value),
    staleTime: 15_000,
    retry: 1,
  })
}

export function useStudent(id: Ref<string> | string) {
  const idVal = computed(() => (typeof id === 'string' ? id : id.value))
  return useQuery({
    queryKey: computed(() => studentKeys.detail(idVal.value)),
    queryFn: () => studentsApi.getById(idVal.value),
    enabled: computed(() => !!idVal.value),
    staleTime: 30_000,
    retry: 1,
  })
}

export function useStudentCourses(id: Ref<string> | string) {
  const idVal = computed(() => (typeof id === 'string' ? id : id.value))
  return useQuery({
    queryKey: computed(() => studentKeys.courses(idVal.value)),
    queryFn: () => studentsApi.getCourses(idVal.value),
    enabled: computed(() => !!idVal.value),
    staleTime: 30_000,
    retry: 1,
  })
}

export function useStudentProgress(id: Ref<string> | string) {
  const idVal = computed(() => (typeof id === 'string' ? id : id.value))
  return useQuery({
    queryKey: computed(() => studentKeys.progress(idVal.value)),
    queryFn: () => studentsApi.getProgress(idVal.value),
    enabled: computed(() => !!idVal.value),
    staleTime: 30_000,
    retry: 1,
  })
}

export function useStudentCourseProgress(studentId: Ref<string> | string, courseId: Ref<string> | string) {
  const sId = computed(() => (typeof studentId === 'string' ? studentId : studentId.value))
  const cId = computed(() => (typeof courseId === 'string' ? courseId : courseId.value))
  return useQuery({
    queryKey: computed(() => studentKeys.courseProgress(sId.value, cId.value)),
    queryFn: () => studentsApi.getCourseProgress(sId.value, cId.value),
    enabled: computed(() => !!sId.value && !!cId.value),
    staleTime: 30_000,
    retry: 1,
  })
}

export function useStudentActivity(id: Ref<string> | string) {
  const idVal = computed(() => (typeof id === 'string' ? id : id.value))
  return useQuery({
    queryKey: computed(() => studentKeys.activity(idVal.value)),
    queryFn: () => studentsApi.getActivity(idVal.value),
    enabled: computed(() => !!idVal.value),
    staleTime: 30_000,
    retry: 1,
  })
}

export function useCourseStudents(courseId: Ref<string> | string) {
  const idVal = computed(() => (typeof courseId === 'string' ? courseId : courseId.value))
  return useQuery({
    queryKey: computed(() => studentKeys.courseStudents(idVal.value)),
    queryFn: () => studentsApi.getCourseStudents(idVal.value),
    enabled: computed(() => !!idVal.value),
    staleTime: 30_000,
    retry: 1,
  })
}

export function useStudentEnrollments(id: Ref<string> | string) {
  const idVal = computed(() => (typeof id === 'string' ? id : id.value))
  return useQuery({
    queryKey: computed(() => [...studentKeys.detail(idVal.value), 'enrollments']),
    queryFn: () => enrollmentsApi.getByStudent(idVal.value),
    enabled: computed(() => !!idVal.value),
    staleTime: 30_000,
    retry: 1,
  })
}

export function useSuspendStudent() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (studentId: string) => studentsApi.suspend(studentId),
    onSuccess: (_data, studentId) => {
      queryClient.invalidateQueries({ queryKey: studentKeys.detail(studentId) })
      queryClient.invalidateQueries({ queryKey: studentKeys.all })
    },
  })
}

export function useActivateStudent() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (studentId: string) => studentsApi.activate(studentId),
    onSuccess: (_data, studentId) => {
      queryClient.invalidateQueries({ queryKey: studentKeys.detail(studentId) })
      queryClient.invalidateQueries({ queryKey: studentKeys.all })
    },
  })
}

export function useLessonProgress(studentId: Ref<string> | string, lessonId: Ref<string> | string) {
  const sId = computed(() => (typeof studentId === 'string' ? studentId : studentId.value))
  const lId = computed(() => (typeof lessonId === 'string' ? lessonId : lessonId.value))
  return useQuery({
    queryKey: computed(() => [...lessonKeys.detail(lId.value), 'progress', sId.value]),
    queryFn: () => progressApi.getLesson(sId.value, lId.value),
    enabled: computed(() => !!sId.value && !!lId.value),
    staleTime: 15_000,
    retry: 1,
  })
}
