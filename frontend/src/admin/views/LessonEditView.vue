<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowRight, BookOpen, Paperclip, Video } from '@lucide/vue'
import { useLesson, useUpdateLesson } from '@/admin/queries/useLessons'
import { useSections } from '@/admin/queries/useSections'
import { useLessonAttachments, useAddAttachment, useDeleteAttachment } from '@/admin/queries/useAttachments'
import LessonForm from '@/admin/components/LessonForm.vue'
import VideoUploader from '@/admin/components/VideoUploader.vue'
import AttachmentList from '@/admin/components/AttachmentList.vue'
import StatePanel from '@/admin/components/StatePanel.vue'
import { useToast } from '@/admin/composables/useToast'
import type { LessonInput } from '@/admin/api/panelApi'

const route = useRoute()
const router = useRouter()
const toast = useToast()

const lessonId = computed(() => String(route.params.id ?? ''))
const { data: lesson, isLoading: lessonLoading, isError, error, refetch } = useLesson(lessonId)

const courseId = computed(() => lesson.value?.course_id ?? '')
const { data: sections, isLoading: sectionsLoading } = useSections(courseId)
const { data: attachments, isLoading: attachmentsLoading, refetch: refetchAttachments } = useLessonAttachments(lessonId)

const updateLessonMutation = useUpdateLesson()
const addAttachmentMutation = useAddAttachment()
const deleteAttachmentMutation = useDeleteAttachment()

function goBack() {
  if (courseId.value) {
    router.push({ name: 'admin.course-detail', params: { id: courseId.value } })
  } else {
    router.back()
  }
}

function handleFormSubmit(payload: Record<string, unknown>) {
  const input: Partial<LessonInput> = {
    title: payload.title as string,
    description: payload.description as string | undefined,
    section_id: payload.section_id as string | undefined,
    thumbnail: payload.thumbnail as string | undefined,
    status: payload.status as string,
    is_free: Boolean(payload.is_free),
  }

  updateLessonMutation.mutate(
    { lessonId: lessonId.value, input },
    {
      onSuccess: () => {
        toast.success('تم تحديث المحاضرة بنجاح')
        refetch()
      },
      onError: (err: unknown) => {
        toast.error('فشل تحديث المحاضرة', (err as Error)?.message)
      },
    }
  )
}

function onUploadFile(file: File) {
  addAttachmentMutation.mutate(
    { lessonId: lessonId.value, file },
    {
      onSuccess: () => {
        toast.success('تم إضافة المرفق بنجاح')
        refetchAttachments()
      },
      onError: (err: unknown) => toast.error('فشلت إضافة المرفق', (err as Error)?.message),
    }
  )
}

function onAddLink(link: { url: string; name?: string }) {
  addAttachmentMutation.mutate(
    { lessonId: lessonId.value, url: link.url, name: link.name },
    {
      onSuccess: () => {
        toast.success('تم إضافة الرابط بنجاح')
        refetchAttachments()
      },
      onError: (err: unknown) => toast.error('فشلت إضافة الرابط', (err as Error)?.message),
    }
  )
}

function onDeleteAttachment(attachmentId: string) {
  deleteAttachmentMutation.mutate(attachmentId, {
    onSuccess: () => {
      toast.success('تم حذف المرفق')
      refetchAttachments()
    },
    onError: (err: unknown) => toast.error('فشل حذف المرفق', (err as Error)?.message),
  })
}

function onVideoUploaded() {
  toast.success('تم تحديث الفيديو وتثبيته على المحاضرة')
  refetch()
}
</script>

<template>
  <div class="max-w-4xl mx-auto space-y-6">
    <button class="btn-ghost text-sm -mb-2" @click="goBack">
      <ArrowRight class="h-4 w-4" />
      العودة
    </button>

    <StatePanel
      :loading="lessonLoading || sectionsLoading"
      :error="isError ? (error as Error)?.message : null"
      :empty="!lesson"
      empty-title="المحاضرة غير موجودة"
      :skeleton-count="5"
      @retry="() => refetch()"
    >
      <template v-if="lesson">
        <div class="flex items-center justify-between gap-3">
          <div>
            <h1 class="text-xl font-extrabold flex items-center gap-2">
              <BookOpen class="h-5 w-5 text-orange-500" />
              تعديل المحاضرة: {{ lesson.title }}
            </h1>
            <p class="text-xs text-slate-500 dark:text-slate-400 mt-1">
              تعديل البيانات الأساسية، الفيديو المرفق، والمرفقات
            </p>
          </div>
        </div>

        <!-- Section 1: Basic Info -->
        <div class="card p-6 space-y-4">
          <h3 class="text-sm font-extrabold flex items-center gap-2 text-slate-800 dark:text-slate-100">
            <BookOpen class="h-4 w-4 text-orange-500" />
            بيانات المحاضرة
          </h3>
          <LessonForm
            :sections="sections ?? []"
            :initial="lesson"
            :loading="updateLessonMutation.isPending.value"
            submit-label="حفظ التعديلات"
            @submit="handleFormSubmit"
          />
        </div>

        <!-- Section 2: Video Upload -->
        <div class="card p-6 space-y-4">
          <h3 class="text-sm font-extrabold flex items-center gap-2 text-slate-800 dark:text-slate-100">
            <Video class="h-4 w-4 text-red-500" />
            فيديو المحاضرة (YouTube)
          </h3>
          <VideoUploader
            :lesson-id="lesson.id"
            :existing-video-id="lesson.youtube_video_id"
            @uploaded="onVideoUploaded"
          />
        </div>

        <!-- Section 3: Attachments -->
        <div class="card p-6 space-y-4">
          <h3 class="text-sm font-extrabold flex items-center gap-2 text-slate-800 dark:text-slate-100">
            <Paperclip class="h-4 w-4 text-orange-500" />
            مرفقات المحاضرة
          </h3>
          <StatePanel
            :loading="attachmentsLoading"
            :error="null"
            :empty="false"
            :skeleton-count="2"
            @retry="() => refetchAttachments()"
          >
            <AttachmentList
              :attachments="attachments ?? []"
              :lesson-id="lessonId"
              :loading="addAttachmentMutation.isPending.value"
              :can-manage="true"
              @upload-file="onUploadFile"
              @add-link="onAddLink"
              @delete="onDeleteAttachment"
            />
          </StatePanel>
        </div>
      </template>
    </StatePanel>
  </div>
</template>
