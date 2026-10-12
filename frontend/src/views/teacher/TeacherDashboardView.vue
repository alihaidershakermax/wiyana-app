<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import {
  Users,
  BookOpen,
  MessageSquareQuote,
  Upload,
  Sparkles,
  Play,
} from "@lucide/vue";
import { apiClient } from "@/api/client";
import { useRouter } from "vue-router";

const router = useRouter();

interface DashboardOverview {
  my_courses: number;
  my_students: number;
  my_lessons: number;
  total_views: number;
  completion_rate: number;
}

interface TeacherCourse {
  id: string;
  title: string;
  status: string;
  students_count: number;
  lessons_count: number;
  created_at?: string;
}

interface QuestionItem {
  id: string;
  user_id?: string;
  lesson_id?: string;
  question_text?: string;
  status?: string;
  created_at?: string;
}

const isLoading = ref(true);

const teacherName = ref("");
const overview = ref<DashboardOverview>({
  my_courses: 0,
  my_students: 0,
  my_lessons: 0,
  total_views: 0,
  completion_rate: 0,
});
const courses = ref<TeacherCourse[]>([]);
const openQuestions = ref<QuestionItem[]>([]);

const stats = computed(() => [
  {
    label: "الطلاب المشتركين في كورساتي",
    value: String(overview.value.my_students || 0),
    icon: Users,
    change: "عدد الطلاب النشطين حالياً",
  },
  {
    label: "إجمالي المحاضرات المرفوعة",
    value: String(overview.value.my_lessons || 0),
    icon: BookOpen,
    change: `في ${overview.value.my_courses || 0} مساق`,
  },
  {
    label: "أسئلة الطلاب غير المجابة",
    value: String(openQuestions.value.length || 0),
    icon: MessageSquareQuote,
    change: "تتطلب الإجابة",
  },
]);

const activeCourse = computed(() => courses.value[0] ?? null);

onMounted(async () => {
  try {
    const [dashRes, coursesRes, meRes, questionsRes] = await Promise.all([
      apiClient.get<DashboardOverview>("/teacher/dashboard"),
      apiClient.get<{ courses: TeacherCourse[] }>("/teacher/courses"),
      apiClient.get<any>("/users/me"),
      apiClient.get<QuestionItem[]>("/questions/teacher"),
    ]);

    overview.value = { ...overview.value, ...(dashRes.data || {}) };
    courses.value = coursesRes.data?.courses || [];
    const u = meRes.data || {};
    teacherName.value = u.name || `${u.first_name || ""} ${u.last_name || ""}`.trim() || u.username || "";
    const questions = Array.isArray(questionsRes.data) ? questionsRes.data : [];
    openQuestions.value = questions.filter((q) => q.status !== "answered");
  } catch {
    // keep zero defaults
  } finally {
    isLoading.value = false;
  }
});
</script>

