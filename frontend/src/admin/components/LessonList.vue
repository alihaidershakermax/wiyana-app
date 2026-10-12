<script setup lang="ts">
import { computed } from 'vue'
import { ChevronDown, Plus, Layers } from '@lucide/vue'
import type { Lesson, Section } from '@/admin/types/panel'
import LessonCard from './LessonCard.vue'

const props = defineProps<{
  sections: Section[]
  lessons: Lesson[]
  loading?: boolean
  canCreate?: boolean
}>()

const emit = defineEmits<{
  edit: [lesson: Lesson]
  delete: [lesson: Lesson]
  publish: [lesson: Lesson]
  archive: [lesson: Lesson]
  create: [sectionId: string]
}>()

const grouped = computed(() => {
  return props.sections.map((section) => ({
    section,
    lessons: props.lessons
      .filter((l) => l.section_id === section.id)
      .sort((a, b) => a.order - b.order),
  }))
})

const unassigned = computed(() =>
  props.lessons.filter((l) => !l.section_id || !props.sections.some((s) => s.id === l.section_id))
)

const openSections = new Set<string>()
function toggleSection(id: string) {
  if (openSections.has(id)) openSections.delete(id)
  else openSections.add(id)
}
function isOpen(id: string) {
  return !openSections.has(id)
}
</script>

<template>
  <div class="space-y-6">
    <div v-for="{ section, lessons: sectionLessons } in grouped" :key="section.id" class="space-y-2">
      <div class="flex items-center justify-between px-1">
        <button class="flex items-center gap-2 group" @click="toggleSection(section.id)">
          <ChevronDown class="h-4 w-4 text-slate-400 transition-transform group-hover:text-slate-600" :class="isOpen(section.id) ? '' : '-rotate-90'" />
          <span class="text-sm font-bold text-slate-700 dark:text-slate-200">{{ section.title }}</span>
          <span class="text-[10px] font-bold bg-orange-50 dark:bg-orange-950/40 text-orange-600 dark:text-orange-400 px-2 py-0.5 rounded-full">
            {{ sectionLessons.length }} محاضرة
          </span>
        </button>
        <button
          v-if="canCreate"
          class="flex items-center gap-1 text-xs font-bold text-orange-600 dark:text-orange-400 hover:text-orange-700 transition-colors"
          @click="emit('create', section.id)"
        >
          <Plus class="h-3.5 w-3.5" />
          إضافة محاضرة
        </button>
      </div>

      <div v-if="isOpen(section.id)" class="space-y-2">
        <LessonCard
          v-for="lesson in sectionLessons"
          :key="lesson.id"
          :lesson="lesson"
          :draggable="true"
          @edit="emit('edit', $event)"
          @delete="emit('delete', $event)"
          @publish="emit('publish', $event)"
          @archive="emit('archive', $event)"
        />
        <div
          v-if="!sectionLessons.length"
          class="rounded-2xl border border-dashed border-slate-200 dark:border-slate-800 py-6 text-center"
        >
          <Layers class="h-5 w-5 text-slate-300 dark:text-slate-600 mx-auto mb-1" />
          <p class="text-xs text-slate-400">لا توجد محاضرات في هذا القسم</p>
        </div>
      </div>
    </div>

    <div v-if="unassigned.length" class="space-y-2">
      <div class="px-1">
        <span class="text-sm font-bold text-slate-700 dark:text-slate-200">بدون قسم</span>
        <span class="text-[10px] font-bold bg-slate-100 dark:bg-slate-800 text-slate-500 px-2 py-0.5 rounded-full mr-2">
          {{ unassigned.length }}
        </span>
      </div>
      <LessonCard
        v-for="lesson in unassigned"
        :key="lesson.id"
        :lesson="lesson"
        :draggable="true"
        @edit="emit('edit', $event)"
        @delete="emit('delete', $event)"
        @publish="emit('publish', $event)"
        @archive="emit('archive', $event)"
      />
    </div>

    <div v-if="!loading && !props.lessons.length" class="rounded-2xl border border-dashed border-slate-200 dark:border-slate-800 py-10 text-center">
      <Layers class="h-8 w-8 text-slate-300 dark:text-slate-600 mx-auto mb-2" />
      <p class="text-sm text-slate-500 dark:text-slate-400">لا توجد محاضرات بعد</p>
      <p v-if="canCreate" class="text-xs text-slate-400 mt-1">ابدأ بإضافة أول محاضرة</p>
    </div>
  </div>
</template>
