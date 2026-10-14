<script setup lang="ts">
import { computed, onMounted, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { Lock } from "@lucide/vue";
import { useLessonPlayer } from "@/composables/useLessonPlayer";
import LessonHeader from "@/components/lesson/LessonHeader.vue";
import VideoPlayer from "@/components/lesson/VideoPlayer.vue";
import ProtectedVideoPlayer from "@/components/lesson/ProtectedVideoPlayer.vue";
import DashPlayer from "@/components/video/DashPlayer.vue";
import LessonInfo from "@/components/lesson/LessonInfo.vue";
import CourseProgress from "@/components/lesson/CourseProgress.vue";
import LessonNavigation from "@/components/lesson/LessonNavigation.vue";
import LessonTabs, { type LessonTabKey } from "@/components/lesson/LessonTabs.vue";
import LessonsTab from "@/components/lesson/LessonsTab.vue";
import AttachmentsTab from "@/components/lesson/AttachmentsTab.vue";
import NotesTab from "@/components/lesson/NotesTab.vue";
import SkeletonPlayer from "@/components/lesson/SkeletonPlayer.vue";
import OfflineState from "@/components/ui/OfflineState.vue";
import ErrorState from "@/components/ui/ErrorState.vue";
import { checkAccess } from "@/api/protection";
import type { AccessStatus } from "@/types/protection";

const route = useRoute();
const router = useRouter();

const lessonId = computed(() => String(route.params.id || ""));

const { bundle, isLoading, isError, refetch } = useLessonPlayer(lessonId);

const activeTab = ref<LessonTabKey>("lessons");

const currentLesson = computed(() => bundle.value?.lesson ?? null);
const course = computed(() => bundle.value?.course ?? null);
const lessons = computed(() => bundle.value?.lessons ?? []);
const isEnrolled = computed(() => bundle.value?.isEnrolled ?? false);
const isLocked = computed(() => currentLesson.value?.isLocked ?? false);

// DASH encrypted video detection (file_path ending in .mpd)
const isDashVideo = computed(
  () => !!currentLesson.value?.file_path?.endsWith(".mpd")
);
const dashUrl = computed(() => currentLesson.value?.file_path ?? "");

// Content-protection gate: file/HLS videos use the protected player; YouTube uses standard player
const isProtectedVideo = computed(
  () =>
    !!currentLesson.value?.video &&
    (currentLesson.value.video.kind === "file" || currentLesson.value.video.kind === "hls")
);
const accessStatus = ref<AccessStatus | null>(null);
const accessChecked = ref(false);
const accessBanned = computed(() => accessStatus.value?.banned ?? false);

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
  activeTab.value = "lessons";
  window.scrollTo({ top: 0 });
});

