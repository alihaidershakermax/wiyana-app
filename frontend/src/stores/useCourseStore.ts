import { defineStore } from "pinia";
import { ref, computed } from "vue";
import { CourseTabCategory, CourseSortOption } from "@/types/course";

export const useCourseStore = defineStore("course", () => {
  const searchQuery = ref<string>("");
  const activeTab = ref<CourseTabCategory>("all");
  const selectedDepartment = ref<string>("");
  const selectedGrade = ref<string>("");
  const selectedSubject = ref<string>("");
  const selectedInstructor = ref<string>("");
  const isFreeOnly = ref<boolean>(false);
  const isPaidOnly = ref<boolean>(false);
  const isSubscribedOnly = ref<boolean>(false);
  const sortBy = ref<CourseSortOption>("newest");
  const isFilterSheetOpen = ref<boolean>(false);

  const activeFiltersCount = computed(() => {
    let count = 0;
    if (selectedDepartment.value) count++;
    if (selectedGrade.value) count++;
    if (selectedSubject.value) count++;
    if (selectedInstructor.value) count++;
    if (isFreeOnly.value) count++;
    if (isPaidOnly.value) count++;
    if (isSubscribedOnly.value) count++;
    if (sortBy.value !== "newest") count++;
    return count;
  });

  function setSearchQuery(query: string) {
    searchQuery.value = query;
  }

  function setActiveTab(tab: CourseTabCategory) {
    activeTab.value = tab;
  }

  function setSortBy(sort: CourseSortOption) {
    sortBy.value = sort;
  }

  function toggleFilterSheet(open?: boolean) {
    isFilterSheetOpen.value = open !== undefined ? open : !isFilterSheetOpen.value;
  }

  function resetFilters() {
    searchQuery.value = "";
    activeTab.value = "all";
    selectedDepartment.value = "";
    selectedGrade.value = "";
    selectedSubject.value = "";
    selectedInstructor.value = "";
    isFreeOnly.value = false;
    isPaidOnly.value = false;
    isSubscribedOnly.value = false;
    sortBy.value = "newest";
  }

  return {
    searchQuery,
    activeTab,
    selectedDepartment,
    selectedGrade,
    selectedSubject,
    selectedInstructor,
    isFreeOnly,
    isPaidOnly,
    isSubscribedOnly,
    sortBy,
    isFilterSheetOpen,
    activeFiltersCount,
    setSearchQuery,
    setActiveTab,
    setSortBy,
    toggleFilterSheet,
    resetFilters,
  };
});
