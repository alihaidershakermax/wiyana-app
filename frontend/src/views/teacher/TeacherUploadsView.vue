<script setup lang="ts">
import { ref, computed, onMounted, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import {
  Upload,
  FileText,
  CheckCircle2,
  AlertCircle,
  Loader2,
  BookOpen,
  Trash2,
  Edit,
  Play,
  X,
  Paperclip,
  Film,
  UserCheck,
  UserX,
  GraduationCap,
  Search,
  Eye,
  BarChart3,
  Clock,
  Users,
  Lock,
  Unlock,
  Video,
  ChevronDown,
  Plus,
} from "@lucide/vue";
import { apiClient } from "@/api/client";
import { adminApiClient, getAccessToken as getAdminToken } from "@/admin/api/adminApi";
import { useAdminStore } from "@/admin/stores/useAdminStore";

function teacherClient() {
  return getAdminToken() ? adminApiClient : apiClient;
}

interface CourseItem {
  id: string;
  title: string;
  description?: string;
  stage?: string;
  department?: string;
  price?: number;
  teacher_id?: string;
  teacher_name?: string;
  teacher_avatar?: string;
  teacher_username?: string;
  lessons_count?: number;
  students_count?: number;
}

interface LectureItem {
  id: string;
  title: string;
  description?: string;
  course_id: string;
  course_title?: string;
  section_id?: string;
  order: number;
  duration?: number;
  duration_minutes?: number;
  is_free?: boolean;
  is_free_preview?: boolean;
  status: string;
  youtube_video_id?: string;
  youtube_url?: string;
  video_url?: string;
  thumbnail?: string;
  has_pdf?: boolean;
  attachments?: any[];
  created_at?: string;
  views_count?: number;
  created_by?: string;
  uploader?: { id: string; name: string; username?: string; avatar?: string } | null;
}

const adminStore = useAdminStore();
const route = useRoute();
const router = useRouter();

// ─── State ────────────────────────────────────────────────────────────────────
const courses = ref<CourseItem[]>([]);
const selectedCourseId = ref<string>("");
const lectures = ref<LectureItem[]>([]);
const isLoading = ref(true);
const lecturesLoading = ref(false);

// Role check
const isOwnerOrAdmin = computed(() => {
  const role = String(adminStore.userRole || "").toUpperCase();
  return role === "OWNER" || role === "ADMIN" || role === "MASTER_ADMIN";
});

// Toast
const toastMessage = ref<{ text: string; type: "success" | "error" | "info" } | null>(null);
function showToast(text: string, type: "success" | "error" | "info" = "success") {
  toastMessage.value = { text, type };
  setTimeout(() => { toastMessage.value = null; }, 4000);
}

// ─── Upload Form State (inline, not modal) ────────────────────────────────────
const uploadTitle = ref("");
const uploadDescription = ref("");
const uploadYoutubeUrl = ref("");
const uploadIsFreePreview = ref(false);
const uploadVideoFile = ref<File | null>(null);
const uploadPdfFile = ref<File | null>(null);
const isSubmittingUpload = ref(false);
const uploadError = ref("");
const uploadSuccess = ref(false);
const uploadProgress = ref(0);

// ─── Edit Modal ───────────────────────────────────────────────────────────────
const isEditModalOpen = ref(false);
const editingLecture = ref<LectureItem | null>(null);
const editTitle = ref("");
const editDescription = ref("");
const editYoutubeUrl = ref("");
const editIsFreePreview = ref(false);
const isSubmittingEdit = ref(false);

// ─── Delete ───────────────────────────────────────────────────────────────────
const lectureToDelete = ref<LectureItem | null>(null);
const deletingLectureLoading = ref(false);

// ─── Lectures search ──────────────────────────────────────────────────────────
const lectureSearchQuery = ref("");
const filteredLectures = computed(() => {
  const q = lectureSearchQuery.value.trim().toLowerCase();
  if (!q) return lectures.value;
  return lectures.value.filter(
    (l) =>
      (l.title || "").toLowerCase().includes(q) ||
      (l.course_title || "").toLowerCase().includes(q)
  );
});

// ─── Preview ──────────────────────────────────────────────────────────────────
const previewVideoUrl = ref<string | null>(null);
const previewVideoTitle = ref<string>("");

// ─── Stats Modal ──────────────────────────────────────────────────────────────
const statsLecture = ref<LectureItem | null>(null);
const isStatsModalOpen = ref(false);

// ─── Assign Teacher (Admin only) ──────────────────────────────────────────────
const isAssignTeacherModalOpen = ref(false);
const teacherSearchQuery = ref("");
const availableTeachers = ref<any[]>([]);
const loadingTeachers = ref(false);
const selectedTeacherForAssign = ref<any | null>(null);
const manualTeacherIdInput = ref("");
const assigningTeacherLoading = ref(false);
const removingTeacherLoading = ref(false);

// ─── Load Courses ─────────────────────────────────────────────────────────────
async function loadCourses() {
  isLoading.value = true;
  try {
    let coursesList: CourseItem[] = [];
    const res1 = await teacherClient().get<{ courses: CourseItem[] }>("/teacher/courses").catch(() => null);
    if (res1?.data?.courses && Array.isArray(res1.data.courses)) {
      coursesList = res1.data.courses;
    }
    if (coursesList.length === 0) {
      const res2 = await teacherClient().get<any>("/courses").catch(() => null);
      if (res2?.data) {
        coursesList = Array.isArray(res2.data) ? res2.data : (res2.data.items || res2.data.courses || []);
      }
    }
    courses.value = coursesList;
    const qCid = String(route.query.course_id || "");
    if (qCid && courses.value.some((c) => c.id === qCid)) {
      selectedCourseId.value = qCid;
    } else if (courses.value.length > 0 && (!selectedCourseId.value || !courses.value.some((c) => c.id === selectedCourseId.value))) {
      selectedCourseId.value = courses.value[0].id;
    }
    await loadSelectedCourseLectures();
  } catch {
    courses.value = [];
  } finally {
    isLoading.value = false;
  }
}

async function loadSelectedCourseLectures() {
  if (!selectedCourseId.value) { lectures.value = []; return; }
  lecturesLoading.value = true;
  try {
    let rawLessons: any[] = [];
    if (isOwnerOrAdmin.value) {
      const allRes = await teacherClient().get<any[]>("/teacher/lessons").catch(() => null);
      if (allRes?.data && Array.isArray(allRes.data)) {
        rawLessons = allRes.data;
      }
    } else {
      const lecRes = await teacherClient().get<any[]>(`/courses/${selectedCourseId.value}/lessons`).catch(() => ({ data: [] }));
      rawLessons = Array.isArray(lecRes.data) ? lecRes.data : [];
    }
    lectures.value = rawLessons.map((l: any, idx: number) => {
      const atts = l.attachments || [];
      const hasPdf = atts.some((a: any) => a.type === "pdf" || String(a.url || a || "").toLowerCase().endsWith(".pdf"));
      return {
        id: l.id || l._id,
        title: l.title || `محاضرة #${idx + 1}`,
        description: l.description || "",
        course_id: l.course_id || selectedCourseId.value,
        course_title: l.course_title || "",
        section_id: l.section_id,
        order: l.order || idx + 1,
        duration: l.duration || 0,
        duration_minutes: l.duration ? Math.round(l.duration / 60) : 0,
        is_free: Boolean(l.is_free || l.is_free_preview),
        is_free_preview: Boolean(l.is_free || l.is_free_preview),
        status: l.status || "published",
        youtube_video_id: l.youtube_video_id,
        youtube_url: l.youtube_url,
        video_url: l.video_url || l.youtube_url,
        thumbnail: l.thumbnail,
        has_pdf: hasPdf,
        attachments: atts,
        created_at: l.created_at,
        views_count: l.views_count || 0,
        created_by: l.created_by,
        uploader: l.uploader || null,
      };
    }).sort((a, b) => {
      const tA = new Date(a.created_at || 0).getTime();
      const tB = new Date(b.created_at || 0).getTime();
      return tB - tA || (b.order || 0) - (a.order || 0);
    });
  } catch {
    lectures.value = [];
  } finally {
    lecturesLoading.value = false;
  }
}

watch(selectedCourseId, () => { loadSelectedCourseLectures(); });

onMounted(async () => {
  if (!adminStore.user) await adminStore.initializeAuth().catch(() => {});
  await loadCourses();
});

const currentCourse = computed(() => courses.value.find((c) => c.id === selectedCourseId.value));

// ─── Upload ────────────────────────────────────────────────────────────────────
function onVideoFilePicked(e: Event) {
  const t = e.target as HTMLInputElement;
  if (t.files?.[0]) uploadVideoFile.value = t.files[0];
}
function onPdfFilePicked(e: Event) {
  const t = e.target as HTMLInputElement;
  if (t.files?.[0]) uploadPdfFile.value = t.files[0];
}
function resetUpload() {
  uploadTitle.value = "";
  uploadDescription.value = "";
  uploadYoutubeUrl.value = "";
  uploadIsFreePreview.value = false;
  uploadVideoFile.value = null;
  uploadPdfFile.value = null;
  uploadError.value = "";
  uploadSuccess.value = false;
  uploadProgress.value = 0;
}

async function submitLectureUpload() {
  if (!selectedCourseId.value) { uploadError.value = "يرجى اختيار المساق أولاً"; return; }
  if (!uploadTitle.value.trim()) { uploadError.value = "يرجى إدخال عنوان المحاضرة"; return; }
  if (!uploadYoutubeUrl.value.trim() && !uploadVideoFile.value) {
    uploadError.value = "يرجى إدخال رابط يوتيوب أو اختيار ملف فيديو";
    return;
  }

  isSubmittingUpload.value = true;
  uploadError.value = "";
  uploadProgress.value = 10;

  try {
    const formData = new FormData();
    formData.append("course_id", selectedCourseId.value);
    formData.append("title", uploadTitle.value.trim());
    if (uploadDescription.value.trim()) formData.append("description", uploadDescription.value.trim());
    formData.append("is_free_preview", String(uploadIsFreePreview.value));
    formData.append("is_free", String(uploadIsFreePreview.value));
    formData.append("status", "published");
    if (uploadYoutubeUrl.value.trim()) formData.append("youtube_url", uploadYoutubeUrl.value.trim());
    if (uploadVideoFile.value) formData.append("video", uploadVideoFile.value);
    if (uploadPdfFile.value) formData.append("pdf", uploadPdfFile.value);

    uploadProgress.value = 10;
    await teacherClient().post("/teacher/upload", formData, {
      headers: { "Content-Type": "multipart/form-data" },
      timeout: 600000,
      onUploadProgress: (progressEvent) => {
        if (progressEvent.total) {
          uploadProgress.value = Math.min(95, Math.round((progressEvent.loaded * 95) / progressEvent.total));
        }
      },
    });
    uploadProgress.value = 100;
    uploadSuccess.value = true;
    showToast("✅ تم رفع ونشر المحاضرة في المساق بنجاح!");
    setTimeout(() => { resetUpload(); loadSelectedCourseLectures(); }, 1500);
  } catch (err: any) {
    uploadError.value = err?.response?.data?.detail || "فشل رفع المحاضرة. يرجى التحقق والمحاولة ثانية.";
  } finally {
    isSubmittingUpload.value = false;
  }
}

// ─── Edit ──────────────────────────────────────────────────────────────────────
function openEditLecture(lecture: LectureItem) {
  if (isOwnerOrAdmin.value) {
    router.push({ name: "admin.lesson-edit", params: { id: lecture.id } });
    return;
  }
  editingLecture.value = lecture;
  editTitle.value = lecture.title;
  editDescription.value = lecture.description || "";
  editYoutubeUrl.value = lecture.youtube_url || "";
  editIsFreePreview.value = Boolean(lecture.is_free_preview || lecture.is_free);
  isEditModalOpen.value = true;
}
async function submitLectureEdit() {
  if (!editingLecture.value || !editTitle.value.trim()) return;
  isSubmittingEdit.value = true;
  try {
    await teacherClient().patch(`/teacher/lessons/${editingLecture.value.id}`, {
      title: editTitle.value.trim(),
      description: editDescription.value.trim() || undefined,
      youtube_url: editYoutubeUrl.value.trim() || undefined,
      is_free_preview: editIsFreePreview.value,
      is_free: editIsFreePreview.value,
    });
    showToast("تم تحديث المحاضرة بنجاح");
    isEditModalOpen.value = false;
    editingLecture.value = null;
    await loadSelectedCourseLectures();
  } catch (err: any) {
    showToast(err?.response?.data?.detail || "فشل تحديث المحاضرة", "error");
  } finally {
    isSubmittingEdit.value = false;
  }
}

// ─── Delete ────────────────────────────────────────────────────────────────────
async function handleConfirmDeleteLecture() {
  if (!lectureToDelete.value) return;
  deletingLectureLoading.value = true;
  try {
    await teacherClient().delete(`/teacher/lessons/${lectureToDelete.value.id}`);
    showToast("تم حذف المحاضرة بنجاح");
    lectureToDelete.value = null;
    await loadSelectedCourseLectures();
  } catch (err: any) {
    showToast(err?.response?.data?.detail || "فشل حذف المحاضرة", "error");
  } finally {
    deletingLectureLoading.value = false;
  }
}

// ─── Preview ───────────────────────────────────────────────────────────────────
function playVideo(lec: LectureItem) {
  let vid = lec.youtube_video_id;
  if (!vid && lec.youtube_url) {
    const match = lec.youtube_url.match(/(?:v=|\/embed\/|youtu\.be\/|\/v\/|\/shorts\/)([0-9A-Za-z_-]{11})/);
    if (match) vid = match[1];
  }
  if (vid) {
    previewVideoUrl.value = `https://www.youtube-nocookie.com/embed/${vid}?autoplay=1&controls=1&modestbranding=1&rel=0&showinfo=0&iv_load_policy=3&fs=1`;
    previewVideoTitle.value = lec.title;
  } else if (lec.video_url) {
    previewVideoUrl.value = lec.video_url;
    previewVideoTitle.value = lec.title;
  } else {
    showToast("فيديو المحاضرة غير متوفر للمعاينة المباشرة", "info");
  }
}

// ─── Stats ─────────────────────────────────────────────────────────────────────
function openStats(lec: LectureItem) {
  statsLecture.value = lec;
  isStatsModalOpen.value = true;
}

// ─── Assign Teacher (Owner/Admin only) ─────────────────────────────────────────
async function openAssignTeacherModal() {
  isAssignTeacherModalOpen.value = true;
  selectedTeacherForAssign.value = null;
  manualTeacherIdInput.value = "";
  teacherSearchQuery.value = "";
  loadingTeachers.value = true;
  try {
    const res = await teacherClient().get<any>("/users/", { params: { limit: 100 } });
    availableTeachers.value = Array.isArray(res.data) ? res.data : (res.data?.items || []);
  } catch { availableTeachers.value = []; }
  finally { loadingTeachers.value = false; }
}

const filteredTeachersList = computed(() => {
  const q = teacherSearchQuery.value.trim().toLowerCase();
  if (!q) return availableTeachers.value.slice(0, 30);
  return availableTeachers.value.filter((u: any) => {
    const name = `${u.first_name || ""} ${u.last_name || ""} ${u.username || ""} ${u.name || ""}`.toLowerCase();
    const idStr = String(u.id || u._id || u.telegram_id || "");
    return name.includes(q) || idStr.includes(q);
  }).slice(0, 30);
});

function selectTeacher(user: any) {
  selectedTeacherForAssign.value = user;
  manualTeacherIdInput.value = user.id || user._id || String(user.telegram_id || "");
}

async function handleConfirmAssignTeacher() {
  if (!selectedCourseId.value) return;
  const tId = manualTeacherIdInput.value.trim() || selectedTeacherForAssign.value?.id || selectedTeacherForAssign.value?.telegram_id;
  if (!tId) { showToast("يرجى اختيار أستاذ من القائمة", "error"); return; }
  assigningTeacherLoading.value = true;
  try {
    const res = await teacherClient().post(`/courses/${selectedCourseId.value}/teacher`, { teacher_id: String(tId) });
    const teacherName = res.data?.teacher_name || [selectedTeacherForAssign.value?.first_name, selectedTeacherForAssign.value?.last_name].filter(Boolean).join(" ") || selectedTeacherForAssign.value?.username || "أستاذ المادة";
    const idx = courses.value.findIndex((c) => c.id === selectedCourseId.value);
    if (idx !== -1) {
      courses.value[idx].teacher_id = String(tId);
      courses.value[idx].teacher_name = teacherName;
    }
    showToast("تم تعيين الأستاذ للمساق بنجاح!");
    isAssignTeacherModalOpen.value = false;
    await loadCourses();
  } catch (err: any) {
    showToast(err?.response?.data?.detail || "فشل تعيين الأستاذ", "error");
  } finally {
    assigningTeacherLoading.value = false;
  }
}

async function handleRemoveTeacher() {
  if (!selectedCourseId.value) return;
  if (!confirm("هل أنت متأكد من إزالة الأستاذ المعيّن من هذا المساق؟")) return;
  removingTeacherLoading.value = true;
  try {
    await teacherClient().delete(`/courses/${selectedCourseId.value}/teacher`);
    const idx = courses.value.findIndex((c) => c.id === selectedCourseId.value);
    if (idx !== -1) { courses.value[idx].teacher_id = undefined; courses.value[idx].teacher_name = undefined; }
    showToast("تمت إزالة الأستاذ من المساق بنجاح");
    await loadCourses();
  } catch (err: any) {
    showToast(err?.response?.data?.detail || "فشل إزالة الأستاذ", "error");
  } finally {
    removingTeacherLoading.value = false;
  }
}

function formatDate(d?: string) {
  if (!d) return "—";
  return new Date(d).toLocaleDateString("ar-IQ", { year: "numeric", month: "short", day: "numeric" });
}
function formatDuration(mins?: number) {
  if (!mins) return "—";
  if (mins < 60) return `${mins} د`;
  const h = Math.floor(mins / 60);
  const m = mins % 60;
  return m > 0 ? `${h}س ${m}د` : `${h}س`;
}
</script>

<template>
  <div class="space-y-0 max-w-6xl mx-auto text-start pb-20">

    <!-- ── Toast ─────────────────────────────────────────────────────────────── -->
    <transition name="toast">
      <div
        v-if="toastMessage"
        :class="[
          'fixed top-5 left-1/2 -translate-x-1/2 z-[9999] px-5 py-3 rounded-2xl shadow-2xl font-bold text-xs sm:text-sm flex items-center gap-2 border',
          toastMessage.type === 'success'
            ? 'bg-emerald-600 text-white border-emerald-500'
            : toastMessage.type === 'error'
            ? 'bg-red-600 text-white border-red-500'
            : 'bg-orange-500 text-white border-orange-400'
        ]"
      >
        <CheckCircle2 v-if="toastMessage.type === 'success'" class="w-4 h-4" />
        <AlertCircle v-else class="w-4 h-4" />
        <span>{{ toastMessage.text }}</span>
      </div>
    </transition>

    <!-- ── Loading ────────────────────────────────────────────────────────────── -->
    <div v-if="isLoading" class="flex items-center justify-center py-24">
      <div class="flex flex-col items-center gap-3">
        <Loader2 class="w-10 h-10 text-orange-500 animate-spin" />
        <p class="text-sm font-bold text-[#6B7280] dark:text-[#A1A1AA]">جاري تحميل بيانات المساقات...</p>
      </div>
    </div>

    <template v-else>

      <!-- ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
           SECTION 1 ▸ UPLOAD NEW VIDEO (hidden for owner/admin)
      ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ -->
      <div v-if="!isOwnerOrAdmin" class="bg-white dark:bg-[#111111] border border-[#E5E7EB] dark:border-[#27272A] rounded-3xl overflow-hidden">

        <!-- Header -->
        <div class="bg-gradient-to-l from-orange-500/10 via-orange-500/5 to-transparent border-b border-[#E5E7EB] dark:border-[#27272A] px-6 py-4 flex items-center justify-between">
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 rounded-2xl bg-orange-500 flex items-center justify-center shadow-md shadow-orange-500/30">
              <Upload class="w-5 h-5 text-white" />
            </div>
            <div>
              <h2 class="text-base font-black text-[#111111] dark:text-white">رفع محاضرة جديدة</h2>
              <p class="text-[11px] text-[#6B7280] dark:text-[#A1A1AA] font-medium">أدخل عنوان الفيديو ورابط يوتيوب أو ملف مع مرفق ثم انشر</p>
            </div>
          </div>
        </div>

        <div class="p-6 space-y-5">

          <!-- Course Selector (preselected card or dropdown) -->
          <div>
            <label class="block text-xs font-black text-[#111111] dark:text-white mb-2">المساق المستهدف</label>
            <!-- No courses -->
            <div v-if="courses.length === 0" class="p-4 rounded-2xl bg-amber-500/10 border border-amber-500/20 flex items-center gap-3">
              <AlertCircle class="w-5 h-5 text-amber-500 shrink-0" />
              <p class="text-xs font-bold text-amber-700 dark:text-amber-400">لا توجد مساقات متاحة حالياً. تواصل مع مالك المنصة.</p>
            </div>
            <!-- Single course card -->
            <div v-else-if="courses.length === 1 || !isOwnerOrAdmin" class="p-4 rounded-2xl bg-orange-500/10 border border-orange-500/30 flex items-center justify-between gap-3">
              <div class="flex items-center gap-3 min-w-0">
                <div class="w-9 h-9 rounded-xl bg-orange-500 flex items-center justify-center shrink-0">
                  <BookOpen class="w-4 h-4 text-white" />
                </div>
                <div class="min-w-0">
                  <p class="text-sm font-black text-[#111111] dark:text-white truncate">{{ currentCourse?.title || courses[0]?.title }}</p>
                  <p class="text-[11px] text-[#6B7280] dark:text-[#A1A1AA]">{{ currentCourse?.department || courses[0]?.department || 'القسم العام' }} • {{ currentCourse?.stage || courses[0]?.stage || 'المرحلة الدراسية' }}</p>
                </div>
              </div>
              <span class="text-[10px] px-2.5 py-1 rounded-full bg-orange-500 text-white font-extrabold shrink-0 whitespace-nowrap">محدد مسبقاً</span>
            </div>
            <!-- Multi-course dropdown -->
            <div v-else class="relative">
              <select
                v-model="selectedCourseId"
                class="w-full appearance-none px-4 py-3 rounded-2xl border border-[#E5E7EB] dark:border-[#27272A] bg-[#F8F8F8] dark:bg-[#0F0F0F] text-[#111111] dark:text-white font-bold text-sm focus:outline-none focus:ring-2 focus:ring-orange-500 pr-10 cursor-pointer"
              >
                <option v-for="c in courses" :key="c.id" :value="c.id">{{ c.title }} — {{ c.stage || 'عام' }}</option>
              </select>
              <ChevronDown class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-[#9CA3AF] pointer-events-none" />
            </div>
          </div>

          <!-- Two column layout for form fields -->
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <!-- Video Title -->
            <div class="md:col-span-2">
              <label class="block text-xs font-black text-[#111111] dark:text-white mb-1.5">عنوان الفيديو / المحاضرة *</label>
              <input
                v-model="uploadTitle"
                type="text"
                placeholder="مثال: المحاضرة 01 — شرح الأساسيات والمفاهيم"
                class="w-full px-4 py-3 rounded-2xl border border-[#E5E7EB] dark:border-[#27272A] bg-[#F8F8F8] dark:bg-[#0F0F0F] text-[#111111] dark:text-white text-sm focus:outline-none focus:ring-2 focus:ring-orange-500 font-medium"
              />
            </div>

            <!-- Description -->
            <div class="md:col-span-2">
              <label class="block text-xs font-black text-[#111111] dark:text-white mb-1.5">وصف المحاضرة (اختياري)</label>
              <textarea
                v-model="uploadDescription"
                rows="2"
                placeholder="اكتب نبذة أو ملاحظات للطلاب..."
                class="w-full px-4 py-3 rounded-2xl border border-[#E5E7EB] dark:border-[#27272A] bg-[#F8F8F8] dark:bg-[#0F0F0F] text-[#111111] dark:text-white text-sm focus:outline-none focus:ring-2 focus:ring-orange-500 font-medium resize-none"
              />
            </div>
          </div>

          <!-- Video Source -->
          <div class="p-4 rounded-2xl border border-[#E5E7EB] dark:border-[#27272A] bg-[#F9F9F9] dark:bg-[#0D0D0D] space-y-3">
            <p class="text-xs font-black text-[#111111] dark:text-white flex items-center gap-1.5">
              <Film class="w-4 h-4 text-orange-500" />
              <span>مصدر الفيديو *</span>
            </p>
            <!-- YouTube URL -->
            <div>
              <label class="text-[11px] font-bold text-[#6B7280] dark:text-[#A1A1AA] block mb-1">رابط يوتيوب</label>
              <input
                v-model="uploadYoutubeUrl"
                type="url"
                placeholder="https://www.youtube.com/watch?v=..."
                class="w-full px-3.5 py-2.5 rounded-xl border border-[#E5E7EB] dark:border-[#27272A] bg-white dark:bg-[#0F0F0F] text-[#111111] dark:text-white text-xs font-medium focus:outline-none focus:ring-2 focus:ring-orange-500"
              />
            </div>
            <!-- Divider -->
            <div class="flex items-center gap-2">
              <div class="h-px flex-1 bg-[#E5E7EB] dark:bg-[#27272A]" />
              <span class="text-[10px] font-bold text-[#9CA3AF]">أو رفع ملف مباشر</span>
              <div class="h-px flex-1 bg-[#E5E7EB] dark:bg-[#27272A]" />
            </div>
            <!-- Video file upload -->
            <label class="flex items-center gap-3 border border-dashed border-[#E5E7EB] dark:border-[#27272A] rounded-xl p-3 cursor-pointer hover:border-orange-400 transition-colors bg-white dark:bg-[#0F0F0F]">
              <input type="file" accept="video/*" class="hidden" @change="onVideoFilePicked" />
              <div class="w-8 h-8 rounded-lg bg-orange-500/10 flex items-center justify-center shrink-0">
                <Video class="w-4 h-4 text-orange-500" />
              </div>
              <div class="min-w-0">
                <p class="text-xs font-bold text-[#111111] dark:text-white truncate">
                  {{ uploadVideoFile ? uploadVideoFile.name : 'اضغط لاختيار ملف فيديو' }}
                </p>
                <p class="text-[10px] text-[#9CA3AF]">MP4, MKV, MOV</p>
              </div>
            </label>
          </div>

          <!-- PDF Attachment -->
          <div class="p-4 rounded-2xl border border-[#E5E7EB] dark:border-[#27272A] bg-[#F9F9F9] dark:bg-[#0D0D0D]">
            <p class="text-xs font-black text-[#111111] dark:text-white flex items-center gap-1.5 mb-2.5">
              <Paperclip class="w-4 h-4 text-orange-500" />
              <span>مرفق PDF / ملزمة (اختياري)</span>
            </p>
            <label class="flex items-center gap-3 border border-dashed border-[#E5E7EB] dark:border-[#27272A] rounded-xl p-3 cursor-pointer hover:border-orange-400 transition-colors bg-white dark:bg-[#0F0F0F]">
              <input type="file" accept=".pdf,.doc,.docx" class="hidden" @change="onPdfFilePicked" />
              <div class="w-8 h-8 rounded-lg bg-blue-500/10 flex items-center justify-center shrink-0">
                <FileText class="w-4 h-4 text-blue-500" />
              </div>
              <div class="min-w-0">
                <p class="text-xs font-bold text-[#111111] dark:text-white truncate">
                  {{ uploadPdfFile ? uploadPdfFile.name : 'اضغط لاختيار ملف PDF أو Word' }}
                </p>
                <p class="text-[10px] text-[#9CA3AF]">PDF, DOC, DOCX</p>
              </div>
            </label>
          </div>

          <!-- Free preview toggle -->
          <label class="flex items-center gap-3 cursor-pointer p-3 rounded-2xl border border-[#E5E7EB] dark:border-[#27272A] hover:bg-[#F9F9F9] dark:hover:bg-[#0D0D0D] transition-colors">
            <div
              @click="uploadIsFreePreview = !uploadIsFreePreview"
              :class="[
                'relative w-11 h-6 rounded-full transition-all duration-200 shrink-0',
                uploadIsFreePreview ? 'bg-orange-500' : 'bg-[#E5E7EB] dark:bg-[#27272A]'
              ]"
            >
              <div :class="['absolute top-1 w-4 h-4 bg-white rounded-full shadow transition-all duration-200', uploadIsFreePreview ? 'left-6' : 'left-1']" />
            </div>
            <div>
              <p class="text-xs font-black text-[#111111] dark:text-white">إتاحة كمعاينة مجانية</p>
              <p class="text-[10px] text-[#6B7280] dark:text-[#A1A1AA]">تسمح لجميع الطلاب بمشاهدة هذه المحاضرة مجاناً</p>
            </div>
          </label>

          <!-- Progress bar during upload -->
          <div v-if="isSubmittingUpload" class="space-y-1">
            <div class="flex items-center justify-between text-[11px] font-bold text-[#6B7280]">
              <span>جاري الرفع والمعالجة...</span>
              <span>{{ uploadProgress }}%</span>
            </div>
            <div class="h-2 rounded-full bg-[#F3F4F6] dark:bg-[#27272A] overflow-hidden">
              <div
                class="h-full rounded-full bg-gradient-to-r from-orange-400 to-orange-600 transition-all duration-500"
                :style="{ width: `${uploadProgress}%` }"
              />
            </div>
          </div>

          <!-- Error / Success -->
          <p v-if="uploadError" class="p-3 rounded-xl bg-red-500/10 border border-red-500/20 text-red-600 dark:text-red-400 text-xs font-bold flex items-center gap-2">
            <AlertCircle class="w-4 h-4 shrink-0" />
            {{ uploadError }}
          </p>
          <p v-if="uploadSuccess" class="p-3 rounded-xl bg-emerald-500/10 border border-emerald-500/20 text-emerald-600 dark:text-emerald-400 text-xs font-bold flex items-center gap-2">
            <CheckCircle2 class="w-4 h-4 shrink-0" />
            تم رفع ونشر المحاضرة بنجاح!
          </p>

          <!-- Action Buttons -->
          <div class="flex items-center justify-between gap-3 pt-1">
            <button
              type="button"
              @click="resetUpload"
              class="px-5 py-2.5 rounded-2xl border border-[#E5E7EB] dark:border-[#27272A] text-xs font-bold text-[#6B7280] dark:text-[#A1A1AA] hover:bg-[#F8F8F8] dark:hover:bg-[#1A1A1A] transition-colors cursor-pointer"
            >
              مسح الحقول
            </button>
            <button
              type="button"
              @click="submitLectureUpload"
              :disabled="isSubmittingUpload || courses.length === 0"
              class="flex-1 px-6 py-3 rounded-2xl bg-orange-500 hover:bg-orange-600 text-white font-black text-sm flex items-center justify-center gap-2 shadow-md shadow-orange-500/30 disabled:opacity-50 disabled:cursor-not-allowed active:scale-95 transition-all cursor-pointer"
            >
              <Loader2 v-if="isSubmittingUpload" class="w-4 h-4 animate-spin" />
              <Upload v-else class="w-4 h-4" />
              <span>{{ isSubmittingUpload ? 'جاري الرفع والنشر...' : 'نشر المحاضرة الآن' }}</span>
            </button>
          </div>
        </div>
      </div>

      <!-- ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
           SECTION 2 ▸ LECTURES MANAGEMENT
      ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ -->
      <div :class="['bg-white dark:bg-[#111111] border border-[#E5E7EB] dark:border-[#27272A] rounded-3xl overflow-hidden', !isOwnerOrAdmin ? 'mt-5' : '']">

        <!-- Header -->
        <div class="border-b border-[#E5E7EB] dark:border-[#27272A] px-6 py-4 flex flex-col sm:flex-row sm:items-center gap-3 justify-between">
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 rounded-2xl bg-[#111111] dark:bg-white flex items-center justify-center shadow">
              <GraduationCap class="w-5 h-5 text-white dark:text-[#111111]" />
            </div>
            <div>
              <h2 class="text-base font-black text-[#111111] dark:text-white">قسم المحاضرات</h2>
              <p class="text-[11px] text-[#6B7280] dark:text-[#A1A1AA] font-medium">
                <template v-if="isOwnerOrAdmin">
                  كل المحاضرات في المنصة
                  <span class="mx-1">•</span>
                  <span>{{ filteredLectures.length }} محاضرة</span>
                </template>
                <template v-else>
                  {{ currentCourse?.title || 'المساق المحدد' }}
                  <span class="mx-1">•</span>
                  <span>{{ filteredLectures.length }} محاضرة</span>
                </template>
              </p>
            </div>
          </div>

          <!-- Course tabs if multiple -->
          <div v-if="courses.length > 1 && !isOwnerOrAdmin" class="flex items-center gap-2 overflow-x-auto pb-1 scrollbar-none">
            <button
              v-for="c in courses"
              :key="c.id"
              @click="selectedCourseId = c.id"
              :class="[
                'px-3 py-1.5 rounded-xl text-xs font-bold whitespace-nowrap transition-all cursor-pointer',
                selectedCourseId === c.id
                  ? 'bg-orange-500 text-white shadow shadow-orange-500/20'
                  : 'bg-[#F8F8F8] dark:bg-[#1A1A1A] text-[#6B7280] dark:text-[#A1A1AA] hover:bg-[#F0F0F0] dark:hover:bg-[#222]'
              ]"
            >{{ c.title }}</button>
          </div>

          <!-- Owner: Assign teacher -->
          <div v-if="isOwnerOrAdmin" class="flex items-center gap-2">
            <div v-if="currentCourse?.teacher_name" class="flex items-center gap-2 px-3 py-1.5 rounded-xl bg-emerald-500/10 border border-emerald-500/20">
              <UserCheck class="w-3.5 h-3.5 text-emerald-500" />
              <span class="text-xs font-bold text-emerald-700 dark:text-emerald-400">{{ currentCourse.teacher_name }}</span>
              <button
                @click="handleRemoveTeacher"
                :disabled="removingTeacherLoading"
                class="text-red-400 hover:text-red-600 transition-colors cursor-pointer"
              >
                <UserX class="w-3.5 h-3.5" />
              </button>
            </div>
            <button
              @click="openAssignTeacherModal"
              class="px-3 py-1.5 rounded-xl border border-[#E5E7EB] dark:border-[#27272A] text-xs font-bold text-[#6B7280] dark:text-[#A1A1AA] hover:border-orange-400 hover:text-orange-500 transition-all cursor-pointer flex items-center gap-1"
            >
              <UserCheck class="w-3.5 h-3.5" />
              <span>{{ currentCourse?.teacher_name ? 'تغيير' : 'تعيين أستاذ' }}</span>
            </button>
          </div>
        </div>

        <!-- Lectures search -->
        <div v-if="lectures.length > 0" class="border-b border-[#E5E7EB] dark:border-[#27272A] px-6 py-3">
          <div class="relative">
            <Search class="absolute right-3 top-1/2 -translate-y-1/2 w-4 h-4 text-[#9CA3AF]" />
            <input
              v-model="lectureSearchQuery"
              type="text"
              placeholder="ابحث عن محاضرة بالعنوان..."
              class="w-full pr-10 px-4 py-2.5 rounded-xl border border-[#E5E7EB] dark:border-[#27272A] bg-[#F8F8F8] dark:bg-[#0F0F0F] text-[#111111] dark:text-white text-xs focus:outline-none focus:ring-2 focus:ring-orange-500"
            />
          </div>
        </div>

        <!-- Stats bar -->
        <div v-if="filteredLectures.length > 0" class="grid grid-cols-3 border-b border-[#E5E7EB] dark:border-[#27272A]">
          <div class="px-6 py-3 text-center border-l border-[#E5E7EB] dark:border-[#27272A]">
            <p class="text-xl font-black text-orange-500">{{ filteredLectures.length }}</p>
            <p class="text-[10px] font-bold text-[#6B7280] dark:text-[#A1A1AA]">محاضرة</p>
          </div>
          <div class="px-6 py-3 text-center border-l border-[#E5E7EB] dark:border-[#27272A]">
            <p class="text-xl font-black text-emerald-500">{{ filteredLectures.filter(l => l.status === 'published').length }}</p>
            <p class="text-[10px] font-bold text-[#6B7280] dark:text-[#A1A1AA]">منشورة</p>
          </div>
          <div class="px-6 py-3 text-center">
            <p class="text-xl font-black text-blue-500">{{ filteredLectures.filter(l => l.is_free).length }}</p>
            <p class="text-[10px] font-bold text-[#6B7280] dark:text-[#A1A1AA]">مجانية</p>
          </div>
        </div>

        <!-- Lectures loading -->
        <div v-if="lecturesLoading" class="flex items-center justify-center py-16">
          <Loader2 class="w-8 h-8 text-orange-500 animate-spin" />
        </div>

        <!-- Empty state -->
        <div v-else-if="filteredLectures.length === 0" class="flex flex-col items-center justify-center py-16 gap-3">
          <div class="w-16 h-16 rounded-3xl bg-[#F8F8F8] dark:bg-[#1A1A1A] flex items-center justify-center">
            <Video class="w-8 h-8 text-[#9CA3AF]" />
          </div>
          <p class="text-sm font-black text-[#6B7280] dark:text-[#A1A1AA]">{{ lectureSearchQuery ? 'لا توجد نتائج مطابقة للبحث' : 'لا توجد محاضرات بعد' }}</p>
          <p class="text-xs text-[#9CA3AF]">{{ lectureSearchQuery ? 'جرب البحث بكلمات أخرى' : 'ارفع محاضرتك الأولى باستخدام النموذج أعلاه' }}</p>
        </div>

        <!-- Lectures list -->
        <div v-else class="divide-y divide-[#F3F4F6] dark:divide-[#1A1A1A]">
          <div
            v-for="(lec, idx) in filteredLectures"
            :key="lec.id"
            class="flex items-center gap-4 px-6 py-4 hover:bg-[#FAFAFA] dark:hover:bg-[#0D0D0D] transition-colors group"
          >
            <!-- Index + Preview -->
            <div class="relative shrink-0">
              <div
                class="w-14 h-10 rounded-xl bg-gradient-to-br from-[#1A1A1A] to-[#2A2A2A] dark:from-[#27272A] dark:to-[#111] flex items-center justify-center cursor-pointer overflow-hidden"
                @click="playVideo(lec)"
              >
                <img v-if="lec.thumbnail" :src="lec.thumbnail" class="w-full h-full object-cover" />
                <Play v-else class="w-4 h-4 text-white/60" />
                <div class="absolute inset-0 bg-black/0 group-hover:bg-black/30 transition-all flex items-center justify-center">
                  <Play class="w-4 h-4 text-white opacity-0 group-hover:opacity-100 transition-opacity drop-shadow" />
                </div>
              </div>
              <div class="absolute -top-1.5 -right-1.5 w-5 h-5 rounded-full bg-[#F8F8F8] dark:bg-[#1A1A1A] border border-[#E5E7EB] dark:border-[#27272A] text-[10px] font-black text-[#6B7280] dark:text-[#A1A1AA] flex items-center justify-center">
                {{ idx + 1 }}
              </div>
            </div>

            <!-- Info -->
            <div class="flex-1 min-w-0">
              <p class="text-sm font-black text-[#111111] dark:text-white truncate">{{ lec.title }}</p>
              <div class="flex items-center gap-3 mt-0.5 flex-wrap">
                <span v-if="isOwnerOrAdmin && lec.course_title" class="text-[10px] font-bold px-2 py-0.5 rounded-full bg-purple-500/10 text-purple-600 dark:text-purple-400">
                  <BookOpen class="w-2.5 h-2.5 inline-block" /> {{ lec.course_title }}
                </span>
                <span v-if="isOwnerOrAdmin && lec.uploader" class="text-[10px] text-[#9CA3AF] flex items-center gap-1">
                  <UserCheck class="w-2.5 h-2.5" />{{ lec.uploader.name }}
                </span>
                <span class="text-[10px] text-[#9CA3AF]">{{ formatDate(lec.created_at) }}</span>
                <span v-if="lec.duration_minutes" class="text-[10px] text-[#9CA3AF] flex items-center gap-0.5">
                  <Clock class="w-2.5 h-2.5" />{{ formatDuration(lec.duration_minutes) }}
                </span>
                <span
                  :class="[
                    'text-[10px] px-2 py-0.5 rounded-full font-bold',
                    lec.is_free ? 'bg-emerald-500/10 text-emerald-600 dark:text-emerald-400' : 'bg-[#F3F4F6] dark:bg-[#27272A] text-[#6B7280] dark:text-[#A1A1AA]'
                  ]"
                >
                  {{ lec.is_free ? 'مجانية' : 'مدفوعة' }}
                </span>
                <span v-if="lec.has_pdf" class="text-[10px] px-2 py-0.5 rounded-full font-bold bg-blue-500/10 text-blue-600 dark:text-blue-400 flex items-center gap-0.5">
                  <Paperclip class="w-2.5 h-2.5" />PDF
                </span>
                <span
                  :class="[
                    'text-[10px] px-2 py-0.5 rounded-full font-bold',
                    lec.status === 'published' ? 'bg-orange-500/10 text-orange-600 dark:text-orange-400' : 'bg-gray-500/10 text-gray-500'
                  ]"
                >
                  {{ lec.status === 'published' ? 'منشور' : 'مسودة' }}
                </span>
              </div>
            </div>

            <!-- Actions -->
            <div :class="['flex items-center gap-1 shrink-0 transition-opacity', isOwnerOrAdmin ? 'opacity-100' : 'opacity-0 group-hover:opacity-100']">
              <button
                @click="playVideo(lec)"
                class="w-8 h-8 rounded-xl flex items-center justify-center text-[#6B7280] dark:text-[#A1A1AA] hover:bg-orange-500/10 hover:text-orange-500 transition-colors cursor-pointer"
                title="معاينة الفيديو"
              >
                <Play class="w-4 h-4" />
              </button>
              <button
                @click="openStats(lec)"
                class="w-8 h-8 rounded-xl flex items-center justify-center text-[#6B7280] dark:text-[#A1A1AA] hover:bg-blue-500/10 hover:text-blue-500 transition-colors cursor-pointer"
                title="الإحصائيات"
              >
                <BarChart3 class="w-4 h-4" />
              </button>
              <button
                @click="openEditLecture(lec)"
                class="w-8 h-8 rounded-xl flex items-center justify-center text-[#6B7280] dark:text-[#A1A1AA] hover:bg-emerald-500/10 hover:text-emerald-500 transition-colors cursor-pointer"
                title="تعديل"
              >
                <Edit class="w-4 h-4" />
              </button>
              <button
                @click="lectureToDelete = lec"
                class="w-8 h-8 rounded-xl flex items-center justify-center text-[#6B7280] dark:text-[#A1A1AA] hover:bg-red-500/10 hover:text-red-500 transition-colors cursor-pointer"
                title="حذف"
              >
                <Trash2 class="w-4 h-4" />
              </button>
            </div>
          </div>
        </div>
      </div>

    </template>

    <!-- ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
         MODALS
    ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ -->

    <!-- Video Preview Modal -->
    <div v-if="previewVideoUrl" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/90 backdrop-blur-md" @click.self="previewVideoUrl = null">
      <div class="bg-black w-full max-w-3xl rounded-3xl overflow-hidden shadow-2xl">
        <div class="flex items-center justify-between p-4 border-b border-white/10">
          <p class="text-sm font-black text-white truncate">{{ previewVideoTitle }}</p>
          <button @click="previewVideoUrl = null" class="p-1.5 rounded-xl bg-white/10 hover:bg-white/20 text-white transition-colors cursor-pointer"><X class="w-4 h-4" /></button>
        </div>
        <iframe
          :src="previewVideoUrl"
          class="w-full aspect-video"
          allowfullscreen
          referrerpolicy="strict-origin"
          allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
        />
      </div>
    </div>

    <!-- Edit Lecture Modal -->
    <div v-if="isEditModalOpen" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/75 backdrop-blur-md overflow-y-auto">
      <div class="bg-white dark:bg-[#181818] border border-[#E5E7EB] dark:border-[#27272A] w-full max-w-lg rounded-3xl p-6 shadow-2xl space-y-4 my-8">
        <div class="flex items-center justify-between border-b border-[#E5E7EB] dark:border-[#27272A] pb-3">
          <h3 class="text-base font-black text-[#111111] dark:text-white flex items-center gap-2">
            <Edit class="w-5 h-5 text-orange-500" />تعديل المحاضرة
          </h3>
          <button @click="isEditModalOpen = false" class="p-2 rounded-xl text-[#6B7280] hover:text-[#111111] dark:hover:text-white hover:bg-[#F8F8F8] dark:hover:bg-[#222222] transition-colors cursor-pointer"><X class="w-5 h-5" /></button>
        </div>
        <div class="space-y-3 text-xs">
          <div>
            <label class="block mb-1 font-black text-[#111111] dark:text-white">عنوان المحاضرة *</label>
            <input v-model="editTitle" type="text" class="w-full px-4 py-2.5 rounded-xl border border-[#E5E7EB] dark:border-[#27272A] bg-[#F8F8F8] dark:bg-[#0F0F0F] text-[#111111] dark:text-white focus:outline-none focus:ring-2 focus:ring-orange-500" />
          </div>
          <div>
            <label class="block mb-1 font-bold text-[#6B7280] dark:text-[#A1A1AA]">وصف المحاضرة (اختياري)</label>
            <textarea v-model="editDescription" rows="2" class="w-full px-4 py-2.5 rounded-xl border border-[#E5E7EB] dark:border-[#27272A] bg-[#F8F8F8] dark:bg-[#0F0F0F] text-[#111111] dark:text-white focus:outline-none focus:ring-2 focus:ring-orange-500 resize-none" />
          </div>
          <div>
            <label class="block mb-1 font-bold text-[#6B7280] dark:text-[#A1A1AA]">رابط يوتيوب (اختياري)</label>
            <input v-model="editYoutubeUrl" type="url" class="w-full px-4 py-2.5 rounded-xl border border-[#E5E7EB] dark:border-[#27272A] bg-[#F8F8F8] dark:bg-[#0F0F0F] text-[#111111] dark:text-white focus:outline-none focus:ring-2 focus:ring-orange-500" />
          </div>
          <label class="flex items-center gap-2.5 cursor-pointer">
            <input v-model="editIsFreePreview" type="checkbox" class="w-4 h-4 rounded text-orange-500 focus:ring-orange-400" />
            <span class="font-bold text-[#111111] dark:text-white">إتاحة كمعاينة مجانية</span>
          </label>
        </div>
        <div class="flex items-center justify-end gap-3 pt-2 border-t border-[#E5E7EB] dark:border-[#27272A]">
          <button @click="isEditModalOpen = false" class="px-5 py-2.5 rounded-xl border border-[#E5E7EB] dark:border-[#27272A] font-bold text-xs text-[#6B7280] dark:text-[#A1A1AA] hover:bg-[#F8F8F8] dark:hover:bg-[#222] cursor-pointer">إلغاء</button>
          <button @click="submitLectureEdit" :disabled="isSubmittingEdit" class="px-6 py-2.5 rounded-xl bg-orange-500 hover:bg-orange-600 text-white font-black text-xs flex items-center gap-2 disabled:opacity-50 cursor-pointer">
            <Loader2 v-if="isSubmittingEdit" class="w-4 h-4 animate-spin" />
            <span>{{ isSubmittingEdit ? 'جاري الحفظ...' : 'حفظ التعديلات' }}</span>
          </button>
        </div>
      </div>
    </div>

    <!-- Delete Confirm Modal -->
    <div v-if="lectureToDelete" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/75 backdrop-blur-md">
      <div class="bg-white dark:bg-[#181818] border border-[#E5E7EB] dark:border-[#27272A] w-full max-w-sm rounded-3xl p-6 shadow-2xl space-y-4">
        <div class="flex items-start gap-3">
          <div class="w-10 h-10 rounded-2xl bg-red-500/10 flex items-center justify-center shrink-0">
            <Trash2 class="w-5 h-5 text-red-500" />
          </div>
          <div>
            <h3 class="text-sm font-black text-[#111111] dark:text-white">حذف المحاضرة</h3>
            <p class="text-xs text-[#6B7280] dark:text-[#A1A1AA] mt-1">هل أنت متأكد من حذف "<strong>{{ lectureToDelete.title }}</strong>"؟ لا يمكن التراجع عن هذا الإجراء.</p>
          </div>
        </div>
        <div class="flex items-center justify-end gap-3">
          <button @click="lectureToDelete = null" class="px-5 py-2.5 rounded-xl border border-[#E5E7EB] dark:border-[#27272A] font-bold text-xs text-[#6B7280] dark:text-[#A1A1AA] hover:bg-[#F8F8F8] dark:hover:bg-[#222] cursor-pointer">إلغاء</button>
          <button @click="handleConfirmDeleteLecture" :disabled="deletingLectureLoading" class="px-6 py-2.5 rounded-xl bg-red-500 hover:bg-red-600 text-white font-black text-xs flex items-center gap-2 disabled:opacity-50 cursor-pointer">
            <Loader2 v-if="deletingLectureLoading" class="w-4 h-4 animate-spin" />
            <span>{{ deletingLectureLoading ? 'جاري الحذف...' : 'نعم، حذف' }}</span>
          </button>
        </div>
      </div>
    </div>

    <!-- Stats Modal -->
    <div v-if="isStatsModalOpen && statsLecture" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/75 backdrop-blur-md">
      <div class="bg-white dark:bg-[#181818] border border-[#E5E7EB] dark:border-[#27272A] w-full max-w-sm rounded-3xl p-6 shadow-2xl space-y-4">
        <div class="flex items-center justify-between">
          <h3 class="text-sm font-black text-[#111111] dark:text-white flex items-center gap-2">
            <BarChart3 class="w-5 h-5 text-blue-500" />
            إحصائيات المحاضرة
          </h3>
          <button @click="isStatsModalOpen = false; statsLecture = null" class="p-2 rounded-xl text-[#6B7280] hover:text-[#111111] dark:hover:text-white hover:bg-[#F8F8F8] dark:hover:bg-[#222] transition-colors cursor-pointer"><X class="w-4 h-4" /></button>
        </div>
        <p class="text-sm font-black text-[#111111] dark:text-white truncate border-b border-[#E5E7EB] dark:border-[#27272A] pb-3">{{ statsLecture.title }}</p>
        <div class="grid grid-cols-2 gap-3">
          <div class="p-3.5 rounded-2xl bg-blue-500/10 border border-blue-500/20 text-center">
            <p class="text-2xl font-black text-blue-600 dark:text-blue-400">{{ statsLecture.views_count || 0 }}</p>
            <p class="text-[11px] font-bold text-[#6B7280] dark:text-[#A1A1AA] mt-0.5">مشاهدة</p>
          </div>
          <div class="p-3.5 rounded-2xl bg-emerald-500/10 border border-emerald-500/20 text-center">
            <p class="text-2xl font-black text-emerald-600 dark:text-emerald-400">{{ statsLecture.duration_minutes || 0 }}</p>
            <p class="text-[11px] font-bold text-[#6B7280] dark:text-[#A1A1AA] mt-0.5">دقيقة مدة</p>
          </div>
          <div class="p-3.5 rounded-2xl bg-orange-500/10 border border-orange-500/20 text-center col-span-2">
            <p class="text-sm font-bold text-[#111111] dark:text-white">{{ formatDate(statsLecture.created_at) }}</p>
            <p class="text-[11px] font-bold text-[#6B7280] dark:text-[#A1A1AA] mt-0.5">تاريخ النشر</p>
          </div>
        </div>
        <div class="flex items-center gap-2">
          <span :class="['flex-1 py-2 rounded-xl text-center text-xs font-bold', statsLecture.status === 'published' ? 'bg-orange-500/10 text-orange-600 dark:text-orange-400' : 'bg-gray-500/10 text-gray-500']">
            {{ statsLecture.status === 'published' ? '✓ منشورة' : '• مسودة' }}
          </span>
          <span :class="['flex-1 py-2 rounded-xl text-center text-xs font-bold', statsLecture.is_free ? 'bg-emerald-500/10 text-emerald-600 dark:text-emerald-400' : 'bg-red-500/10 text-red-500']">
            {{ statsLecture.is_free ? '🔓 مجانية' : '🔒 مدفوعة' }}
          </span>
        </div>
      </div>
    </div>

    <!-- Assign Teacher Modal (Owner/Admin only) -->
    <div v-if="isAssignTeacherModalOpen && isOwnerOrAdmin" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/75 backdrop-blur-md overflow-y-auto">
      <div class="bg-white dark:bg-[#181818] border border-[#E5E7EB] dark:border-[#27272A] w-full max-w-lg rounded-3xl p-6 shadow-2xl space-y-4 my-8">
        <div class="flex items-center justify-between border-b border-[#E5E7EB] dark:border-[#27272A] pb-3">
          <h3 class="text-base font-black text-[#111111] dark:text-white flex items-center gap-2">
            <UserCheck class="w-5 h-5 text-orange-500" />تعيين أستاذ للمساق
          </h3>
          <button @click="isAssignTeacherModalOpen = false" class="p-2 rounded-xl text-[#6B7280] hover:text-[#111111] dark:hover:text-white hover:bg-[#F8F8F8] dark:hover:bg-[#222] transition-colors cursor-pointer"><X class="w-5 h-5" /></button>
        </div>
        <p class="text-xs font-bold text-[#6B7280] dark:text-[#A1A1AA]">
          المساق: <span class="text-[#111111] dark:text-white">{{ currentCourse?.title }}</span>
        </p>
        <div class="relative">
          <Search class="absolute right-3 top-1/2 -translate-y-1/2 w-4 h-4 text-[#9CA3AF]" />
          <input v-model="teacherSearchQuery" type="text" placeholder="ابحث عن أستاذ بالاسم أو المعرف..." class="w-full pr-10 px-4 py-2.5 rounded-xl border border-[#E5E7EB] dark:border-[#27272A] bg-[#F8F8F8] dark:bg-[#0F0F0F] text-[#111111] dark:text-white text-xs focus:outline-none focus:ring-2 focus:ring-orange-500" />
        </div>
        <div v-if="loadingTeachers" class="flex items-center justify-center py-8"><Loader2 class="w-6 h-6 text-orange-500 animate-spin" /></div>
        <div v-else class="space-y-1.5 max-h-60 overflow-y-auto scrollbar-thin scrollbar-thumb-[#E5E7EB] dark:scrollbar-thumb-[#27272A]">
          <button
            v-for="user in filteredTeachersList"
            :key="user.id || user._id"
            @click="selectTeacher(user)"
            :class="[
              'w-full flex items-center gap-3 px-3 py-2.5 rounded-xl border text-right transition-all cursor-pointer',
              selectedTeacherForAssign?.id === user.id
                ? 'border-orange-500 bg-orange-500/10'
                : 'border-[#E5E7EB] dark:border-[#27272A] hover:border-orange-300 bg-white dark:bg-[#111]'
            ]"
          >
            <div class="w-8 h-8 rounded-full bg-gradient-to-br from-orange-400 to-orange-600 flex items-center justify-center text-white text-xs font-black shrink-0">
              {{ (user.first_name || user.username || "?")[0].toUpperCase() }}
            </div>
            <div class="min-w-0 flex-1 text-right">
              <p class="text-xs font-black text-[#111111] dark:text-white truncate">{{ user.first_name || "" }} {{ user.last_name || "" }}</p>
              <p class="text-[10px] text-[#9CA3AF] truncate">@{{ user.username || user.telegram_id || user.id }}</p>
            </div>
            <CheckCircle2 v-if="selectedTeacherForAssign?.id === user.id" class="w-4 h-4 text-orange-500 shrink-0" />
          </button>
        </div>
        <div v-if="selectedTeacherForAssign" class="p-3 rounded-xl bg-orange-500/10 border border-orange-500/20 text-xs font-bold text-orange-700 dark:text-orange-400">
          تم اختيار: {{ selectedTeacherForAssign.first_name || "" }} {{ selectedTeacherForAssign.last_name || "" }}
        </div>
        <div class="flex items-center justify-end gap-3 pt-2 border-t border-[#E5E7EB] dark:border-[#27272A]">
          <button @click="isAssignTeacherModalOpen = false" class="px-5 py-2.5 rounded-xl border border-[#E5E7EB] dark:border-[#27272A] font-bold text-xs text-[#6B7280] cursor-pointer">إلغاء</button>
          <button @click="handleConfirmAssignTeacher" :disabled="assigningTeacherLoading || !selectedTeacherForAssign" class="px-6 py-2.5 rounded-xl bg-orange-500 hover:bg-orange-600 text-white font-black text-xs flex items-center gap-2 disabled:opacity-50 cursor-pointer">
            <Loader2 v-if="assigningTeacherLoading" class="w-4 h-4 animate-spin" />
            <span>{{ assigningTeacherLoading ? 'جاري التعيين...' : 'تعيين الأستاذ' }}</span>
          </button>
        </div>
      </div>
    </div>

  </div>
</template>

<style scoped>
.toast-enter-active, .toast-leave-active { transition: all 0.3s ease; }
.toast-enter-from, .toast-leave-to { opacity: 0; transform: translate(-50%, -20px); }
</style>
