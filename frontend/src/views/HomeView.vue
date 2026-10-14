<script setup lang="ts">
import { computed, ref, onMounted } from "vue";
import { useRouter } from "vue-router";
import { useHomeData } from "@/composables/useHomeData";
import { useHomeStore } from "@/stores/useHomeStore";
import { Course } from "@/types/course";
import { apiClient } from "@/api/client";

// Components
import OfflineState from "@/components/ui/OfflineState.vue";
import HomeHeader from "@/components/home/HomeHeader.vue";
import LatestPublishBanner from "@/components/home/LatestPublishBanner.vue";
import SubjectsGrid from "@/components/home/SubjectsGrid.vue";
import HomeSkeleton from "@/components/ui/HomeSkeleton.vue";
import ErrorState from "@/components/ui/ErrorState.vue";
import SearchModal from "@/components/home/SearchModal.vue";

// TanStack Query & Pinia Store
const { data, isLoading, isError, refetch } = useHomeData();
const homeStore = useHomeStore();
const router = useRouter();

// Subjects (المواد الدراسية)
const subjects = ref<any[]>([]);
onMounted(async () => {
  try {
    const res = await apiClient.get("/subjects/");
    const d = res.data;
    subjects.value = Array.isArray(d) ? d : d?.data || [];
  } catch {}
});

function handleSubjectClick(subjectId: string) {
  if (subjectId) router.push({ name: "subject-detail", params: { id: subjectId } });
}

// Find the very latest published item as fallback
const latestPublishedItem = computed<Course | null>(() => {
  if (!data.value) return null;
  return data.value.latestCourses[0] || data.value.recommendedCourses[0] || null;
});

// All courses combined for search modal
const allCoursesCombined = computed<Course[]>(() => {
  if (!data.value) return [];
  const map = new Map<string, Course>();
  [
    ...data.value.continueLearningCourses,
    ...data.value.recommendedCourses,
    ...data.value.latestCourses,
    ...data.value.popularCourses,
  ].forEach((c) => map.set(c.id, c));
  return Array.from(map.values());
});

function handleCourseClick(courseId: string) {
  if (courseId) {
    router.push(`/courses/${courseId}`);
  }
}

function handleViewAllCourses() {
  router.push("/courses");
}

function handleSearchClick() {
  router.push("/search");
}
</script>

<template>
  <div class="min-h-screen bg-white dark:bg-[#0F0F0F] text-[#111111] dark:text-[#FFFFFF] flex flex-col font-sans transition-colors duration-200" data-page="home">
    <!-- Offline Connection Warning -->
    <OfflineState @refresh="refetch()" />

    <!-- Fixed Header (only when data loaded) -->
    <HomeHeader
      v-if="data"
      :user="data.user"
      @click-search="handleSearchClick"
      @click-notifications="router.push('/notifications')"
      @click-profile="router.push('/profile')"
    />

    <!-- Main Content Area -->
    <main class="flex-1 max-w-7xl mx-auto page-gutter py-4 sm:py-6 w-full flex flex-col gap-5 sm:gap-6" :aria-busy="isLoading ? 'true' : 'false'">
      <!-- Loading Skeleton State -->
      <div v-if="isLoading" role="status" aria-live="polite" aria-label="جاري تحميل الصفحة الرئيسية">
        <HomeSkeleton />
      </div>

      <!-- Error State -->
      <section v-else-if="isError" class="state-panel" aria-live="polite">
        <ErrorState @retry="refetch()" />
      </section>

      <!-- Loaded Dashboard Content -->
      <div v-else-if="data" class="flex flex-col gap-5 sm:gap-6 animate-fade-in" aria-label="محتوى الصفحة الرئيسية">
        <!-- 1. Animated Latest Published Banner (بنر متحرك أحدث ما تم نشره) -->
        <LatestPublishBanner
          :items="data.latestCourses"
          :latest-item="latestPublishedItem"
          @select-course="handleCourseClick"
        />

        <!-- 2. Subjects Grid (المواد الدراسية مباشرة) -->
        <SubjectsGrid
          :subjects="subjects"
          @select-subject="handleSubjectClick"
        />
      </div>
    </main>

    <!-- Search Modal -->
    <SearchModal
      :is-open="homeStore.isSearchOpen"
      :courses="allCoursesCombined"
      @close="homeStore.toggleSearch(false)"
      @select-course="handleCourseClick"
    />
  </div>
</template>

<style scoped>
.animate-fade-in {
  animation: fadeIn 0.25s cubic-bezier(0.16, 1, 0.3, 1) forwards;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(4px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
