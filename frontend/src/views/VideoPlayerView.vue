<script setup lang="ts">
import { computed, onMounted, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { Lock, BookOpen, ThumbsUp, Share2, Download, ChevronDown, ChevronUp, FileText, StickyNote } from "@lucide/vue";
import { useLessonPlayer } from "@/composables/useLessonPlayer";
import LessonHeader from "@/components/lesson/LessonHeader.vue";
import VideoPlayer from "@/components/lesson/VideoPlayer.vue";
import ProtectedVideoPlayer from "@/components/lesson/ProtectedVideoPlayer.vue";
import DashPlayer from "@/components/video/DashPlayer.vue";
import LessonInfo from "@/components/lesson/LessonInfo.vue";
import CourseProgress from "@/components/lesson/CourseProgress.vue";
import LessonNavigation from "@/components/lesson/LessonNavigation.vue";
import SkeletonPlayer from "@/components/lesson/SkeletonPlayer.vue";
import OfflineState from "@/components/ui/OfflineState.vue";
import ErrorState from "@/components/ui/ErrorState.vue";
import { checkAccess } from "@/api/protection";
import type { AccessStatus } from "@/types/protection";
const route = useRoute();
const router = useRouter();

const lessonId = computed(() => String(route.params.id || ""));

const { bundle, isLoading, isError, refetch } = useLessonPlayer(lessonId);

const currentLesson = computed(() => bundle.value?.lesson ?? null);
const course = computed(() => bundle.value?.course ?? null);
const lessons = computed(() => bundle.value?.lessons ?? []);
const isLocked = computed(() => false); // Subscriptions removed — all content free

// Plain video (mp4 — no protection, no watermark)
const isPlainVideo = computed(() => {
  const fp = currentLesson.value?.file_path ?? "";
  return fp.endsWith(".mp4") || fp.endsWith(".webm") || fp.endsWith(".mov");
});
const plainUrl = computed(() => currentLesson.value?.file_path ?? "");

// DASH video (mpd — dash.js reference player, no DRM)
const isDashVideo = computed(() => {
  const fp = currentLesson.value?.file_path ?? "";
  return fp.endsWith(".mpd");
});
const dashUrl = computed(() => currentLesson.value?.file_path ?? "");

// YouTube-style UI state
const descExpanded = ref(false);
const showAttachments = ref(false);
const showNotes = ref(false);
const liked = ref(false);

function toggleLike() { liked.value = !liked.value; }

function shareLesson() {
  const url = window.location.href;
  if (navigator.share) {
    navigator.share({ title: currentLesson.value?.title, url }).catch(() => {});
  } else {
    navigator.clipboard?.writeText(url);
  }
}

function downloadVideo() {
  const a = document.createElement("a");
  a.href = plainUrl.value;
  a.download = `${currentLesson.value?.title || "lecture"}.mp4`;
  a.target = "_blank";
  document.body.appendChild(a);
  a.click();
  a.remove();
}

function navigateToLesson(id: string) {
  if (id === lessonId.value) return;
  router.push({ name: "lesson-player", params: { id } });
}

function goToCourse() {
  if (course.value?.id) router.push({ name: "course-detail", params: { id: course.value.id } });
}

function closeApp() {
  try {
    (window as any).Telegram?.WebApp?.close?.();
  } catch {}
}

watch(lessonId, () => {
  descExpanded.value = false;
  showAttachments.value = false;
  showNotes.value = false;
  window.scrollTo({ top: 0 });
});
</script>

<template>
  <div class="min-h-screen bg-white dark:bg-[#0F0F0F] text-[#111111] dark:text-[#FFFFFF] flex flex-col font-sans transition-colors duration-200 pb-24 text-start">
    <!-- Loading State -->
    <SkeletonPlayer v-if="isLoading && !bundle" />

    <!-- Error State -->
    <div v-else-if="isError || !currentLesson" class="min-h-screen flex flex-col">
      <LessonHeader title="مشغل المحاضرة" />
      <div class="flex-1 max-w-7xl mx-auto w-full px-4">
        <ErrorState
          title="تعذر تحميل المحاضرة"
          message="حدثت مشكلة أثناء تحميل بيانات المحاضرة. تحقق من اتصالك وحاول مرة أخرى."
          @retry="refetch()"
        />
      </div>
    </div>

    <!-- Player Screen -->
    <template v-else>
      <!-- Normal content -->
      <template>
      <!-- Offline Connection Warning -->
      <OfflineState @refresh="refetch()" />

      <!-- Fixed Header -->
      <LessonHeader
        :title="currentLesson.title"
        :course-title="course?.title"
        :course-id="course?.id"
      />

      <!-- YouTube mobile style: video sticky on top -->
      <!-- Video Player -->
      <div class="w-full bg-black sticky top-0 z-20">
        <div class="w-full aspect-video">
          <!-- DASH player (dash.js reference, no DRM) -->
          <DashPlayer
            v-if="isDashVideo"
            :key="currentLesson.id"
            :src="dashUrl"
          />

          <!-- Plain video (no protection, no watermark) -->
          <video
            v-else-if="isPlainVideo"
            :key="currentLesson.id"
            :src="plainUrl"
            controls
            playsinline
            preload="metadata"
            class="w-full h-full"
          />

          <!-- Fallback: protected player -->
          <ProtectedVideoPlayer
            v-else-if="currentLesson.video && (currentLesson.video.kind === 'file' || currentLesson.video.kind === 'hls')"
            :key="currentLesson.id"
            :video="currentLesson.video"
            :lesson-id="currentLesson.id"
            :poster="course?.thumbnailUrl"
          />

          <!-- Fallback: regular player -->
          <VideoPlayer
            v-else
            :key="currentLesson.id"
            :video="currentLesson.video || { id: currentLesson.id, kind: 'youtube', url: currentLesson.youtube_url ?? '', youtube_video_id: currentLesson.youtube_video_id, youtube_url: currentLesson.youtube_url }"
            :lesson-id="currentLesson.id"
            :poster="course?.thumbnailUrl"
          />
        </div>
      </div>

      <!-- YouTube-style info -->
      <main class="flex-1 w-full px-3 py-3 flex flex-col gap-3">
        <!-- Title -->
        <h1 class="text-[15px] font-bold leading-snug line-clamp-2">{{ currentLesson.title }}</h1>

        <!-- Channel row + actions (YouTube style) -->
        <div class="flex items-center justify-between gap-2">
          <div class="flex items-center gap-2.5 min-w-0">
            <div class="w-9 h-9 rounded-full bg-orange-500/15 flex items-center justify-center overflow-hidden shrink-0">
              <img v-if="course?.thumbnailUrl" :src="course.thumbnailUrl" class="w-full h-full object-cover" alt="" />
              <BookOpen v-else class="w-4 h-4 text-orange-500" />
            </div>
            <div class="flex flex-col min-w-0">
              <span class="text-[13px] font-bold truncate">{{ course?.subjectName || course?.title }}</span>
              <span class="text-[11px] opacity-60">{{ lessons.length }} محاضرة</span>
            </div>
          </div>
        </div>

        <!-- Action buttons row (YouTube style pills) -->
        <div class="flex items-center gap-2 overflow-x-auto pb-1">
          <button
            type="button"
            @click="toggleLike"
            class="flex items-center gap-1.5 px-4 py-2 rounded-full text-[13px] font-bold shrink-0 transition-colors"
            :class="liked ? 'bg-orange-500/20 text-orange-500' : 'bg-black/5 dark:bg-white/10'"
          >
            <ThumbsUp class="w-4 h-4" :fill="liked ? 'currentColor' : 'none'" />
            أعجبني
          </button>
          <button
            type="button"
            @click="shareLesson"
            class="flex items-center gap-1.5 px-4 py-2 rounded-full bg-black/5 dark:bg-white/10 text-[13px] font-bold shrink-0"
          >
            <Share2 class="w-4 h-4" />
            مشاركة
          </button>
          <button
            v-if="isPlainVideo"
            type="button"
            @click="downloadVideo"
            class="flex items-center gap-1.5 px-4 py-2 rounded-full bg-black/5 dark:bg-white/10 text-[13px] font-bold shrink-0"
          >
            <Download class="w-4 h-4" />
            تحميل
          </button>
        </div>

        <!-- Description box (YouTube style, expandable) -->
        <button
          v-if="currentLesson.description"
          type="button"
          @click="descExpanded = !descExpanded"
          class="rounded-xl bg-black/5 dark:bg-white/10 p-3 text-start w-full"
        >
          <p class="text-[13px] leading-relaxed" :class="{ 'line-clamp-2': !descExpanded }">
            {{ currentLesson.description }}
          </p>
          <span class="text-[11px] font-bold opacity-60 mt-1 flex items-center gap-1">
            {{ descExpanded ? "عرض أقل" : "عرض المزيد" }}
            <ChevronDown v-if="!descExpanded" class="w-3 h-3" />
            <ChevronUp v-else class="w-3 h-3" />
          </span>
        </button>

        <!-- Attachments (collapsible) -->
        <div v-if="currentLesson.attachments?.length" class="rounded-xl bg-black/5 dark:bg-white/10 overflow-hidden">
          <button
            type="button"
            @click="showAttachments = !showAttachments"
            class="w-full flex items-center justify-between p-3"
          >
            <span class="text-[13px] font-bold flex items-center gap-2">
              <FileText class="w-4 h-4" />
              المرفقات ({{ currentLesson.attachments.length }})
            </span>
            <ChevronDown v-if="!showAttachments" class="w-4 h-4 opacity-60" />
            <ChevronUp v-else class="w-4 h-4 opacity-60" />
          </button>
          <div v-if="showAttachments" class="px-3 pb-3 flex flex-col gap-2">
            <a
              v-for="att in currentLesson.attachments"
              :key="att.id"
              :href="att.url"
              target="_blank"
              class="flex items-center gap-2 p-2.5 rounded-lg bg-black/5 dark:bg-white/5 text-[13px] font-medium"
            >
              <FileText class="w-4 h-4 shrink-0 opacity-60" />
              <span class="truncate">{{ att.name || att.title }}</span>
            </a>
          </div>
        </div>

        <!-- Up next -->
        <div class="flex flex-col gap-2 mt-1">
          <h2 class="text-[15px] font-bold">محاضرات أخرى</h2>
          <button
            v-for="l in lessons.filter((x: any) => x.id !== currentLesson.id)"
            :key="l.id"
            type="button"
            @click="navigateToLesson(l.id)"
            class="flex gap-3 p-2 rounded-xl hover:bg-black/5 dark:hover:bg-white/5 transition-colors text-start w-full"
          >
            <div class="w-32 aspect-video rounded-lg bg-black/10 dark:bg-white/10 overflow-hidden shrink-0 relative">
              <img v-if="l.thumbnail" :src="l.thumbnail" class="w-full h-full object-cover" alt="" />
              <div v-else class="w-full h-full flex items-center justify-center">
                <BookOpen class="w-6 h-6 opacity-30" />
              </div>
              <span v-if="l.duration" class="absolute bottom-1 end-1 px-1.5 py-0.5 rounded bg-black/80 text-white text-[10px] font-bold">
                {{ Math.floor(l.duration / 60) }}:{{ String(l.duration % 60).padStart(2, "0") }}
              </span>
            </div>
            <div class="flex flex-col gap-1 min-w-0 py-0.5">
              <h3 class="text-[13px] font-bold leading-snug line-clamp-2">{{ l.title }}</h3>
              <span class="text-[11px] opacity-60">{{ course?.subjectName || course?.title }}</span>
            </div>
          </button>
          <p v-if="!lessons.filter((x: any) => x.id !== currentLesson.id).length" class="text-[13px] opacity-60 text-center py-4">
            لا توجد محاضرات أخرى بعد
          </p>
        </div>
      </main>
      </template>
    </template>
  </div>
</template>

<style scoped>
/* YouTube mobile player styles */
</style>
