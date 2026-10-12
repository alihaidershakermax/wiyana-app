<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { CreditCard, RefreshCw, Eye, Trash2, Power } from '@lucide/vue'
import { adminApiClient } from '@/admin/api/adminApi'
import StatePanel from '@/admin/components/StatePanel.vue'
import ActionDropdown, { type MenuItem } from '@/admin/components/ui/ActionDropdown.vue'
import ToastContainer from '@/admin/components/ui/ToastContainer.vue'
import { useToast } from '@/admin/composables/useToast'
import { formatDate } from '@/utils/format'

interface SubscriptionItem {
  id?: string
  _id?: string
  student_id: string
  course_id: string
  code_used?: string
  status: string
  expires_at?: string
  created_at: string
}

const toast = useToast()
const items = ref<SubscriptionItem[]>([])
const loading = ref(false)
const errorMsg = ref<string | null>(null)

async function fetchSubscriptions() {
  loading.value = true
  errorMsg.value = null
  try {
    const { data } = await adminApiClient.get('/panel/subscriptions')
    items.value = data.data ?? data
  } catch (err: unknown) {
    errorMsg.value = (err as Error)?.message || 'فشل تحميل الاشتراكات'
  } finally {
    loading.value = false
  }
}

function getActionItems(s: SubscriptionItem): MenuItem[] {
  return [
    {
      id: 'view',
      label: 'تفاصيل الاشتراك',
      icon: Eye,
      onClick: () => toast.info('تفاصيل الاشتراك', `الاشتراك الخاص بالطالب ${s.student_id}`),
    },
    {
      id: 'cancel',
      label: s.status === 'active' ? 'إلغاء الاشتراك' : 'تجديد الاشتراك',
      icon: Power,
      danger: s.status === 'active',
      onClick: async () => {
        try {
          const subId = s.id || s._id
          const newStatus = s.status === 'active' ? 'cancelled' : 'active'
          await adminApiClient.patch(`/panel/subscriptions/${subId}`, { status: newStatus })
          s.status = newStatus
          toast.success(newStatus === 'cancelled' ? 'تم إلغاء الاشتراك' : 'تم تجديد الاشتراك')
        } catch {
          toast.error('خطأ', 'هذه العملية غير متوفرة حالياً')
        }
      },
    },
    {
      id: 'delete',
      label: 'حذف السجل',
      icon: Trash2,
      danger: true,
      onClick: async () => {
        try {
          const subId = s.id || s._id
          await adminApiClient.delete(`/panel/subscriptions/${subId}`)
          items.value = items.value.filter(i => (i.id || i._id) !== subId)
          toast.success('تم حذف السجل بنجاح')
        } catch {
          toast.error('خطأ', 'هذه العملية غير متوفرة حالياً')
        }
      },
    },
  ]
}

onMounted(() => {
  fetchSubscriptions()
})
</script>

<template>
  <div class="space-y-6 text-[#F5F5F5]">
    <ToastContainer />
    <div class="flex flex-wrap items-center justify-between gap-3">
      <div>
        <h1 class="text-xl font-extrabold flex items-center gap-2">
          <CreditCard class="h-5 w-5 text-[#F97316]" />
          اشتراكات الطلاب
        </h1>
        <p class="text-xs text-[#A1A1AA] mt-1">
          عرض اشتراكات المساقات المفعلة وحالة الصلاحية لكل طالب
        </p>
      </div>

      <button class="btn-secondary text-xs" @click="fetchSubscriptions">
        <RefreshCw class="h-4 w-4" />
        تحديث
      </button>
    </div>

    <div class="card overflow-hidden">
      <StatePanel
        :loading="loading"
        :error="errorMsg"
        :empty="!items.length"
        empty-title="لا توجد اشتراكات مسجلة"
        empty-message="سيظهر هنا سجل اشتراكات الطلاب بمجرد تفعيل أكواد المساقات"
        :skeleton-count="5"
        @retry="fetchSubscriptions"
      >
        <div class="overflow-x-auto">
          <table class="w-full text-right border-collapse text-sm text-[#F5F5F5]">
            <thead>
              <tr class="bg-[#0F0F0F] text-[#A1A1AA] border-b border-[#242424] text-xs font-extrabold">
                <th class="py-3.5 px-4">الطالب ID</th>
                <th class="py-3.5 px-4">المساق ID</th>
                <th class="py-3.5 px-4">الكود المستخدم</th>
                <th class="py-3.5 px-4">الحالة</th>
                <th class="py-3.5 px-4">تاريخ الانتهاء</th>
                <th class="py-3.5 px-4">تاريخ التفعيل</th>
                <th class="py-3.5 px-4 text-center w-16">إجراءات</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-[#242424] bg-[#111111]">
              <tr v-for="s in items" :key="s.id || s._id" class="hover:bg-[#161616] transition-colors">
                <td class="py-3.5 px-4 font-mono text-xs text-white">
                  {{ s.student_id }}
                </td>
                <td class="py-3.5 px-4 font-mono text-xs text-[#A1A1AA]">
                  {{ s.course_id }}
                </td>
                <td class="py-3.5 px-4 font-mono text-xs font-bold text-[#F97316]">
                  {{ s.code_used || 'مباشر' }}
                </td>
                <td class="py-3.5 px-4">
                  <span
                    :class="[
                      'text-[10px] font-extrabold px-2.5 py-1 rounded-full border inline-block',
                      s.status === 'active' ? 'bg-[#052e16] text-[#22C55E] border-[#14532d]' : 'bg-[#450a0a] text-[#EF4444] border-[#7f1d1d]'
                    ]"
                  >
                    {{ s.status === 'active' ? 'نشط' : 'منتهي' }}
                  </span>
                </td>
                <td class="py-3.5 px-4 text-xs text-[#71717A]">
                  {{ formatDate(s.expires_at) ?? 'غير محدد' }}
                </td>
                <td class="py-3.5 px-4 text-xs text-[#71717A]">
                  {{ formatDate(s.created_at) }}
                </td>
                <td class="py-3.5 px-4 text-center">
                  <ActionDropdown :items="getActionItems(s)" />
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </StatePanel>
    </div>
  </div>
</template>
