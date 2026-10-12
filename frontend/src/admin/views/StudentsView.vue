<script setup lang="ts">
import { ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { Search, Users, Filter, RefreshCw, ChevronRight, ChevronLeft, Power } from '@lucide/vue'
import { useStudents, useSuspendStudent, useActivateStudent } from '@/admin/queries/useStudents'
import StudentTable from '@/admin/components/StudentTable.vue'
import StatePanel from '@/admin/components/StatePanel.vue'
import ConfirmDialog from '@/admin/components/ui/ConfirmDialog.vue'
import { useToast } from '@/admin/composables/useToast'
import type { Student } from '@/admin/types/panel'

const router = useRouter()
const toast = useToast()

const search = ref('')
const status = ref<'active' | 'suspended' | ''>('')
const page = ref(1)
const limit = 20

const filters = ref({ search: '', status: '' })
const { data, isLoading, isError, error, refetch } = useStudents(filters)

const suspend = useSuspendStudent()
const activate = useActivateStudent()
const pendingAction = ref<Student | null>(null)
const actionMode = ref<'suspend' | 'activate'>('suspend')
const actionLoading = ref(false)

watch([search, status], () => {
  page.value = 1
})

function applyFilters() {
  filters.value = { search: search.value.trim(), status: status.value }
  page.value = 1
}

function goToPage(p: number) {
  page.value = p
}

function openAction(s: Student) {
  actionMode.value = s.is_active ? 'suspend' : 'activate'
  pendingAction.value = s
}

function confirmAction() {
  if (!pendingAction.value) return
  actionLoading.value = true
  const id = pendingAction.value.id
  const isSuspend = actionMode.value === 'suspend'
  const mutation = isSuspend ? suspend : activate
  mutation.mutate(id, {
      onSuccess: () => {
        toast.success(isSuspend ? 'تم إيقاف حساب الطالب' : 'تم تفعيل حساب الطالب')
        pendingAction.value = null
        actionLoading.value = false
      },
      onError: (err: unknown) => {
        toast.error('فشلت العملية', (err as Error)?.message)
        pendingAction.value = null
        actionLoading.value = false
      },
    }
  )
}

function viewStudent(s: Student) {
  router.push({ name: 'admin.student-detail', params: { id: s.id } })
}
</script>

<template>
  <div class="space-y-6">
    <div class="flex flex-wrap items-center justify-between gap-3">
      <div>
        <h1 class="text-xl font-extrabold flex items-center gap-2">
          <Users class="h-5 w-5 text-orange-500" />
          الطلاب
        </h1>
        <p class="text-xs text-slate-500 dark:text-slate-400 mt-1">
          إدارة حسابات الطلاب والاشتراكات والنشاط
        </p>
      </div>
      <button class="btn-secondary text-sm" @click="() => refetch()">
        <RefreshCw class="h-4 w-4" />
        تحديث
      </button>
    </div>

    <!-- Filters -->
    <div class="card p-4">
      <div class="flex flex-wrap items-center gap-3">
        <div class="relative flex-1 min-w-[220px]">
          <Search class="absolute right-3 top-1/2 -translate-y-1/2 h-4 w-4 text-slate-400" />
          <input
            v-model="search"
            type="text"
            placeholder="ابحث بالاسم أو اسم المستخدم أو Telegram ID"
            class="input w-full pr-10 text-sm"
            @keyup.enter="applyFilters"
          />
        </div>
        <select v-model="status" class="input text-sm w-40">
          <option value="">كل الحالات</option>
          <option value="active">نشط</option>
          <option value="suspended">موقوف</option>
        </select>
        <button class="btn-primary text-sm" @click="applyFilters">
          <Filter class="h-4 w-4" />
          تصفية
        </button>
      </div>
    </div>

    <!-- Table -->
    <div class="card overflow-hidden">
      <StatePanel
        :loading="isLoading"
        :error="isError ? (error as Error)?.message : null"
        :empty="!data?.data?.length"
        empty-title="لا يوجد طلاب"
        empty-message="لم يتم العثور على طلاب مطابقين"
        :skeleton-count="6"
        @retry="() => refetch()"
      >
        <StudentTable
          :students="data?.data ?? []"
          @view="viewStudent"
          @toggle-status="openAction"
        />

        <!-- Pagination -->
        <div v-if="data && data.meta && data.meta.total_pages > 1" class="flex items-center justify-between px-4 py-3 border-t border-slate-100 dark:border-slate-800">
          <p class="text-xs text-slate-400">
            عرض {{ data.data.length }} من {{ data.meta.total }} طالب
          </p>
          <div class="flex items-center gap-2">
            <button
              class="btn-secondary text-xs py-1.5"
              :disabled="page <= 1"
              @click="goToPage(page - 1)"
            >
              <ChevronRight class="h-3.5 w-3.5" />
              السابق
            </button>
            <span class="text-xs font-bold text-slate-500">{{ page }} / {{ data.meta.total_pages }}</span>
            <button
              class="btn-secondary text-xs py-1.5"
              :disabled="page >= data.meta.total_pages"
              @click="goToPage(page + 1)"
            >
              التالي
              <ChevronLeft class="h-3.5 w-3.5" />
            </button>
          </div>
        </div>
      </StatePanel>
    </div>

    <ConfirmDialog
      :open="!!pendingAction"
      :title="actionMode === 'suspend' ? 'إيقاف حساب الطالب' : 'تفعيل حساب الطالب'"
      :message="actionMode === 'suspend' ? 'سيتم منع الطالب من الوصول إلى المنصة. هل أنت متأكد؟' : 'سيتمكن الطالب من الوصول إلى المنصة مرة أخرى. هل أنت متأكد؟'"
      confirm-label="تأكيد"
      cancel-label="إلغاء"
      variant="danger"
      :loading="actionLoading"
      @confirm="confirmAction"
      @cancel="pendingAction = null"
    />
  </div>
</template>
