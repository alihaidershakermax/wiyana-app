import { ref } from 'vue'
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { computed, type Ref } from 'vue'
import { lessonsApi, lessonKeys } from '@/admin/api/panelApi'

export function useUploadLessonVideo() {
  const queryClient = useQueryClient()
  const progress = ref(0)
  const uploadState = ref<'idle' | 'uploading' | 'processing' | 'completed' | 'failed'>('idle')
  const errorMessage = ref<string | null>(null)

  const mutation = useMutation({
    mutationFn: ({
      lessonId,
      file,
      title,
      description,
      privacyStatus,
    }: {
      lessonId: string
      file: File
      title?: string
      description?: string
      privacyStatus?: string
    }) => {
      uploadState.value = 'uploading'
      progress.value = 0
      errorMessage.value = null
      return lessonsApi.uploadVideo(
        lessonId,
        file,
        { title, description, privacyStatus },
        (p) => (progress.value = p)
      )
    },
    onSuccess: (data, variables) => {
      uploadState.value = 'completed'
      progress.value = 100
      queryClient.invalidateQueries({ queryKey: lessonKeys.detail(variables.lessonId) })
      queryClient.invalidateQueries({ queryKey: lessonKeys.all })
    },
    onError: (err: unknown) => {
      uploadState.value = 'failed'
      errorMessage.value =
        (err as { response?: { data?: { detail?: string } } })?.response?.data?.detail ??
        (err as Error)?.message ??
        'فشل رفع الفيديو إلى YouTube'
    },
  })

  return { ...mutation, progress, uploadState, errorMessage }
}

export function useYouTubeStatus() {
  return useQuery({
    queryKey: ['youtube-status'] as const,
    queryFn: () => lessonsApi.getYouTubeStatus(),
    staleTime: 60_000,
    retry: 1,
  })
}

// Re-export for consumers that need keys
export { lessonKeys }

export function useLessonPlayer(lessonId: Ref<string> | string) {
  const idVal = computed(() => (typeof lessonId === 'string' ? lessonId : lessonId.value))
  return useQuery({
    queryKey: computed(() => lessonKeys.player(idVal.value)),
    queryFn: () => lessonsApi.getPlayer(idVal.value),
    enabled: computed(() => !!idVal.value),
    staleTime: 30_000,
    retry: 1,
  })
}
