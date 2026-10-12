<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useCourse, usePublishCourse, useArchiveCourse, useAssignTeacher, useRemoveTeacher } from '@/admin/queries/useCourses'
import { useSections, useCreateSection, useUpdateSection, useDeleteSection, useReorderSections } from '@/admin/queries/useSections'
import { teacherManagementApi } from '@/admin/api/teacherManagementApi'
import { useAdminStore } from '@/admin/stores/useAdminStore'
import { useToast } from '@/admin/composables/useToast'
import ConfirmDialog from '@/admin/components/ui/ConfirmDialog.vue'
import ToastContainer from '@/admin/components/ui/ToastContainer.vue'
import {
  ArrowRight, Edit, BookOpen, Users, FileText, CheckCircle, Archive,
  Plus, Trash2, GripVertical, Loader2, UserCheck, UserX, Globe
} from '@lucide/vue'
import { COURSE_STATUS_LABELS, COURSE_STATUS_COLORS } from '@/admin/types/course'
import type { Section } from '@/admin/types/section'

const router = useRouter()
const route = useRoute()
const toast = useToast()
const adminStore = useAdminStore()

const courseId = ref(route.params.id as string)
const { data: course, isLoading, isError, refetch } = useCourse(courseId)
const { data: sections, isLoading: sectionsLoading } = useSections(courseId)

const publishCourse = usePublishCourse()
const archiveCourse = useArchiveCourse()
const assignTeacher = useAssignTeacher()
const removeTeacher = useRemoveTeacher()
const createSection = useCreateSection()
const updateSection = useUpdateSection()
const deleteSection = useDeleteSection()
const reorderSections = useReorderSections()

// Section management
const showNewSectionForm = ref(false)
const newSectionTitle = ref('')
const addingSectionLoading = ref(false)
const editingSection = ref<Section | null>(null)
const editingTitle = ref('')
const confirmDeleteSection = ref<Section | null>(null)
const isDeletingSection = ref(false)
const processingAction = ref<string | null>(null)

// Teacher assignment form (Platform Owner only)
const showAssignForm = ref(false)
const teacherIdInput = ref('')
const assigningTeacher = ref(false)

// Teacher search from DB
const teacherSearch = ref('')
const allUsers = ref<any[]>([])
const usersLoading = ref(false)
const selectedUser = ref<any | null>(null)

async function loadUsers() {
  if (usersLoading.value) return
  usersLoading.value = true
  try {
    // Load both teachers and all students (backend auto-promotes)
    const { adminApiClient } = await import('@/admin/api/adminApi')
    const { data } = await adminApiClient.get('/users/', { params: { limit: 200 } })
    allUsers.value = Array.isArray(data) ? data : (data?.items ?? [])
  } catch {
    // fallback: try teachers endpoint
    try {
      const teachers = await teacherManagementApi.getTeachers({ limit: 200 })
      allUsers.value = teachers
    } catch { allUsers.value = [] }
  } finally {
    usersLoading.value = false
  }
}

watch(showAssignForm, (val) => { if (val && adminStore.isOwner) loadUsers() })

const filteredUsers = computed(() => {
  const q = teacherSearch.value.trim().toLowerCase()
  if (!q) return allUsers.value.slice(0, 50)
  return allUsers.value.filter(u => {
    const name = `${u.first_name ?? ''} ${u.last_name ?? ''} ${u.username ?? ''}`.toLowerCase()
    const tgId = String(u.telegram_id ?? u.user_id ?? '')
    return name.includes(q) || tgId.includes(q)
  }).slice(0, 50)
})

function selectUser(user: any) {
  selectedUser.value = user
  const uid = user.id || user._id || String(user.telegram_id || user.user_id || '')
  teacherIdInput.value = uid
}

async function handlePublish() {
  processingAction.value = 'publish'
  try {
    await publishCourse.mutateAsync(courseId.value)
    toast.success('تم النشر', 'تم نشر المساق بنجاح')
  } catch (err: any) {
    toast.error('خطأ', err?.response?.data?.detail ?? 'فشل نشر المساق')
  } finally {
    processingAction.value = null
  }
}

