<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useSubjects, useDeleteSubject } from '@/admin/queries/useSubjects'
import { useToast } from '@/admin/composables/useToast'
import ConfirmDialog from '@/admin/components/ui/ConfirmDialog.vue'
import ToastContainer from '@/admin/components/ui/ToastContainer.vue'
import {
  Plus, Search, Filter, FolderOpen, Edit, Trash2,
  ChevronLeft, ChevronRight, RefreshCw, BookOpen, Users,
  AlertTriangle, Archive
} from '@lucide/vue'
import type { SubjectListFilters } from '@/admin/types/subject'
import type { Subject } from '@/admin/types/subject'

const router = useRouter()
const toast = useToast()

const filters = reactive<SubjectListFilters>({
  search: '',
  department: '',
  stage: '',
  status: undefined,
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

const { data, isLoading, isError, refetch } = useSubjects(computed(() => ({ ...filters })))
const deleteSubject = useDeleteSubject()

const confirmDelete = ref<Subject | null>(null)
const isDeleting = ref(false)

const stages = ['الأول', 'الثاني', 'الثالث', 'الرابع', 'الخامس', 'السادس']
const statuses = [
  { value: 'active', label: 'نشط' },
  { value: 'inactive', label: 'غير نشط' },
  { value: 'archived', label: 'مؤرشف' },
]

const statusClasses: Record<string, string> = {
  active: 'bg-emerald-100 text-emerald-700 dark:bg-emerald-950/40 dark:text-emerald-400',
  inactive: 'bg-amber-100 text-amber-700 dark:bg-amber-950/40 dark:text-amber-400',
  archived: 'bg-slate-100 text-slate-600 dark:bg-slate-800 dark:text-slate-400',
}
const statusLabels: Record<string, string> = {
  active: 'نشط',
  inactive: 'غير نشط',
  archived: 'مؤرشف',
}

async function handleDelete() {
  if (!confirmDelete.value) return
  isDeleting.value = true
  try {
    const result = await deleteSubject.mutateAsync(confirmDelete.value.id)
    if (result.archived) {
      toast.warning('تم الأرشفة', result.message)
    } else {
      toast.success('تم الحذف', 'تم حذف المادة بنجاح')
    }
    confirmDelete.value = null
  } catch (err: any) {
    toast.error('خطأ', err?.response?.data?.detail ?? 'فشل حذف المادة')
  } finally {
    isDeleting.value = false
  }
}

function goPage(page: number) {
  filters.page = page
}

const subjects = computed(() => data.value?.data ?? [])
const meta = computed(() => data.value?.meta)
</script>

<template>
  <div dir="rtl" class="space-y-6 animate-fade-in">
    <ToastContainer />

    <!-- Header -->
    <div class="flex items-center justify-between flex-wrap gap-4">
      <div>
        <h1 class="text-2xl font-extrabold text-slate-900 dark:text-white">المواد الدراسية</h1>
        <p class="mt-1 text-sm text-slate-500 dark:text-slate-400">
          إدارة المواد الدراسية المتاحة في المنصة
        </p>
      </div>
      <button @click="router.push({ name: 'admin.subject-create' })" class="btn-primary gap-2">
        <Plus class="h-4 w-4" />
        مادة جديدة
      </button>
    </div>

    <!-- Filters -->
    <div class="card p-4 space-y-4">
      <div class="flex flex-wrap gap-3">
        <!-- Search -->
        <div class="relative flex-1 min-w-48">
          <Search class="absolute right-3 top-1/2 -translate-y-1/2 h-4 w-4 text-slate-400" />
          <input
            v-model="searchInput"
            @input="onSearch"
            type="text"
            placeholder="البحث باسم المادة أو الكود..."
            class="w-full pr-9 pl-4 py-2.5 rounded-xl border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 text-sm focus:outline-none focus:ring-2 focus:ring-orange-400/50 focus:border-orange-400 transition-colors"
          />
        </div>

        <!-- Stage filter -->
        <select
          v-model="filters.stage"
          @change="filters.page = 1"
          class="px-4 py-2.5 rounded-xl border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 text-sm focus:outline-none focus:ring-2 focus:ring-orange-400/50 text-slate-700 dark:text-slate-200"
        >
          <option value="">كل المراحل</option>
          <option v-for="stage in stages" :key="stage" :value="stage">{{ stage }}</option>
        </select>

        <!-- Status filter -->
        <select
          v-model="filters.status"
          @change="filters.page = 1"
          class="px-4 py-2.5 rounded-xl border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 text-sm focus:outline-none focus:ring-2 focus:ring-orange-400/50 text-slate-700 dark:text-slate-200"
        >
          <option :value="undefined">كل الحالات</option>
          <option v-for="s in statuses" :key="s.value" :value="s.value">{{ s.label }}</option>
        </select>

        <button @click="refetch()" :disabled="isLoading" class="btn-ghost border border-slate-200 dark:border-slate-700 gap-2">
          <RefreshCw :class="['h-4 w-4', isLoading ? 'animate-spin' : '']" />
        </button>
      </div>
    </div>

    <!-- Error -->
    <div v-if="isError" class="card p-8 text-center border-red-200 bg-red-50 dark:bg-red-950/20">
      <p class="text-red-600 dark:text-red-400 font-medium mb-3">حدث خطأ أثناء تحميل المواد</p>
      <button @click="refetch()" class="btn-secondary text-red-500 border-red-200">إعادة المحاولة</button>
    </div>

    <!-- Loading Skeleton -->
    <div v-else-if="isLoading" class="card overflow-hidden">
      <div class="divide-y divide-slate-100 dark:divide-slate-800">
        <div v-for="i in 8" :key="i" class="p-4 animate-pulse flex gap-4">
          <div class="h-10 w-10 bg-slate-200 dark:bg-slate-700 rounded-xl shrink-0" />
          <div class="flex-1 space-y-2 pt-1">
            <div class="h-3.5 bg-slate-200 dark:bg-slate-700 rounded w-1/3" />
            <div class="h-3 bg-slate-200 dark:bg-slate-700 rounded w-1/2" />
          </div>
          <div class="h-6 w-16 bg-slate-200 dark:bg-slate-700 rounded-full" />
        </div>
      </div>
    </div>

    <!-- Table -->
    <div v-else-if="subjects.length" class="card overflow-hidden">
      <!-- Desktop Table -->
      <div class="hidden md:block overflow-x-auto">
        <table class="w-full text-sm">
          <thead>
            <tr class="border-b border-slate-100 dark:border-slate-800 bg-slate-50 dark:bg-slate-800/50">
              <th class="text-right font-bold text-slate-500 dark:text-slate-400 px-4 py-3">المادة</th>
              <th class="text-right font-bold text-slate-500 dark:text-slate-400 px-4 py-3">الكود</th>
              <th class="text-right font-bold text-slate-500 dark:text-slate-400 px-4 py-3">القسم</th>
              <th class="text-right font-bold text-slate-500 dark:text-slate-400 px-4 py-3">المرحلة</th>
              <th class="text-right font-bold text-slate-500 dark:text-slate-400 px-4 py-3">المساقات</th>
              <th class="text-right font-bold text-slate-500 dark:text-slate-400 px-4 py-3">الحالة</th>
              <th class="text-right font-bold text-slate-500 dark:text-slate-400 px-4 py-3">إجراءات</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100 dark:divide-slate-800">
            <tr
              v-for="subject in subjects"
              :key="subject.id"
              class="hover:bg-slate-50 dark:hover:bg-slate-800/30 transition-colors"
            >
              <td class="px-4 py-3">
                <div class="flex items-center gap-3">
                  <div class="h-9 w-9 rounded-xl bg-purple-50 dark:bg-purple-950/40 border border-purple-200 dark:border-purple-900/60 flex items-center justify-center shrink-0">
                    <FolderOpen class="h-4 w-4 text-purple-500" />
                  </div>
                  <div>
                    <p class="font-semibold text-slate-900 dark:text-white">{{ subject.name }}</p>
                    <p v-if="subject.description" class="text-xs text-slate-400 truncate max-w-xs">{{ subject.description }}</p>
                  </div>
                </div>
              </td>
              <td class="px-4 py-3">
                <span class="font-mono text-xs bg-slate-100 dark:bg-slate-800 px-2 py-1 rounded-lg text-slate-600 dark:text-slate-300" dir="ltr">
                  {{ subject.code }}
                </span>
              </td>
              <td class="px-4 py-3 text-slate-600 dark:text-slate-300">{{ subject.department ?? '—' }}</td>
              <td class="px-4 py-3 text-slate-600 dark:text-slate-300">{{ subject.stage ?? '—' }}</td>
              <td class="px-4 py-3">
                <span class="flex items-center gap-1.5 text-slate-600 dark:text-slate-300">
                  <BookOpen class="h-3.5 w-3.5 text-slate-400" />
                  {{ subject.courses_count }}
                </span>
              </td>
              <td class="px-4 py-3">
                <span :class="['text-xs font-bold px-2.5 py-1 rounded-full', statusClasses[subject.status]]">
                  {{ statusLabels[subject.status] }}
                </span>
              </td>
              <td class="px-4 py-3">
                <div class="flex items-center gap-2">
                  <button
                    @click="router.push({ name: 'admin.subject-edit', params: { id: subject.id } })"
                    class="p-2 rounded-lg hover:bg-blue-50 dark:hover:bg-blue-950/30 text-slate-400 hover:text-blue-500 transition-colors"
                    title="تعديل"
                  >
                    <Edit class="h-4 w-4" />
                  </button>
                  <button
                    @click="confirmDelete = subject"
                    class="p-2 rounded-lg hover:bg-red-50 dark:hover:bg-red-950/30 text-slate-400 hover:text-red-500 transition-colors"
                    title="حذف"
                  >
                    <Trash2 class="h-4 w-4" />
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Mobile Cards -->
      <div class="md:hidden divide-y divide-slate-100 dark:divide-slate-800">
        <div
          v-for="subject in subjects"
          :key="subject.id"
          class="p-4 space-y-3"
        >
          <div class="flex items-start gap-3">
            <div class="h-10 w-10 rounded-xl bg-purple-50 dark:bg-purple-950/40 border border-purple-200 dark:border-purple-900/60 flex items-center justify-center shrink-0">
              <FolderOpen class="h-5 w-5 text-purple-500" />
            </div>
            <div class="flex-1 min-w-0">
              <p class="font-semibold text-slate-900 dark:text-white">{{ subject.name }}</p>
              <p class="text-xs text-slate-400 font-mono" dir="ltr">{{ subject.code }}</p>
            </div>
            <span :class="['text-xs font-bold px-2 py-0.5 rounded-full shrink-0', statusClasses[subject.status]]">
              {{ statusLabels[subject.status] }}
            </span>
          </div>
          <div class="flex items-center justify-between text-xs text-slate-500">
            <div class="flex items-center gap-3">
              <span>{{ subject.department ?? '—' }}</span>
              <span v-if="subject.stage">• {{ subject.stage }}</span>
              <span class="flex items-center gap-1"><BookOpen class="h-3 w-3" /> {{ subject.courses_count }}</span>
            </div>
            <div class="flex gap-2">
              <button @click="router.push({ name: 'admin.subject-edit', params: { id: subject.id } })" class="p-1.5 rounded-lg hover:bg-blue-50 text-slate-400 hover:text-blue-500">
                <Edit class="h-4 w-4" />
              </button>
              <button @click="confirmDelete = subject" class="p-1.5 rounded-lg hover:bg-red-50 text-slate-400 hover:text-red-500">
                <Trash2 class="h-4 w-4" />
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Pagination -->
      <div v-if="meta && meta.total_pages > 1" class="border-t border-slate-100 dark:border-slate-800 p-4 flex items-center justify-between">
        <p class="text-sm text-slate-500">
          {{ meta.total }} مادة • صفحة {{ meta.page }} من {{ meta.total_pages }}
        </p>
        <div class="flex items-center gap-2">
          <button
            :disabled="filters.page === 1"
            @click="goPage(filters.page! - 1)"
            class="btn-ghost border border-slate-200 dark:border-slate-700 px-3 py-2 disabled:opacity-40"
          >
            <ChevronRight class="h-4 w-4" />
          </button>
          <button
            :disabled="filters.page === meta.total_pages"
            @click="goPage(filters.page! + 1)"
            class="btn-ghost border border-slate-200 dark:border-slate-700 px-3 py-2 disabled:opacity-40"
          >
            <ChevronLeft class="h-4 w-4" />
          </button>
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <div v-else class="card p-16 text-center">
      <FolderOpen class="h-14 w-14 text-slate-200 dark:text-slate-700 mx-auto mb-4" />
      <h3 class="text-lg font-bold text-slate-700 dark:text-slate-300 mb-2">لا توجد مواد دراسية</h3>
      <p class="text-sm text-slate-400 mb-6">ابدأ بإنشاء أول مادة دراسية في المنصة</p>
      <button @click="router.push({ name: 'admin.subject-create' })" class="btn-primary gap-2">
        <Plus class="h-4 w-4" />
        إنشاء مادة
      </button>
    </div>

    <!-- Delete Confirm -->
    <ConfirmDialog
      :open="!!confirmDelete"
      title="حذف المادة الدراسية"
      :message="confirmDelete?.courses_count
        ? `مادة '${confirmDelete?.name}' مرتبطة بـ ${confirmDelete?.courses_count} مساق. سيتم أرشفتها بدلاً من حذفها.`
        : `هل أنت متأكد من حذف مادة '${confirmDelete?.name}'؟ لا يمكن التراجع عن هذا الإجراء.`"
      :confirm-label="confirmDelete?.courses_count ? 'أرشفة' : 'حذف'"
      cancel-label="إلغاء"
      variant="danger"
      :loading="isDeleting"
      @confirm="handleDelete"
      @cancel="confirmDelete = null"
    />
  </div>
</template>
