<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useCourses, useDeleteCourse, usePublishCourse, useArchiveCourse } from '@/admin/queries/useCourses'
import { useToast } from '@/admin/composables/useToast'
import ConfirmDialog from '@/admin/components/ui/ConfirmDialog.vue'
import ActionDropdown, { type MenuItem } from '@/admin/components/ui/ActionDropdown.vue'
import ToastContainer from '@/admin/components/ui/ToastContainer.vue'
import {
  Plus, Search, BookOpen, Edit, Trash2, ChevronLeft, ChevronRight,
  RefreshCw, Users, FileText, CheckCircle, Archive, Eye
} from '@lucide/vue'
import type { CourseListFilters, Course } from '@/admin/types/course'
import { COURSE_STATUS_LABELS } from '@/admin/types/course'

const router = useRouter()
const toast = useToast()

const filters = reactive<CourseListFilters>({
  page: 1,
  limit: 20,
})
const searchInput = ref('')
let searchTimer: ReturnType<typeof setTimeout>

function onSearch() {
  clearTimeout(searchTimer)
  searchTimer = setTimeout(() => {
    filters.search = searchInput.value || undefined
    filters.page = 1
  }, 400)
}

const { data, isLoading, isError, refetch } = useCourses(computed(() => ({ ...filters })))
const deleteCourse = useDeleteCourse()
const publishCourse = usePublishCourse()
const archiveCourse = useArchiveCourse()

const confirmDelete = ref<Course | null>(null)
const isDeleting = ref(false)
const processingId = ref<string | null>(null)

const courses = computed(() => data.value?.data ?? [])
const meta = computed(() => data.value?.meta)

async function handleDelete() {
  if (!confirmDelete.value) return
  isDeleting.value = true
  try {
    await deleteCourse.mutateAsync(confirmDelete.value.id)
    toast.success('تم الحذف', 'تم حذف المساق بنجاح')
    confirmDelete.value = null
  } catch (err: any) {
    toast.error('خطأ', err?.response?.data?.detail ?? 'فشل حذف المساق')
  } finally {
    isDeleting.value = false
  }
}

async function handlePublish(course: Course) {
  processingId.value = course.id
  try {
    await publishCourse.mutateAsync(course.id)
    toast.success('تم النشر', `تم نشر مساق "${course.title}" بنجاح`)
  } catch (err: any) {
    toast.error('خطأ', err?.response?.data?.detail ?? 'فشل نشر المساق')
  } finally {
    processingId.value = null
  }
}

async function handleArchive(course: Course) {
  processingId.value = course.id
  try {
    await archiveCourse.mutateAsync(course.id)
    toast.success('تم الأرشفة', `تم أرشفة مساق "${course.title}"`)
  } catch (err: any) {
    toast.error('خطأ', err?.response?.data?.detail ?? 'فشل أرشفة المساق')
  } finally {
    processingId.value = null
  }
}

function getCourseActionItems(course: Course): MenuItem[] {
  const items: MenuItem[] = [
    {
      id: 'detail',
      label: 'عرض التفاصيل',
      icon: Eye,
      onClick: () => router.push({ name: 'admin.course-detail', params: { id: course.id } }),
    },
    {
      id: 'edit',
      label: 'تعديل المساق',
      icon: Edit,
      onClick: () => router.push({ name: 'admin.course-edit', params: { id: course.id } }),
    },
  ]

  if (course.status !== 'published') {
    items.push({
      id: 'publish',
      label: 'نشر المساق',
      icon: CheckCircle,
      onClick: () => handlePublish(course),
    })
  }

  if (course.status !== 'archived') {
    items.push({
      id: 'archive',
      label: 'أرشفة المساق',
      icon: Archive,
      onClick: () => handleArchive(course),
    })
  }

  items.push({
    id: 'delete',
    label: 'حذف المساق',
    icon: Trash2,
    danger: true,
    onClick: () => { confirmDelete.value = course },
  })

  return items
}

function goPage(page: number) {
  filters.page = page
}
</script>

