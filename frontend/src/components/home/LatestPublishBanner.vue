<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from "vue";
import { Course } from "@/types/course";
import {
  Play,
  Sparkles,
  Flame,
  ArrowLeft,
  ChevronRight,
  ChevronLeft,
  BookOpen,
  User,
  GraduationCap,
  Clock,
} from "@lucide/vue";

const props = defineProps<{
  items?: Course[];
  latestItem?: Course | null;
}>();

const emit = defineEmits<{
  (e: "selectCourse", courseId: string): void;
}>();

// Collect all banner courses
const bannerCourses = computed<Course[]>(() => {
  if (props.items && props.items.length > 0) {
    return props.items.slice(0, 5);
  }
  if (props.latestItem) {
    return [props.latestItem];
  }
  return [];
});

const activeIndex = ref(0);
const isPaused = ref(false);
let timer: ReturnType<typeof setInterval> | null = null;

const currentCourse = computed<Course | null>(() => {
  if (bannerCourses.value.length === 0) return null;
  return bannerCourses.value[activeIndex.value] || bannerCourses.value[0];
});

function nextSlide() {
  if (bannerCourses.value.length <= 1) return;
  activeIndex.value = (activeIndex.value + 1) % bannerCourses.value.length;
}

function prevSlide() {
  if (bannerCourses.value.length <= 1) return;
  activeIndex.value =
    (activeIndex.value - 1 + bannerCourses.value.length) %
    bannerCourses.value.length;
}

function goToSlide(index: number) {
  activeIndex.value = index;
}

function startAutoSlide() {
  stopAutoSlide();
  if (bannerCourses.value.length > 1) {
    timer = setInterval(() => {
      if (!isPaused.value) {
        nextSlide();
      }
    }, 4500);
  }
}

function stopAutoSlide() {
  if (timer) {
    clearInterval(timer);
    timer = null;
  }
}

// Touch swipe gestures for mobile
const touchStartX = ref(0);
const touchEndX = ref(0);

function handleTouchStart(e: TouchEvent) {
  touchStartX.value = e.touches[0].clientX;
}

function handleTouchEnd(e: TouchEvent) {
  touchEndX.value = e.changedTouches[0].clientX;
  const diff = touchStartX.value - touchEndX.value;
  // RTL aware: swipe left means next, swipe right means prev
  if (Math.abs(diff) > 45) {
    if (diff > 0) {
      nextSlide();
    } else {
      prevSlide();
    }
  }
}

onMounted(() => {
  startAutoSlide();
});

onUnmounted(() => {
  stopAutoSlide();
});
</script>

