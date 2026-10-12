<script setup lang="ts">
import { ref, computed } from "vue";
import { useRouter } from "vue-router";
import { useCoursesData } from "@/composables/useCoursesData";
import { useCourseStore } from "@/stores/useCourseStore";
import { Course } from "@/types/course";

// Components
import OfflineState from "@/components/ui/OfflineState.vue";
import CoursesHeader from "@/components/courses/CoursesHeader.vue";
import SearchBar from "@/components/courses/SearchBar.vue";
import CategoryTabs from "@/components/courses/CategoryTabs.vue";
import FilterSheet from "@/components/courses/FilterSheet.vue";
import CourseGrid from "@/components/courses/CourseGrid.vue";
import SkeletonCard from "@/components/ui/SkeletonCard.vue";
import EmptyState from "@/components/ui/EmptyState.vue";
import ErrorState from "@/components/ui/ErrorState.vue";

const router = useRouter();

// TanStack Query & Pinia Store
const { courses, isLoading, isError, refetch, metadata } = useCoursesData();
const courseStore = useCourseStore();

const isSearchVisible = ref<boolean>(true);

// Reactive Filtered & Sorted Courses
const filteredCourses = computed<Course[]>(() => {
  if (!courses.value) return [];
  let result = [...courses.value];

  // 1. Search Query Filter
  if (courseStore.searchQuery.trim()) {
    const q = courseStore.searchQuery.trim().toLowerCase();
    result = result.filter(
      (c) =>
        c.title.toLowerCase().includes(q) ||
        (c.departmentName && c.departmentName.toLowerCase().includes(q)) ||
        (c.gradeName && c.gradeName.toLowerCase().includes(q)) ||
        (c.instructor?.name && c.instructor.name.toLowerCase().includes(q))
    );
  }

  // 2. Tab Filter
  const tab = courseStore.activeTab;
  if (tab === "my_courses") {
    result = result.filter((c) => c.enrollmentStatus !== "not_enrolled");
  } else if (tab === "newest") {
    result = result.filter((c) => c.categoryId === "newest");
  } else if (tab === "free") {
    result = result.filter((c) => c.isFree);
  } else if (tab === "paid") {
    result = result.filter((c) => !c.isFree);
  } else if (tab === "popular") {
    result = result.filter((c) => c.categoryId === "popular" || c.viewsCount > 8000);
  }

  // 3. Department Filter
  if (courseStore.selectedDepartment) {
    const deptObj = metadata.value?.departments.find((d) => d.id === courseStore.selectedDepartment);
    if (deptObj) {
      result = result.filter((c) => c.departmentName === deptObj.name);
    }
  }

  // 4. Grade Filter
  if (courseStore.selectedGrade) {
    const gradeObj = metadata.value?.grades.find((g) => g.id === courseStore.selectedGrade);
    if (gradeObj) {
      result = result.filter((c) => c.gradeName === gradeObj.name);
    }
  }

  // 5. Instructor Filter
  if (courseStore.selectedInstructor) {
    const instObj = metadata.value?.instructors.find((i) => i.id === courseStore.selectedInstructor);
    if (instObj) {
      result = result.filter((c) => c.instructor.name === instObj.name);
    }
  }

  // 7. Free/Paid/Subscription Filters from FilterSheet
  if (courseStore.isFreeOnly) {
    result = result.filter((c) => c.isFree);
  }
  if (courseStore.isPaidOnly) {
    result = result.filter((c) => !c.isFree);
  }
  if (courseStore.isSubscribedOnly) {
    result = result.filter((c) => c.enrollmentStatus === "enrolled");
  }

  // 8. Sorting
  const sort = courseStore.sortBy;
  result.sort((a, b) => {
    if (sort === "newest") {
      return new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime();
    }
    if (sort === "oldest") {
      return new Date(a.createdAt).getTime() - new Date(b.createdAt).getTime();
    }
    if (sort === "popular") {
      return b.viewsCount - a.viewsCount;
    }
    if (sort === "top_rated") {
      return b.rating - a.rating;
    }
    if (sort === "alphabetical") {
      return a.title.localeCompare(b.title, "ar");
    }
    return 0;
  });

  return result;
});

