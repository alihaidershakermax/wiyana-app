<script setup lang="ts">
import { BookOpen } from "@lucide/vue";

interface Subject {
  id: string;
  name: string;
  description?: string;
  courses_count?: number;
}

defineProps<{
  subjects: Subject[];
}>();

const emit = defineEmits<{
  selectSubject: [id: string];
}>();
</script>

<template>
  <div v-if="subjects && subjects.length > 0" class="flex flex-col gap-4">
    <div class="flex items-center justify-between px-1">
      <h2 class="text-base sm:text-lg font-extrabold tracking-tight">المواد الدراسية</h2>
    </div>
    <div class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-3">
      <div
        v-for="subject in subjects"
        :key="subject.id"
        data-glass="card"
        class="glass-card cursor-pointer p-4 flex flex-col items-center gap-3 text-center transition-transform hover:scale-[1.02] active:scale-95"
        @click="emit('selectSubject', subject.id)"
      >
        <div class="w-14 h-14 rounded-2xl bg-orange-500/15 flex items-center justify-center">
          <BookOpen class="w-7 h-7 text-orange-500" />
        </div>
        <div class="flex flex-col gap-1">
          <h3 class="text-sm font-bold leading-snug">{{ subject.name }}</h3>
          <span v-if="subject.courses_count" class="text-xs opacity-60">{{ subject.courses_count }} مساق</span>
        </div>
      </div>
    </div>
  </div>
</template>
