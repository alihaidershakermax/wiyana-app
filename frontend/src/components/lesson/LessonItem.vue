<script setup lang="ts">
import { computed } from "vue";
import { CheckCircle2, Play, Lock, ChevronLeft, Clock } from "@lucide/vue";
import { formatShortDuration } from "@/utils/format";
import { usePlayerStore } from "@/stores/usePlayerStore";
import type { Lesson } from "@/types/lesson";

const props = withDefaults(
  defineProps<{
    lesson: Lesson;
    isCurrent?: boolean;
    isEnrolled?: boolean;
  }>(),
  {
    isCurrent: false,
    isEnrolled: false,
  }
);

const emit = defineEmits<{
  (e: "select", lessonId: string): void;
}>();

const playerStore = usePlayerStore();

const isDone = computed(
  () => props.lesson.isCompleted === true || playerStore.isCompleted(props.lesson.id)
);

const isLocked = computed(
  () => !props.lesson.isFreePreview && !props.isEnrolled && !props.isCurrent
);

const durationText = computed(() => formatShortDuration(props.lesson.durationSeconds));

function onSelect() {
  if (isLocked.value) return;
  emit("select", props.lesson.id);
}
</script>

<template>
  <button
    @click="onSelect"
    type="button"
    :disabled="isLocked"
    :class="[
      'w-full flex items-center justify-between gap-3 p-3.5 rounded-2xl border transition-all text-right',
      isCurrent
        ? 'bg-orange-500/10 dark:bg-orange-500/15 border-orange-500 ring-1 ring-orange-500/40'
        : 'bg-white dark:bg-[#181818] border-[#E5E7EB] dark:border-[#27272A] hover:border-orange-500/40',
      isLocked ? 'opacity-60 cursor-not-allowed' : 'cursor-pointer active:scale-[0.98]'
    ]"
  >
    <div class="flex items-center gap-3 min-w-0">
      <div
        :class="[
          'w-9 h-9 rounded-xl flex items-center justify-center shrink-0',
          isDone
            ? 'bg-orange-500/10 text-orange-500'
            : isCurrent
              ? 'bg-orange-500 text-white'
              : isLocked
                ? 'bg-[#F8F8F8] dark:bg-[#222222] text-[#9CA3AF] dark:text-[#71717A]'
                : 'bg-orange-500/10 text-orange-500'
        ]"
      >
        <CheckCircle2 v-if="isDone" class="w-4 h-4" />
        <Play v-else-if="isCurrent" class="w-4 h-4 fill-white ml-0.5" />
        <Lock v-else-if="isLocked" class="w-4 h-4" />
        <Play v-else class="w-4 h-4 fill-orange-500 ml-0.5" />
      </div>

      <div class="flex flex-col gap-0.5 min-w-0">
        <span
          :class="[
            'text-xs sm:text-sm font-bold truncate',
            isCurrent ? 'text-orange-600 dark:text-orange-400' : 'text-[#111111] dark:text-[#FFFFFF]'
          ]"
        >
          {{ lesson.title }}
        </span>
        <span class="inline-flex items-center gap-1.5 text-[10px] text-[#6B7280] dark:text-[#A1A1AA]">
          <span v-if="durationText" class="inline-flex items-center gap-1">
            <Clock class="w-3 h-3" />
            {{ durationText }}
          </span>
          <span
            v-if="isLocked"
            class="px-1.5 py-0.5 rounded-full text-[9px] font-bold bg-[#F8F8F8] dark:bg-[#222222] text-[#6B7280] dark:text-[#A1A1AA]"
          >
            مقفلة
          </span>
          <span
            v-else-if="lesson.isFreePreview"
            class="px-1.5 py-0.5 rounded-full text-[9px] font-bold bg-orange-500/10 text-orange-600 dark:text-orange-400"
          >
            معاينة مجانية
          </span>
          <span
            v-if="isCurrent"
            class="px-1.5 py-0.5 rounded-full text-[9px] font-bold bg-orange-500/15 text-orange-600 dark:text-orange-400"
          >
            الحالية
          </span>
        </span>
      </div>
    </div>

    <ChevronLeft class="w-4 h-4 text-[#9CA3AF] dark:text-[#71717A] shrink-0" />
  </button>
</template>