function handleCourseClick(courseId: string) {
  if (courseId) {
    router.push(`/courses/${courseId}`);
  }
}

function handleCourseAction(courseId: string) {
  if (courseId) {
    router.push(`/courses/${courseId}`);
  }
}
</script>

<template>
  <div class="min-h-screen bg-white dark:bg-[#0F0F0F] text-[#111111] dark:text-[#FFFFFF] flex flex-col font-sans transition-colors duration-200" data-page="courses">
    <!-- Offline Connection Notice -->
    <OfflineState @refresh="refetch()" />

    <!-- Courses Header -->
    <CoursesHeader
      :total-count="filteredCourses.length"
      :active-filters-count="courseStore.activeFiltersCount"
      :is-search-visible="isSearchVisible"
      @toggle-search="isSearchVisible = !isSearchVisible"
      @open-filters="courseStore.toggleFilterSheet(true)"
    />

    <!-- Search Input Bar -->
    <Transition
      enter-active-class="transition duration-200 ease-out"
      enter-from-class="opacity-0 -translate-y-2"
      enter-to-class="opacity-100 translate-y-0"
      leave-active-class="transition duration-150 ease-in"
      leave-from-class="opacity-100 translate-y-0"
      leave-to-class="opacity-0 -translate-y-2"
    >
      <SearchBar
        v-if="isSearchVisible"
        :model-value="courseStore.searchQuery"
        @update:model-value="courseStore.setSearchQuery"
        @clear="courseStore.setSearchQuery('')"
      />
    </Transition>

    <!-- Main Content Container -->
    <main class="flex-1 max-w-7xl mx-auto page-gutter py-6 w-full flex flex-col gap-6" :aria-busy="isLoading ? 'true' : 'false'">
      <!-- Category Tabs Bar -->
      <CategoryTabs
        :active-tab="courseStore.activeTab"
        @select-tab="courseStore.setActiveTab"
      />

      <!-- Loading State Skeletons -->
      <div v-if="isLoading" class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4" role="status" aria-live="polite" aria-label="جاري تحميل المساقات">
        <SkeletonCard v-for="i in 6" :key="i" />
      </div>

      <!-- Error State -->
      <section v-else-if="isError" class="state-panel" aria-live="polite">
        <ErrorState @retry="refetch()" />
      </section>

      <!-- Empty State -->
      <section v-else-if="filteredCourses.length === 0" class="state-panel">
        <EmptyState
        title="لم نجد أي مساقات مطابقة"
        description="جرب البحث بكلمات مختلفة أو قم بإعادة تعيين فلاتر البحث الحالية."
        action-text="إعادة تعيين الفلاتر"
        @click-action="courseStore.resetFilters()"
        />
      </section>

      <!-- Course Grid Results -->
      <CourseGrid
        v-else
        :courses="filteredCourses"
        @select-course="handleCourseClick"
        @course-action="handleCourseAction"
      />
    </main>

    <!-- Filter Bottom Sheet Modal -->
    <FilterSheet
      :is-open="courseStore.isFilterSheetOpen"
      :metadata="metadata"
      :selected-department="courseStore.selectedDepartment"
      :selected-grade="courseStore.selectedGrade"
      :selected-instructor="courseStore.selectedInstructor"
      :is-free-only="courseStore.isFreeOnly"
      :is-paid-only="courseStore.isPaidOnly"
      :is-subscribed-only="courseStore.isSubscribedOnly"
      :sort-by="courseStore.sortBy"
      @close="courseStore.toggleFilterSheet(false)"
      @update:department="(v) => (courseStore.selectedDepartment = v)"
      @update:grade="(v) => (courseStore.selectedGrade = v)"
      @update:instructor="(v) => (courseStore.selectedInstructor = v)"
      @update:free-only="(v) => (courseStore.isFreeOnly = v)"
      @update:paid-only="(v) => (courseStore.isPaidOnly = v)"
      @update:subscribed-only="(v) => (courseStore.isSubscribedOnly = v)"
      @update:sort-by="(v) => courseStore.setSortBy(v)"
      @reset="courseStore.resetFilters()"
    />
  </div>
</template>