<template>
  <div class="space-y-6 text-start">
    <!-- Loading -->
    <div v-if="isLoading" class="space-y-6">
      <div class="h-32 bg-[#F8F8F8] dark:bg-[#222222] rounded-3xl animate-pulse" />
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 sm:gap-6">
        <div v-for="i in 4" :key="i" class="h-28 bg-[#F8F8F8] dark:bg-[#222222] rounded-3xl animate-pulse" />
      </div>
    </div>

    <template v-else>
      <!-- Welcome Card -->
      <div class="bg-orange-500 rounded-3xl p-6 sm:p-8 text-white shadow-sm flex flex-col md:flex-row items-start md:items-center justify-between gap-6">
        <div>
          <div class="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-white/20 text-xs font-bold mb-3 border border-white/20">
            <Sparkles class="w-3.5 h-3.5" />
            <span>بوابة الأستاذ والمحاضر المعتمد</span>
          </div>
          <h2 class="text-2xl sm:text-3xl font-black tracking-tight mb-2">
            مرحباً بك {{ teacherName || 'أستاذنا الكريم' }}
          </h2>
          <p class="text-xs sm:text-sm text-orange-100 max-w-xl leading-relaxed">
            يمكنك من هنا متابعة أداء مساقاتك، رفع محاضرات فيديو وملازم جديدة، والإجابة على استفسارات الطلاب.
          </p>
        </div>

        <button
          type="button"
          @click="router.push('/teacher/uploads')"
          class="px-5 py-3 rounded-2xl bg-white text-orange-600 hover:bg-orange-50 font-black text-xs sm:text-sm flex items-center justify-center gap-2 shadow-sm transition-all active:scale-95 cursor-pointer shrink-0"
        >
          <Upload class="w-4 h-4" />
          <span>رفع محاضرة جديدة</span>
        </button>
      </div>

      <!-- Stats Grid -->
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4 sm:gap-6">
        <div
          v-for="item in stats"
          :key="item.label"
          class="bg-white dark:bg-[#181818] border border-[#E5E7EB] dark:border-[#27272A] rounded-3xl p-5 shadow-xs flex flex-col justify-between"
        >
          <div class="flex items-start justify-between mb-3">
            <span class="text-xs font-bold text-[#6B7280] dark:text-[#A1A1AA]">{{ item.label }}</span>
            <div
              class="w-10 h-10 rounded-2xl flex items-center justify-center bg-orange-500/10 text-orange-600 dark:text-orange-400"
            >
              <component :is="item.icon" class="w-5 h-5" />
            </div>
          </div>

          <div>
            <div class="text-2xl font-black text-[#111111] dark:text-[#FFFFFF]">
              {{ item.value }}
            </div>
            <p class="text-[11px] text-[#6B7280] dark:text-[#A1A1AA] mt-1 font-semibold">
              {{ item.change }}
            </p>
          </div>
        </div>
      </div>

      <!-- Active Course Spotlight -->
      <div v-if="activeCourse" class="bg-white dark:bg-[#181818] border border-[#E5E7EB] dark:border-[#27272A] rounded-3xl p-6 shadow-xs space-y-4">
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-3">
            <div class="w-12 h-12 rounded-2xl bg-orange-500/10 text-orange-500 flex items-center justify-center">
              <Play class="w-6 h-6" />
            </div>
            <div>
              <span class="text-xs font-bold text-orange-500">المقرر الأكثر تفاعلاً</span>
              <h3 class="text-base font-black text-[#111111] dark:text-[#FFFFFF]">{{ activeCourse.title }}</h3>
            </div>
          </div>

          <button
            type="button"
            @click="router.push('/teacher/uploads')"
            class="text-xs font-bold text-orange-500 hover:underline cursor-pointer"
          >
            إدارة المقرر
          </button>
        </div>

        <div class="grid grid-cols-2 sm:grid-cols-4 gap-3 pt-2 border-t border-[#E5E7EB] dark:border-[#27272A]">
          <div class="p-3 bg-[#F8F8F8] dark:bg-[#222222] rounded-2xl text-center">
            <div class="text-xs text-[#6B7280] dark:text-[#A1A1AA]">الطلاب</div>
            <div class="text-sm font-black text-[#111111] dark:text-[#FFFFFF] mt-0.5">{{ activeCourse.students_count }}</div>
          </div>
          <div class="p-3 bg-[#F8F8F8] dark:bg-[#222222] rounded-2xl text-center">
            <div class="text-xs text-[#6B7280] dark:text-[#A1A1AA]">المحاضرات</div>
            <div class="text-sm font-black text-[#111111] dark:text-[#FFFFFF] mt-0.5">{{ activeCourse.lessons_count }}</div>
          </div>
          <div class="p-3 bg-[#F8F8F8] dark:bg-[#222222] rounded-2xl text-center">
            <div class="text-xs text-[#6B7280] dark:text-[#A1A1AA]">الحالة</div>
            <div class="text-sm font-black text-orange-500 mt-0.5">منشور</div>
          </div>
          <div class="p-3 bg-[#F8F8F8] dark:bg-[#222222] rounded-2xl text-center">
            <div class="text-xs text-[#6B7280] dark:text-[#A1A1AA]">إضافة محاضرة</div>
            <button
              type="button"
              @click="router.push('/teacher/uploads')"
              class="text-xs font-bold text-orange-500 hover:underline mt-0.5 block cursor-pointer"
            >
              + إضافة
            </button>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>