onMounted(async () => {
  try {
    accessStatus.value = await checkAccess();
  } catch {
    accessStatus.value = null;
  } finally {
    accessChecked.value = true;
  }
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
      <!-- Banned: full page block gate -->
      <div
        v-if="accessChecked && accessBanned"
        class="min-h-screen flex flex-col items-center justify-center gap-4 px-6 text-center"
      >
        <div class="w-16 h-16 rounded-full bg-red-500/15 flex items-center justify-center">
          <Lock class="w-8 h-8 text-red-500" />
        </div>
        <div class="flex flex-col gap-1.5 max-w-md">
          <h3 class="text-base sm:text-lg font-black">حسابك محظور مؤقتاً</h3>
          <p class="text-xs sm:text-sm font-medium opacity-70 leading-relaxed">
            تم تجاوز الحد المسموح من محاولات التقاط المحتوى. يمكنك إعادة التحقق لاحقاً بعد انتهاء مدة الحظر.
          </p>
        </div>
        <div class="flex flex-col gap-2 w-full max-w-xs">
          <button
            type="button"
            class="px-5 py-2.5 rounded-xl bg-orange-500 hover:bg-orange-600 active:scale-95 transition-all text-white text-sm font-black cursor-pointer"
            @click="router.go(0)"
          >
            إعادة التحقق من الوصول
          </button>
          <button
            type="button"
            class="px-5 py-2.5 rounded-xl bg-white/10 hover:bg-white/20 active:scale-95 transition-all text-sm font-black cursor-pointer"
            @click="closeApp"
          >
            إغلاق
          </button>
        </div>
      </div>

      <!-- Normal content -->
      <template v-else>
      <!-- Offline Connection Warning -->
      <OfflineState @refresh="refetch()" />

      <!-- Fixed Header -->
      <LessonHeader
        :title="currentLesson.title"
        :course-title="course?.title"
        :course-id="course?.id"
      />

      <!-- Video Player Container with Clean Fit -->
      <div class="w-full max-w-7xl mx-auto sm:px-4 sm:pt-3">
        <div class="w-full aspect-video bg-slate-950 sm:rounded-2xl overflow-hidden shadow-2xl border border-white/5 relative">
          <!-- Locked lesson: paywall before the player -->
          <div
            v-if="isLocked"
            class="relative w-full h-full bg-slate-950 flex flex-col items-center justify-center gap-4 px-6 text-center"
          >
            <div class="w-16 h-16 rounded-full bg-orange-500/15 flex items-center justify-center">
              <Lock class="w-7 h-7 text-orange-500" />
            </div>
            <div class="flex flex-col gap-1.5 max-w-md">
              <h3 class="text-base sm:text-lg font-black text-white">هذه المحاضرة مقفلة</h3>
              <p class="text-xs sm:text-sm font-medium text-white/60 leading-relaxed">
                للوصول إلى هذه المحاضرة وجميع محاضرات المساق، يرجى تفعيل اشتراكك أولاً.
              </p>
            </div>
            <button
              @click="goToCourse"
              type="button"
              class="px-5 py-2.5 rounded-xl bg-orange-500 hover:bg-orange-600 active:scale-95 transition-all text-white text-sm font-black cursor-pointer shadow-lg shadow-orange-500/25"
            >
              عرض المساق والاشتراك
            </button>
          </div>

          <!-- DASH encrypted player -->
          <DashPlayer
            v-if="isDashVideo"
            :key="currentLesson.id"
            :src="dashUrl"
          />

          <!-- Protected player for file/HLS videos -->
          <ProtectedVideoPlayer
            v-else-if="isProtectedVideo && currentLesson.video"
            :key="currentLesson.id"
            :video="currentLesson.video"
            :lesson-id="currentLesson.id"
            :poster="course?.thumbnailUrl"
          />

          <!-- Regular player (YouTube / Embed) -->
          <VideoPlayer
            v-else
            :key="currentLesson.id"
            :video="currentLesson.video || { id: currentLesson.id, kind: 'youtube', url: currentLesson.youtube_url ?? '', youtube_video_id: currentLesson.youtube_video_id, youtube_url: currentLesson.youtube_url }"
            :lesson-id="currentLesson.id"
            :poster="course?.thumbnailUrl"
          />
        </div>
      </div>

      <!-- Main Content Area -->
      <main class="flex-1 max-w-7xl mx-auto px-4 py-5 w-full flex flex-col gap-5">
        <!-- Lesson Info -->
        <LessonInfo v-if="course" :lesson="currentLesson" :course="course" />

        <!-- Course Progress -->
        <CourseProgress :lessons="lessons" />

        <!-- Previous / Next Navigation -->
        <LessonNavigation
          :lessons="lessons"
          :current-lesson-id="currentLesson.id"
          @navigate="navigateToLesson"
        />

        <!-- Tabs -->
        <LessonTabs
          v-model="activeTab"
          :lessons-count="lessons.length"
          :attachments-count="currentLesson.attachments.length"
        />

        <!-- Tab Content -->
        <Transition name="tab-fade" mode="out-in">
          <LessonsTab
            v-if="activeTab === 'lessons'"
            :key="`lessons-${currentLesson.id}`"
            :lessons="lessons"
            :current-lesson-id="currentLesson.id"
            :is-enrolled="isEnrolled"
            @select="navigateToLesson"
          />
          <AttachmentsTab
            v-else-if="activeTab === 'attachments'"
            :key="`attachments-${currentLesson.id}`"
            :attachments="currentLesson.attachments"
          />
          <NotesTab
            v-else-if="activeTab === 'notes'"
            :key="`notes-${currentLesson.id}`"
            :lesson-id="currentLesson.id"
          />
        </Transition>
      </main>
      </template>
    </template>
  </div>
</template>

<style scoped>
.tab-fade-enter-active,
.tab-fade-leave-active {
  transition: opacity 0.15s ease, transform 0.15s ease;
}
.tab-fade-enter-from {
  opacity: 0;
  transform: translateY(4px);
}
.tab-fade-leave-to {
  opacity: 0;
  transform: translateY(-4px);
}
</style>
