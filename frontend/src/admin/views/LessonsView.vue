<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowRight, BookOpen, Plus, Send, Archive, Trash2, RefreshCw } from '@lucide/vue'
import { useCourseLessons, usePublishLesson, useArchiveLesson, useDeleteLesson } from '@/admin/queries/useLessons'
import { useSections } from '@/admin/queries/useSections'
import LessonList from '@/admin/components/LessonList.vue'
import StatePanel from '@/admin/components/StatePanel.vue'
import ConfirmDialog from '@/admin/components/ui/ConfirmDialog.vue'
import { useToast } from '@/admin/composables/useToast'
import type { Lesson } from '@/admin/types/panel'

const route = useRoute()
const router = useRouter()
const toast = useToast()

const courseId = computed(() => String(route.params.courseId ?? ''))
const {
  data: lessons,
  isLoading: lessonsLoading,
  isError: lessonsError,
  error: lessonsErrorObj,
  refetch: refetchLessons,
} = useCourseLessons(courseId)

const {
  data: sections,
  isLoading: sectionsLoading,
} = useSections(courseId)

const publish = usePublishLesson()
const archive = useArchiveLesson()
const remove = useDeleteLesson()

const pendingDelete = ref<Lesson | null>(null)
const deleting = ref(false)

const loading = computed(() => lessonsLoading.value || sectionsLoading.value)

function goBack() {
  router.back()
}

function createLesson(sectionId: string) {
  router.push({ name: 'admin.lesson-create', params: { courseId: courseId.value, sectionId } })
}

function editLesson(lesson: Lesson) {
  router.push({ name: 'admin.lesson-edit', params: { id: lesson.id } })
}

function onPublish(lesson: Lesson) {
  publish.mutate(lesson.id, {
    onSuccess: () => toast.success('تم نشر المحاضرة'),
    onError: (err: unknown) => toast.error('فشل النشر', (err as Error)?.message),
  })
}

function onArchive(lesson: Lesson) {
  archive.mutate(lesson.id, {
    onSuccess: () => toast.success('تمت أرشفة المحاضرة'),
    onError: (err: unknown) => toast.error('فشلت الأرشفة', (err as Error)?.message),
  })
}

function confirmDelete() {
  if (!pendingDelete.value) return
  deleting.value = true
  remove.mutate(pendingDelete.value.id, {
    onSuccess: () => {
      toast.success('تم حذف المحاضرة')
      pendingDelete.value = null
      deleting.value = false
    },
    onError: (err: unknown) => {
      toast.error('فشل الحذف', (err as Error)?.message)
      pendingDelete.value = null
      deleting.value = false
    },
  })
}
</script>

<template>
  <div class="space-y-6">
    <button class="btn-ghost text-sm -mb-2" @click="goBack">
      <ArrowRight class="h-4 w-4" />
      العودة
    </button>

    <div class="flex flex-wrap items-center justify-between gap-3">
      <div>
        <h1 class="text-xl font-extrabold flex items-center gap-2">
          <BookOpen class="h-5 w-5 text-orange-500" />
          محاضرات المساق
        </h1>
        <p class="text-xs text-slate-500 dark:text-slate-400 mt-1">
          إدارة المحاضرات حسب الأقسام — مرتبطة بـ YouTube والمرفقات
        </p>
      </div>
      <div class="flex items-center gap-2">
        <button class="btn-secondary text-sm" @click="() => refetchLessons()">
          <RefreshCw class="h-4 w-4" />
          تحديث
        </button>
      </div>
    </div>

    <StatePanel
      :loading="loading"
      :error="lessonsError ? (lessonsErrorObj as Error)?.message : null"
      :empty="false"
      :skeleton-count="5"
      @retry="() => refetchLessons()"
    >
      <LessonList
        :sections="sections ?? []"
        :lessons="lessons ?? []"
        :loading="loading"
        :can-create="true"
        @create="createLesson"
        @edit="editLesson"
        @publish="onPublish"
        @archive="onArchive"
        @delete="pendingDelete = $event"
      />
    </StatePanel>

    <ConfirmDialog
      :open="!!pendingDelete"
      title="حذف المحاضرة"
      message="سيتم حذف المحاضرة نهائياً. لا يمكن التراجع عن هذا الإجراء. هل أنت متأكد؟"
      confirm-label="حذف"
      cancel-label="إلغاء"
      variant="danger"
      :loading="deleting"
      @confirm="confirmDelete"
      @cancel="pendingDelete = null"
    />
  </div>
</template>