async function handleArchive() {
  processingAction.value = 'archive'
  try {
    await archiveCourse.mutateAsync(courseId.value)
    toast.success('تم الأرشفة', 'تم أرشفة المساق')
  } catch (err: any) {
    toast.error('خطأ', err?.response?.data?.detail ?? 'فشل الأرشفة')
  } finally {
    processingAction.value = null
  }
}

async function handleAssignTeacher() {
  const teacherId = selectedUser.value?.id || selectedUser.value?._id || selectedUser.value?.telegram_id || teacherIdInput.value.trim()
  if (!teacherId) {
    toast.error('خطأ', 'يرجى اختيار أستاذ من القائمة أو إدخال المعرف')
    return
  }
  assigningTeacher.value = true
  try {
    await assignTeacher.mutateAsync({
      id: courseId.value,
      payload: { teacher_id: String(teacherId) },
    })
    toast.success('تم التعيين', 'تم تعيين الأستاذ بنجاح')
    showAssignForm.value = false
    teacherIdInput.value = ''
    selectedUser.value = null
    refetch()
  } catch (err: any) {
    const detail = err?.response?.data?.detail || err?.message || 'فشل تعيين الأستاذ'
    toast.error('خطأ', detail)
  } finally {
    assigningTeacher.value = false
  }
}

async function handleRemoveTeacher() {
  processingAction.value = 'remove-teacher'
  try {
    await removeTeacher.mutateAsync(courseId.value)
    toast.success('تم الإزالة', 'تم إزالة الأستاذ من المساق بنجاح')
    refetch()
  } catch (err: any) {
    toast.error('خطأ', err?.response?.data?.detail ?? 'فشل إزالة الأستاذ')
  } finally {
    processingAction.value = null
  }
}

async function handleAddSection() {
  if (!newSectionTitle.value.trim()) return
  addingSectionLoading.value = true
  try {
    await createSection.mutateAsync({
      courseId: courseId.value,
      payload: { title: newSectionTitle.value.trim() },
    })
    toast.success('تم الإضافة', 'تم إضافة القسم بنجاح')
    newSectionTitle.value = ''
    showNewSectionForm.value = false
  } catch (err: any) {
    toast.error('خطأ', err?.response?.data?.detail ?? 'فشل إضافة القسم')
  } finally {
    addingSectionLoading.value = false
  }
}

function startEditSection(section: Section) {
  editingSection.value = section
  editingTitle.value = section.title
}

async function handleUpdateSection() {
  if (!editingSection.value || !editingTitle.value.trim()) return
  try {
    await updateSection.mutateAsync({
      id: editingSection.value.id,
      courseId: courseId.value,
      payload: { title: editingTitle.value.trim() },
    })
    toast.success('تم التحديث', 'تم تحديث عنوان القسم')
    editingSection.value = null
  } catch (err: any) {
    toast.error('خطأ', err?.response?.data?.detail ?? 'فشل تحديث القسم')
  }
}

async function handleDeleteSection() {
  if (!confirmDeleteSection.value) return
  isDeletingSection.value = true
  try {
    await deleteSection.mutateAsync({
      id: confirmDeleteSection.value.id,
      courseId: courseId.value,
    })
    toast.success('تم الحذف', 'تم حذف القسم بنجاح')
    confirmDeleteSection.value = null
  } catch (err: any) {
    toast.error('خطأ', err?.response?.data?.detail ?? 'فشل حذف القسم')
  } finally {
    isDeletingSection.value = false
  }
}

const formattedDate = (dateStr: string) =>
  new Date(dateStr).toLocaleDateString('ar', { year: 'numeric', month: 'long', day: 'numeric' })
</script>

