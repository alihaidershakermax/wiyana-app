<script setup lang="ts">
import { ref } from 'vue'
import { Paperclip, UploadCloud, Link as LinkIcon, Loader2 } from '@lucide/vue'
import { useAttachments, useUploadAttachment, useAddAttachmentLink, useDeleteAttachment } from '@/admin/queries/useAttachments'
import AttachmentCard from './AttachmentCard.vue'
import StatePanel from './StatePanel.vue'
import { useToast } from '@/admin/composables/useToast'

const props = defineProps<{ lessonId: string }>()

const toast = useToast()
const fileInput = ref<HTMLInputElement | null>(null)
const linkUrl = ref('')
const linkName = ref('')
const showLinkForm = ref(false)
const deletingId = ref<string | null>(null)

const { data: attachments, isLoading, isError, error, refetch } = useAttachments(props.lessonId)
const upload = useUploadAttachment()
const addLink = useAddAttachmentLink()
const remove = useDeleteAttachment()

function onFileChange(e: Event) {
  const target = e.target as HTMLInputElement
  const file = target.files?.[0]
  if (!file) return
  upload.mutate(
    { lessonId: props.lessonId, file },
    {
      onSuccess: () => toast.success('تم رفع المرفق'),
      onError: (err: unknown) =>
        toast.error('فشل رفع المرفق', (err as Error)?.message ?? 'حاول مرة أخرى'),
    }
  )
  target.value = ''
}

function addLinkSubmit() {
  if (!linkUrl.value.trim()) return
  addLink.mutate(
    { lessonId: props.lessonId, url: linkUrl.value.trim(), name: linkName.value.trim() || undefined },
    {
      onSuccess: () => {
        toast.success('تمت إضافة الرابط')
        linkUrl.value = ''
        linkName.value = ''
        showLinkForm.value = false
      },
      onError: (err: unknown) => toast.error('فشل إضافة الرابط', (err as Error)?.message),
    }
  )
}

function onDelete(att: { id: string }) {
  deletingId.value = att.id
  remove.mutate(
    { lessonId: props.lessonId, attachmentId: att.id },
    {
      onSuccess: () => {
        toast.success('تم حذف المرفق')
        deletingId.value = null
      },
      onError: (err: unknown) => {
        toast.error('فشل حذف المرفق', (err as Error)?.message)
        deletingId.value = null
      },
    }
  )
}
</script>

<template>
  <div class="rounded-2xl border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 p-4 space-y-4">
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-2">
        <Paperclip class="h-4 w-4 text-orange-500" />
        <h4 class="text-sm font-bold">المرفقات</h4>
        <span v-if="attachments?.length" class="text-[10px] font-bold bg-orange-50 dark:bg-orange-950/40 text-orange-600 dark:text-orange-400 px-2 py-0.5 rounded-full">
          {{ attachments.length }}
        </span>
      </div>
      <div class="flex items-center gap-2">
        <button class="btn-secondary text-xs py-1.5" @click="showLinkForm = !showLinkForm">
          <LinkIcon class="h-3.5 w-3.5" />
          إضافة رابط
        </button>
        <button class="btn-primary text-xs py-1.5" @click="fileInput?.click()">
          <UploadCloud v-if="!upload.isPending.value" class="h-3.5 w-3.5" />
          <Loader2 v-else class="h-3.5 w-3.5 animate-spin" />
          رفع ملف
        </button>
        <input ref="fileInput" type="file" class="hidden" @change="onFileChange" />
      </div>
    </div>

    <div v-if="showLinkForm" class="grid sm:grid-cols-2 gap-3">
      <input
        v-model="linkUrl"
        type="url"
        placeholder="https://..."
        dir="ltr"
        class="input w-full text-sm"
        @keyup.enter="addLinkSubmit"
      />
      <input
        v-model="linkName"
        type="text"
        placeholder="اسم الرابط (اختياري)"
        class="input w-full text-sm"
        @keyup.enter="addLinkSubmit"
      />
    </div>

    <StatePanel
      :loading="isLoading"
      :error="isError ? (error as Error)?.message : null"
      :empty="!!attachments && attachments.length === 0"
      empty-title="لا توجد مرفقات"
      empty-message="أضف ملفات PDF أو صور أو روابط للمحاضرة"
      @retry="() => refetch()"
    >
      <AttachmentCard
        :attachments="attachments ?? []"
        :deleting-id="deletingId"
        @delete="onDelete"
      />
    </StatePanel>
  </div>
</template>
