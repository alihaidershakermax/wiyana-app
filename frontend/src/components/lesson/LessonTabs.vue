<script setup lang="ts">
import { computed } from "vue";
import { ListVideo, Paperclip, StickyNote } from "@lucide/vue";

export type LessonTabKey = "lessons" | "attachments" | "notes";

const props = withDefaults(
  defineProps<{
    modelValue: LessonTabKey;
    lessonsCount?: number;
    attachmentsCount?: number;
  }>(),
  {
    lessonsCount: 0,
    attachmentsCount: 0,
  }
);

const emit = defineEmits<{
  (e: "update:modelValue", value: LessonTabKey): void;
}>();

const tabs = computed<
  Array<{
    key: LessonTabKey;
    label: string;
    icon: any;
    count?: number;
  }>
>(() => [
  { key: "lessons", label: "المحاضرات", icon: ListVideo, count: props.lessonsCount },
  { key: "attachments", label: "الملازم والمرفقات", icon: Paperclip, count: props.attachmentsCount },
  { key: "notes", label: "ملاحظاتي", icon: StickyNote },
]);

function select(key: LessonTabKey) {
  emit("update:modelValue", key);
}
</script>

<template>
  <div class="flex border-b border-[#E5E7EB] dark:border-[#27272A]" role="tablist">
    <button
      v-for="tab in tabs"
      :key="tab.key"
      @click="select(tab.key)"
      type="button"
      role="tab"
      :aria-selected="modelValue === tab.key"
      :class="[
        'flex-1 py-3 text-xs sm:text-sm font-bold text-center border-b-2 transition-all inline-flex items-center justify-center gap-1.5 min-w-0 cursor-pointer',
        modelValue === tab.key
          ? 'border-orange-500 text-orange-500'
          : 'border-transparent text-[#6B7280] hover:text-[#111111] dark:text-[#A1A1AA] dark:hover:text-white'
      ]"
    >
      <component :is="tab.icon" class="w-4 h-4 shrink-0" />
      <span class="truncate">{{ tab.label }}</span>
      <span
        v-if="tab.count !== undefined && tab.count > 0"
        class="px-1.5 py-0.5 rounded-full text-[9px] font-black bg-orange-500/10 text-orange-500"
      >
        {{ tab.count }}
      </span>
    </button>
  </div>
</template>
