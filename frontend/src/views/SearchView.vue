<script setup lang="ts">
import { ref, watch } from "vue";
import { useRouter } from "vue-router";
import { useQuery } from "@tanstack/vue-query";
import { Search, X, ChevronRight, BookOpen, Loader2 } from "@lucide/vue";
import { searchPlatform } from "@/api/search";

const router = useRouter();
const searchQuery = ref("");
let searchTimeout: ReturnType<typeof setTimeout> | null = null;

const { data: searchResults, isLoading } = useQuery({
  queryKey: ["platformSearch", searchQuery],
  queryFn: () => searchPlatform(searchQuery.value.trim()),
  enabled: ref(false),
  staleTime: 1000 * 60 * 2,
});

// Debounced search trigger
watch(searchQuery, (val) => {
  if (searchTimeout) clearTimeout(searchTimeout);
  if (!val.trim()) return;
  searchTimeout = setTimeout(() => {
    searchPlatform(val.trim()).then((res) => {
      (searchResults as any).value = res;
    });
  }, 300);
});

function goBack() {
  router.back();
}

function openCourse(courseId: string) {
  router.push(`/courses/${courseId}`);
}
</script>

<template>
  <div class="min-h-screen bg-white dark:bg-[#0F0F0F] text-[#111111] dark:text-[#FFFFFF] flex flex-col font-sans transition-colors pb-20 text-start">
    <!-- Top Search Header -->
    <header class="sticky top-0 z-30 w-full backdrop-blur-xl bg-white/90 dark:bg-[#0F0F0F]/90 border-b border-[#E5E7EB] dark:border-[#27272A] p-3.5">
      <div class="max-w-7xl mx-auto flex items-center gap-3">
        <button
          @click="goBack"
          type="button"
          aria-label="الرجوع"
          class="w-11 h-11 rounded-2xl bg-[#F8F8F8] dark:bg-[#181818] text-[#111111] dark:text-[#FFFFFF] border border-[#E5E7EB] dark:border-[#27272A] flex items-center justify-center active:scale-90 transition-transform shrink-0 cursor-pointer"
        >
          <ChevronRight class="w-6 h-6" />
        </button>

        <!-- Search Bar Input -->
        <div class="flex-1 relative flex items-center">
          <Search class="w-5 h-5 text-orange-500 absolute right-3 pointer-events-none" />
          <input
            v-model="searchQuery"
            type="text"
            placeholder="ابحث عن مساق، مرحلة، أو مدرس..."
            class="w-full h-11 bg-[#F8F8F8] dark:bg-[#181818] text-[#111111] dark:text-[#FFFFFF] pr-10 pl-10 rounded-2xl text-xs sm:text-sm font-semibold focus:outline-none focus:ring-2 focus:ring-orange-500/30 placeholder:text-[#9CA3AF] dark:placeholder:text-[#71717A] border border-transparent focus:border-orange-500"
            autoFocus
          />
          <button
            v-if="searchQuery.length > 0"
            @click="searchQuery = ''"
            type="button"
            class="absolute left-3 text-[#6B7280] dark:text-[#A1A1AA] p-1 cursor-pointer"
          >
            <X class="w-4 h-4" />
          </button>
        </div>
      </div>
    </header>

    <!-- Search Results Content -->
    <main class="flex-1 max-w-7xl mx-auto px-4 py-4 w-full flex flex-col gap-3">
      <div v-if="!searchQuery.trim()" class="py-12 text-center text-[#6B7280] dark:text-[#A1A1AA] text-xs sm:text-sm">
        <BookOpen class="w-12 h-12 text-orange-500/30 mx-auto mb-3" />
        <p>اكتب اسم المساق، المرحلة أو المدرس للبدء بالبحث.</p>
      </div>

      <div v-else-if="isLoading" class="py-12 text-center text-[#6B7280] dark:text-[#A1A1AA] text-xs sm:text-sm flex items-center justify-center gap-2">
        <Loader2 class="w-5 h-5 animate-spin text-orange-500" />
        <span>جاري البحث...</span>
      </div>

      <div v-else-if="!searchResults?.courses?.length && searchResults !== undefined" class="py-12 text-center text-[#6B7280] dark:text-[#A1A1AA] text-xs sm:text-sm">
        لم نجد نتائج مطابقة لـ "{{ searchQuery }}"
      </div>

      <div v-else-if="searchResults?.courses">
        <div
          v-for="course in searchResults.courses"
          :key="course.id || course._id"
          @click="openCourse(String(course.id || course._id))"
          class="flex items-center justify-between p-3 bg-white dark:bg-[#181818] border border-[#E5E7EB] dark:border-[#27272A] rounded-2xl cursor-pointer hover:border-orange-500/40 active:scale-98 transition-all"
        >
          <div class="flex items-center gap-3 min-w-0">
            <img :src="course.thumbnail_url || course.coverImageUrl" :alt="course.title" class="w-16 h-16 rounded-xl object-cover shrink-0" />
            <div class="flex flex-col min-w-0">
              <span class="text-[10px] font-bold text-orange-500 uppercase">{{ course.stage || course.department || 'مساق تعليمي' }}</span>
              <h3 class="text-xs sm:text-sm font-bold text-[#111111] dark:text-[#FFFFFF] truncate">{{ course.title }}</h3>
              <span class="text-[11px] text-[#6B7280] dark:text-[#A1A1AA]">{{ course.teacher_name || course.instructor?.name }}</span>
            </div>
          </div>
          <ChevronRight class="w-4 h-4 text-[#9CA3AF] dark:text-[#71717A] rotate-180 shrink-0" />
        </div>
      </div>
    </main>
  </div>
</template>
