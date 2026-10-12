import { defineStore } from "pinia";
import { ref } from "vue";

export const useHomeStore = defineStore("home", () => {
  const selectedCategory = ref<string>("all");
  const isAnnouncementDismissed = ref<boolean>(false);
  const isSearchOpen = ref<boolean>(false);
  const searchQuery = ref<string>("");

  function setCategory(category: string) {
    selectedCategory.value = category;
  }

  function dismissAnnouncement() {
    isAnnouncementDismissed.value = true;
  }

  function toggleSearch(open?: boolean) {
    isSearchOpen.value = open !== undefined ? open : !isSearchOpen.value;
  }

  return {
    selectedCategory,
    isAnnouncementDismissed,
    isSearchOpen,
    searchQuery,
    setCategory,
    dismissAnnouncement,
    toggleSearch,
  };
});
