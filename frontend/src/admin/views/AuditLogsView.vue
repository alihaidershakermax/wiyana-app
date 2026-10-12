<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Shield, RefreshCw } from '@lucide/vue'
import { adminApiClient } from '@/admin/api/adminApi'
import StatePanel from '@/admin/components/StatePanel.vue'
import { formatDate } from '@/utils/format'

interface AuditLog {
  id?: string
  _id?: string
  user_id: string
  action: string
  resource: string
  ip_address?: string
  timestamp: string
}

const logs = ref<AuditLog[]>([])
const loading = ref(false)
const errorMsg = ref<string | null>(null)

async function fetchAuditLogs() {
  loading.value = true
  errorMsg.value = null
  try {
    const { data } = await adminApiClient.get('/audit-logs/')
    logs.value = data.data ?? data
  } catch (err: unknown) {
    errorMsg.value = (err as Error)?.message || 'فشل تحميل سجل الأمان'
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchAuditLogs()
})
</script>

<template>
  <div class="space-y-6 text-[#F5F5F5]">
    <div class="flex flex-wrap items-center justify-between gap-3">
      <div>
        <h1 class="text-xl font-extrabold flex items-center gap-2">
          <Shield class="h-5 w-5 text-[#F97316]" />
          سجل الأمان والتدقيق (Audit Logs)
        </h1>
        <p class="text-xs text-[#A1A1AA] mt-1">
          تسجيل جميع العمليات الحساسة وتغييرات النظام بالمنصة
        </p>
      </div>

      <button class="btn-secondary text-xs" @click="fetchAuditLogs">
        <RefreshCw class="h-4 w-4" />
        تحديث
      </button>
    </div>

    <div class="card overflow-hidden">
      <StatePanel
        :loading="loading"
        :error="errorMsg"
        :empty="!logs.length"
        empty-title="لا يوجد سجل أحداث حتى الآن"
        :skeleton-count="5"
        @retry="fetchAuditLogs"
      >
        <div class="overflow-x-auto">
          <table class="w-full text-right border-collapse text-sm text-[#F5F5F5]">
            <thead>
              <tr class="bg-[#0F0F0F] text-[#A1A1AA] border-b border-[#242424] text-xs font-extrabold">
                <th class="py-3.5 px-4">المستخدم ID</th>
                <th class="py-3.5 px-4">الإجراء</th>
                <th class="py-3.5 px-4">المورد</th>
                <th class="py-3.5 px-4">عنوان IP</th>
                <th class="py-3.5 px-4">التاريخ والوقت</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-[#242424] bg-[#111111]">
              <tr v-for="l in logs" :key="l.id || l._id" class="hover:bg-[#161616] transition-colors">
                <td class="py-3.5 px-4 font-mono text-xs text-white">
                  {{ l.user_id }}
                </td>
                <td class="py-3.5 px-4">
                  <span class="text-[10px] font-extrabold text-[#F97316] bg-[#1F0E03] border border-orange-900/40 px-2.5 py-0.5 rounded-full">
                    {{ l.action }}
                  </span>
                </td>
                <td class="py-3.5 px-4 text-xs font-mono text-[#A1A1AA]">
                  {{ l.resource }}
                </td>
                <td class="py-3.5 px-4 font-mono text-xs text-[#71717A]" dir="ltr">
                  {{ l.ip_address || '—' }}
                </td>
                <td class="py-3.5 px-4 text-xs text-[#71717A]">
                  {{ formatDate(l.timestamp) }}
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </StatePanel>
    </div>
  </div>
</template>
