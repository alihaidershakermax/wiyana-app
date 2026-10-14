<script setup lang="ts">
import { computed, ref, onMounted } from "vue";
import { useRoute, useRouter } from "vue-router";
import { ArrowRight, PlayCircle } from "@lucide/vue";
import { apiClient } from "@/api/client";

const route = useRoute();
const router = useRouter();
const subjectId = computed(() => String(route.params.id || ""));

const subject = ref<any>(null);
const lessons = ref<any[]>([]);
const isLoading = ref(true);

async function load() {
  isLoading.value = true;
  try {
    // Get subject
    const sRes = await apiClient.get(`/subjects/${subjectId.value}`);
    subject.value = sRes.data?.data || sRes.data;

    // Get courses in this subject
    const cRes = await apiClient.get(`/subjects/${subjectId.value}/courses`);
    const courses = cRes.data?.data || cRes.data || [];

    // Get lessons for each course
    const allLessons: any[] = [];
    for (const c of courses) {
      try {
        const lRes = await apiClient.get(`/courses/${c.id}/lessons`);
        const ls = lRes.data?.data || lRes.data || [];
        allLessons.push(...ls.map((l: any) => ({ ...l, courseTitle: c.title })));
      } catch {}
    }
    lessons.value = allLessons;
  } catch (e) {
    console.error(e);
  } finally {
    isLoading.value = false;
  }
}

function openLesson(id: string) {
  router.push({ name: "lesson-player", params: { id } });
}

function goBack() {
  router.back();
}

onMounted(load);
</script>

<template>
  <div class="min-h-screen bg-white dark:bg-[#0F0F0F] pb-24">
    <!-- Header -->
    <div class="sticky top-0 z-10 bg-white/80 dark:bg-[#0F0F0F]/80 backdrop-blur-lg border-b border-black/5 dark:border-white/10">
      <div class="max-w-7xl mx-auto px-4 py-3 flex items-center gap-3">
        <button @click="goBack" class="p-2 rounded-full hover:bg-black/5 dark:hover:bg-white/10">
          <ArrowRight class="w-5 h-5" />
        </button>
        <h1 class="text-lg font-bold">{{ subject?.name || "المادة" }}</h1>
      </div>
    </div>

    <!-- Lessons list like YouTube -->
    <div class="max-w-7xl mx-auto px-4 py-4">
      <div v-if="isLoading" class="text-center py-10 opacity-60">جاري التحميل...</div>
      <div v-else-if="lessons.length === 0" class="text-center py-10 opacity-60">لا توجد محاضرات بعد</div>
      <div v-else class="flex flex-col gap-3">
        <div
          v-for="lesson in lessons"
          :key="lesson.id"
          @click="openLesson(lesson.id)"
          class="flex gap-3 cursor-pointer group"
        >
          <!-- Thumbnail like YouTube -->
          <div class="relative w-40 sm:w-56 aspect-video rounded-xl overflow-hidden bg-black/10 dark:bg-white/10 shrink-0">
            <img
              v-if="lesson.thumbnail"
              :src="lesson.thumbnail"
              class="w-full h-full object-cover"
              alt=""
            />
            <div class="absolute inset-0 flex items-center justify-center">
              <PlayCircle class="w-10 h-10 text-white/80 group-hover:text-white group-hover:scale-110 transition-all" />
            </div>
            <span v-if="lesson.duration" class="absolute bottom-1 left-1 bg-black/80 text-white text-[11px] px-1.5 py-0.5 rounded">
              {{ Math.floor(lesson.duration / 60) }}:{{ String(lesson.duration % 60).padStart(2, "0") }}
            </span>
          </div>
          <!-- Info like YouTube -->
          <div class="flex flex-col gap-1 py-1 min-w-0">
            <h3 class="text-sm font-bold leading-snug line-clamp-2">{{ lesson.title }}</h3>
            <span class="text-xs opacity-60">{{ lesson.courseTitle }}</span>
            <span v-if="lesson.description" class="text-xs opacity-50 line-clamp-1">{{ lesson.description }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
