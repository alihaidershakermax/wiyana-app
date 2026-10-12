<script setup lang="ts">
import { Paperclip } from "@lucide/vue";
import AttachmentCard from "@/components/lesson/AttachmentCard.vue";
import EmptyState from "@/components/ui/EmptyState.vue";
import type { LessonAttachment } from "@/types/lesson";

defineProps<{
  attachments: LessonAttachment[];
}>();
</script>

<template>
  <div v-if="attachments.length === 0" v-motion :initial="{ opacity: 0, y: 12 }" :enter="{ opacity: 1, y: 0 }">
    <EmptyState
      title="لا توجد مرفقات"
      description="لا توجد ملفات أو روابط مرفقة بهذه المحاضرة."
    >
      <template #icon>
        <Paperclip class="w-8 h-8" />
      </template>
    </EmptyState>
  </div>

  <div v-else class="flex flex-col gap-2.5">
    <div
      v-for="(attachment, index) in attachments"
      :key="attachment.id"
      v-motion
      :initial="{ opacity: 0, y: 10 }"
      :enter="{ opacity: 1, y: 0, transition: { delay: Math.min(index * 40, 300) } }"
    >
      <AttachmentCard :attachment="attachment" />
    </div>
  </div>
</template>
