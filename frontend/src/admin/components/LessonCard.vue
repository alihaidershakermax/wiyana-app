<script setup lang="ts">
import { computed } from 'vue'
import { Play, Clock, Eye, Lock, Unlock, GripVertical, FileText, Edit, Trash2, Send, Archive } from '@lucide/vue'
import type { Lesson } from '@/admin/types/panel'
import { formatShortDuration, formatDate } from '@/utils/format'
import ActionDropdown, { type MenuItem } from '@/admin/components/ui/ActionDropdown.vue'

const props = defineProps<{
  lesson: Lesson
  draggable?: boolean
}>()

const emit = defineEmits<{
  edit: [lesson: Lesson]
  delete: [lesson: Lesson]
  publish: [lesson: Lesson]
  archive: [lesson: Lesson]
  view: [lesson: Lesson]
}>()

const statusMeta = computed(() => {
  const map = {
    draft: { label: 'مسودة', cls: 'bg-[#18181b] text-[#A1A1AA] border-[#27272a]' },
    published: { label: 'منشور', cls: 'bg-[#052e16] text-[#22C55E] border-[#14532d]' },
    archived: { label: 'مؤرشف', cls: 'bg-[#451a03] text-[#F59E0B] border-[#78350f]' },
  }
  return map[props.lesson.status] ?? map.draft
})

const hasVideo = computed(() => !!props.lesson.youtube_video_id && !props.lesson.youtube_video_id?.startsWith('yt_sim_'))
const isSimulated = computed(() => !!props.lesson.youtube_video_id?.startsWith('yt_sim_'))

function getActionItems(): MenuItem[] {
  const items: MenuItem[] = [
    {
      id: 'edit',
      label: 'تعديل المحاضرة',
      icon: Edit,
      onClick: () => emit('edit', props.lesson),
    },
  ]

  if (props.lesson.status !== 'published') {
    items.push({
      id: 'publish',
      label: 'نشر المحاضرة',
      icon: Send,
      onClick: () => emit('publish', props.lesson),
    })
  } else {
    items.push({
      id: 'archive',
      label: 'أرشفة المحاضرة',
      icon: Archive,
      onClick: () => emit('archive', props.lesson),
    })
  }

  items.push({
    id: 'delete',
    label: 'حذف المحاضرة',
    icon: Trash2,
    danger: true,
    onClick: () => emit('delete', props.lesson),
  })

  return items
}
</script>

<template>
  <div
    class="group relative flex items-center gap-4 rounded-2xl border border-[#242424] bg-[#111111] p-3.5 hover:border-[#F97316]/40 transition-colors text-[#F5F5F5]"
  >
    <span v-if="draggable" class="cursor-grab text-[#71717A] hover:text-white shrink-0">
      <GripVertical class="h-5 w-5" />
    </span>

    <!-- Thumbnail -->
    <div class="relative h-16 w-28 shrink-0 rounded-xl overflow-hidden bg-[#080808] border border-[#242424]">
      <img
        v-if="lesson.thumbnail"
        :src="lesson.thumbnail"
        :alt="lesson.title"
        class="h-full w-full object-cover"
        loading="lazy"
      />
      <div v-else class="h-full w-full flex items-center justify-center">
        <Play class="h-6 w-6 text-[#F97316]" />
      </div>
      <div class="absolute inset-0 bg-gradient-to-t from-black/60 to-transparent" />
      <span class="absolute bottom-1 left-1 rounded bg-black/80 text-white text-[10px] font-bold px-1.5 py-0.5 flex items-center gap-1">
        <Clock class="h-3 w-3" />
        {{ formatShortDuration(lesson.duration ?? 0) ?? '—' }}
      </span>
      <span v-if="isSimulated" class="absolute top-1 right-1 rounded bg-[#F97316] text-white text-[9px] font-bold px-1 py-0.5">
        تجريبي
      </span>
    </div>

    <!-- Info -->
    <div class="flex-1 min-w-0">
      <div class="flex items-center gap-2">
        <h3 class="text-sm font-extrabold text-white truncate">{{ lesson.title }}</h3>
        <span :class="['text-[10px] font-extrabold px-2.5 py-0.5 rounded-full border shrink-0', statusMeta.cls]">
          {{ statusMeta.label }}
        </span>
        <span
          v-if="!hasVideo"
          :class="[
            'text-[10px] font-extrabold px-2.5 py-0.5 rounded-full border shrink-0',
            lesson.youtube_video_id ? 'bg-[#451a03] text-[#F59E0B] border-[#78350f]' : 'bg-[#450a0a] text-[#EF4444] border-[#7f1d1d]'
          ]"
        >
          {{ lesson.youtube_video_id ? 'قيد المعالجة' : 'بدون فيديو' }}
        </span>
      </div>
      <p class="text-xs text-[#A1A1AA] mt-1 truncate">{{ lesson.description || 'لا يوجد وصف' }}</p>
      <div class="flex items-center gap-4 mt-2 text-[11px] text-[#71717A]">
        <span class="flex items-center gap-1">
          <Eye class="h-3.5 w-3.5" />
          {{ lesson.views }} مشاهدة
        </span>
        <span class="flex items-center gap-1">
          <FileText class="h-3.5 w-3.5" />
          {{ lesson.attachments?.length ?? 0 }} مرفقات
        </span>
        <span class="flex items-center gap-1">
          <template v-if="lesson.is_free">
            <Unlock class="h-3.5 w-3.5 text-[#22C55E]" />
            مجانية
          </template>
          <template v-else>
            <Lock class="h-3.5 w-3.5 text-[#71717A]" />
            مدفوعة
          </template>
        </span>
        <span class="hidden sm:inline">{{ formatDate(lesson.created_at) }}</span>
      </div>
    </div>

    <!-- Teleported Action Menu -->
    <div class="shrink-0">
      <ActionDropdown :items="getActionItems()" />
    </div>
  </div>
</template>
