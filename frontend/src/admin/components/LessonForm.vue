<script setup lang="ts">
import { ref, reactive, watch } from 'vue'
import { Save, Loader2, Image, Lock, Unlock } from '@lucide/vue'
import type { Lesson, Section } from '@/admin/types/panel'

const props = defineProps<{
  sections: Section[]
  initial?: Partial<Lesson> | null
  loading?: boolean
  submitLabel?: string
}>()

const emit = defineEmits<{ submit: [payload: Record<string, unknown>] }>()

const form = reactive({
  title: '',
  description: '',
  section_id: '',
  thumbnail: '',
  status: 'draft',
  is_free: false,
})

watch(
  () => props.initial,
  (v) => {
    form.title = v?.title ?? ''
    form.description = v?.description ?? ''
    form.section_id = v?.section_id ?? props.sections[0]?.id ?? ''
    form.thumbnail = v?.thumbnail ?? ''
    form.status = v?.status ?? 'draft'
    form.is_free = v?.is_free ?? false
  },
  { immediate: true }
)

function submit() {
  if (!form.title.trim()) return
  emit('submit', {
    title: form.title.trim(),
    description: form.description.trim() || undefined,
    section_id: form.section_id || undefined,
    thumbnail: form.thumbnail.trim() || undefined,
    status: form.status,
    is_free: form.is_free,
  })
}
</script>

<template>
  <form class="space-y-4" @submit.prevent="submit">
    <div class="space-y-1.5">
      <label class="text-xs font-bold text-slate-600 dark:text-slate-300">عنوان المحاضرة *</label>
      <input
        v-model="form.title"
        type="text"
        placeholder="مثال: شرح درس الحركة المتسارعة"
        class="input w-full text-sm"
        required
      />
    </div>

    <div class="space-y-1.5">
      <label class="text-xs font-bold text-slate-600 dark:text-slate-300">الوصف</label>
      <textarea
        v-model="form.description"
        rows="3"
        placeholder="وصف مختصر لمحتوى المحاضرة"
        class="input w-full text-sm resize-none"
      />
    </div>

    <div class="grid sm:grid-cols-2 gap-4">
      <div class="space-y-1.5">
        <label class="text-xs font-bold text-slate-600 dark:text-slate-300">القسم</label>
        <select v-model="form.section_id" class="input w-full text-sm">
          <option value="">بدون قسم</option>
          <option v-for="s in sections" :key="s.id" :value="s.id">{{ s.title }}</option>
        </select>
      </div>

      <div class="space-y-1.5">
        <label class="text-xs font-bold text-slate-600 dark:text-slate-300">الحالة</label>
        <select v-model="form.status" class="input w-full text-sm">
          <option value="draft">مسودة</option>
          <option value="published">منشور</option>
          <option value="archived">مؤرشف</option>
        </select>
      </div>
    </div>

    <div class="space-y-1.5">
      <label class="text-xs font-bold text-slate-600 dark:text-slate-300 flex items-center gap-1.5">
        <Image class="h-3.5 w-3.5" />
        رابط الصورة المصغرة
      </label>
      <input
        v-model="form.thumbnail"
        type="url"
        placeholder="https://..."
        dir="ltr"
        class="input w-full text-sm"
      />
    </div>

    <label class="flex items-center justify-between rounded-xl border border-slate-200 dark:border-slate-800 p-3 cursor-pointer hover:border-orange-300 dark:hover:border-orange-800 transition-colors">
      <span class="flex items-center gap-2 text-sm font-bold text-slate-700 dark:text-slate-200">
        <Unlock class="h-4 w-4 text-emerald-500" />
        محاضرة مجانية (تظهر بدون اشتراك)
      </span>
      <button
        type="button"
        role="switch"
        :aria-checked="form.is_free"
        class="relative h-6 w-11 rounded-full transition-colors"
        :class="form.is_free ? 'bg-orange-500' : 'bg-slate-300 dark:bg-slate-700'"
        @click="form.is_free = !form.is_free"
      >
        <span
          class="absolute top-0.5 h-5 w-5 rounded-full bg-white shadow transition-all"
          :class="form.is_free ? 'left-0.5' : 'left-[22px]'"
        />
      </button>
    </label>

    <button type="submit" class="btn-primary w-full py-2.5" :disabled="loading || !form.title.trim()">
      <Loader2 v-if="loading" class="h-4 w-4 animate-spin" />
      <Save v-else class="h-4 w-4" />
      {{ submitLabel ?? 'حفظ المحاضرة' }}
    </button>
  </form>
</template>