<template>
  <div
    v-if="bannerCourses.length > 0 && currentCourse"
    class="w-full flex flex-col gap-2.5 text-start select-none"
    @mouseenter="isPaused = true"
    @mouseleave="isPaused = false"
  >
    <!-- Section Title & Status Indicator -->
    <div class="flex items-center justify-between px-1">
      <div class="flex items-center gap-2">
        <div class="relative flex h-3 w-3 items-center justify-center">
          <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-orange-400 opacity-75"></span>
          <span class="relative inline-flex rounded-full h-2 w-2 bg-orange-500"></span>
        </div>
        <h2 class="text-sm sm:text-base font-extrabold text-[#111111] dark:text-[#FFFFFF] tracking-tight">
          أحدث ما تم نشره
        </h2>
      </div>

      <div class="flex items-center gap-2">
        <span class="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full bg-orange-500/10 text-orange-500 text-[11px] font-bold">
          <Sparkles class="w-3 h-3" />
          <span>جديد المنصة</span>
        </span>
      </div>
    </div>

    <!-- Animated Moving Banner Card -->
    <div
      data-glass="banner"
      class="glass-thick relative w-full rounded-3xl overflow-hidden text-white cursor-pointer group transition-all duration-300"
      @touchstart="handleTouchStart"
      @touchend="handleTouchEnd"
      @click="emit('selectCourse', currentCourse.id)"
    >
      <!-- Aspect Ratio Canvas with Background Images -->
      <div class="relative w-full aspect-[16/9] sm:aspect-[21/8] min-h-[210px] overflow-hidden bg-[#181818]">
        <!-- Smooth Background Image Slide -->
        <transition name="slide-fade" mode="out-in">
          <div :key="currentCourse.id" class="absolute inset-0 w-full h-full">
            <img
              v-if="currentCourse.coverImageUrl"
              :src="currentCourse.coverImageUrl"
              :alt="currentCourse.title"
              class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-700 ease-out brightness-90"
            />
            <div
              v-else
              class="w-full h-full flex items-center justify-center bg-gradient-to-br from-orange-600 via-zinc-900 to-black"
            >
              <GraduationCap class="w-20 h-20 text-white/30 animate-pulse" />
            </div>
          </div>
        </transition>

        <!-- Cinematic Gradient Overlays -->
        <div class="absolute inset-0 bg-gradient-to-t from-black/95 via-black/45 to-black/20"></div>
        <div class="absolute inset-0 bg-gradient-to-r from-black/80 via-transparent to-transparent"></div>

        <!-- Top Floating Badges & Navigation Arrows -->
        <div class="absolute top-3.5 right-3.5 left-3.5 flex items-center justify-between z-10 pointer-events-none">
          <div class="flex items-center gap-2">
            <span class="inline-flex items-center gap-1.5 text-[11px] font-black uppercase tracking-wider bg-orange-500 text-white px-3 py-1 rounded-full shadow-lg shadow-orange-500/30">
              <Flame class="w-3.5 h-3.5 fill-white" />
              <span>نُشر حديثاً</span>
            </span>

            <span
              v-if="currentCourse.gradeName"
              class="text-[11px] font-bold bg-black/60 backdrop-blur-md text-white/90 px-3 py-1 rounded-full border border-white/10"
            >
              {{ currentCourse.gradeName }}
            </span>
          </div>

          <!-- Slider Arrow Controls (Desktop) -->
          <div
            v-if="bannerCourses.length > 1"
            class="hidden sm:flex items-center gap-1.5 pointer-events-auto"
            @click.stop
          >
            <button
              type="button"
              aria-label="السابق"
              @click.stop="prevSlide"
              class="w-8 h-8 rounded-full bg-black/50 hover:bg-orange-500 backdrop-blur-md text-white flex items-center justify-center transition-colors cursor-pointer active:scale-90"
            >
              <ChevronRight class="w-4 h-4" />
            </button>
            <button
              type="button"
              aria-label="التالي"
              @click.stop="nextSlide"
              class="w-8 h-8 rounded-full bg-black/50 hover:bg-orange-500 backdrop-blur-md text-white flex items-center justify-center transition-colors cursor-pointer active:scale-90"
            >
              <ChevronLeft class="w-4 h-4" />
            </button>
          </div>
        </div>

        <!-- Center Glowing Play Button -->
        <div class="absolute inset-0 flex items-center justify-center pointer-events-none z-10">
          <div class="w-14 h-14 rounded-full bg-orange-500/90 text-white flex items-center justify-center shadow-xl shadow-orange-500/50 group-hover:scale-110 group-hover:bg-orange-500 transition-all duration-300 backdrop-blur-xs">
            <Play class="w-6 h-6 fill-white ml-0.5" />
          </div>
        </div>

        <!-- Bottom Content Overlay -->
        <div class="absolute bottom-0 inset-x-0 p-4 sm:p-5 flex flex-col gap-2 z-10">
          <transition name="slide-up" mode="out-in">
            <div :key="currentCourse.id" class="flex flex-col gap-1.5">
              <h3 class="text-base sm:text-xl font-black text-white line-clamp-1 group-hover:text-orange-400 transition-colors drop-shadow-md">
                {{ currentCourse.title }}
              </h3>

              <div class="flex items-center justify-between text-xs text-zinc-300 font-medium">
                <div class="flex items-center gap-3">
                  <span v-if="currentCourse.instructor?.name" class="flex items-center gap-1 font-bold text-white">
                    <User class="w-3.5 h-3.5 text-orange-400" />
                    <span>{{ currentCourse.instructor.name }}</span>
                  </span>
                  <span v-if="currentCourse.lecturesCount > 0" class="flex items-center gap-1 text-zinc-400">
                    <BookOpen class="w-3.5 h-3.5" />
                    <span>{{ currentCourse.lecturesCount }} محاضرة</span>
                  </span>
                </div>

                <div class="inline-flex items-center gap-1 text-orange-400 font-bold text-xs group-hover:translate-x-[-4px] transition-transform">
                  <span>مشاهدة المساق</span>
                  <ArrowLeft class="w-3.5 h-3.5" />
                </div>
              </div>
            </div>
          </transition>

          <!-- Slider Pagination Dots / Progress Indicator -->
          <div
            v-if="bannerCourses.length > 1"
            class="flex items-center justify-center gap-1.5 pt-1"
            @click.stop
          >
            <button
              v-for="(_, idx) in bannerCourses"
              :key="idx"
              type="button"
              :aria-label="`شريحة ${idx + 1}`"
              @click.stop="goToSlide(idx)"
              :class="[
                'h-1.5 rounded-full transition-all duration-300 cursor-pointer',
                activeIndex === idx
                  ? 'w-6 bg-orange-500 shadow-sm shadow-orange-500/50'
                  : 'w-2 bg-white/40 hover:bg-white/70'
              ]"
            />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.slide-fade-enter-active,
.slide-fade-leave-active {
  transition: opacity 0.4s ease, transform 0.4s ease;
}

.slide-fade-enter-from {
  opacity: 0;
  transform: scale(1.02);
}

.slide-fade-leave-to {
  opacity: 0;
  transform: scale(0.98);
}

.slide-up-enter-active,
.slide-up-leave-active {
  transition: opacity 0.3s ease, transform 0.3s ease;
}

.slide-up-enter-from {
  opacity: 0;
  transform: translateY(6px);
}

.slide-up-leave-to {
  opacity: 0;
  transform: translateY(-6px);
}
</style>
