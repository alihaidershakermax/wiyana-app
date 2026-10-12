import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query'
import { computed, type Ref } from 'vue'
import { attachmentsApi, lessonKeys } from '@/admin/api/panelApi'

export function useAttachments(lessonId: Ref<string> | string) {
  const idVal = computed(() => (typeof lessonId === 'string' ? lessonId : lessonId.value))
  return useQuery({
    queryKey: computed(() => lessonKeys.attachments(idVal.value)),
    queryFn: () => attachmentsApi.getByLesson(idVal.value),
    enabled: computed(() => !!idVal.value),
    staleTime: 30_000,
    retry: 1,
  })
}

export const useLessonAttachments = useAttachments

export function useUploadAttachment() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ lessonId, file, name }: { lessonId: string; file: File; name?: string }) =>
      attachmentsApi.upload(lessonId, file, name),
    onSuccess: (_data, variables) => {
      queryClient.invalidateQueries({ queryKey: lessonKeys.attachments(variables.lessonId) })
    },
  })
}

export function useAddAttachmentLink() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ lessonId, url, name }: { lessonId: string; url: string; name?: string }) =>
      attachmentsApi.addLink(lessonId, url, name),
    onSuccess: (_data, variables) => {
      queryClient.invalidateQueries({ queryKey: lessonKeys.attachments(variables.lessonId) })
    },
  })
}

export function useAddAttachment() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ lessonId, file, url, name }: { lessonId: string; file?: File; url?: string; name?: string }) => {
      if (file) return attachmentsApi.upload(lessonId, file, name)
      if (url) return attachmentsApi.addLink(lessonId, url, name)
      throw new Error('يرجى تحديد ملف أو رابط للمرفق')
    },
    onSuccess: (_data, variables) => {
      queryClient.invalidateQueries({ queryKey: lessonKeys.attachments(variables.lessonId) })
    },
  })
}

export function useDeleteAttachment() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (arg: string | { lessonId: string; attachmentId: string }) => {
      const attachmentId = typeof arg === 'string' ? arg : arg.attachmentId
      return attachmentsApi.remove(attachmentId)
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['attachments'] })
    },
  })
}
