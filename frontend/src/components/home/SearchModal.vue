<script setup lang="ts">
import { ref, computed } from "vue";
import { Search, X, BookOpen, ChevronLeft } from "@lucide/vue";
import { Course } from "@/types/course";

const props = defineProps<{
  isOpen: boolean;
  courses: Course[];
}>();

const emit = defineEmits<{
  (e: "close"): void;
  (e: "selectCourse", courseId: string): void;
}>();

const query = ref("");

const filteredCourses = computed(() => {
  if (!query.value.trim()) return [];
  const q = query.value.trim().toLowerCase();
  return props.courses.filter(
    (c) =>
      c.title.toLowerCase().includes(q) ||
      (c.departmentName && c.departmentName.toLowerCase().includes(q)) ||
      (c.gradeName && c.gradeName.toLowerCase().includes(q)) ||
      (c.instructor?.name && c.instructor.name.toLowerCase().includes(q))
  );
});
</script>

<template>
  <Transition
    enter-active-class="transition duration-200 ease-out"
    enter-from-class="opacity-0"
    enter-to-class="opacity-100"
    leave-active-class="transition duration-150 ease-in"
    leave-from-class="opacity-100"
    leave-to-class="opacity-0"
  >
    <div
      v-if="isOpen"
      class="fixed inset-0 z-50 flex items-start justify-center p-4 sm:p-6 pt-16 sm:pt-20 bg-black/70 backdrop-blur-md"
      @click.self="emit('close')"
    >
      <div
        class="w-full max-w-2xl bg-white dark:bg-[#181818] border border-[#E5E7EB] dark:border-[#27272A] rounded-3xl shadow-2xl overflow-hidden flex flex-col"
      >
        <!-- Search Input Bar -->
        <div class="p-4 border-b border-[#E5E7EB] dark:border-[#27272A] flex items-center gap-3">
          <Search class="w-5 h-5 text-orange-500 shrink-0" />
          <input
            v-model="query"
            type="text"
            placeholder="ابحث عن مساق، مرحلة، أو أستاذ..."
            class="w-full bg-transparent text-[#111111] dark:text-[#FFFFFF] text-sm font-semibold focus:outline-none placeholder:text-[#9CA3AF] dark:placeholder:text-[#71717A]"
            autoFocus
          />
          <button
            @click="emit('close')"
            type="button"
            class="p-2 rounded-xl text-[#6B7280] hover:text-[#111111] dark:text-[#A1A1AA] dark:hover:text-[#FFFFFF] hover:bg-[#F8F8F8] dark:hover:bg-[#222222] transition-colors cursor-pointer"
          >
            <X class="w-5 h-5" />
          </button>
        </div>

        <!-- Search Results List -->
        <div class="p-4 max-h-[60vh] overflow-y-auto flex flex-col gap-2">
          <div v-if="query.trim() === ''" class="py-8 text-center text-[#6B7280] dark:text-[#A1A1AA] text-xs sm:text-sm">
            اكتب اسم المساق أو المادة للبدء بالبحث السريع.
          </div>

          <div v-else-if="filteredCourses.length === 0" class="py-8 text-center text-[#6B7280] dark:text-[#A1A1AA] text-xs sm:text-sm">
            لم نجد نتائج مطابقة لـ "{{ query }}"
          </div>

          <button
            v-for="course in filteredCourses"
            :key="course.id"
            @click="emit('selectCourse', course.id); emit('close')"
            type="button"
            class="flex items-center justify-between gap-3 p-3 rounded-2xl hover:bg-[#F8F8F8] dark:hover:bg-[#222222] transition-colors text-right w-full group cursor-pointer"
          >
            <div class="flex items-center gap-3 min-w-0">
              <div class="w-12 h-12 rounded-xl bg-[#F8F8F8] dark:bg-[#222222] overflow-hidden shrink-0">
                <img :src="course.coverImageUrl" :alt="course.title" class="w-full h-full object-cover" />
              </div>
              <div class="flex flex-col min-w-0">
                <span class="text-xs font-bold text-orange-500">{{ course.subjectName }}</span>
                <span class="text-sm font-bold text-[#111111] dark:text-[#FFFFFF] truncate group-hover:text-orange-500 transition-colors">
                  {{ course.title }}
                </span>
                <span class="text-[11px] text-[#6B7280] dark:text-[#A1A1AA]">{{ course.instructor.name }}</span>
              </div>
            </div>

            <ChevronLeft class="w-4 h-4 text-[#9CA3AF] dark:text-[#71717A] group-hover:text-orange-500 transition-colors" />
          </button>
        </div>
      </div>
    </div>
  </Transition>
</template>
