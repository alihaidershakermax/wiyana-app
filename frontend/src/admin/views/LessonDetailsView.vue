<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowRight, BookOpen, Edit, Video, Paperclip, Send, Archive, Trash2 } from '@lucide/vue'
import { useLesson, usePublishLesson, useArchiveLesson, useDeleteLesson } from '@/admin/queries/useLessons'
import { useLessonAttachments } from '@/admin/queries/useAttachments'
import AttachmentList from '@/admin/components/AttachmentList.vue'
import StatePanel from '@/admin/components/StatePanel.vue'
import ConfirmDialog from '@/admin/components/ui/ConfirmDialog.vue'
import { useToast } from '@/admin/composables/useToast'
import { formatDate } from '@/utils/format'

const route = useRoute()
const router = useRouter()
const toast = useToast()

const lessonId = computed(() => String(route.params.id ?? ''))
const { data: lesson, isLoading: lessonLoading, isError, error, refetch } = useLesson(lessonId)
const { data: attachments, isLoading: attachmentsLoading } = useLessonAttachments(lessonId)

const publish = usePublishLesson()
const archive = useArchiveLesson()
const remove = useDeleteLesson()

const showDeleteModal = ref(false)
const actionLoading = ref(false)

function goBack() {
  if (lesson.value?.course_id) {
    router.push({ name: 'admin.course-detail', params: { id: lesson.value.course_id } })
  } else {
    router.back()
  }
}

function goToEdit() {
  router.push({ name: 'admin.lesson-edit', params: { id: lessonId.value } })
}

function onPublish() {
  if (!lesson.value) return
  publish.mutate(lesson.value.id, {
    onSuccess: () => {
      toast.success('تم نشر المحاضرة بنجاح')
      refetch()
    },
    onError: (err: unknown) => toast.error('فشل النشر', (err as Error)?.message),
  })
}

function onArchive() {
  if (!lesson.value) return
  archive.mutate(lesson.value.id, {
    onSuccess: () => {
      toast.success('تم أرشفة المحاضرة بنجاح')
      refetch()
    },
    onError: (err: unknown) => toast.error('فشلت الأرشفة', (err as Error)?.message),
  })
}

function confirmDelete() {
  if (!lesson.value) return
  actionLoading.value = true
  remove.mutate(lesson.value.id, {
    onSuccess: () => {
      toast.success('تم حذف المحاضرة بنجاح')
      showDeleteModal.value = false
      actionLoading.value = false
      goBack()
    },
    onError: (err: unknown) => {
      toast.error('فشل الحذف', (err as Error)?.message)
      showDeleteModal.value = false
      actionLoading.value = false
    },
  })
}
</script>

<template>
  <div class="max-w-4xl mx-auto space-y-6">
    <button class="btn-ghost text-sm -mb-2" @click="goBack">
      <ArrowRight class="h-4 w-4" />
      العودة
    </button>

    <StatePanel
      :loading="lessonLoading"
      :error="isError ? (error as Error)?.message : null"
      :empty="!lesson"
      empty-title="المحاضرة غير موجودة"
      :skeleton-count="5"
      @retry="() => refetch()"
    >
      <template v-if="lesson">
        <div class="flex flex-wrap items-center justify-between gap-3">
          <div>
            <h1 class="text-xl font-extrabold flex items-center gap-2">
              <BookOpen class="h-5 w-5 text-orange-500" />
              {{ lesson.title }}
            </h1>
            <p class="text-xs text-slate-500 dark:text-slate-400 mt-1">
              تاريخ الإنشاء: {{ formatDate(lesson.created_at) }}
            </p>
          </div>

          <div class="flex items-center gap-2">
            <button class="btn-secondary text-sm" @click="goToEdit">
              <Edit class="h-4 w-4" />
              تعديل
            </button>
            <button
              v-if="lesson.status === 'draft'"
              class="btn-primary text-sm"
              @click="onPublish"
            >
              <Send class="h-4 w-4" />
              نشر
            </button>
            <button
              v-if="lesson.status === 'published'"
              class="btn-secondary text-sm text-amber-600"
              @click="onArchive"
            >
              <Archive class="h-4 w-4" />
              أرشفة
            </button>
            <button class="btn-secondary text-sm text-red-500" @click="showDeleteModal = true">
              <Trash2 class="h-4 w-4" />
              حذف
            </button>
          </div>
        </div>

        <!-- Video Player preview -->
        <div class="card overflow-hidden">
          <div v-if="lesson.youtube_video_id" class="aspect-video w-full bg-black">
            <iframe
              :src="`https://www.youtube.com/embed/${lesson.youtube_video_id}?controls=0&modestbranding=1&rel=0&showinfo=0&iv_load_policy=3&fs=0`"
              title="YouTube video player"
              class="w-full h-full border-0"
              allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
              allowfullscreen
            />
          </div>
          <div v-else class="p-12 text-center space-y-3 bg-slate-50 dark:bg-slate-900/50">
            <Video class="h-12 w-12 text-slate-300 dark:text-slate-600 mx-auto" />
            <p class="text-sm font-bold text-slate-500">لا يوجد فيديو مرفع لهذه المحاضرة</p>
            <button class="btn-primary text-xs mx-auto" @click="goToEdit">
              رفع فيديو الآن
            </button>
          </div>

          <div class="p-5 space-y-3">
            <div class="flex items-center justify-between">
              <span class="text-xs font-bold px-2.5 py-1 rounded-full" :class="lesson.status === 'published' ? 'bg-emerald-100 text-emerald-700' : 'bg-slate-100 text-slate-600'">
                {{ lesson.status === 'published' ? 'منشور' : lesson.status === 'archived' ? 'مؤرشف' : 'مسودة' }}
              </span>
              <span v-if="lesson.is_free" class="text-xs font-bold text-emerald-600">مجانية</span>
            </div>
            <p v-if="lesson.description" class="text-sm text-slate-600 dark:text-slate-300 leading-relaxed">
              {{ lesson.description }}
            </p>
          </div>
        </div>

        <!-- Attachments -->
        <div class="card p-6 space-y-4">
          <h3 class="text-sm font-extrabold flex items-center gap-2">
            <Paperclip class="h-4 w-4 text-orange-500" />
            المرفقات ({{ attachments?.length ?? 0 }})
          </h3>
          <StatePanel
            :loading="attachmentsLoading"
            :error="null"
            :empty="!attachments?.length"
            empty-title="لا توجد مرفقات"
            :skeleton-count="2"
          >
            <AttachmentList
              :attachments="attachments ?? []"
              :lesson-id="lessonId"
              :can-manage="false"
            />
          </StatePanel>
        </div>
      </template>
    </StatePanel>

    <ConfirmDialog
      :open="showDeleteModal"
      title="حذف المحاضرة"
      message="هل أنت متأكد من حذف هذه المحاضرة نهائياً؟ لا يمكن التراجع."
      confirm-label="حذف"
      cancel-label="إلغاء"
      variant="danger"
      :loading="actionLoading"
      @confirm="confirmDelete"
      @cancel="showDeleteModal = false"
    />
  </div>
</template>
