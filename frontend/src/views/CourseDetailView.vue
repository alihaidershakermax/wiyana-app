<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useQuery, useMutation, useQueryClient } from "@tanstack/vue-query";
import {
  ChevronRight,
  Play,
  Lock,
  FileText,
  Users,
  Clock,
  BookOpen,
  CheckCircle2,
  Bookmark,
  KeyRound,
  Sparkles,
} from "@lucide/vue";
import { fetchCourseById } from "@/api/courses";
import { checkCourseFavorite, addCourseToFavorites, removeCourseFromFavorites } from "@/api/favorites";
import { redeemSubscriptionCode } from "@/api/subscription";
import { useUserStore } from "@/stores/useUserStore";
import ProgressBar from "@/components/ui/ProgressBar.vue";
import SkeletonCard from "@/components/ui/SkeletonCard.vue";
import ErrorState from "@/components/ui/ErrorState.vue";

const route = useRoute();
const router = useRouter();
const queryClient = useQueryClient();
const userStore = useUserStore();

const courseId = computed(() => String(route.params.id || ""));
const activeTab = ref<"lectures" | "description" | "files">("lectures");
const isBookmarked = ref(false);
const isRedeemModalOpen = ref(false);
const redeemCodeInput = ref("");
const redeemError = ref("");
const isRedeeming = ref(false);

const TELEGRAM_BOT_USERNAME = "WenaAppBot";

function formatFileSize(bytes: number): string {
  if (!bytes || bytes === 0) return "0 B";
  const k = 1024;
  const sizes = ["B", "KB", "MB", "GB"];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + " " + sizes[i];
}

function downloadViaTelegram(file: any) {
  if (!file || !file.telegram_file_id) {
    userStore.showToast("الملف غير متاح للتنزيل عبر تليجرام", "error");
    return;
  }
  
  const telegramLink = `https://t.me/${TELEGRAM_BOT_USERNAME}?start=file_${file.telegram_file_id}`;
  window.open(telegramLink, "_blank");
  userStore.showToast("جاري تحويلك إلى بوت تليجرام لتنزيل الملف...");
}

// Query course details from real FastAPI backend
const { data: course, isLoading, isError, refetch } = useQuery({
  queryKey: ["courseDetail", courseId],
  queryFn: () => fetchCourseById(courseId.value),
  enabled: computed(() => Boolean(courseId.value)),
  staleTime: 1000 * 60 * 3,
});

// Check if favorite
onMounted(async () => {
  if (courseId.value) {
    try {
      isBookmarked.value = await checkCourseFavorite(courseId.value);
    } catch {
      isBookmarked.value = false;
    }
  }
});

async function toggleFavorite() {
  if (!courseId.value) return;
  try {
    if (isBookmarked.value) {
      await removeCourseFromFavorites(courseId.value);
      isBookmarked.value = false;
      userStore.showToast("تم إزالة المساق من المحفوظات");
    } else {
      await addCourseToFavorites(courseId.value);
      isBookmarked.value = true;
      userStore.showToast("تمت إضافة المساق إلى المحفوظات");
    }
    queryClient.invalidateQueries({ queryKey: ["favorites"] });
  } catch (err) {
    userStore.showToast("تعذر تحديث المحفوظات", "error");
  }
}

async function handleRedeemCode() {
  if (!redeemCodeInput.value.trim()) {
    redeemError.value = "يرجى إدخال رمز الاشتراك أولاً";
    return;
  }
  isRedeeming.value = true;
  redeemError.value = "";
  try {
    const res = await redeemSubscriptionCode(redeemCodeInput.value.trim());
    userStore.showToast(res.message || "تم تفعيل الاشتراك بنجاح!");
    isRedeemModalOpen.value = false;
    redeemCodeInput.value = "";
    queryClient.invalidateQueries({ queryKey: ["courseDetail", courseId] });
    queryClient.invalidateQueries({ queryKey: ["homeDashboardData"] });
    queryClient.invalidateQueries({ queryKey: ["userProfileData"] });
    queryClient.invalidateQueries({ queryKey: ["userSubscriptionInfo"] });
    queryClient.invalidateQueries({ queryKey: ["mySubscriptions"] });
    queryClient.invalidateQueries({ queryKey: ["coursesList"] });
    refetch();
  } catch (err: any) {
    redeemError.value = err.response?.data?.detail || "رمز الاشتراك غير صالح أو منتهي";
  } finally {
    isRedeeming.value = false;
  }
}

