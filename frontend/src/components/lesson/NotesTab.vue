<script setup lang="ts">
import NoteEditor from "@/components/lesson/NoteEditor.vue";

defineProps<{
  lessonId: string;
}>();
</script>

<template>
  <div v-motion :initial="{ opacity: 0, y: 12 }" :enter="{ opacity: 1, y: 0 }">
    <NoteEditor :lesson-id="lessonId" />
  </div>
</template>
