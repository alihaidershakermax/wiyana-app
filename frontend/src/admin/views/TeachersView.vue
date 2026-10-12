<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { GraduationCap, Plus, RefreshCw, Power } from '@lucide/vue'
import { useTeachers, useUpdateTeacherStatus } from '@/admin/queries/useTeachers'
import StatePanel from '@/admin/components/StatePanel.vue'
import ConfirmDialog from '@/admin/components/ui/ConfirmDialog.vue'
import { useToast } from '@/admin/composables/useToast'
import type { Teacher } from '@/admin/types/teacher'

const router = useRouter()
const toast = useToast()

const { data, isLoading, isError, error, refetch } = useTeachers()

const updateStatus = useUpdateTeacherStatus()
const pendingAction = ref<Teacher | null>(null)
const actionLoading = ref(false)

function openAction(t: Teacher) {
  pendingAction.value = t
}

function confirmAction() {
  if (!pendingAction.value) return
  actionLoading.value = true
  const teacher = pendingAction.value
  updateStatus.mutate(
    { id: teacher.id, isActive: !teacher.is_active },
    {
      onSuccess: () => {
        toast.success(teacher.is_active ? 'تم إيقاف حساب الأستاذ' : 'تم تفعيل حساب الأستاذ')
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
</script>

<template>
  <div class="space-y-6">
    <div class="flex flex-wrap items-center justify-between gap-3">
      <div>
        <h1 class="text-xl font-extrabold flex items-center gap-2">
          <GraduationCap class="h-5 w-5 text-orange-500" />
          الأساتذة
        </h1>
        <p class="text-xs text-slate-500 dark:text-slate-400 mt-1">
          إدارة حسابات الأساتذة وتخصيص المساقات
        </p>
      </div>
      <div class="flex items-center gap-2">
        <button class="btn-secondary text-sm" @click="() => refetch()">
          <RefreshCw class="h-4 w-4" />
          تحديث
        </button>
        <button class="btn-primary text-sm" @click="router.push({ name: 'admin.teacher-create' })">
          <Plus class="h-4 w-4" />
          إضافة أستاذ
        </button>
      </div>
    </div>

    <!-- Table -->
    <div class="card overflow-hidden">
      <StatePanel
        :loading="isLoading"
        :error="isError ? (error as Error)?.message : null"
        :empty="!data?.length"
        empty-title="لا يوجد أساتذة"
        empty-message="ابدأ بإضافة أول أستاذ"
        :skeleton-count="6"
        @retry="() => refetch()"
      >
        <template #empty-action>
          <button class="btn-primary text-sm" @click="router.push({ name: 'admin.teacher-create' })">
            <Plus class="h-4 w-4" />
            إضافة أستاذ
          </button>
        </template>

        <div class="overflow-x-auto">
          <table class="w-full text-sm">
            <thead>
              <tr class="text-right border-b border-slate-100 dark:border-slate-800 text-xs text-slate-500 dark:text-slate-400">
                <th class="px-4 py-3 font-bold">الأستاذ</th>
                <th class="px-4 py-3 font-bold">يوزر نيم</th>
                <th class="px-4 py-3 font-bold">معرف تيليجرام</th>
                <th class="px-4 py-3 font-bold">التخصص</th>
                <th class="px-4 py-3 font-bold">الحالة</th>
                <th class="px-4 py-3 font-bold">إجراء</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="teacher in data"
                :key="teacher.id"
                class="border-b border-slate-50 dark:border-slate-800/60 hover:bg-slate-50 dark:hover:bg-slate-800/40 transition-colors"
              >
                <td class="px-4 py-3 font-bold text-slate-800 dark:text-slate-200">
                  {{ teacher.first_name || 'مدرس' }}
                </td>
                <td class="px-4 py-3 font-mono text-slate-600 dark:text-slate-300" dir="ltr">
                  {{ teacher.username || '—' }}
                </td>
                <td class="px-4 py-3 font-mono text-slate-600 dark:text-slate-300" dir="ltr">
                  {{ teacher.telegram_id || '—' }}
                </td>
                <td class="px-4 py-3 text-slate-500 dark:text-slate-400">
                  {{ teacher.bio || '—' }}
                </td>
                <td class="px-4 py-3">
                  <span
                    :class="[
                      'inline-flex items-center px-2.5 py-1 rounded-full text-xs font-bold',
                      teacher.is_active
                        ? 'bg-emerald-100 text-emerald-700 dark:bg-emerald-950/40 dark:text-emerald-400'
                        : 'bg-slate-100 text-slate-600 dark:bg-slate-800 dark:text-slate-400',
                    ]"
                  >
                    {{ teacher.is_active ? 'نشط' : 'موقوف' }}
                  </span>
                </td>
                <td class="px-4 py-3">
                  <button
                    class="btn-ghost text-xs py-1.5 px-2"
                    @click="openAction(teacher)"
                  >
                    <Power class="h-3.5 w-3.5" />
                    {{ teacher.is_active ? 'إيقاف' : 'تفعيل' }}
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </StatePanel>
    </div>

    <ConfirmDialog
      :open="!!pendingAction"
      :title="pendingAction?.is_active ? 'إيقاف حساب الأستاذ' : 'تفعيل حساب الأستاذ'"
      :message="pendingAction?.is_active ? 'سيتم منع الأستاذ من الوصول إلى لوحة التحكم. هل أنت متأكد؟' : 'سيتمكن الأستاذ من الوصول إلى لوحة التحكم مرة أخرى. هل أنت متأكد؟'"
      confirm-label="تأكيد"
      cancel-label="إلغاء"
      variant="danger"
      :loading="actionLoading"
      @confirm="confirmAction"
      @cancel="pendingAction = null"
    />
  </div>
</template>