function goBack() {
  router.back();
}

function openLesson(lesson: any) {
  if (lesson.is_locked && !course.value?.is_enrolled) {
    if (!userStore.isAuthenticated) {
      userStore.showToast("يرجى تسجيل الدخول عبر تليجرام أولاً", "error");
      return;
    }
    isRedeemModalOpen.value = true;
    return;
  }
  router.push(`/lessons/${lesson.id}`);
}
</script>

<template>
  <div class="min-h-screen bg-white dark:bg-[#0F0F0F] text-[#111111] dark:text-[#FFFFFF] flex flex-col font-sans transition-colors pb-28 text-start">
    <!-- Loading State -->
    <div v-if="isLoading" class="max-w-7xl mx-auto px-4 py-8 w-full flex flex-col gap-6">
      <div class="aspect-video w-full rounded-3xl bg-[#F8F8F8] dark:bg-[#222222] animate-pulse"></div>
      <div class="grid grid-cols-4 gap-3">
        <div v-for="i in 4" :key="i" class="h-16 rounded-2xl bg-[#F8F8F8] dark:bg-[#222222] animate-pulse"></div>
      </div>
    </div>

    <!-- Error State -->
    <div v-else-if="isError || !course" class="max-w-7xl mx-auto px-4 py-12 w-full">
      <ErrorState title="تعذر تحميل تفاصيل المساق" message="المساق غير موجود أو حدث خطأ أثناء الاتصال بالخادم." @retry="refetch()" />
    </div>

    <!-- Loaded Course Content -->
    <template v-else>
      <!-- Top Bar Header with Cover -->
      <div class="relative w-full aspect-video sm:aspect-[21/9] bg-[#0F0F0F] overflow-hidden">
        <img
          :src="course.thumbnail_url || course.coverImageUrl || 'https://images.unsplash.com/photo-1636466497217-26a8cbeaf0aa?auto=format&fit=crop&q=80&w=800'"
          :alt="course.title"
          class="w-full h-full object-cover"
        />
        <div class="absolute inset-0 bg-gradient-to-t from-black/80 via-black/40 to-black/60"></div>

        <!-- Header Action Buttons -->
        <div class="absolute top-4 left-4 right-4 flex items-center justify-between z-10">
          <button
            @click="goBack"
            type="button"
            aria-label="الرجوع"
            class="w-11 h-11 rounded-2xl bg-black/60 backdrop-blur-md text-white flex items-center justify-center active:scale-90 transition-transform cursor-pointer"
          >
            <ChevronRight class="w-6 h-6" />
          </button>

          <div class="flex items-center gap-2">
            <button
              @click="toggleFavorite"
              type="button"
              class="w-11 h-11 rounded-2xl bg-black/60 backdrop-blur-md text-white flex items-center justify-center active:scale-90 transition-transform cursor-pointer"
            >
              <Bookmark :class="['w-5 h-5 transition-colors', isBookmarked ? 'fill-orange-500 text-orange-500' : 'text-white']" />
            </button>
          </div>
        </div>

        <!-- Quick Info Badge on Cover -->
        <div class="absolute bottom-4 right-4 left-4 flex flex-col gap-1 text-white">
          <span class="text-xs font-bold text-orange-400 uppercase tracking-wider">
            {{ course.department ? `${course.department} • ` : '' }}{{ course.stage || course.gradeName || 'المرحلة الدراسية' }}
          </span>
          <h1 class="text-lg sm:text-2xl font-black leading-snug line-clamp-2">
            {{ course.title }}
          </h1>
        </div>
      </div>

      <!-- Main Content Container -->
      <main class="flex-1 max-w-7xl mx-auto px-4 py-5 w-full flex flex-col gap-5">
        <!-- Course Meta Stats Row -->
        <div class="grid grid-cols-3 gap-2 bg-white dark:bg-[#181818] border border-[#E5E7EB] dark:border-[#27272A] rounded-2xl p-3 text-center shadow-xs">
          <div class="flex flex-col items-center gap-0.5">
            <Users class="w-4 h-4 text-orange-500" />
            <span class="text-xs font-black">{{ course.students_count || 0 }}</span>
            <span class="text-[9px] text-[#6B7280] dark:text-[#A1A1AA]">طالب مشترك</span>
          </div>

          <div class="flex flex-col items-center gap-0.5 border-r border-[#E5E7EB] dark:border-[#27272A]">
            <BookOpen class="w-4 h-4 text-orange-500" />
            <span class="text-xs font-black">{{ course.lessons_count || 0 }}</span>
            <span class="text-[9px] text-[#6B7280] dark:text-[#A1A1AA]">محاضرة</span>
          </div>

          <div class="flex flex-col items-center gap-0.5 border-r border-[#E5E7EB] dark:border-[#27272A]">
            <Clock class="w-4 h-4 text-orange-500" />
            <span class="text-xs font-black">{{ course.total_duration_minutes || 0 }}س</span>
            <span class="text-[9px] text-[#6B7280] dark:text-[#A1A1AA]">ساعات الفيديو</span>
          </div>
        </div>

        <!-- Instructor Card -->
        <div class="flex items-center gap-3.5 bg-white dark:bg-[#181818] border border-[#E5E7EB] dark:border-[#27272A] rounded-3xl p-4 shadow-xs">
          <img
            :src="course.teacher_avatar || 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?auto=format&fit=crop&q=80&w=150'"
            :alt="course.teacher_name"
            class="w-12 h-12 rounded-2xl object-cover shrink-0"
          />
          <div class="flex flex-col min-w-0">
            <span class="text-[10px] font-bold text-orange-500 uppercase">أستاذ المساق</span>
            <h3 class="text-sm font-bold text-[#111111] dark:text-[#FFFFFF] truncate">{{ course.teacher_name || 'أستاذ المساق' }}</h3>
            <p class="text-xs text-[#6B7280] dark:text-[#A1A1AA] truncate">{{ course.department || course.departmentName || 'الفرع الأكاديمي' }}</p>
          </div>
        </div>

        <!-- Progress Bar if Enrolled -->
        <div v-if="course.is_enrolled" class="bg-white dark:bg-[#181818] border border-[#E5E7EB] dark:border-[#27272A] rounded-2xl p-4 shadow-xs">
          <div class="flex items-center justify-between text-xs font-bold mb-2">
            <span class="text-[#111111] dark:text-[#FFFFFF]">نسبة إنجازك في المساق</span>
            <span class="text-orange-500">{{ course.completion_percentage }}%</span>
          </div>
          <ProgressBar :percentage="course.completion_percentage" />
        </div>

        <!-- Locked Course Alert Banner if Not Enrolled -->
        <div
          v-else
          class="flex items-center justify-between p-4 rounded-3xl bg-orange-500/10 border border-orange-500/20 text-[#111111] dark:text-orange-200"
        >
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 rounded-2xl bg-orange-500/20 text-orange-500 flex items-center justify-center shrink-0">
              <Lock class="w-5 h-5" />
            </div>
            <div>
              <h4 class="text-xs font-bold">المساق يحتاج إلى اشتراك</h4>
              <p class="text-[11px] text-[#6B7280] dark:text-[#A1A1AA] mt-0.5">
                أدخل كود الاشتراك لتفعيل كافة محاضرات وملحقات المساق فوراً.
              </p>
            </div>
          </div>

          <button
            @click="isRedeemModalOpen = true"
            type="button"
            class="px-3.5 py-2 rounded-xl bg-orange-500 hover:bg-[#EA580C] active:scale-95 text-white font-bold text-xs shadow-xs transition-all shrink-0 cursor-pointer"
          >
            تفعيل الكود
          </button>
        </div>

        <!-- Navigation Tabs -->
        <div class="flex border-b border-[#E5E7EB] dark:border-[#27272A]">
          <button
            @click="activeTab = 'lectures'"
            type="button"
            :class="[
              'flex-1 py-3 text-xs sm:text-sm font-bold text-center border-b-2 transition-all cursor-pointer',
              activeTab === 'lectures'
                ? 'border-orange-500 text-orange-500'
                : 'border-transparent text-[#6B7280] dark:text-[#A1A1AA] hover:text-[#111111] dark:hover:text-[#FFFFFF]'
            ]"
          >
            المحاضرات ({{ course.lessons_count || 0 }})
          </button>

          <button
            @click="activeTab = 'description'"
            type="button"
            :class="[
              'flex-1 py-3 text-xs sm:text-sm font-bold text-center border-b-2 transition-all cursor-pointer',
              activeTab === 'description'
                ? 'border-orange-500 text-orange-500'
                : 'border-transparent text-[#6B7280] dark:text-[#A1A1AA] hover:text-[#111111] dark:hover:text-[#FFFFFF]'
            ]"
          >
            تفاصيل المساق
          </button>

          <button
            @click="activeTab = 'files'"
            type="button"
            :class="[
              'flex-1 py-3 text-xs sm:text-sm font-bold text-center border-b-2 transition-all cursor-pointer',
              activeTab === 'files'
                ? 'border-orange-500 text-orange-500'
                : 'border-transparent text-[#6B7280] dark:text-[#A1A1AA] hover:text-[#111111] dark:hover:text-[#FFFFFF]'
            ]"
          >
            الملفات والملازم
          </button>
        </div>

        <!-- Tab Content: Lectures List -->
        <div v-if="activeTab === 'lectures'" class="flex flex-col gap-2.5">
          <div
            v-for="lec in course.lessons"
            :key="lec.id"
            @click="openLesson(lec)"
            :class="[
              'flex items-center justify-between gap-3 p-3.5 bg-white dark:bg-[#181818] border rounded-2xl transition-all cursor-pointer group active:scale-[0.99] shadow-xs',
              lec.is_locked
                ? 'border-[#E5E7EB] dark:border-[#27272A] opacity-85'
                : 'border-[#E5E7EB] dark:border-[#27272A] hover:border-orange-500/50'
            ]"
          >
            <div class="flex items-center gap-3 min-w-0">
              <div
                :class="[
                  'w-9 h-9 rounded-xl flex items-center justify-center shrink-0',
                  lec.completed
                    ? 'bg-orange-500/10 text-orange-500'
                    : lec.is_locked
                    ? 'bg-[#F8F8F8] dark:bg-[#222222] text-[#9CA3AF] dark:text-[#71717A]'
                    : 'bg-orange-500/10 text-orange-500'
                ]"
              >
                <CheckCircle2 v-if="lec.completed" class="w-4 h-4" />
                <Lock v-else-if="lec.is_locked" class="w-4 h-4 text-[#9CA3AF] dark:text-[#71717A]" />
                <Play v-else class="w-4 h-4 fill-orange-500 ml-0.5" />
              </div>

              <div class="flex flex-col min-w-0">
                <span class="text-xs sm:text-sm font-bold text-[#111111] dark:text-[#FFFFFF] truncate group-hover:text-orange-500 transition-colors">
                  {{ lec.title }}
                </span>
                <span class="text-[10px] text-[#6B7280] dark:text-[#A1A1AA]">
                  {{ lec.duration ? `${Math.round(lec.duration / 60)} دقيقة` : 'محاضرة تفاعلية' }}
                  <template v-if="lec.progress > 0"> • تم إكمال {{ lec.progress }}%</template>
                </span>
              </div>
            </div>

            <div class="flex items-center gap-2 shrink-0">
              <span v-if="lec.is_free || lec.is_free_preview" class="px-2 py-0.5 rounded-full text-[10px] font-bold bg-orange-500/10 text-orange-500">
                معاينة مجانية
              </span>
              <Lock v-else-if="lec.is_locked" class="w-4 h-4 text-[#9CA3AF] dark:text-[#71717A]" />
              <ChevronRight class="w-4 h-4 text-[#9CA3AF] dark:text-[#71717A] group-hover:text-orange-500 transition-colors" />
            </div>
          </div>
        </div>

        <!-- Tab Content: Description -->
        <div v-else-if="activeTab === 'description'" class="bg-white dark:bg-[#181818] border border-[#E5E7EB] dark:border-[#27272A] rounded-3xl p-5 text-xs sm:text-sm text-[#6B7280] dark:text-[#A1A1AA] leading-relaxed shadow-xs">
          {{ course.description || 'كورس شامل لكافة المفاهيم والمحاضرات الأكاديمية مع شروحات وحلول نموذجية لجميع الاختبارات.' }}
        </div>

        <!-- Tab Content: Files -->
        <div v-else class="flex flex-col gap-2.5">
          <div v-if="!course.is_enrolled" class="p-6 bg-orange-500/10 border border-orange-500/20 rounded-2xl text-center">
            <Lock class="w-8 h-8 text-orange-500 mx-auto mb-2" />
            <p class="text-sm font-bold text-[#111111] dark:text-orange-400">يجب الاشتراك في المساق أولاً للوصول إلى الملفات</p>
          </div>
          <template v-else>
            <div v-if="course.attachments && course.attachments.length > 0" class="flex flex-col gap-2.5">
              <div
                v-for="file in course.attachments"
                :key="file.id"
                class="p-4 bg-white dark:bg-[#181818] border border-[#E5E7EB] dark:border-[#27272A] rounded-2xl flex items-center justify-between shadow-xs"
              >
                <div class="flex items-center gap-3">
                  <div class="w-9 h-9 rounded-xl bg-orange-500/10 text-orange-500 flex items-center justify-center shrink-0">
                    <FileText class="w-4 h-4" />
                  </div>
                  <div>
                    <h4 class="text-xs font-bold text-[#111111] dark:text-[#FFFFFF]">{{ file.name || 'ملف المرفق' }}</h4>
                    <span class="text-[10px] text-[#6B7280] dark:text-[#A1A1AA]">{{ (file.type || 'file').toUpperCase() }} • {{ formatFileSize(file.size) }}</span>
                  </div>
                </div>
                <button
                  @click="downloadViaTelegram(file)"
                  type="button"
                  class="px-3 py-1.5 rounded-xl bg-orange-500 hover:bg-[#EA580C] active:scale-95 text-white font-bold text-xs transition-all cursor-pointer flex items-center gap-1"
                >
                  <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="currentColor">
                    <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm4.64 6.8c-.15 1.58-.8 5.42-1.13 7.19-.14.75-.42 1-.68 1.03-.58.05-1.02-.38-1.58-.75-.88-.58-1.38-.94-2.23-1.5-.99-.65-.35-1.01.22-1.59.15-.15 2.71-2.48 2.76-2.69a.2.2 0 00-.05-.18c-.06-.05-.14-.03-.21-.02-.09.02-1.49.95-4.22 2.79-.4.27-.76.41-1.08.4-.36-.01-1.04-.2-1.55-.37-.63-.2-1.12-.31-1.07-.66.02-.18.27-.36.74-.55 2.92-1.27 4.86-2.11 5.83-2.51 2.78-1.16 3.35-1.36 3.73-1.36.08 0 .27.02.39.12.1.08.13.19.14.27-.01.06.01.24 0 .24z"/>
                  </svg>
                  تنزيل
                </button>
              </div>
            </div>
            <div v-else class="p-8 bg-[#F8F8F8] dark:bg-[#181818] border border-[#E5E7EB] dark:border-[#27272A] rounded-2xl text-center">
              <FileText class="w-10 h-10 text-[#9CA3AF] dark:text-[#71717A] mx-auto mb-3" />
              <p class="text-sm text-[#6B7280] dark:text-[#A1A1AA]">لا توجد ملفات مرفقة لهذا المساق بعد</p>
            </div>
          </template>
        </div>
      </main>

      <!-- Bottom Fixed Action Bar -->
      <div class="fixed bottom-0 left-0 right-0 z-40 bg-white/95 dark:bg-[#0F0F0F]/95 backdrop-blur-xl border-t border-[#E5E7EB] dark:border-[#27272A] p-4">
        <div class="max-w-md mx-auto flex items-center justify-between gap-4">
          <div class="flex flex-col">
            <span class="text-[10px] font-bold text-[#6B7280] dark:text-[#A1A1AA]">سعر الاشتراك</span>
            <span class="text-base font-black text-orange-500">
              {{ course.price > 0 ? `${course.price.toLocaleString('ar-IQ')} د.ع` : 'مجاني بالكامل' }}
            </span>
          </div>

          <button
            v-if="course.is_enrolled"
            @click="course.lessons?.[0] ? openLesson(course.lessons[0]) : null"
            type="button"
            class="flex-1 h-12 inline-flex items-center justify-center gap-2 rounded-2xl bg-orange-500 hover:bg-[#EA580C] active:scale-95 text-white font-bold text-sm shadow-sm shadow-orange-500/20 transition-all cursor-pointer"
          >
            <Play class="w-4 h-4 fill-white" />
            <span>متابعة مشاهدة الدروس</span>
          </button>

          <button
            v-else
            @click="isRedeemModalOpen = true"
            type="button"
            class="flex-1 h-12 inline-flex items-center justify-center gap-2 rounded-2xl bg-orange-500 hover:bg-[#EA580C] active:scale-95 text-white font-bold text-sm shadow-sm shadow-orange-500/20 transition-all cursor-pointer"
          >
            <KeyRound class="w-4 h-4" />
            <span>أدخل كود الاشتراك للتفعيل</span>
          </button>
        </div>
      </div>
    </template>

    <!-- Subscription Code Redemption Modal -->
    <div
      v-if="isRedeemModalOpen"
      class="fixed inset-0 z-[9999] bg-black/70 backdrop-blur-md flex items-center justify-center p-4 overflow-y-auto"
    >
      <div class="w-full max-w-md bg-white dark:bg-[#181818] rounded-3xl p-6 flex flex-col gap-4 border border-[#E5E7EB] dark:border-[#27272A] shadow-2xl animate-scale-in">
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-2.5">
            <div class="w-10 h-10 rounded-2xl bg-orange-500/10 text-orange-500 flex items-center justify-center">
              <Sparkles class="w-5 h-5" />
            </div>
            <div>
              <h3 class="text-sm font-bold text-[#111111] dark:text-[#FFFFFF]">تفعيل كود الاشتراك</h3>
              <p class="text-xs text-[#6B7280] dark:text-[#A1A1AA]">أدخل رمز الاشتراك لفتح المساق فوراً</p>
            </div>
          </div>
          <button
            @click="isRedeemModalOpen = false"
            type="button"
            class="text-[#6B7280] hover:text-[#111111] dark:text-[#A1A1AA] dark:hover:text-white text-lg font-bold cursor-pointer"
          >
            ✕
          </button>
        </div>

        <div class="flex flex-col gap-2">
          <input
            v-model="redeemCodeInput"
            type="text"
            placeholder="مثال: WYN-A2B3-C4D5"
            dir="ltr"
            class="w-full h-12 px-4 rounded-2xl bg-[#F8F8F8] dark:bg-[#0F0F0F] border border-[#E5E7EB] dark:border-[#27272A] text-[#111111] dark:text-[#FFFFFF] font-mono text-center uppercase tracking-wider text-sm focus:outline-none focus:border-orange-500 transition-colors"
          />
          <span v-if="redeemError" class="text-xs font-bold text-red-500 text-center">{{ redeemError }}</span>
        </div>

        <div class="flex items-center gap-3 mt-2">
          <button
            @click="isRedeemModalOpen = false"
            type="button"
            class="flex-1 h-11 rounded-2xl border border-[#E5E7EB] dark:border-[#27272A] font-bold text-xs text-[#6B7280] dark:text-[#A1A1AA] hover:bg-[#F8F8F8] dark:hover:bg-[#222222] transition-colors cursor-pointer"
          >
            إلغاء
          </button>

          <button
            @click="handleRedeemCode"
            :disabled="isRedeeming"
            type="button"
            class="flex-1 h-11 rounded-2xl bg-orange-500 hover:bg-[#EA580C] active:scale-95 disabled:opacity-50 text-white font-bold text-xs shadow-sm transition-all cursor-pointer"
          >
            {{ isRedeeming ? 'جاري التحقق...' : 'تفعيل الاشتراك' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
@keyframes scaleIn {
  from {
    opacity: 0;
    transform: scale(0.95);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}
.animate-scale-in {
  animation: scaleIn 0.2s cubic-bezier(0.16, 1, 0.3, 1) forwards;
}
</style>
