<script setup lang="ts">
import { computed, ref } from 'vue'
import { FileText, FileImage, Link, FileArchive, File, Download, Trash2, Loader2, Send } from '@lucide/vue'
import type { Attachment } from '@/admin/types/panel'
import { formatDate } from '@/utils/format'

const props = defineProps<{
  attachments: Attachment[]
  loading?: boolean
  deletingId?: string | null
}>()

const emit = defineEmits<{ delete: [attachment: Attachment] }>()

const downloadingId = ref<string | null>(null)

function iconFor(type: string) {
  const map: Record<string, unknown> = {
    pdf: FileText,
    image: FileImage,
    link: Link,
    document: File,
    file: FileArchive,
    telegram_file: Send,
  }
  return map[type] ?? File
}

const typeLabels: Record<string, string> = {
  pdf: 'PDF',
  image: 'صورة',
  link: 'رابط',
  document: 'مستند',
  file: 'ملف',
  telegram_file: 'تليجرام',
}

function formatSize(size?: number | null): string {
  if (!size || size <= 0) return ''
  if (size < 1024) return `${size} بايت`
  if (size < 1024 * 1024) return `${(size / 1024).toFixed(0)} ك.ب`
  return `${(size / 1024 / 1024).toFixed(1)} م.ب`
}

async function openAttachment(att: Attachment) {
  if (!att.url || downloadingId.value) return

  // External links are intentional content (not platform URLs), open directly.
  if (att.type === 'link') {
    window.open(att.url, '_blank', 'noopener')
    return
  }

  downloadingId.value = att.id
  try {
    // Fetch as a blob so the real storage URL is never exposed.
    const res = await fetch(att.url, { credentials: 'same-origin' })
    if (!res.ok) throw new Error('fetch failed')
    const blob = await res.blob()
    const objectUrl = URL.createObjectURL(blob)
    if (att.type === 'image' || att.type === 'pdf') {
      window.open(objectUrl, '_blank')
    } else {
      const a = document.createElement('a')
      a.href = objectUrl
      a.download = att.name || 'ملف'
      document.body.appendChild(a)
      a.click()
      a.remove()
    }
    setTimeout(() => URL.revokeObjectURL(objectUrl), 60000)
  } catch {
    // Do not fall back to revealing the raw URL.
  } finally {
    downloadingId.value = null
  }
}
</script>

<template>
  <div class="space-y-2">
    <div
      v-for="att in attachments"
      :key="att.id"
      class="flex items-center gap-3 rounded-xl border border-slate-200 dark:border-slate-800 p-3 hover:border-orange-300 dark:hover:border-orange-800 transition-colors"
    >
      <div class="h-10 w-10 rounded-lg bg-slate-100 dark:bg-slate-800 flex items-center justify-center shrink-0">
        <component :is="iconFor(att.type)" class="h-5 w-5 text-orange-500" />
      </div>
      <div class="flex-1 min-w-0">
        <p class="text-sm font-bold text-slate-700 dark:text-slate-200 truncate">{{ att.name }}</p>
        <p class="text-[11px] text-slate-400">
          {{ typeLabels[att.type] ?? att.type }}
          <span v-if="(att.type as string) === 'telegram_file' && (att as any).telegram_channel_id"> · قناة: {{ (att as any).telegram_channel_id }}</span>
          <span v-else-if="formatSize(att.size)"> · {{ formatSize(att.size) }}</span>
          <span v-if="att.created_at"> · {{ formatDate(att.created_at) }}</span>
        </p>
      </div>
      <button
        v-if="att.url"
        type="button"
        :disabled="downloadingId === att.id"
        @click="openAttachment(att)"
        class="p-2 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800 text-slate-400 hover:text-orange-500 transition-colors disabled:opacity-50 cursor-pointer"
        :title="att.type === 'link' ? 'فتح' : 'تحميل'"
      >
        <Loader2 v-if="downloadingId === att.id" class="h-4 w-4 animate-spin" />
        <Download v-else class="h-4 w-4" />
      </button>
      <button
        :disabled="deletingId === att.id"
        class="p-2 rounded-lg hover:bg-red-50 dark:hover:bg-red-950/40 text-slate-400 hover:text-red-600 transition-colors disabled:opacity-50"
        title="حذف"
        @click="emit('delete', att)"
      >
        <Loader2 v-if="deletingId === att.id" class="h-4 w-4 animate-spin" />
        <Trash2 v-else class="h-4 w-4" />
      </button>
    </div>

    <div
      v-if="!loading && !attachments.length"
      class="rounded-xl border border-dashed border-slate-200 dark:border-slate-800 py-6 text-center"
    >
      <FileText class="h-5 w-5 text-slate-300 dark:text-slate-600 mx-auto mb-1" />
      <p class="text-xs text-slate-400">لا توجد مرفقات</p>
    </div>
  </div>
</template>
