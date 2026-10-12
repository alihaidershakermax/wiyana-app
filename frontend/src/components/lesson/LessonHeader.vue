<script setup lang="ts">
import { onMounted, onUnmounted, ref } from "vue";
import { useRouter } from "vue-router";
import { ChevronRight, MoreVertical, BookOpen } from "@lucide/vue";

const props = withDefaults(
  defineProps<{
    title: string;
    courseTitle?: string;
    courseId?: string;
  }>(),
  {
    courseTitle: "",
    courseId: "",
  }
);

const router = useRouter();
const menuOpen = ref(false);
const menuRef = ref<HTMLDivElement | null>(null);

function goBack() {
  if (window.history.length > 1) {
    router.back();
  } else {
    router.push("/");
  }
}

function toggleMenu() {
  menuOpen.value = !menuOpen.value;
}

function openCourse() {
  menuOpen.value = false;
  if (props.courseId) {
    router.push(`/courses/${props.courseId}`);
  }
}

function onDocumentClick(event: MouseEvent) {
  if (menuRef.value && !menuRef.value.contains(event.target as Node)) {
    menuOpen.value = false;
  }
}

onMounted(() => {
  document.addEventListener("click", onDocumentClick);
});

onUnmounted(() => {
  document.removeEventListener("click", onDocumentClick);
});
</script>

<template>
  <header class="sticky top-0 z-30 w-full backdrop-blur-xl bg-white/90 dark:bg-[#0F0F0F]/90 border-b border-[#E5E7EB] dark:border-[#27272A] px-4 py-3 flex items-center gap-3">
    <button
      @click="goBack"
      type="button"
      aria-label="الرجوع"
      class="w-11 h-11 rounded-2xl bg-[#F8F8F8] dark:bg-[#181818] text-[#111111] dark:text-white border border-[#E5E7EB] dark:border-[#27272A] flex items-center justify-center active:scale-90 transition-transform shrink-0 cursor-pointer"
    >
      <ChevronRight class="w-6 h-6" />
    </button>

    <div class="flex flex-col min-w-0 flex-1 text-start">
      <span v-if="courseTitle" class="text-[10px] font-bold text-orange-500 uppercase tracking-wider truncate">
        {{ courseTitle }}
      </span>
      <h1 class="text-xs sm:text-sm font-bold text-[#111111] dark:text-white truncate">
        {{ title }}
      </h1>
    </div>

    <div ref="menuRef" class="relative shrink-0">
      <button
        @click.stop="toggleMenu"
        type="button"
        aria-label="المزيد"
        class="w-11 h-11 rounded-2xl bg-[#F8F8F8] dark:bg-[#181818] text-[#111111] dark:text-white border border-[#E5E7EB] dark:border-[#27272A] flex items-center justify-center active:scale-90 transition-transform cursor-pointer"
      >
        <MoreVertical class="w-5 h-5" />
      </button>

      <Transition
        enter-active-class="transition-all duration-150 ease-out"
        enter-from-class="opacity-0 scale-95 translate-y-1"
        enter-to-class="opacity-100 scale-100 translate-y-0"
        leave-active-class="transition-all duration-100 ease-in"
        leave-from-class="opacity-100 scale-100 translate-y-0"
        leave-to-class="opacity-0 scale-95 translate-y-1"
      >
        <div
          v-if="menuOpen"
          class="absolute left-0 top-14 w-52 bg-white dark:bg-[#181818] border border-[#E5E7EB] dark:border-[#27272A] rounded-2xl shadow-xl p-1.5 z-50 text-start"
        >
          <button
            v-if="courseId"
            @click="openCourse"
            type="button"
            class="w-full flex items-center gap-2.5 px-3.5 py-2.5 rounded-xl text-xs font-bold text-[#111111] dark:text-[#A1A1AA] hover:bg-[#F8F8F8] dark:hover:bg-[#222222] transition-colors cursor-pointer"
          >
            <BookOpen class="w-4 h-4 text-orange-500" />
            <span>فتح صفحة الكورس</span>
          </button>
        </div>
      </Transition>
    </div>
  </header>
</template>