<template>
  <div dir="rtl" class="space-y-6 animate-fade-in">
    <ToastContainer />

    <!-- Loading -->
    <div v-if="isLoading" class="space-y-6">
      <div class="h-8 bg-slate-200 dark:bg-slate-700 rounded w-1/3 animate-pulse" />
      <div class="card p-6 animate-pulse space-y-4">
        <div class="h-6 bg-slate-200 dark:bg-slate-700 rounded w-1/2" />
        <div class="h-4 bg-slate-200 dark:bg-slate-700 rounded w-2/3" />
        <div class="grid grid-cols-3 gap-4 mt-4">
          <div v-for="i in 3" :key="i" class="h-20 bg-slate-200 dark:bg-slate-700 rounded-xl" />
        </div>
      </div>
    </div>

    <!-- Error -->
    <div v-else-if="isError" class="card p-12 text-center">
      <p class="text-slate-500 dark:text-slate-400 mb-4">لم يتم العثور على المساق</p>
      <button @click="router.push({ name: 'admin.courses' })" class="btn-primary">العودة</button>
    </div>

    <template v-else-if="course">
      <!-- Header -->
      <div class="flex items-start justify-between flex-wrap gap-4">
        <div class="flex items-center gap-4">
          <button @click="router.push({ name: 'admin.courses' })" class="btn-ghost border border-slate-200 dark:border-slate-700 p-2">
            <ArrowRight class="h-5 w-5" />
          </button>
          <div>
            <div class="flex items-center gap-3 flex-wrap">
              <h1 class="text-2xl font-extrabold text-slate-900 dark:text-white">{{ course.title }}</h1>
              <span :class="['text-xs font-bold px-2.5 py-1 rounded-full', COURSE_STATUS_COLORS[course.status]]">
                {{ COURSE_STATUS_LABELS[course.status] }}
              </span>
            </div>
            <p class="text-sm text-slate-500 dark:text-slate-400 mt-1">
              {{ course.subject_name ?? 'بدون مادة' }} • {{ course.stage ?? 'بدون مرحلة' }}
            </p>
          </div>
        </div>

        <!-- Actions -->
        <div class="flex items-center gap-2 flex-wrap">
          <button
            @click="router.push({ name: 'admin.course-edit', params: { id: course.id } })"
            class="btn-ghost border border-slate-200 dark:border-slate-700 gap-2"
          >
            <Edit class="h-4 w-4" /> تعديل
          </button>
          <button
            v-if="course.status !== 'published'"
            @click="handlePublish"
            :disabled="processingAction === 'publish'"
            class="btn-primary gap-2 bg-emerald-500 hover:bg-emerald-600"
          >
            <Loader2 v-if="processingAction === 'publish'" class="h-4 w-4 animate-spin" />
            <CheckCircle v-else class="h-4 w-4" />
            نشر
          </button>
          <button
            v-if="course.status !== 'archived'"
            @click="handleArchive"
            :disabled="processingAction === 'archive'"
            class="btn-ghost border border-amber-200 dark:border-amber-800 text-amber-600 gap-2"
          >
            <Loader2 v-if="processingAction === 'archive'" class="h-4 w-4 animate-spin" />
            <Archive v-else class="h-4 w-4" />
            أرشفة
          </button>
        </div>
      </div>

      <!-- Stats Row -->
      <div class="grid grid-cols-2 sm:grid-cols-4 gap-4">
        <div class="card p-4 flex items-center gap-3">
          <div class="h-10 w-10 rounded-xl bg-blue-50 dark:bg-blue-950/40 border border-blue-200 dark:border-blue-900/60 flex items-center justify-center">
            <Users class="h-5 w-5 text-blue-500" />
          </div>
          <div>
            <p class="text-xs text-slate-500">الطلاب</p>
            <p class="text-xl font-extrabold text-slate-900 dark:text-white">{{ course.students_count }}</p>
          </div>
        </div>
        <div class="card p-4 flex items-center gap-3">
          <div class="h-10 w-10 rounded-xl bg-orange-50 dark:bg-orange-950/40 border border-orange-200 dark:border-orange-900/60 flex items-center justify-center">
            <FileText class="h-5 w-5 text-orange-500" />
          </div>
          <div>
            <p class="text-xs text-slate-500">المحاضرات</p>
            <p class="text-xl font-extrabold text-slate-900 dark:text-white">{{ course.lessons_count }}</p>
          </div>
        </div>
        <div class="card p-4 flex items-center gap-3">
          <div class="h-10 w-10 rounded-xl bg-green-50 dark:bg-green-950/40 border border-green-200 dark:border-green-900/60 flex items-center justify-center">
            <BookOpen class="h-5 w-5 text-green-500" />
          </div>
          <div>
            <p class="text-xs text-slate-500">الأقسام</p>
            <p class="text-xl font-extrabold text-slate-900 dark:text-white">{{ sections?.length ?? 0 }}</p>
          </div>
        </div>
        <div class="card p-4 flex items-center gap-3">
          <div class="h-10 w-10 rounded-xl bg-purple-50 dark:bg-purple-950/40 border border-purple-200 dark:border-purple-900/60 flex items-center justify-center">
            <Globe class="h-5 w-5 text-purple-500" />
          </div>
          <div>
            <p class="text-xs text-slate-500">السعر</p>
            <p class="text-xl font-extrabold text-slate-900 dark:text-white">{{ course.price }} د.ع</p>
          </div>
        </div>
      </div>

      <!-- Quick Navigation -->
      <div class="card p-4 flex items-center justify-between flex-wrap gap-4">
        <span class="text-sm font-semibold text-slate-700 dark:text-slate-200">إدارة محتوى المساق:</span>
        <div class="flex gap-3">
          <button
            @click="router.push({ name: 'admin.course-lessons', params: { courseId: course.id } })"
            class="btn-primary gap-2 text-xs"
          >
            <BookOpen class="h-3.5 w-3.5" /> عرض المحاضرات ({{ course.lessons_count }})
          </button>
          <button
            @click="router.push({ name: 'admin.lesson-create', params: { courseId: course.id } })"
            class="btn-ghost border border-slate-200 dark:border-slate-700 gap-2 text-xs"
          >
            <Plus class="h-3.5 w-3.5" /> إضافة محاضرة
          </button>
          <button
            @click="router.push({ name: 'admin.course-students', params: { courseId: course.id } })"
            class="btn-ghost border border-slate-200 dark:border-slate-700 gap-2 text-xs"
          >
            <Users class="h-3.5 w-3.5" /> الطلاب المشتركون ({{ course.students_count }})
          </button>
        </div>
      </div>

      <!-- Teacher Assignment Section -->
      <div class="card p-6 space-y-4">
        <div class="flex items-center justify-between">
          <h2 class="text-sm font-bold text-slate-700 dark:text-slate-200 flex items-center gap-2">
            <UserCheck class="h-4 w-4 text-orange-500" />
            الأستاذ المعيّن
          </h2>
          <!-- Platform Owner Only controls -->
          <div v-if="adminStore.isOwner" class="flex items-center gap-2">
            <button
              v-if="course.teacher_id"
              @click="handleRemoveTeacher"
              :disabled="processingAction === 'remove-teacher'"
              class="btn-ghost border border-red-200 dark:border-red-900/50 text-red-500 gap-2 text-xs py-1.5"
            >
              <UserX class="h-3.5 w-3.5" /> إزالة
            </button>
            <button @click="showAssignForm = !showAssignForm; selectedUser = null; teacherSearch = ''; teacherIdInput = ''" class="btn-ghost border border-slate-200 dark:border-slate-700 gap-2 text-xs py-1.5">
              <Plus class="h-3.5 w-3.5" /> تعيين أستاذ
            </button>
          </div>
        </div>

        <!-- Current Teacher Card -->
        <div v-if="course.teacher_id" class="flex items-center gap-3 p-3 rounded-xl bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700">
          <img
            v-if="course.teacher_avatar"
            :src="course.teacher_avatar"
            :alt="course.teacher_name ?? ''"
            class="h-10 w-10 rounded-full object-cover shrink-0 ring-2 ring-orange-500/30"
          />
          <div v-else class="h-10 w-10 rounded-full bg-orange-100 dark:bg-orange-950/40 flex items-center justify-center shrink-0">
            <Users class="h-5 w-5 text-orange-500" />
          </div>
          <div>
            <p class="font-semibold text-slate-900 dark:text-white">{{ course.teacher_name ?? course.teacher_id }}</p>
            <p class="text-xs text-slate-400">{{ course.teacher_bio ?? 'أستاذ المساق' }}</p>
          </div>
        </div>
        <div v-else class="text-sm text-slate-400 italic">لم يتم تعيين أستاذ لهذا المساق</div>

        <!-- Search & Select Assign Form - Platform Owner Only -->
        <div v-if="adminStore.isOwner && showAssignForm" class="space-y-3">
          <!-- Search box -->
          <input
            v-model="teacherSearch"
            type="text"
            placeholder="ابحث باسم الأستاذ أو معرف التيليجرام..."
            class="w-full px-4 py-2.5 rounded-xl border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 text-sm focus:outline-none focus:ring-2 focus:ring-orange-400/50"
          />
          <!-- Manual ID input fallback -->
          <div class="space-y-1.5">
            <label class="text-xs text-slate-500 dark:text-slate-400">أو أدخل المعرف مباشرة (UUID / Telegram ID / Username):</label>
            <input
              v-model="teacherIdInput"
              type="text"
              placeholder="مثال: 123456789 أو @username أو uuid-xxxx..."
              class="w-full px-3 py-2 rounded-xl border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 text-sm font-mono focus:outline-none focus:ring-2 focus:ring-orange-400/50"
            />
          </div>
          <!-- User list -->
          <div v-if="usersLoading" class="text-xs text-slate-400 text-center py-3">جارٍ تحميل المستخدمين...</div>
          <div v-else-if="filteredUsers.length" class="max-h-52 overflow-y-auto rounded-xl border border-slate-200 dark:border-slate-700 divide-y divide-slate-100 dark:divide-slate-800">
            <button
              v-for="u in filteredUsers" :key="u.id ?? u.telegram_id"
              type="button"
              @click="selectUser(u)"
              :class="[
                'w-full flex items-center gap-3 p-3 text-right hover:bg-orange-50 dark:hover:bg-orange-950/20 transition-colors',
                selectedUser?.id === u.id ? 'bg-orange-50 dark:bg-orange-950/30' : ''
              ]"
            >
              <img
                v-if="u.photo_url || u.avatarUrl"
                :src="u.photo_url || u.avatarUrl"
                class="h-8 w-8 rounded-full object-cover shrink-0"
              />
              <div v-else class="h-8 w-8 rounded-full bg-slate-200 dark:bg-slate-700 flex items-center justify-center shrink-0 text-xs font-bold text-slate-500">
                {{ (u.first_name || u.username || '?').charAt(0).toUpperCase() }}
              </div>
              <div class="text-right flex-1">
                <p class="text-sm font-semibold text-slate-900 dark:text-white leading-tight">
                  {{ [u.first_name, u.last_name].filter(Boolean).join(' ') || u.username || u.name || '—' }}
                </p>
                <p class="text-[11px] text-slate-400 font-mono">{{ u.telegram_id || u.user_id || u.id }}</p>
              </div>
              <span v-if="u.role" :class="[
                'text-[10px] px-2 py-0.5 rounded-full font-bold',
                u.role === 'teacher' ? 'bg-green-100 text-green-700 dark:bg-green-950/40 dark:text-green-400' :
                u.role === 'master_admin' ? 'bg-orange-100 text-orange-700' : 'bg-slate-100 text-slate-500'
              ]">{{ u.role }}</span>
            </button>
          </div>
          <p v-else class="text-xs text-slate-400 text-center py-2">لا يوجد مستخدمون</p>

          <!-- Selected user preview + assign button -->
          <div v-if="selectedUser || teacherIdInput.trim()" class="flex items-center justify-between gap-3 p-3 rounded-xl bg-orange-50 dark:bg-orange-950/20 border border-orange-200 dark:border-orange-900/40">
            <p class="text-sm font-semibold text-orange-700 dark:text-orange-300">
              {{ selectedUser
                ? 'تم اختيار: ' + ([selectedUser.first_name, selectedUser.last_name].filter(Boolean).join(' ') || selectedUser.username || selectedUser.name)
                : 'تعيين بمعرف: ' + teacherIdInput.trim()
              }}
            </p>
            <button
              @click="handleAssignTeacher"
              :disabled="assigningTeacher"
              class="btn-primary gap-2 py-2 text-sm shrink-0"
            >
              <Loader2 v-if="assigningTeacher" class="h-4 w-4 animate-spin" />
              تعيين كأستاذ
            </button>
          </div>
        </div>
      </div>

      <!-- Sections -->
      <div class="card p-6 space-y-4">
        <div class="flex items-center justify-between">
          <h2 class="text-sm font-bold text-slate-700 dark:text-slate-200 flex items-center gap-2">
            <BookOpen class="h-4 w-4 text-orange-500" />
            أقسام المساق والمحاضرات
          </h2>
          <button
            @click="showNewSectionForm = !showNewSectionForm"
            class="btn-ghost border border-slate-200 dark:border-slate-700 gap-2 text-xs py-1.5"
          >
            <Plus class="h-3.5 w-3.5" /> قسم جديد
          </button>
        </div>

        <!-- New Section Form -->
        <div v-if="showNewSectionForm" class="flex gap-3">
          <input
            v-model="newSectionTitle"
            type="text"
            placeholder="عنوان القسم (مثال: الفصل الأول، الوحدة الأولى)..."
            @keyup.enter="handleAddSection"
            class="flex-1 px-4 py-2.5 rounded-xl border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 text-sm focus:outline-none focus:ring-2 focus:ring-orange-400/50"
          />
          <button
            @click="handleAddSection"
            :disabled="addingSectionLoading || !newSectionTitle.trim()"
            class="btn-primary gap-2 py-2.5 disabled:opacity-50"
          >
            <Loader2 v-if="addingSectionLoading" class="h-4 w-4 animate-spin" />
            <Plus v-else class="h-4 w-4" />
            إضافة
          </button>
          <button @click="showNewSectionForm = false; newSectionTitle = ''" class="btn-ghost border border-slate-200 dark:border-slate-700 py-2.5 px-4">
            إلغاء
          </button>
        </div>

        <!-- Sections List with Direct Upload Button for Each Section -->
        <div v-if="sectionsLoading" class="space-y-3">
          <div v-for="i in 3" :key="i" class="h-14 bg-slate-100 dark:bg-slate-800 rounded-xl animate-pulse" />
        </div>

        <div v-else-if="sections?.length" class="space-y-3">
          <div
            v-for="(section, idx) in sections"
            :key="section.id"
            class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 p-4 rounded-2xl border border-slate-200 dark:border-slate-700 bg-white/50 dark:bg-slate-900/40 hover:border-orange-500/50 transition-all group"
          >
            <div class="flex items-center gap-3 min-w-0">
              <div class="h-7 w-7 rounded-xl bg-orange-500/10 text-orange-500 flex items-center justify-center shrink-0 text-xs font-bold">
                {{ idx + 1 }}
              </div>

              <template v-if="editingSection?.id === section.id">
                <input
                  v-model="editingTitle"
                  type="text"
                  @keyup.enter="handleUpdateSection"
                  @keyup.escape="editingSection = null"
                  class="px-3 py-1.5 rounded-lg border border-orange-300 dark:border-orange-700 bg-white dark:bg-slate-800 text-sm focus:outline-none focus:ring-1 focus:ring-orange-400"
                  autofocus
                />
                <button @click="handleUpdateSection" class="btn-primary text-xs py-1.5 px-3">حفظ</button>
                <button @click="editingSection = null" class="btn-ghost text-xs py-1.5 px-3 border border-slate-200 dark:border-slate-700">إلغاء</button>
              </template>

              <template v-else>
                <div class="min-w-0">
                  <p class="text-sm font-extrabold text-slate-900 dark:text-white truncate">{{ section.title }}</p>
                  <p class="text-xs text-slate-400 mt-0.5">
                    {{ section.lessons_count > 0 ? `${section.lessons_count} محاضرة` : 'لا توجد محاضرات بعد' }}
                  </p>
                </div>
              </template>
            </div>

            <!-- Action buttons: Direct Upload to this Section + Edit / Delete -->
            <div class="flex items-center gap-2 shrink-0 justify-end pt-2 sm:pt-0 border-t sm:border-t-0 border-slate-100 dark:border-slate-800">
              <button
                type="button"
                @click="router.push({ name: 'admin.lesson-create', query: { courseId: course.id, sectionId: section.id } })"
                class="btn-primary text-xs py-1.5 px-3 flex items-center gap-1.5 shadow-sm shadow-orange-500/20 active:scale-95 transition-all cursor-pointer"
                title="رفع محاضرة لهذا القسم مباشرة"
              >
                <Plus class="h-3.5 w-3.5" />
                <span>رفع محاضرة للقسم</span>
              </button>

              <button
                @click="startEditSection(section)"
                class="p-2 rounded-xl hover:bg-blue-50 dark:hover:bg-blue-950/30 text-slate-400 hover:text-blue-500 transition-colors cursor-pointer"
                title="تعديل اسم القسم"
              >
                <Edit class="h-3.5 w-3.5" />
              </button>

              <button
                @click="confirmDeleteSection = section"
                class="p-2 rounded-xl hover:bg-red-50 dark:hover:bg-red-950/30 text-slate-400 hover:text-red-500 transition-colors cursor-pointer"
                title="حذف القسم"
              >
                <Trash2 class="h-3.5 w-3.5" />
              </button>
            </div>
          </div>
        </div>

        <div v-else class="text-center py-8">
          <BookOpen class="h-10 w-10 text-slate-200 dark:text-slate-700 mx-auto mb-3" />
          <p class="text-sm text-slate-400 mb-3">لا توجد أقسام في هذا المساق بعد. أنشئ قسماً لتتمكن من رفع المحاضرات داخله.</p>
          <button @click="showNewSectionForm = true" class="btn-primary text-sm gap-2">
            <Plus class="h-4 w-4" /> إضافة أول قسم للمساق
          </button>
        </div>
      </div>

      <!-- Course Info -->
      <div class="card p-6">
        <h2 class="text-sm font-bold text-slate-700 dark:text-slate-200 mb-4">معلومات المساق</h2>
        <dl class="space-y-3">
          <div class="flex justify-between text-sm">
            <dt class="text-slate-500">الوصف</dt>
            <dd class="text-slate-900 dark:text-white max-w-sm text-right">{{ course.description ?? '—' }}</dd>
          </div>
          <div class="flex justify-between text-sm">
            <dt class="text-slate-500">القسم</dt>
            <dd class="text-slate-900 dark:text-white">{{ course.department ?? '—' }}</dd>
          </div>
          <div class="flex justify-between text-sm">
            <dt class="text-slate-500">تاريخ الإنشاء</dt>
            <dd class="text-slate-900 dark:text-white" dir="ltr">{{ formattedDate(course.created_at) }}</dd>
          </div>
          <div v-if="course.published_at" class="flex justify-between text-sm">
            <dt class="text-slate-500">تاريخ النشر</dt>
            <dd class="text-slate-900 dark:text-white" dir="ltr">{{ formattedDate(course.published_at) }}</dd>
          </div>
        </dl>
      </div>
    </template>

    <!-- Delete Section Confirm -->
    <ConfirmDialog
      :open="!!confirmDeleteSection"
      title="حذف القسم"
      :message="`هل أنت متأكد من حذف قسم '${confirmDeleteSection?.title}'? ستُفصل المحاضرات عن هذا القسم.`"
      confirm-label="حذف"
      cancel-label="إلغاء"
      variant="danger"
      :loading="isDeletingSection"
      @confirm="handleDeleteSection"
      @cancel="confirmDeleteSection = null"
    />
  </div>
</template>