<template>
  <div dir="rtl" class="space-y-6 animate-fade-in text-[#F5F5F5]">
    <ToastContainer />
    <!-- Header -->
    <div class="flex items-center justify-between flex-wrap gap-4">
      <div>
        <h1 class="text-2xl font-extrabold text-white flex items-center gap-2">
          <BookOpen class="h-6 w-6 text-[#F97316]" />
          المساقات التعليمية
        </h1>
        <p class="mt-1 text-xs text-[#A1A1AA]">إدارة جميع المساقات الدراسية بالمنصة</p>
      </div>
      <button @click="router.push({ name: 'admin.course-create' })" class="btn-primary gap-2 text-xs">
        <Plus class="h-4 w-4" />
        إنشاء مساق جديد
      </button>
    </div>

    <!-- Filters -->
    <div class="card p-4">
      <div class="flex flex-wrap gap-3">
        <div class="relative flex-1 min-w-48">
          <Search class="absolute right-3 top-1/2 -translate-y-1/2 h-4 w-4 text-[#71717A]" />
          <input
            v-model="searchInput"
            @input="onSearch"
            type="text"
            placeholder="البحث باسم المساق..."
            class="input w-full pr-9 text-sm"
          />
        </div>
        <select
          v-model="filters.status"
          @change="filters.page = 1"
          class="input text-sm font-semibold w-40"
        >
          <option :value="undefined">كل الحالات</option>
          <option value="draft">مسودة</option>
          <option value="published">منشور</option>
          <option value="archived">مؤرشف</option>
        </select>
        <button @click="refetch()" :disabled="isLoading" class="btn-secondary text-xs">
          <RefreshCw :class="['h-4 w-4', isLoading ? 'animate-spin' : '']" />
        </button>
      </div>
    </div>

    <!-- Error State -->
    <div v-if="isError" class="card p-8 text-center border-red-900/50 bg-red-950/20">
      <p class="text-red-400 font-medium mb-3 text-xs">حدث خطأ أثناء تحميل المساقات</p>
      <button @click="refetch()" class="btn-secondary text-xs text-red-400">إعادة المحاولة</button>
    </div>

    <!-- Loading -->
    <div v-else-if="isLoading" class="card overflow-hidden">
      <div class="divide-y divide-[#242424]">
        <div v-for="i in 5" :key="i" class="p-4 animate-pulse flex gap-4">
          <div class="h-12 w-12 bg-[#181818] rounded-xl shrink-0" />
          <div class="flex-1 space-y-2 pt-1.5">
            <div class="h-3.5 bg-[#181818] rounded w-1/2" />
            <div class="h-3 bg-[#181818] rounded w-1/3" />
          </div>
        </div>
      </div>
    </div>

    <!-- Course List -->
    <div v-else-if="courses.length" class="card overflow-hidden">
      <!-- Desktop Table -->
      <div class="hidden lg:block overflow-x-auto">
        <table class="w-full text-sm text-[#F5F5F5]">
          <thead>
            <tr class="bg-[#0F0F0F] border-b border-[#242424] text-right">
              <th class="font-bold text-[#A1A1AA] px-4 py-3.5 text-xs">المساق</th>
              <th class="font-bold text-[#A1A1AA] px-4 py-3.5 text-xs">المادة</th>
              <th class="font-bold text-[#A1A1AA] px-4 py-3.5 text-xs">الأستاذ</th>
              <th class="font-bold text-[#A1A1AA] px-4 py-3.5 text-xs">الطلاب</th>
              <th class="font-bold text-[#A1A1AA] px-4 py-3.5 text-xs">المحاضرات</th>
              <th class="font-bold text-[#A1A1AA] px-4 py-3.5 text-xs">الحالة</th>
              <th class="font-bold text-[#A1A1AA] px-4 py-3.5 text-xs text-center w-16">إجراءات</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-[#242424] bg-[#111111]">
            <tr
              v-for="course in courses"
              :key="course.id"
              class="hover:bg-[#161616] transition-colors"
            >
              <td class="px-4 py-3.5">
                <div class="flex items-center gap-3">
                  <div class="h-10 w-10 rounded-xl bg-[#1F0E03] border border-orange-900/40 flex items-center justify-center shrink-0 overflow-hidden text-[#F97316]">
                    <img v-if="course.thumbnail_url" :src="course.thumbnail_url" class="h-full w-full object-cover" />
                    <BookOpen v-else class="h-5 w-5" />
                  </div>
                  <div class="max-w-xs">
                    <p
                      class="font-bold text-white truncate cursor-pointer hover:text-[#F97316] transition-colors"
                      @click="router.push({ name: 'admin.course-detail', params: { id: course.id } })"
                    >{{ course.title }}</p>
                    <p v-if="course.stage" class="text-[11px] text-[#A1A1AA]">{{ course.stage }}</p>
                  </div>
                </div>
              </td>
              <td class="px-4 py-3.5 text-[#A1A1AA] text-xs">{{ course.subject_name ?? '—' }}</td>
              <td class="px-4 py-3.5 text-[#A1A1AA] text-xs">{{ course.teacher_name ?? '—' }}</td>
              <td class="px-4 py-3.5 text-xs text-[#A1A1AA]">
                <span class="flex items-center gap-1.5">
                  <Users class="h-3.5 w-3.5 text-[#71717A]" /> {{ course.students_count }}
                </span>
              </td>
              <td class="px-4 py-3.5 text-xs text-[#A1A1AA]">
                <span class="flex items-center gap-1.5">
                  <FileText class="h-3.5 w-3.5 text-[#71717A]" /> {{ course.lessons_count }}
                </span>
              </td>
              <td class="px-4 py-3.5">
                <span
                  :class="[
                    'text-[10px] font-extrabold px-2.5 py-1 rounded-full border inline-block',
                    course.status === 'published' ? 'bg-[#052e16] text-[#22C55E] border-[#14532d]' :
                    course.status === 'archived' ? 'bg-[#451a03] text-[#F59E0B] border-[#78350f]' :
                    'bg-[#18181b] text-[#A1A1AA] border-[#27272a]',
                  ]"
                >
                  {{ COURSE_STATUS_LABELS[course.status] || course.status }}
                </span>
              </td>
              <td class="px-4 py-3.5 text-center">
                <ActionDropdown :items="getCourseActionItems(course)" />
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Mobile Cards -->
      <div class="lg:hidden divide-y divide-[#242424]">
        <div
          v-for="course in courses"
          :key="course.id"
          class="p-4 space-y-3 hover:bg-[#161616]"
        >
          <div class="flex items-start justify-between gap-3">
            <div class="flex items-center gap-3">
              <div class="h-10 w-10 rounded-xl bg-[#1F0E03] border border-orange-900/40 flex items-center justify-center shrink-0 text-[#F97316]">
                <BookOpen class="h-5 w-5" />
              </div>
              <div>
                <p class="font-bold text-white text-sm">{{ course.title }}</p>
                <p class="text-[11px] text-[#A1A1AA]">{{ course.teacher_name ?? '—' }} • {{ course.subject_name ?? '—' }}</p>
              </div>
            </div>
            <ActionDropdown :items="getCourseActionItems(course)" />
          </div>
        </div>
      </div>

      <!-- Pagination -->
      <div v-if="meta && meta.total_pages > 1" class="border-t border-[#242424] p-4 flex items-center justify-between">
        <p class="text-xs text-[#A1A1AA]">{{ meta.total }} مساق • صفحة {{ meta.page }} من {{ meta.total_pages }}</p>
        <div class="flex items-center gap-2">
          <button :disabled="filters.page === 1" @click="goPage(filters.page! - 1)" class="btn-secondary text-xs px-2.5 py-1.5 disabled:opacity-40">
            <ChevronRight class="h-4 w-4" />
          </button>
          <button :disabled="filters.page === meta.total_pages" @click="goPage(filters.page! + 1)" class="btn-secondary text-xs px-2.5 py-1.5 disabled:opacity-40">
            <ChevronLeft class="h-4 w-4" />
          </button>
        </div>
      </div>
    </div>

    <!-- Empty -->
    <div v-else class="card p-16 text-center space-y-3">
      <BookOpen class="h-12 w-12 text-[#71717A] mx-auto" />
      <h3 class="text-base font-bold text-white">لا توجد مساقات</h3>
      <p class="text-xs text-[#A1A1AA]">ابدأ بإنشاء أول مساق تعليمي في المنصة</p>
      <button @click="router.push({ name: 'admin.course-create' })" class="btn-primary text-xs gap-2">
        <Plus class="h-4 w-4" /> إنشاء مساق
      </button>
    </div>

    <!-- Delete Confirm -->
    <ConfirmDialog
      :open="!!confirmDelete"
      title="حذف المساق"
      :message="`هل أنت متأكد من حذف مساق '${confirmDelete?.title}'؟ لا يمكن التراجع عن هذا الإجراء.`"
      confirm-label="حذف"
      cancel-label="إلغاء"
      variant="danger"
      :loading="isDeleting"
      @confirm="handleDelete"
      @cancel="confirmDelete = null"
    />
  </div>
</template>
